package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// Torg12BatchItem — рейс с реквизитами ТОРГ-12 для пакетной печати.
type Torg12BatchItem struct {
	Details  *dto.SessionDetails
	Document dto.Torg12Document
}

// renderTorg12Batch собирает пакет ТОРГ-12 за период: документ на страницу.
func renderTorg12Batch(
	items []Torg12BatchItem,
	filter journalFilter,
	settings dto.ReportSettings,
) []byte {
	var body strings.Builder

	body.WriteString("<h1>ТОРГ-12 за период</h1>")
	body.WriteString("<p class=\"subtitle\">Пакетная печать товарных накладных</p>")
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
		return wrapHTML("ТОРГ-12 за период", false, body.String())
	}

	for _, item := range items {
		body.WriteString("<section class=\"talon-page\">")
		body.WriteString(torg12Body(item.Details, item.Document, settings))
		body.WriteString("</section>")
	}

	body.WriteString("<p class=\"note\">Реквизиты берутся сохранённые на рейсе; для незаполненных " +
		"документов подставляются значения по умолчанию.</p>")

	return wrapHTML("ТОРГ-12 за период", false, body.String())
}
