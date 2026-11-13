package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/reports"
	"go_weighing_service/internal/repository"
)

const maxInlinePhotos = 4

// ReportFile — готовый файл отчёта для отдачи в HTTP.
type ReportFile struct {
	Name        string
	ContentType string
	Data        []byte
}

// ReportService — каталог отчётов, рендер бланков и реквизиты организации.
type ReportService struct {
	sessions *repository.SessionRepository
	docs     *repository.DocumentRepository
	photos   *repository.PhotoRepository
	photoDir string
}

func NewReportService(
	sessions *repository.SessionRepository,
	docs *repository.DocumentRepository,
	photos *repository.PhotoRepository,
	photoDir string,
) *ReportService {
	return &ReportService{sessions: sessions, docs: docs, photos: photos, photoDir: photoDir}
}

// Catalog возвращает доступные отчёты.
func (s *ReportService) Catalog() []dto.ReportDefinition {
	return reports.Catalog()
}

// Run формирует отчёт по коду и параметрам.
func (s *ReportService) Run(ctx context.Context, code string, query url.Values) (*ReportFile, error) {
	if _, ok := reports.Definition(code); !ok {
		return nil, &ValidationError{Message: "Отчёт не найден"}
	}

	format := strings.ToLower(strings.TrimSpace(query.Get("format")))
	if format == "" {
		format = "html"
	}

	switch code {
	case "journal":
		return s.runJournal(ctx, query, format)
	case "ts-card", "talon", "ttn", "torg12":
		return s.runDocument(ctx, code, query, format)
	}

	if _, ok := aggregateDimensions[code]; ok {
		return s.runAggregate(ctx, code, query, format)
	}

	return nil, &ValidationError{Message: "Отчёт не найден"}
}

type aggregateDimension struct {
	Dimension    string
	Title        string
	Subtitle     string
	ShowSubGroup bool
}

var aggregateDimensions = map[string]aggregateDimension{
	"cargo-summary": {
		Dimension: "cargo",
		Title:     "Отчёт по грузам",
		Subtitle:  "Нетто и количество закрытых рейсов по грузам",
	},
	"counterparty-summary": {
		Dimension:    "counterparty",
		Title:        "Отчёт по контрагентам",
		Subtitle:     "Нетто по контрагентам с разбивкой по грузам",
		ShowSubGroup: true,
	},
	"carrier-summary": {
		Dimension:    "carrier",
		Title:        "Отчёт по перевозчикам",
		Subtitle:     "Нетто по перевозчикам с разбивкой по грузам",
		ShowSubGroup: true,
	},
	"driver-summary": {
		Dimension:    "driver",
		Title:        "Отчёт по водителям",
		Subtitle:     "Нетто по водителям с разбивкой по грузам",
		ShowSubGroup: true,
	},
	"storage-summary": {
		Dimension:    "storage",
		Title:        "Взвешивания по складам",
		Subtitle:     "Нетто по складам с разбивкой по типу операции",
		ShowSubGroup: true,
	},
}

func (s *ReportService) runAggregate(
	ctx context.Context,
	code string,
	query url.Values,
	format string,
) (*ReportFile, error) {
	meta := aggregateDimensions[code]

	if format != "html" && format != "xlsx" {
		return nil, &ValidationError{Message: "Агрегатный отчёт поддерживает форматы HTML и XLSX"}
	}

	dateFrom, err := parseReportDate(query.Get("dateFrom"), false)
	if err != nil {
		return nil, err
	}
	dateTo, err := parseReportDate(query.Get("dateTo"), true)
	if err != nil {
		return nil, err
	}

	items, err := s.sessions.Aggregate(ctx, meta.Dimension, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}

	filter := reports.AggregateFilter{
		Title:        meta.Title,
		Subtitle:     meta.Subtitle,
		ShowSubGroup: meta.ShowSubGroup,
		DateFrom:     strings.TrimSpace(query.Get("dateFrom")),
		DateTo:       strings.TrimSpace(query.Get("dateTo")),
		Organization: settings.OrganizationName,
	}

	stamp := time.Now().Format("20060102-150405")
	if format == "xlsx" {
		data, err := reports.RenderAggregateXLSX(items, filter)
		if err != nil {
			return nil, err
		}
		return &ReportFile{
			Name:        code + "-" + stamp + ".xlsx",
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Data:        data,
		}, nil
	}

	return &ReportFile{
		Name:        code + "-" + stamp + ".html",
		ContentType: "text/html; charset=utf-8",
		Data:        reports.RenderAggregate(items, filter),
	}, nil
}

