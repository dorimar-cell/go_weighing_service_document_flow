package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// renderTalonsBatch собирает талоны за период: один талон на страницу.
func renderTalonsBatch(items []dto.SessionDTO, filter journalFilter, settings dto.ReportSettings) []byte {
	var body strings.Builder

	body.WriteString("<h1>Талоны взвешивания</h1>")
	body.WriteString("<p class=\"subtitle\">Пакетная печать за период</p>")
	body.WriteString("<div class=\"meta\">")
	if settings.OrganizationName != "" {
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
	}
	body.WriteString("<span>" + esc(aggregateMeta(aggregateFilter{
		DateFrom: filter.DateFrom,
		DateTo:   filter.DateTo,
	})) + "</span>")
	body.WriteString("<span>Талонов: " + fmt.Sprintf("%d", len(items)) + "</span>")
	body.WriteString("</div>")

	if len(items) == 0 {
		body.WriteString("<p class=\"muted\">За выбранный период закрытых рейсов нет.</p>")
		return wrapHTML("Талоны взвешивания", false, body.String())
	}

	for _, item := range items {
		session := item
		brutto, tare := bruttoTare(&session)

		body.WriteString("<section class=\"talon-page\">")
		body.WriteString("<h1>Талон взвешивания</h1>")
		body.WriteString("<p class=\"subtitle\">№ " + esc(int64Text(session.ID)) + " от " +
			esc(dateOnly(session.FirstMeasuredAt)) + "</p>")

		body.WriteString(pairRows([][2]string{
			{"Транспортное средство", esc(session.TransportNumber)},
			{"Прицеп", orDash(session.TrailerNumber)},
			{"Водитель", orDash(session.DriverName)},
			{"Пропуск", orDash(session.PassNumber)},
			{"Груз", orDash(session.CargoName)},
			{"Контрагент", orDash(session.CounterpartyName)},
			{"Перевозчик", orDash(session.CarrierName)},
			{"Склад", orDash(session.StorageName)},
			{"Тип операции", orDash(session.OperationTypeName)},
		}))

		body.WriteString("<table><thead><tr><th>Показатель</th><th>Кг</th><th>Тонн</th></tr></thead><tbody>")
		body.WriteString("<tr><td>Брутто</td><td class=\"num\">" + formatNumber(brutto, 1) +
			"</td><td class=\"num\">" + formatNumber(brutto/1000, 3) + "</td></tr>")
		body.WriteString("<tr><td>Тара</td><td class=\"num\">" + formatNumber(tare, 1) +
			"</td><td class=\"num\">" + formatNumber(tare/1000, 3) + "</td></tr>")
		body.WriteString("<tr class=\"total\"><td>Нетто</td><td class=\"num\">" +
			formatNumber(nettoValue(&session), 1) + "</td><td class=\"num\">" +
			tonnesPtr(session.NetWeightKg) + "</td></tr>")
		body.WriteString("</tbody></table>")

		body.WriteString("<div class=\"signs\">")
		body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span></div>")
		body.WriteString("<div>Водитель: <span class=\"sign-line\"></span> " + orDash(session.DriverName) + "</div>")
		body.WriteString("</div>")
		body.WriteString("</section>")
	}

	return wrapHTML("Талоны взвешивания", false, body.String())
}
