package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// CargoActBatchItem — рейс для пакетной печати актов.
type CargoActBatchItem struct {
	Details *dto.SessionDetails
}

// renderCargoActBatch собирает пакет актов приёма-передачи за период.
func renderCargoActBatch(
	items []CargoActBatchItem,
	filter journalFilter,
	settings dto.ReportSettings,
) []byte {
	var body strings.Builder

	body.WriteString("<h1>Акты приёма-передачи за период</h1>")
	body.WriteString("<p class=\"subtitle\">Пакетная печать по закрытым рейсам</p>")
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
		return wrapHTML("Акты приёма-передачи за период", false, body.String())
	}

	for _, item := range items {
		body.WriteString("<section class=\"talon-page\">")
		body.WriteString(cargoActBody(item.Details, settings))
		body.WriteString("</section>")
	}

	return wrapHTML("Акты приёма-передачи за период", false, body.String())
}