func (s *ReportService) runJournal(ctx context.Context, query url.Values, format string) (*ReportFile, error) {
	if format != "html" && format != "xlsx" {
		return nil, &ValidationError{Message: "Журнал поддерживает форматы HTML и XLSX"}
	}

	status := strings.ToLower(strings.TrimSpace(query.Get("status")))
	switch status {
	case "", "all", "open", "closed", "cancelled":
	default:
		return nil, &ValidationError{Message: "Неизвестный статус: " + status}
	}
	if status == "" {
		status = "all"
	}

	dateFrom, err := parseReportDate(query.Get("dateFrom"), false)
	if err != nil {
		return nil, err
	}
	dateTo, err := parseReportDate(query.Get("dateTo"), true)
	if err != nil {
		return nil, err
	}

	items, err := s.sessions.Export(ctx, status, strings.TrimSpace(query.Get("search")), dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}

	filter := reports.JournalFilter{
		DateFrom:     strings.TrimSpace(query.Get("dateFrom")),
		DateTo:       strings.TrimSpace(query.Get("dateTo")),
		StatusLabel:  reportStatusLabel(status),
		Search:       strings.TrimSpace(query.Get("search")),
		Organization: settings.OrganizationName,
	}

	stamp := time.Now().Format("20060102-150405")
	if format == "xlsx" {
		data, err := reports.RenderJournalXLSX(items, filter)
		if err != nil {
			return nil, err
		}
		return &ReportFile{
			Name:        "weighing-journal-" + stamp + ".xlsx",
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Data:        data,
		}, nil
	}

	return &ReportFile{
		Name:        "weighing-journal-" + stamp + ".html",
		ContentType: "text/html; charset=utf-8",
		Data:        reports.RenderJournal(items, filter),
	}, nil
}

func (s *ReportService) runDocument(
	ctx context.Context,
	code string,
	query url.Values,
	format string,
) (*ReportFile, error) {
	if format != "html" {
		return nil, &ValidationError{Message: "Документы формируются в HTML (печать из браузера)"}
	}

	sessionID, err := parseSessionID(query.Get("sessionId"))
	if err != nil {
		return nil, err
	}

	details, err := s.docs.SessionDetails(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}

	stamp := time.Now().Format("20060102-150405")

	switch code {
	case "ts-card":
		photos := s.inlinePhotos(ctx, sessionID)
		return &ReportFile{
			Name:        fmt.Sprintf("ts-card-%d-%s.html", sessionID, stamp),
			ContentType: "text/html; charset=utf-8",
			Data:        reports.RenderTSCard(details, settings, photos),
		}, nil
	case "talon":
		return &ReportFile{
			Name:        fmt.Sprintf("talon-%d-%s.html", sessionID, stamp),
			ContentType: "text/html; charset=utf-8",
			Data:        reports.RenderTalon(details, settings),
		}, nil
	case "ttn":
		doc := s.ttnDefaults(details, settings)
		if saved, err := s.docs.TTN(ctx, sessionID); err != nil {
			return nil, err
		} else if saved != nil {
			doc = *saved
		}
		return &ReportFile{
			Name:        fmt.Sprintf("ttn-%d-%s.html", sessionID, stamp),
			ContentType: "text/html; charset=utf-8",
			Data:        reports.RenderTTN(details, doc, settings),
		}, nil
	case "torg12":
		doc := s.torg12Defaults(details, settings)
		if saved, err := s.docs.Torg12(ctx, sessionID); err != nil {
			return nil, err
		} else if saved != nil {
			doc = *saved
		}
		return &ReportFile{
			Name:        fmt.Sprintf("torg12-%d-%s.html", sessionID, stamp),
			ContentType: "text/html; charset=utf-8",
			Data:        reports.RenderTorg12(details, doc, settings),
		}, nil
	}

	return nil, &ValidationError{Message: "Документ не найден"}
}

