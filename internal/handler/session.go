package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xuri/excelize/v2"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/middleware"
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
	if !requirePermission(w, r, permFirstWeight) {
		return
	}

	var request dto.CreateSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	if !allowManualWeight(w, r, request.ManualWeightKg) {
		return
	}

	session, err := h.service.Create(r.Context(), request, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// Flow: POST /spa/api/weighing/sessions/flow
func (h *SessionHandler) Flow(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permFlowWeight) {
		return
	}

	var request dto.FlowSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	if !allowManualWeight(w, r, request.ManualWeightKg) {
		return
	}

	session, err := h.service.CreateFlow(r.Context(), request, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// Close: POST /spa/api/weighing/sessions/{id}/second
func (h *SessionHandler) Close(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSecondWeight) {
		return
	}

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

	session, err := h.service.Close(r.Context(), id, request.Measurement, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// List: GET /spa/api/weighing/sessions?status=open&page=1&pageSize=25
func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

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

	result, err := h.service.List(
		r.Context(),
		status,
		search,
		dateFrom,
		dateTo,
		r.URL.Query().Get("sort"),
		r.URL.Query().Get("dir"),
		page,
		pageSize,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Get: GET /spa/api/weighing/sessions/{id}
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

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

// Export: GET /spa/api/weighing/sessions/export?status=&search=&dateFrom=&dateTo=
func (h *SessionHandler) Export(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permExport) {
		return
	}

	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

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

	items, err := h.service.Export(r.Context(), status, search, dateFrom, dateTo)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="weighing-journal-%s.csv"`,
		time.Now().Format("20060102-150405"),
	))
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	_ = writer.Write([]string{
		"№", "Статус", "Транспорт", "Прицеп", "Водитель", "Пропуск",
		"Груз", "Контрагент", "Перевозчик", "Склад", "Тип операции",
		"Брутто, т", "Тара, т", "Нетто, т",
		"Первое взвешивание", "Второе взвешивание", "Создан",
	})

	for _, item := range items {
		_ = writer.Write([]string{
			strconv.FormatInt(item.ID, 10),
			statusLabel(item.Status),
			item.TransportNumber,
			item.TrailerNumber,
			item.DriverName,
			item.PassNumber,
			item.CargoName,
			item.CounterpartyName,
			item.CarrierName,
			item.StorageName,
			item.OperationTypeName,
			csvWeight(&item.FirstWeightKg),
			csvWeight(item.SecondWeightKg),
			csvWeight(item.NetWeightKg),
			item.FirstMeasuredAt.Format("02.01.2006 15:04"),
			formatCSVTime(item.SecondMeasuredAt),
			item.CreatedAt.Format("02.01.2006 15:04"),
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		slog.Error("export csv", "error", err)
	}
}

// ExportXLSX: GET /spa/api/weighing/sessions/export.xlsx?status=&search=&dateFrom=&dateTo=
func (h *SessionHandler) ExportXLSX(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permExport) {
		return
	}

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

	items, err := h.service.Export(
		r.Context(),
		r.URL.Query().Get("status"),
		r.URL.Query().Get("search"),
		dateFrom,
		dateTo,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	file := excelize.NewFile()
	defer func() { _ = file.Close() }()

	sheet := "Журнал"
	if _, err := file.NewSheet(sheet); err != nil {
		h.handleError(w, err)
		return
	}
	_ = file.DeleteSheet("Sheet1")

	headers := []string{
		"№", "Статус", "Транспорт", "Прицеп", "Водитель", "Пропуск",
		"Груз", "Контрагент", "Перевозчик", "Склад", "Тип операции",
		"Брутто, т", "Тара, т", "Нетто, т",
		"Первое взвешивание", "Второе взвешивание", "Создан",
	}

	headerStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E8EEF7"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	weightFormat := "#,##0.000"
	weightStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &weightFormat,
	})

	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, 1)
		_ = file.SetCellValue(sheet, cell, header)
	}
	lastHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	_ = file.SetCellStyle(sheet, "A1", lastHeaderCell, headerStyle)

	for index, item := range items {
		row := index + 2
		values := []any{
			item.ID,
			statusLabel(item.Status),
			item.TransportNumber,
			item.TrailerNumber,
			item.DriverName,
			item.PassNumber,
			item.CargoName,
			item.CounterpartyName,
			item.CarrierName,
			item.StorageName,
			item.OperationTypeName,
			weightValue(item.FirstWeightKg),
			weightValuePtr(item.SecondWeightKg),
			weightValuePtr(item.NetWeightKg),
			item.FirstMeasuredAt.Format("02.01.2006 15:04"),
			formatCSVTime(item.SecondMeasuredAt),
			item.CreatedAt.Format("02.01.2006 15:04"),
		}

		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			if value != nil {
				_ = file.SetCellValue(sheet, cell, value)
			}
		}

		for column := 12; column <= 14; column++ {
			cell, _ := excelize.CoordinatesToCellName(column, row)
			_ = file.SetCellStyle(sheet, cell, cell, weightStyle)
		}
	}

	widths := map[string]float64{
		"A": 6, "B": 12, "C": 16, "D": 14, "E": 22, "F": 16,
		"G": 22, "H": 24, "I": 20, "J": 20, "K": 16,
		"L": 12, "M": 12, "N": 12, "O": 18, "P": 18, "Q": 18,
	}
	for column, width := range widths {
		_ = file.SetColWidth(sheet, column, column, width)
	}

	_ = file.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="weighing-journal-%s.xlsx"`,
		time.Now().Format("20060102-150405"),
	))

	if err := file.Write(w); err != nil {
		slog.Error("export xlsx", "error", err)
	}
}

func weightValue(value float64) any {
	return value
}

func weightValuePtr(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func statusLabel(status string) string {
	switch status {
	case "open":
		return "Открыт"
	case "closed":
		return "Закрыт"
	case "cancelled":
		return "Отменён"
	default:
		return status
	}
}

func csvWeight(value *float64) string {
	if value == nil {
		return ""
	}
	return strings.ReplaceAll(strconv.FormatFloat(*value, 'f', 3, 64), ".", ",")
}

func formatCSVTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("02.01.2006 15:04")
}

// Bulk: POST /spa/api/weighing/sessions/bulk
func (h *SessionHandler) Bulk(w http.ResponseWriter, r *http.Request) {
	var request dto.BulkSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	switch request.Action {
	case "cancel":
		if !requirePermission(w, r, permCancelWeight) {
			return
		}
	case "restore":
		if !requirePermission(w, r, permEdit) {
			return
		}
	default:
		writeError(w, http.StatusBadRequest, "invalid_action", "Действие должно быть cancel или restore")
		return
	}

	result, err := h.service.Bulk(r.Context(), request.Action, request.IDs, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Restore: POST /spa/api/weighing/sessions/{id}/restore
func (h *SessionHandler) Restore(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permEdit) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	if err := h.service.Restore(r.Context(), id, currentUserID(r)); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Update: PATCH /spa/api/weighing/sessions/{id}
func (h *SessionHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permEdit) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	var request dto.UpdateSessionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	session, err := h.service.Update(r.Context(), id, request, currentUserID(r))
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// Audit: GET /spa/api/weighing/sessions/{id}/audit
func (h *SessionHandler) Audit(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	items, err := h.service.Audit(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Permissions: GET /spa/api/weighing/permissions — права текущего пользователя.
func (h *SessionHandler) Permissions(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Пользователь не определён")
		return
	}

	codes := make([]string, 0, len(user.Permissions))
	for code := range user.Permissions {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	writeJSON(w, http.StatusOK, map[string]any{
		"isRoot":      user.IsRoot,
		"permissions": codes,
	})
}

// AuditExport: GET /spa/api/weighing/audit/export.csv
func (h *SessionHandler) AuditExport(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

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

	items, err := h.service.AuditExport(
		r.Context(),
		r.URL.Query().Get("action"),
		r.URL.Query().Get("search"),
		dateFrom,
		dateTo,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="audit-%s.csv"`, time.Now().Format("20060102-150405"),
	))
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	writer.Comma = ';'
	_ = writer.Write([]string{
		"ID", "Дата", "Рейс", "Транспорт", "Пользователь", "Действие", "Изменения",
	})
	for _, item := range items {
		_ = writer.Write([]string{
			strconv.FormatInt(item.ID, 10),
			item.CreatedAt.Format("02.01.2006 15:04"),
			strconv.FormatInt(item.SessionID, 10),
			item.TransportNumber,
			item.UserName,
			item.Action,
			string(item.Payload),
		})
	}
	writer.Flush()
}

// AuditJournal: GET /spa/api/weighing/audit
func (h *SessionHandler) AuditJournal(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

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

	result, err := h.service.AuditJournal(
		r.Context(),
		r.URL.Query().Get("action"),
		r.URL.Query().Get("search"),
		dateFrom,
		dateTo,
		parseQueryInt(r, "page", 1),
		parseQueryInt(r, "pageSize", 25),
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Measurements: GET /spa/api/weighing/sessions/{id}/measurements
func (h *SessionHandler) Measurements(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

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
	if !requirePermission(w, r, permCancelWeight) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	if err := h.service.Cancel(r.Context(), id, currentUserID(r)); err != nil {
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
