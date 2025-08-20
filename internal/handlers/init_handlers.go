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

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
