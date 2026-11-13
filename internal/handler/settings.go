package handler

import (
	"net/http"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/service"
)

// SettingsHandler — бизнес-настройки полей рейса.
type SettingsHandler struct {
	service *service.SettingsService
	extra   *service.SettingsExtraService
}

func NewSettingsHandler(svc *service.SettingsService, extra *service.SettingsExtraService) *SettingsHandler {
	return &SettingsHandler{service: svc, extra: extra}
}

// GetApp: GET /spa/api/weighing/settings/app
func (h *SettingsHandler) GetApp(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	settings, err := h.service.App(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить настройки")
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

// SaveApp: PUT /spa/api/weighing/settings/app
func (h *SettingsHandler) SaveApp(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	var payload dto.AppSettings
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	settings, err := h.service.SaveApp(r.Context(), payload, currentUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить настройки")
		return
	}

	if h.extra != nil {
		h.extra.Log(r.Context(), "app", 1, "updated", currentUserID(r), payload)
	}

	writeJSON(w, http.StatusOK, settings)
}
