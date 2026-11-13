package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// InvoiceBatchItem — рейс с реквизитами счёта-фактуры.
type InvoiceBatchItem struct {
	Details  *dto.SessionDetails
	Document dto.InvoiceDocument
}

// renderInvoiceBatch собирает пакет счетов-фактур за период.
func renderInvoiceBatch(
	items []InvoiceBatchItem,
	filter journalFilter,
	settings dto.ReportSettings,
) []byte {
	var body strings.Builder

	body.WriteString("<h1>Счета-фактуры за период</h1>")
	body.WriteString("<p class=\"subtitle\">Пакетная печать счетов-фактур</p>")
	body.WriteString("<div class=\"meta\">")
	if settings.OrganizationName != "" {
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
	}
	body.WriteString("<span>" + esc(aggregateMeta(aggregateFilter{
		DateFrom: filter.DateFrom,
		DateTo:   filter.DateTo,
	})) + "</span>")
	body.WriteString("<span>Документов: " + fmt.Sprintf("%d", len(items)) + "</span>")
	body.WriteString("<span>Сформирован: " + nowInLocation().Format("02.01.2006 15:04") + "</span>")
	body.WriteString("</div>")

	if len(items) == 0 {
		body.WriteString("<p class=\"muted\">За выбранный период закрытых рейсов нет.</p>")
		return wrapHTML("Счета-фактуры за период", false, body.String())
	}

	for _, item := range items {
		body.WriteString("<section class=\"talon-page\">")
		body.WriteString(invoiceBody(item.Details, item.Document, settings))
		body.WriteString("</section>")
	}

	body.WriteString("<p class=\"note\">Реквизиты берутся сохранённые на рейсе; для незаполненных " +
		"документов подставляются значения по умолчанию.</p>")

	return wrapHTML("Счета-фактуры за период", false, body.String())
}
