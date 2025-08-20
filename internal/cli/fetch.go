package cli

import (
	"RSSHub/internal/aggregator"
	"RSSHub/internal/config"
	"RSSHub/internal/handlers"
	"RSSHub/internal/logger"
	"RSSHub/internal/service"
	"RSSHub/internal/storage"
	"log"
	"net/http"
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

	go func() {
		svr := http.Server{
			Addr:    ":8080",
			Handler: cliHandler.Router(),
		}

		if err := svr.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	cliAggregator.Start()

}
