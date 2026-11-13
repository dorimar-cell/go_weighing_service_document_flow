package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

// ReportHandler — каталог отчётов, печать бланков и реквизиты документов.
type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{service: svc}
}

// Catalog: GET /spa/api/weighing/reports
func (h *ReportHandler) Catalog(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	writeJSON(w, http.StatusOK, dto.ReportCatalogResponse{Items: h.service.Catalog()})
}

// Run: GET /spa/api/weighing/reports/{code}/run?format=html|xlsx&...
func (h *ReportHandler) Run(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	code := strings.TrimSpace(chi.URLParam(r, "code"))
	file, err := h.service.Run(r.Context(), code, r.URL.Query())
	if err != nil {
		h.handleError(w, err)
		return
	}

	disposition := "inline"
	if strings.HasSuffix(file.Name, ".xlsx") {
		disposition = "attachment"
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", disposition+`; filename="`+file.Name+`"`)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(file.Data)
}

// GetSettings: GET /spa/api/weighing/reports/settings
func (h *ReportHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	settings, err := h.service.Settings(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

// SaveSettings: PUT /spa/api/weighing/reports/settings
func (h *ReportHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permExport) {
		return
	}

	var payload dto.ReportSettings
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	settings, err := h.service.SaveSettings(r.Context(), payload)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

// GetDocument: GET /spa/api/weighing/sessions/{id}/documents/{kind}
func (h *ReportHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	document, err := h.service.Document(r.Context(), id, strings.TrimSpace(chi.URLParam(r, "kind")))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, document)
}

// SaveDocument: PUT /spa/api/weighing/sessions/{id}/documents/{kind}
func (h *ReportHandler) SaveDocument(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	var raw json.RawMessage
	if err := decodeJSON(w, r, &raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	document, err := h.service.SaveDocument(
		r.Context(),
		id,
		strings.TrimSpace(chi.URLParam(r, "kind")),
		raw,
		currentUserID(r),
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, document)
}

func (h *ReportHandler) handleError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Рейс не найден")
		return
	}

	slog.Error("report handler", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
}
