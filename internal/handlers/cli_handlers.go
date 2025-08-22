package handlers

import (
	"RSSHub/internal/apperrors"
	"RSSHub/internal/models"
	"RSSHub/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.Command
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.JsonResponse(w, 401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	err := h.cliService.AddService(ctx, data)
	if err != nil {
		errCode := apperrors.CheckError(err)

		utils.JsonResponse(w, errCode, map[string]string{
			"error": err.Error(),
		})

		return
	}

	utils.JsonResponse(w, 200, map[string]string{
		"status":  "ok",
		"message": "Added successfully!",
	})
}
