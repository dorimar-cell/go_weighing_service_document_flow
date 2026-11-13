package reports

import (
	"math"
	"strings"

	"go_weighing_service/internal/dto"
)

// renderInvoice собирает счёт-фактуру по рейсу.
func renderInvoice(details *dto.SessionDetails, doc dto.InvoiceDocument, settings dto.ReportSettings) []byte {
	session := &details.Session

	quantity := cargoQuantity(session)
	price := cargoPrice(details, 0)
	sumWithoutVat := quantity * price
	nds := doc.NDS
	if nds < 0 {
		nds = 0
	}
	vatSum := sumWithoutVat * nds / 100
	total := sumWithoutVat + vatSum
	totalKopecks := int64(math.Round(total * 100))

	number := doc.Name
	if number == "" {
		number = "б/н"
	}
	documentDateValue := doc.Date
	if documentDateValue == "" {
		documentDateValue = documentDate(session)
	}

	chief := firstNonEmpty(doc.ChiefName, settings.Chief)
	accountant := firstNonEmpty(doc.AccountantName, settings.Accountant)

	var body strings.Builder
	body.WriteString("<h1>Счёт-фактура</h1>")
	body.WriteString("<p class=\"subtitle\">№ " + esc(number) + " от " + esc(prettyDate(documentDateValue)) + "</p>")

	body.WriteString("<h2>Продавец</h2>")
	body.WriteString(pairRows([][2]string{
		{"Наименование", orDash(settings.OrganizationName)},
		{"ИНН / КПП", orDash(joinRequisites(settings.INN, settings.KPP))},
		{"Адрес", orDash(settings.Address)},
	}))

	body.WriteString("<h2>Покупатель</h2>")
	body.WriteString(pairRows([][2]string{
		{"Наименование", orDash(session.CounterpartyName)},
		{"ИНН / КПП", orDash(joinRequisites(details.CounterpartyINN, details.CounterpartyKPP))},
		{"Адрес", orDash(details.CounterpartyAddress)},
	}))

	if doc.DocumentName != "" {
		body.WriteString("<p class=\"words\">Платёжно-расчётный документ: " + esc(doc.DocumentName) + "</p>")
	}

	body.WriteString("<table><thead><tr>" +
		"<th>№</th><th>Наименование товара (груза)</th><th>Ед. изм.</th><th>Количество, т</th>" +
		"<th>Цена, руб.</th><th>Стоимость без налога, руб.</th><th>НДС, %</th>" +
		"<th>Сумма налога, руб.</th><th>Всего с налогом, руб.</th>" +
		"</tr></thead><tbody>")
	body.WriteString("<tr><td class=\"center\">1</td><td>" + orDash(session.CargoName) + "</td>" +
		"<td class=\"center\">т</td><td class=\"num\">" + formatNumber(quantity, 3) + "</td>" +
		"<td class=\"num\">" + formatNumber(price, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(sumWithoutVat, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(nds, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(vatSum, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(total, 2) + "</td></tr>")
	body.WriteString("<tr class=\"total\"><td colspan=\"3\" class=\"num\">Всего к оплате</td>" +
		"<td class=\"num\">" + formatNumber(quantity, 3) + "</td><td></td>" +
		"<td class=\"num\">" + formatNumber(sumWithoutVat, 2) + "</td><td></td>" +
		"<td class=\"num\">" + formatNumber(vatSum, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(total, 2) + "</td></tr>")
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"words\">Всего к оплате: <b>" + rublesWords(totalKopecks) + "</b></p>")
	body.WriteString("<p class=\"words\">В том числе НДС: " + formatNumber(vatSum, 2) +
		" руб. (" + formatNumber(nds, 2) + " %)</p>")

	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Руководитель организации: <span class=\"sign-line\"></span> " + orDash(chief) + "</div>")
	body.WriteString("<div>Главный бухгалтер: <span class=\"sign-line\"></span> " + orDash(accountant) + "</div>")
	if doc.BusinessmanName != "" || doc.RegIPCertificate != "" {
		body.WriteString("<div>Индивидуальный предприниматель: <span class=\"sign-line\"></span> " +
			orDash(doc.BusinessmanName) + "</div>")
		if doc.RegIPCertificate != "" {
			body.WriteString("<div>Свидетельство: " + esc(doc.RegIPCertificate) + "</div>")
		}
	}
	body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Счёт-фактура сформирована веб-системой весовой. " +
		"Суммы рассчитаны по нетто рейса и цене груза; при необходимости скорректируйте реквизиты и цену.</p>")

	return wrapHTML("Счёт-фактура №"+number, false, body.String())
}

func joinRequisites(inn, kpp string) string {
	parts := make([]string, 0, 2)
	if inn != "" {
		parts = append(parts, "ИНН "+inn)
	}
	if kpp != "" {
		parts = append(parts, "КПП "+kpp)
	}
	return strings.Join(parts, ", ")
}
