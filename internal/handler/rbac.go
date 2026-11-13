package handler

import (
	"errors"
	"net/http"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

// RBACHandler — матрица прав групп.
type RBACHandler struct {
	service *service.RBACService
}

func NewRBACHandler(svc *service.RBACService) *RBACHandler {
	return &RBACHandler{service: svc}
}

// Matrix: GET /spa/api/weighing/rbac
func (h *RBACHandler) Matrix(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	matrix, err := h.service.Matrix(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить права")
		return
	}

	writeJSON(w, http.StatusOK, matrix)
}

// SaveGroup: PUT /spa/api/weighing/rbac/groups/{id}
func (h *RBACHandler) SaveGroup(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	id, err := parseNestedID(r, "id")
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректная группа")
		return
	}

	var payload dto.SaveGroupPermissionsRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	matrix, err := h.service.SaveGroup(r.Context(), id, payload.Permissions, currentUserID(r))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Группа не найдена")
			return
		}
		if validation, ok := err.(*service.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить права")
		return
	}

	writeJSON(w, http.StatusOK, matrix)
}
