package handler

import (
	"errors"
	"net/http"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

// StationHandler — весовые и терминалы.
type StationHandler struct {
	service *service.StationService
}

func NewStationHandler(svc *service.StationService) *StationHandler {
	return &StationHandler{service: svc}
}

// List: GET /spa/api/weighing/stations
func (h *StationHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	items, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить станции")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Save: PUT /spa/api/weighing/stations/{id}
func (h *StationHandler) Save(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	id, err := parseNestedID(r, "id")
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректная станция")
		return
	}

	var payload dto.SaveStationRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	items, err := h.service.Save(r.Context(), id, payload, currentUserID(r))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Станция не найдена")
			return
		}
		if validation, ok := err.(*service.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить станцию")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
