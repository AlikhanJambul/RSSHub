package cli

import (
	"RSSHub/internal/adapter/handlers"
	"RSSHub/internal/adapter/postgres"
	"RSSHub/internal/adapter/rss"
	"RSSHub/internal/app"
	"RSSHub/internal/config"
	"RSSHub/internal/logger"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
	if !checkPort() {
		fmt.Fprintln(os.Stderr, "Background process is already running")
		os.Exit(1)
	}

	cfg := config.Load()
	interval, err := time.ParseDuration(cfg.TimerInterval)
	if err != nil {
		log.Fatal(err)
	}

	cliLogger := logger.New()

	db, err := postgres.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	cliParser := rss.NewParser()
	cliRepo := postgres.NewRepo(db)
	cliService := app.NewService(cliRepo)
	cliAggregator := app.InitAggregator(cfg.WorkerCount, interval, cliRepo, *cliLogger, cliParser)
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

	cliLogger.Info(fmt.Sprintf("The background process for fetching feeds has started (interval = %s, workers = %d)"), cfg.TimerInterval, cfg.WorkerCount)

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

func checkPort() bool {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return false
	} else {
		ln.Close()
		return true
	}
}
