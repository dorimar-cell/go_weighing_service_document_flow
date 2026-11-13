package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

const maxBodyBytes = 1 << 20 // 1 MiB

type DictionaryHandler struct {
	service *service.DictionaryService
}

func NewDictionaryHandler(svc *service.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{service: svc}
}

// List: GET /spa/api/weighing/dictionaries/{code}?page=1&pageSize=25&search=
func (h *DictionaryHandler) List(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "pageSize", 25)
	search := r.URL.Query().Get("search")

	result, err := h.service.List(r.Context(), code, page, pageSize, search)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Create: POST /spa/api/weighing/dictionaries/{code}
func (h *DictionaryHandler) Create(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	raw, err := decodePayload(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	item, err := h.service.Create(r.Context(), code, raw)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

// Update: PATCH /spa/api/weighing/dictionaries/{code}/{id}
func (h *DictionaryHandler) Update(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор")
		return
	}

	raw, err := decodePayload(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	item, err := h.service.Update(r.Context(), code, id, raw)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// Delete: DELETE /spa/api/weighing/dictionaries/{code}/{id}
func (h *DictionaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор")
		return
	}

	if err := h.service.Delete(r.Context(), code, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Restore: POST /spa/api/weighing/dictionaries/{code}/{id}/restore
func (h *DictionaryHandler) Restore(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор")
		return
	}

	if err := h.service.Restore(r.Context(), code, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Export: GET /spa/api/weighing/dictionaries/{code}/export.csv?search=
func (h *DictionaryHandler) Export(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permExport) {
		return
	}

	code := chi.URLParam(r, "code")
	export, err := h.service.Export(r.Context(), code, r.URL.Query().Get("search"))
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="%s-%s.csv"`, code, time.Now().Format("20060102-150405"),
	))
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	headers := make([]string, 0, len(export.Fields)+1)
	headers = append(headers, "id")
	for _, field := range export.Fields {
		headers = append(headers, field.Key)
	}
	_ = writer.Write(headers)

	for _, item := range export.Items {
		row := make([]string, 0, len(headers))
		row = append(row, fmt.Sprint(item["id"]))
		for _, field := range export.Fields {
			value := item[field.Key]
			if value == nil {
				row = append(row, "")
				continue
			}
			row = append(row, fmt.Sprint(value))
		}
		_ = writer.Write(row)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		slog.Error("export dictionary csv", "code", code, "error", err)
	}
}

func (h *DictionaryHandler) handleError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Запись не найдена")
		return
	}

	slog.Error("dictionary handler", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
}

func decodePayload(w http.ResponseWriter, r *http.Request) (map[string]json.RawMessage, error) {
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	var raw map[string]json.RawMessage
	if err := decoder.Decode(&raw); err != nil {
		return nil, errors.New("Некорректный JSON")
	}

	return raw, nil
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func parseQueryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return fallback
	}
	return value
}
