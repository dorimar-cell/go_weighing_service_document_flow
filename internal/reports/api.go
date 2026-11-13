package reports

import (
	"time"

	"go_weighing_service/internal/dto"
)

// JournalFilter — фильтры, попадающие в шапку печатного журнала.
type JournalFilter = journalFilter

// AggregateFilter — фильтры агрегатного отчёта.
type AggregateFilter = aggregateFilter

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

// RenderAggregate собирает HTML агрегатного отчёта.
func RenderAggregate(rows []dto.AggregateRow, filter AggregateFilter) []byte {
	return renderAggregate(rows, filter)
}

// RenderAggregateXLSX собирает агрегатный отчёт книгой Excel.
func RenderAggregateXLSX(rows []dto.AggregateRow, filter AggregateFilter) ([]byte, error) {
	return renderAggregateXLSX(rows, filter)
}

// RenderTalon собирает талон взвешивания.
func RenderTalon(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	return renderTalon(details, settings)
}

// RenderTalonVariant собирает талон в выбранном варианте бланка.
func RenderTalonVariant(
	details *dto.SessionDetails,
	settings dto.ReportSettings,
	variant string,
) []byte {
	return renderTalonVariant(details, settings, variant)
}

// RenderTalonsBatch собирает пакет талонов за период.
func RenderTalonsBatch(
	items []dto.SessionDTO,
	filter JournalFilter,
	settings dto.ReportSettings,
) []byte {
	return renderTalonsBatch(items, filter, settings)
}

// RenderFirstWeight собирает отчёт по транспорту на первом взвешивании.
func RenderFirstWeight(items []dto.SessionDTO, filter JournalFilter) []byte {
	return renderFirstWeight(items, filter)
}

// RenderTransportJournal собирает журнал по номеру ТС.
func RenderTransportJournal(items []dto.SessionDTO, filter JournalFilter) []byte {
	return renderTransportJournal(items, filter)
}

// RenderTTNBatch собирает пакет ТТН за период.
func RenderTTNBatch(
	items []TTNBatchItem,
	filter JournalFilter,
	settings dto.ReportSettings,
) []byte {
	return renderTTNBatch(items, filter, settings)
}

// RenderCargoAct собирает акт приёма-передачи товара.
func RenderCargoAct(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	return renderCargoAct(details, settings)
}

// RenderUploadingCertificate собирает свидетельство о разгрузке/погрузке.
func RenderUploadingCertificate(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	return renderUploadingCertificate(details, settings)
}

// RenderFirstWeightXLSX собирает отчёт по первому взвешиванию в Excel.
func RenderFirstWeightXLSX(items []dto.SessionDTO, filter JournalFilter) ([]byte, error) {
	return renderFirstWeightXLSX(items, filter)
}

// RenderTransportJournalXLSX собирает журнал по номеру ТС в Excel.
func RenderTransportJournalXLSX(items []dto.SessionDTO, filter JournalFilter) ([]byte, error) {
	return renderTransportJournalXLSX(items, filter)
}

// RenderInvoiceBatch собирает пакет счетов-фактур за период.
func RenderInvoiceBatch(
	items []InvoiceBatchItem,
	filter JournalFilter,
	settings dto.ReportSettings,
) []byte {
	return renderInvoiceBatch(items, filter, settings)
}

// RenderTorg12Batch собирает пакет ТОРГ-12 за период.
func RenderTorg12Batch(
	items []Torg12BatchItem,
	filter JournalFilter,
	settings dto.ReportSettings,
) []byte {
	return renderTorg12Batch(items, filter, settings)
}

// RenderSessionPack собирает комплект документов по рейсу.
func RenderSessionPack(pack SessionPack, settings dto.ReportSettings) []byte {
	return renderSessionPack(pack, settings)
}

// RenderInvoice собирает счёт-фактуру.
func RenderInvoice(
	details *dto.SessionDetails,
	doc dto.InvoiceDocument,
	settings dto.ReportSettings,
) []byte {
	return renderInvoice(details, doc, settings)
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