// Settings возвращает реквизиты организации для бланков.
func (s *ReportService) Settings(ctx context.Context) (dto.ReportSettings, error) {
	values, err := s.docs.PrintSettings(ctx)
	if err != nil {
		return dto.ReportSettings{}, err
	}

	return dto.ReportSettings{
		OrganizationName: values["report.organization.name"],
		INN:              values["report.organization.inn"],
		KPP:              values["report.organization.kpp"],
		Address:          values["report.organization.address"],
		OKPO:             values["report.organization.okpo"],
		Chief:            values["report.organization.chief"],
		ChiefPosition:    values["report.organization.chiefPosition"],
		Accountant:       values["report.organization.accountant"],
	}, nil
}

// SaveSettings сохраняет реквизиты организации.
func (s *ReportService) SaveSettings(ctx context.Context, settings dto.ReportSettings) (dto.ReportSettings, error) {
	values := map[string]string{
		"report.organization.name":          strings.TrimSpace(settings.OrganizationName),
		"report.organization.inn":           strings.TrimSpace(settings.INN),
		"report.organization.kpp":           strings.TrimSpace(settings.KPP),
		"report.organization.address":       strings.TrimSpace(settings.Address),
		"report.organization.okpo":          strings.TrimSpace(settings.OKPO),
		"report.organization.chief":         strings.TrimSpace(settings.Chief),
		"report.organization.chiefPosition": strings.TrimSpace(settings.ChiefPosition),
		"report.organization.accountant":    strings.TrimSpace(settings.Accountant),
	}

	if err := s.docs.SavePrintSettings(ctx, values); err != nil {
		return dto.ReportSettings{}, err
	}

	return s.Settings(ctx)
}

