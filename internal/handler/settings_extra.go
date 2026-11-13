package handler

import (
	"errors"
	"net/http"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

// SettingsExtraHandler — камеры, лицензии и аудит настроек.
type SettingsExtraHandler struct {
	service *service.SettingsExtraService
}

func NewSettingsExtraHandler(svc *service.SettingsExtraService) *SettingsExtraHandler {
	return &SettingsExtraHandler{service: svc}
}

func (h *SettingsExtraHandler) Cameras(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	items, err := h.service.Cameras(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить камеры")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) SaveCamera(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	id, err := parseNestedID(r, "id")
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректная камера")
		return
	}
	var payload dto.SaveCameraRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	items, err := h.service.SaveCamera(r.Context(), id, payload, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) Licenses(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	items, err := h.service.Licenses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить лицензии")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) CreateLicense(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	var payload dto.LicenseRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	items, err := h.service.CreateLicense(r.Context(), payload, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) UpdateLicense(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	id, err := parseNestedID(r, "id")
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректная лицензия")
		return
	}
	var payload dto.LicenseRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	items, err := h.service.UpdateLicense(r.Context(), id, payload, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) DeleteLicense(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	id, err := parseNestedID(r, "id")
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректная лицензия")
		return
	}
	items, err := h.service.DeleteLicense(r.Context(), id, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *SettingsExtraHandler) SettingsAudit(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}
	response, err := h.service.Audit(
		r.Context(), parseQueryInt(r, "page", 1), parseQueryInt(r, "pageSize", 25),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить аудит")
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *SettingsExtraHandler) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Запись не найдена")
		return
	}
	if validation, ok := err.(*service.ValidationError); ok {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка")
}
