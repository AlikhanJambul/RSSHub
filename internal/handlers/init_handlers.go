package handlers

import (
	"RSSHub/internal/service"
	"RSSHub/internal/storage"
	"net/http"
)

type Handler struct {
	cliRepo    storage.CLIRepo
	cliService service.CLIService
}

func NewHandler(cliRepo storage.CLIRepo, cliService service.CLIService) *Handler {
	return &Handler{cliRepo: cliRepo, cliService: cliService}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
