package cli

import (
	"RSSHub/internal/config"
	"RSSHub/internal/handlers"
	"RSSHub/internal/service"
	"RSSHub/internal/storage"
	"log"
	"net/http"
)

func StartServer() {
	cfg := config.Load()

	db, err := storage.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	cliRepo := storage.NewRepo(db)
	cliService := service.NewService(cliRepo)
	cliHandler := handlers.NewHandler(cliRepo, cliService)

	svr := http.Server{
		Addr:    ":8080",
		Handler: cliHandler.Router(),
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
