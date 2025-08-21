package cli

import (
	"RSSHub/internal/aggregator"
	"RSSHub/internal/config"
	"RSSHub/internal/handlers"
	"RSSHub/internal/logger"
	"RSSHub/internal/service"
	"RSSHub/internal/storage"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
	cfg := config.Load()
	interval, err := time.ParseDuration(cfg.TimerInterval)
	if err != nil {
		log.Fatal(err)
	}

	cliLogger := logger.New()

	db, err := storage.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	cliRepo := storage.NewRepo(db)
	cliService := service.NewService(cliRepo)
	cliAggregator := aggregator.InitAggregator(cfg.WorkerCount, interval, cliRepo, *cliLogger)
	cliHandler := handlers.NewHandler(cliRepo, cliService, cliAggregator, *cliLogger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	svr := http.Server{
		Addr:    ":8080",
		Handler: cliHandler.Router(),
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	cliLogger.Info("Server started")

	go cliAggregator.Start()

	<-stop
	cliLogger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	cliAggregator.Stop()

	cliLogger.Info("Server exited properly")
}
