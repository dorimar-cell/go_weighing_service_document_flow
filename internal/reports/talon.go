package reports

import (
	"strings"

	"go_weighing_service/internal/dto"
)

// renderTalon собирает отвесной талон по рейсу.
func renderTalon(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	session := &details.Session
	brutto, tare := bruttoTare(session)

	var body strings.Builder

	if settings.OrganizationName != "" {
		body.WriteString("<div class=\"meta\"><span><b>" + esc(settings.OrganizationName) + "</b></span>")
		if settings.INN != "" {
			body.WriteString("<span>ИНН: " + esc(settings.INN) + "</span>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<h1>Талон взвешивания</h1>")
	body.WriteString("<p class=\"subtitle\">№ " + esc(int64Text(session.ID)) + " от " +
		esc(documentDate(session)) + "</p>")

	body.WriteString(pairRows([][2]string{
		{"Транспортное средство", esc(session.TransportNumber)},
		{"Прицеп", orDash(session.TrailerNumber)},
		{"Водитель", orDash(session.DriverName)},
		{"Пропуск", orDash(session.PassNumber)},
		{"Груз", orDash(session.CargoName)},
		{"Контрагент", counterpartyText(details)},
		{"Перевозчик", carrierText(details)},
		{"Склад", orDash(session.StorageName)},
		{"Тип операции", orDash(session.OperationTypeName)},
	}))

	body.WriteString("<h2>Результат взвешивания</h2>")
	body.WriteString("<table><thead><tr><th>Показатель</th><th>Кг</th><th>Тонн</th></tr></thead><tbody>")
	body.WriteString("<tr><td>Брутто</td><td class=\"num\">" + formatNumber(brutto, 1) +
		"</td><td class=\"num\">" + formatNumber(brutto/1000, 3) + "</td></tr>")
	body.WriteString("<tr><td>Тара</td><td class=\"num\">" + formatNumber(tare, 1) +
		"</td><td class=\"num\">" + formatNumber(tare/1000, 3) + "</td></tr>")
	body.WriteString("<tr class=\"total\"><td>Нетто</td><td class=\"num\">" +
		formatNumber(nettoValue(session), 1) + "</td><td class=\"num\">" +
		tonnesPtr(session.NetWeightKg) + "</td></tr>")
	body.WriteString("</tbody></table>")

	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span> " +
		orDash(details.OperatorName) + "</div>")
	body.WriteString("<div>Водитель: <span class=\"sign-line\"></span> " + orDash(session.DriverName) + "</div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Талон сформирован веб-системой весовой по данным рейса. " +
		"Предназначен для внутреннего учёта и выдачи водителю.</p>")

	return wrapHTML("Талон №"+int64Text(session.ID), false, body.String())
}

func nettoValue(session *dto.SessionDTO) float64 {
	if session.NetWeightKg == nil {
		return 0
	}
	return *session.NetWeightKg
}
