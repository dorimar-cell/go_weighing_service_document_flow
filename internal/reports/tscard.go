package reports

import (
	"strings"

	"go_weighing_service/internal/dto"
)

// inlinePhoto — фото рейса, встроенное в HTML как data URI.
type inlinePhoto struct {
	First bool
	URI   string
}

func renderTSCard(details *dto.SessionDetails, settings dto.ReportSettings, photos []inlinePhoto) []byte {
	session := &details.Session
	brutto, tare := bruttoTare(session)

	var body strings.Builder
	body.WriteString("<h1>Карточка транспортного средства</h1>")
	body.WriteString("<p class=\"subtitle\">Рейс №" + esc(int64Text(session.ID)) +
		" от " + esc(dateOnly(session.FirstMeasuredAt)) + " · " + esc(statusLabel(session.Status)) + "</p>")

	body.WriteString("<div class=\"meta\">")
	if settings.OrganizationName != "" {
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
	}
	if details.OperatorName != "" {
		body.WriteString("<span>Оператор: " + esc(details.OperatorName) + "</span>")
	}
	body.WriteString("<span>Сформирована: " + esc(dateOnlyNow()) + "</span>")
	body.WriteString("</div>")

	body.WriteString("<table class=\"pairs\"><tbody>")
	pair := func(label, value string) {
		body.WriteString("<tr><td class=\"label\">" + esc(label) + "</td><td>" + value + "</td></tr>")
	}

	pair("Транспортное средство", esc(session.TransportNumber))
	pair("Прицеп", orDash(session.TrailerNumber))
	pair("Водитель", orDash(session.DriverName))
	if details.DriverLicense != "" {
		pair("Водительское удостоверение", esc(details.DriverLicense))
	}
	pair("Пропуск", orDash(session.PassNumber))
	pair("Груз", orDash(session.CargoName))
	pair("Контрагент", counterpartyText(details))
	pair("Перевозчик", carrierText(details))
	pair("Склад", orDash(session.StorageName))
	pair("Тип операции", orDash(session.OperationTypeName))
	body.WriteString("</tbody></table>")

	body.WriteString("<h2>Взвешивания</h2>")
	body.WriteString("<table><thead><tr>" +
		"<th>Взвешивание</th><th>Дата</th><th>Время</th><th>Вес, кг</th><th>Вес, т</th>" +
		"</tr></thead><tbody>")
	body.WriteString("<tr><td>Первое</td><td class=\"center\">" + esc(dateOnly(session.FirstMeasuredAt)) +
		"</td><td class=\"center\">" + esc(timeOnly(session.FirstMeasuredAt)) +
		"</td><td class=\"num\">" + formatNumber(session.FirstWeightKg, 1) +
		"</td><td class=\"num\">" + formatNumber(session.FirstWeightKg/1000, 3) + "</td></tr>")
	if session.SecondMeasuredAt != nil {
		body.WriteString("<tr><td>Второе</td><td class=\"center\">" + esc(dateOnlyPtr(session.SecondMeasuredAt)) +
			"</td><td class=\"center\">" + esc(timeOnlyPtr(session.SecondMeasuredAt)) +
			"</td><td class=\"num\">" + formatNumber(*session.SecondWeightKg, 1) +
			"</td><td class=\"num\">" + formatNumber(*session.SecondWeightKg/1000, 3) + "</td></tr>")
	} else {
		body.WriteString("<tr><td>Второе</td><td class=\"center\">—</td><td class=\"center\">—</td>" +
			"<td class=\"num\">—</td><td class=\"num\">—</td></tr>")
	}
	body.WriteString("<tr class=\"total\"><td colspan=\"3\" class=\"num\">Брутто / Тара, кг</td>" +
		"<td class=\"num\">" + formatNumber(brutto, 1) + " / " + formatNumber(tare, 1) + "</td>" +
		"<td class=\"num\">Нетто: " + tonnesPtr(session.NetWeightKg) + "</td></tr>")
	body.WriteString("</tbody></table>")

	if len(photos) > 0 {
		body.WriteString("<h2>Фотографии</h2><div class=\"photos\">")
		for _, photo := range photos {
			caption := "Первое взвешивание"
			if !photo.First {
				caption = "Второе взвешивание"
			}
			body.WriteString("<figure><img src=\"" + esc(photo.URI) + "\" alt=\"\"><figcaption>" +
				esc(caption) + "</figcaption></figure>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span></div>")
	body.WriteString("<div>Водитель: <span class=\"sign-line\"></span></div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Карточка сформирована веб-системой весовой по данным рейса и " +
		"подписанных измерений Station Agent.</p>")

	return wrapHTML("Карточка ТС №"+int64Text(session.ID), false, body.String())
}

func counterpartyText(details *dto.SessionDetails) string {
	name := details.Session.CounterpartyName
	if name == "" {
		return "—"
	}
	return esc(name) + requisites(details.CounterpartyINN, details.CounterpartyKPP)
}

func carrierText(details *dto.SessionDetails) string {
	name := details.Session.CarrierName
	if name == "" {
		return "—"
	}
	return esc(name) + requisites(details.CarrierINN, details.CarrierKPP)
}

func requisites(inn, kpp string) string {
	parts := make([]string, 0, 2)
	if inn != "" {
		parts = append(parts, "ИНН "+esc(inn))
	}
	if kpp != "" {
		parts = append(parts, "КПП "+esc(kpp))
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, ", ") + ")"
}