// Document возвращает реквизиты документа рейса: сохранённые или подставленные по умолчанию.
func (s *ReportService) Document(ctx context.Context, sessionID int64, kind string) (any, error) {
	if sessionID <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}

	details, err := s.docs.SessionDetails(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}

	switch kind {
	case "ttn":
		saved, err := s.docs.TTN(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if saved != nil {
			return saved, nil
		}
		doc := s.ttnDefaults(details, settings)
		return &doc, nil
	case "torg12":
		saved, err := s.docs.Torg12(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if saved != nil {
			return saved, nil
		}
		doc := s.torg12Defaults(details, settings)
		return &doc, nil
	default:
		return nil, &ValidationError{Message: "Неизвестный документ: " + kind}
	}
}

// SaveDocument сохраняет реквизиты ТТН или ТОРГ-12 на рейсе.
func (s *ReportService) SaveDocument(
	ctx context.Context,
	sessionID int64,
	kind string,
	body []byte,
	userID int64,
) (any, error) {
	if sessionID <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}

	details, err := s.docs.SessionDetails(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	switch kind {
	case "ttn":
		var doc dto.TTNDocument
		if err := json.Unmarshal(body, &doc); err != nil {
			return nil, &ValidationError{Message: "Некорректный JSON документа"}
		}
		doc.SessionID = sessionID
		return s.docs.SaveTTN(ctx, &doc, userID)
	case "torg12":
		var doc dto.Torg12Document
		if err := json.Unmarshal(body, &doc); err != nil {
			return nil, &ValidationError{Message: "Некорректный JSON документа"}
		}
		doc.SessionID = sessionID

		fallback := details.Session.FirstMeasuredAt
		if details.Session.SecondMeasuredAt != nil {
			fallback = *details.Session.SecondMeasuredAt
		}
		date := reports.ParseDocumentDate(doc.Date, fallback)
		return s.docs.SaveTorg12(ctx, &doc, date, reports.ParseOptionalDate(doc.ProxyDate), userID)
	default:
		return nil, &ValidationError{Message: "Неизвестный документ: " + kind}
	}
}

func (s *ReportService) ttnDefaults(details *dto.SessionDetails, settings dto.ReportSettings) dto.TTNDocument {
	session := details.Session

	price := 0.0
	if details.CargoSellingPrice != nil {
		price = *details.CargoSellingPrice
	} else if details.CargoPricePerTonne != nil {
		price = *details.CargoPricePerTonne
	}

	return dto.TTNDocument{
		SessionID:             session.ID,
		TTNNumber:             "",
		CustomerOrganization:  settings.OrganizationName,
		ShipperOrganization:   settings.OrganizationName,
		ConsigneeOrganization: session.CounterpartyName,
		PayerOrganization:     session.CounterpartyName,
		CarrierOrganization:   session.CarrierName,
		TrailerNumber:         session.TrailerNumber,
		LoadPoint:             session.StorageName,
		CargoReleaseName:      firstValue(settings.Chief, details.OperatorName),
		CargoReleaseWorkAs:    settings.ChiefPosition,
		AccountantName:        settings.Accountant,
		CargoPickUpName:       details.CounterpartyDirector,
		CargoPickUpWorkAs:     "Директор",
		ConsigneeName:         details.CounterpartyDirector,
		ConsigneeWorkAs:       "Директор",
		PriceCargo:            price,
	}
}

func (s *ReportService) torg12Defaults(details *dto.SessionDetails, settings dto.ReportSettings) dto.Torg12Document {
	session := details.Session

	date := session.FirstMeasuredAt
	if session.SecondMeasuredAt != nil {
		date = *session.SecondMeasuredAt
	}

	return dto.Torg12Document{
		SessionID:          session.ID,
		Name:               "",
		Date:               date.Format("2006-01-02"),
		ProviderOrg:        settings.OrganizationName,
		ConsigneeOrg:       session.CounterpartyName,
		NDS:                20,
		CargoReleaseName:   firstValue(settings.Chief, details.OperatorName),
		CargoReleaseWorkAs: settings.ChiefPosition,
		AccountantName:     settings.Accountant,
		ShipperName:        firstValue(settings.Chief, details.OperatorName),
		ShipperWorkAs:      settings.ChiefPosition,
		CargoPickUpName:    details.CounterpartyDirector,
		CargoPickUpWorkAs:  "Директор",
		ConsigneeName:      details.CounterpartyDirector,
		ConsigneeWorkAs:    "Директор",
	}
}

func (s *ReportService) inlinePhotos(ctx context.Context, sessionID int64) []reports.InlinePhoto {
	if s.photos == nil || s.photoDir == "" {
		return nil
	}

	items, err := s.photos.List(ctx, sessionID)
	if err != nil {
		return nil
	}

	result := make([]reports.InlinePhoto, 0, maxInlinePhotos)
	base := filepath.Clean(s.photoDir) + string(os.PathSeparator)

	for _, item := range items {
		if len(result) >= maxInlinePhotos {
			break
		}

		full := filepath.Join(s.photoDir, filepath.FromSlash(item.FileName))
		if !strings.HasPrefix(filepath.Clean(full), base) {
			continue
		}

		info, err := os.Stat(full)
		if err != nil || info.IsDir() || info.Size() > 6<<20 {
			continue
		}

		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}

		contentType := "image/jpeg"
		if strings.EqualFold(filepath.Ext(full), ".png") {
			contentType = "image/png"
		}

		result = append(result, reports.InlinePhoto{
			First: item.FirstWeight,
			URI:   "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data),
		})
	}

	return result
}

func parseReportDate(value string, endOfDay bool) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		return nil, &ValidationError{Message: "Дата должна быть в формате ГГГГ-ММ-ДД"}
	}
	if endOfDay {
		parsed = parsed.Add(24 * time.Hour)
	}

	return &parsed, nil
}

func parseSessionID(value string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || id <= 0 {
		return 0, &ValidationError{Message: "Укажите корректный рейс"}
	}
	return id, nil
}

func reportStatusLabel(status string) string {
	switch status {
	case "open":
		return "Открытые"
	case "closed":
		return "Закрытые"
	case "cancelled":
		return "Отменённые"
	default:
		return "Все"
	}
}

func firstValue(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
