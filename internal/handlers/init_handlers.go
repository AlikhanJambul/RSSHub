package handlers

import (
	"RSSHub/internal/aggregator"
	"RSSHub/internal/logger"
	"RSSHub/internal/service"
	"RSSHub/internal/storage"
	"net/http"
)

type Handler struct {
	cliRepo       storage.CLIRepo
	cliService    service.CLIService
	cliAggregator aggregator.Aggregator
	cliLogger     logger.Logger
}

func NewHandler(cliRepo storage.CLIRepo, cliService service.CLIService, cliAggregator aggregator.Aggregator, cliLogger logger.Logger) *Handler {
	return &Handler{cliRepo: cliRepo, cliService: cliService, cliAggregator: cliAggregator, cliLogger: cliLogger}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", h.Add)
	mux.HandleFunc("PUT /set-worker", h.SetWorkersCount)
	mux.HandleFunc("PUT /set-interval", h.SetInterval)

	return mux
}
