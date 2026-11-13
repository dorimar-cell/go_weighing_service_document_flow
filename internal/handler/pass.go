package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

type PassHandler struct {
	service *service.PassService
}

func NewPassHandler(svc *service.PassService) *PassHandler {
	return &PassHandler{service: svc}
}

// Create: POST /spa/api/weighing/passes
func (h *PassHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permPasses) {
		return
	}

	var request dto.CreatePassRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	pass, err := h.service.Create(r.Context(), request, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, pass)
}

// List: GET /spa/api/weighing/passes?status=active&search=&page=1&pageSize=25
func (h *PassHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permPasses) {
		return
	}

	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "pageSize", 25)

	result, err := h.service.List(r.Context(), status, search, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Get: GET /spa/api/weighing/passes/{id}
func (h *PassHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permPasses) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор пропуска")
		return
	}

	pass, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pass)
}

// SetEnabled: PATCH /spa/api/weighing/passes/{id}/enabled
func (h *PassHandler) SetEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permPasses) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор пропуска")
		return
	}

	var request dto.UpdatePassEnabledRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	pass, err := h.service.SetEnabled(r.Context(), id, request.Enabled, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pass)
}

func (h *PassHandler) handleError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Пропуск не найден")
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		writeError(w, http.StatusBadRequest, "invalid_reference", "Выбранный элемент справочника не найден")
		return
	}

	slog.Error("pass handler", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
}
