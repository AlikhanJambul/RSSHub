package handlers

import (
	"RSSHub/internal/apperrors"
	"RSSHub/internal/domain"
	"RSSHub/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data domain.Command
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

		h.cliLogger.Error(err.Error())

		utils.JsonResponse(w, errCode, map[string]string{
			"error": err.Error(),
		})

		return
	}

	h.cliLogger.Info("Succes")

	utils.JsonResponse(w, 200, map[string]string{
		"status":  "ok",
		"message": "Added successfully!",
	})
}

func (h *Handler) SetWorkersCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data domain.Command
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.JsonResponse(w, 401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err := h.cliAggregator.ChangeCountWorker(data.Workers)
	if err != nil {
		errCode := apperrors.CheckError(err)
		h.cliLogger.Error(err.Error())

		utils.JsonResponse(w, errCode, map[string]string{
			"error": err.Error(),
		})

		return
	}

	h.cliLogger.Info("Succes")
	utils.JsonResponse(w, 200, map[string]string{})
}

func (h *Handler) SetInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data domain.Command
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.JsonResponse(w, 401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err := h.cliAggregator.ChangeInterval(data.Interval)
	if err != nil {
		errCode := apperrors.CheckError(err)

		h.cliLogger.Error(err.Error())

		utils.JsonResponse(w, errCode, map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.cliLogger.Info("Succes")
	utils.JsonResponse(w, 200, map[string]string{})
}
