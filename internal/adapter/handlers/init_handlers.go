package handlers

import (
	"RSSHub/internal/adapter/postgres"
	"RSSHub/internal/app"
	"RSSHub/internal/logger"
	"net/http"
)

type Handler struct {
	cliRepo       postgres.CLIRepo
	cliService    app.CLIService
	cliAggregator app.Aggregator
	cliLogger     logger.Logger
}

func NewHandler(cliRepo postgres.CLIRepo, cliService app.CLIService, cliAggregator app.Aggregator, cliLogger logger.Logger) *Handler {
	return &Handler{cliRepo: cliRepo, cliService: cliService, cliAggregator: cliAggregator, cliLogger: cliLogger}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", h.Add)
	mux.HandleFunc("PUT /set-worker", h.SetWorkersCount)
	mux.HandleFunc("PUT /set-interval", h.SetInterval)
	mux.HandleFunc("DELETE /delete", h.Delete)
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
