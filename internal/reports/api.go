package reports

import (
	"time"

	"go_weighing_service/internal/dto"
)

// JournalFilter — фильтры, попадающие в шапку печатного журнала.
type JournalFilter = journalFilter

// InlinePhoto — фотография рейса для карточки ТС.
type InlinePhoto = inlinePhoto

// RenderJournal собирает HTML журнала взвешиваний.
func RenderJournal(items []dto.SessionDTO, filter JournalFilter) []byte {
	return renderJournal(items, filter)
}

// RenderJournalXLSX собирает печатную книгу журнала.
func RenderJournalXLSX(items []dto.SessionDTO, filter JournalFilter) ([]byte, error) {
	return renderJournalXLSX(items, filter)
}

// RenderTSCard собирает карточку транспортного средства.
func RenderTSCard(
	details *dto.SessionDetails,
	settings dto.ReportSettings,
	photos []InlinePhoto,
) []byte {
	return renderTSCard(details, settings, photos)
}

// RenderTTN собирает ТТН.
func RenderTTN(details *dto.SessionDetails, doc dto.TTNDocument, settings dto.ReportSettings) []byte {
	return renderTTN(details, doc, settings)
}

// RenderTorg12 собирает ТОРГ-12.
func RenderTorg12(
	details *dto.SessionDetails,
	doc dto.Torg12Document,
	settings dto.ReportSettings,
) []byte {
	return renderTorg12(details, doc, settings)
}

// ParseDocumentDate разбирает дату из формы; при ошибке возвращает fallback.
func ParseDocumentDate(value string, fallback time.Time) time.Time {
	return parseDocumentDate(value, fallback)
}

// ParseOptionalDate разбирает необязательную дату.
func ParseOptionalDate(value string) *time.Time {
	return parseOptionalDate(value)
}
