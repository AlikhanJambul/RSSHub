package handlers

import (
	"RSSHub/internal/models"
	"encoding/json"
	"net/http"
)

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.Command
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: service logic

}
