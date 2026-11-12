package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

type SessionHandler struct {
	service *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{service: svc}
}

// Create: POST /spa/api/weighing/sessions
func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	session, err := h.service.Create(r.Context(), request)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// Flow: POST /spa/api/weighing/sessions/flow
func (h *SessionHandler) Flow(w http.ResponseWriter, r *http.Request) {
	var request dto.FlowSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	session, err := h.service.CreateFlow(r.Context(), request)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// Close: POST /spa/api/weighing/sessions/{id}/second
func (h *SessionHandler) Close(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	var request dto.CloseSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	session, err := h.service.Close(r.Context(), id, request.Measurement)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// List: GET /spa/api/weighing/sessions?status=open&page=1&pageSize=25
func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "pageSize", 25)

	dateFrom, err := parseDateParam(r.URL.Query().Get("dateFrom"), false)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}
	dateTo, err := parseDateParam(r.URL.Query().Get("dateTo"), true)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}

	result, err := h.service.List(r.Context(), status, search, dateFrom, dateTo, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Get: GET /spa/api/weighing/sessions/{id}
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	session, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// Measurements: GET /spa/api/weighing/sessions/{id}/measurements
func (h *SessionHandler) Measurements(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	items, err := h.service.Measurements(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Cancel: POST /spa/api/weighing/sessions/{id}/cancel
func (h *SessionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	if err := h.service.Cancel(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHandler) handleError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Рейс не найден")
		return
	}

	if errors.Is(err, repository.ErrConflict) {
		writeError(w, http.StatusConflict, "conflict", err.Error())
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		writeError(w, http.StatusBadRequest, "invalid_reference", "Выбранный элемент справочника не найден")
		return
	}

	slog.Error("session handler", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return errors.New("Некорректный JSON: " + err.Error())
	}

	return nil
}

func parseDateParam(value string, endOfDay bool) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		return nil, errors.New("Дата должна быть в формате ГГГГ-ММ-ДД")
	}

	if endOfDay {
		parsed = parsed.Add(24 * time.Hour)
	}

	return &parsed, nil
}
