package reports

import (
	"math"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
)

func pairRows(rows [][2]string) string {
	var builder strings.Builder
	builder.WriteString("<table class=\"pairs\"><tbody>")
	for _, row := range rows {
		builder.WriteString("<tr><td class=\"label\">" + esc(row[0]) + "</td><td>" + row[1] + "</td></tr>")
	}
	builder.WriteString("</tbody></table>")
	return builder.String()
}

func documentDate(session *dto.SessionDTO) string {
	if session.SecondMeasuredAt != nil && !session.SecondMeasuredAt.IsZero() {
		return dateOnlyPtr(session.SecondMeasuredAt)
	}
	return dateOnly(session.FirstMeasuredAt)
}

func cargoQuantity(session *dto.SessionDTO) float64 {
	if session.NetWeightKg == nil {
		return 0
	}
	return *session.NetWeightKg / 1000
}

func cargoPrice(details *dto.SessionDetails, fallback float64) float64 {
	if fallback > 0 {
		return fallback
	}
	if details.CargoSellingPrice != nil && *details.CargoSellingPrice > 0 {
		return *details.CargoSellingPrice
	}
	if details.CargoPricePerTonne != nil && *details.CargoPricePerTonne > 0 {
		return *details.CargoPricePerTonne
	}
	return 0
}

func renderTTN(details *dto.SessionDetails, doc dto.TTNDocument, settings dto.ReportSettings) []byte {
	session := &details.Session

	quantity := cargoQuantity(session)
	price := cargoPrice(details, doc.PriceCargo)
	sum := quantity * price

	organization := firstNonEmpty(doc.ShipperOrganization, settings.OrganizationName)
	consignee := firstNonEmpty(doc.ConsigneeOrganization, session.CounterpartyName)
	payer := firstNonEmpty(doc.PayerOrganization, consignee)
	customer := firstNonEmpty(doc.CustomerOrganization, settings.OrganizationName)
	carrier := firstNonEmpty(doc.CarrierOrganization, session.CarrierName)
	number := doc.TTNNumber
	if number == "" {
		number = "б/н"
	}

	var body strings.Builder
	if settings.OrganizationName != "" || settings.INN != "" {
		body.WriteString("<div class=\"meta\">")
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
		if settings.INN != "" {
			body.WriteString("<span>ИНН: " + esc(settings.INN) + "</span>")
		}
		if settings.KPP != "" {
			body.WriteString("<span>КПП: " + esc(settings.KPP) + "</span>")
		}
		if settings.Address != "" {
			body.WriteString("<span>" + esc(settings.Address) + "</span>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<h1>Товарно-транспортная накладная</h1>")
	body.WriteString("<p class=\"subtitle\">№ " + esc(number) + " от " + esc(documentDate(session)) + "</p>")

	body.WriteString(pairRows([][2]string{
		{"Грузоотправитель", orDash(organization)},
		{"Грузополучатель", orDash(consignee)},
		{"Плательщик", orDash(payer)},
		{"Заказчик (владелец)", orDash(customer)},
		{"Перевозчик", carrierCell(details, carrier)},
	}))

	body.WriteString("<h2>Транспорт и водитель</h2>")
	body.WriteString(pairRows([][2]string{
		{"Автомобиль", esc(session.TransportNumber)},
		{"Прицеп", orDash(firstNonEmpty(doc.TrailerNumber, session.TrailerNumber))},
		{"Тип прицепа", orDash(doc.TypeOfTrailer)},
		{"Водитель", orDash(session.DriverName)},
		{"Водительское удостоверение", orDash(details.DriverLicense)},
		{"Путевой лист / путёвка", orDash(doc.Putevka)},
		{"Договор", orDash(doc.Dogovor)},
		{"Приложение", orDash(doc.Prilozhenie)},
	}))

	body.WriteString("<h2>Груз</h2>")
	body.WriteString("<table><thead><tr>" +
		"<th>№</th><th>Наименование груза</th><th>Ед. изм.</th><th>Количество, т</th>" +
		"<th>Цена, руб.</th><th>Стоимость, руб.</th>" +
		"</tr></thead><tbody>")
	body.WriteString("<tr><td class=\"center\">1</td><td>" + orDash(session.CargoName) + "</td>" +
		"<td class=\"center\">т</td><td class=\"num\">" + formatNumber(quantity, 3) + "</td>" +
		"<td class=\"num\">" + formatNumber(price, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(sum, 2) + "</td></tr>")
	body.WriteString("<tr class=\"total\"><td colspan=\"3\" class=\"num\">Итого</td>" +
		"<td class=\"num\">" + formatNumber(quantity, 3) + "</td><td></td>" +
		"<td class=\"num\">" + formatNumber(sum, 2) + "</td></tr>")
	body.WriteString("</tbody></table>")

	if doc.PriceTransportation > 0 {
		body.WriteString("<p class=\"words\">Стоимость перевозки: " +
			formatNumber(doc.PriceTransportation, 2) + " руб.</p>")
	}

	body.WriteString("<h2>Погрузка и разгрузка</h2>")
	body.WriteString(pairRows([][2]string{
		{"Пункт погрузки", orDash(doc.LoadPoint)},
		{"Пункт разгрузки", orDash(doc.UploadPoint)},
		{"Расстояние перевозки, км", numberOrDash(doc.Distance)},
		{"Грузоподъёмность, т", numberOrDash(doc.CargoCapacity)},
	}))

	if session.NetWeightKg != nil {
		body.WriteString("<p class=\"words\"><b>Масса груза прописью:</b> " +
			kilogramsWords(int64(math.Round(*session.NetWeightKg))) + ".</p>")
	}

	body.WriteString("<h2>Подписи сторон</h2>")
	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Груз отпустил: <span class=\"sign-line\"></span> " +
		signature(doc.CargoReleaseName, doc.CargoReleaseWorkAs) + "</div>")
	body.WriteString("<div>Главный бухгалтер: <span class=\"sign-line\"></span> " +
		signature(doc.AccountantName, "") + "</div>")
	body.WriteString("<div>Груз принял: <span class=\"sign-line\"></span> " +
		signature(doc.CargoPickUpName, doc.CargoPickUpWorkAs) + "</div>")
	body.WriteString("<div>Груз получил: <span class=\"sign-line\"></span> " +
		signature(doc.ConsigneeName, doc.ConsigneeWorkAs) + "</div>")
	body.WriteString("<div>Водитель: <span class=\"sign-line\"></span> " + orDash(session.DriverName) + "</div>")
	body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Бланк ТТН сформирован веб-системой весовой на основании данных рейса " +
		"и сохранённых реквизитов документа. Унифицированная форма 1-Т адаптирована для печати из браузера.</p>")

	return wrapHTML("ТТН №"+number, false, body.String())
}

func renderTorg12(details *dto.SessionDetails, doc dto.Torg12Document, settings dto.ReportSettings) []byte {
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

	provider := firstNonEmpty(doc.ProviderOrg, settings.OrganizationName)
	consignee := firstNonEmpty(doc.ConsigneeOrg, session.CounterpartyName)

	var body strings.Builder
	if settings.OrganizationName != "" {
		body.WriteString("<div class=\"meta\"><span><b>" + esc(settings.OrganizationName) + "</b></span>")
		if settings.INN != "" {
			body.WriteString("<span>ИНН: " + esc(settings.INN) + "</span>")
		}
		if settings.KPP != "" {
			body.WriteString("<span>КПП: " + esc(settings.KPP) + "</span>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<h1>Товарная накладная</h1>")
	body.WriteString("<p class=\"subtitle\">№ " + esc(number) + " от " + esc(documentDateValue) + "</p>")

	body.WriteString(pairRows([][2]string{
		{"Поставщик (грузоотправитель)", orDash(provider)},
		{"Грузополучатель", orDash(consignee)},
		{"Структурное подразделение", orDash(doc.StructName)},
		{"Основание", orDash(doc.Foundation)},
	}))

	body.WriteString("<table><thead><tr>" +
		"<th>№</th><th>Наименование товара</th><th>Ед. изм.</th><th>Количество, т</th>" +
		"<th>Цена, руб.</th><th>Сумма без налога, руб.</th><th>НДС, %</th>" +
		"<th>Сумма НДС, руб.</th><th>Всего с НДС, руб.</th>" +
		"</tr></thead><tbody>")
	body.WriteString("<tr><td class=\"center\">1</td><td>" + orDash(session.CargoName) + "</td>" +
		"<td class=\"center\">т</td><td class=\"num\">" + formatNumber(quantity, 3) + "</td>" +
		"<td class=\"num\">" + formatNumber(price, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(sumWithoutVat, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(nds, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(vatSum, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(total, 2) + "</td></tr>")
	body.WriteString("<tr class=\"total\"><td colspan=\"3\" class=\"num\">Всего</td>" +
		"<td class=\"num\">" + formatNumber(quantity, 3) + "</td><td></td>" +
		"<td class=\"num\">" + formatNumber(sumWithoutVat, 2) + "</td><td></td>" +
		"<td class=\"num\">" + formatNumber(vatSum, 2) + "</td>" +
		"<td class=\"num\">" + formatNumber(total, 2) + "</td></tr>")
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"words\">Всего отпущено на сумму: <b>" + rublesWords(totalKopecks) + "</b></p>")
	body.WriteString("<p class=\"words\">НДС: " + formatNumber(vatSum, 2) + " руб. (" +
		formatNumber(nds, 2) + " %)</p>")

	body.WriteString("<h2>Подписи сторон</h2>")
	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Отпуск груза разрешил: <span class=\"sign-line\"></span> " +
		signature(doc.CargoReleaseName, doc.CargoReleaseWorkAs) + "</div>")
	body.WriteString("<div>Главный бухгалтер: <span class=\"sign-line\"></span> " +
		signature(doc.AccountantName, "") + "</div>")
	body.WriteString("<div>Груз отпустил: <span class=\"sign-line\"></span> " +
		signature(doc.ShipperName, doc.ShipperWorkAs) + "</div>")
	body.WriteString("<div>Груз принял: <span class=\"sign-line\"></span> " +
		signature(doc.CargoPickUpName, doc.CargoPickUpWorkAs) + "</div>")
	body.WriteString("<div>Груз получил: <span class=\"sign-line\"></span> " +
		signature(doc.ConsigneeName, doc.ConsigneeWorkAs) + "</div>")
	body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
	body.WriteString("</div>")

	if doc.ProxyName != "" || doc.ProxyGivenBy != "" || doc.ProxyDate != "" {
		body.WriteString("<p class=\"words\">Доверенность: № " + orDash(doc.ProxyName) +
			" от " + orDash(doc.ProxyDate) + ", выдана " + orDash(doc.ProxyGivenBy) + ".</p>")
	}

	body.WriteString("<p class=\"note\">Бланк ТОРГ-12 сформирован веб-системой весовой. " +
		"Одна строка товара соответствует грузу рейса; цена берётся из справочника груза или реквизитов документа.</p>")

	return wrapHTML("ТОРГ-12 №"+number, false, body.String())
}

func signature(name, position string) string {
	parts := make([]string, 0, 2)
	if name != "" {
		parts = append(parts, esc(name))
	}
	if position != "" {
		parts = append(parts, "("+esc(position)+")")
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func numberOrDash(value float64) string {
	if value <= 0 {
		return "—"
	}
	return formatNumber(value, 2)
}

func carrierCell(details *dto.SessionDetails, name string) string {
	if name == "" {
		return "—"
	}
	return esc(name) + requisites(details.CarrierINN, details.CarrierKPP)
}

// parseDocumentDate разбирает дату из формы; пустая строка даёт fallback.
func parseDocumentDate(value string, fallback time.Time) time.Time {
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(value), time.UTC)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseOptionalDate(value string) *time.Time {
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(value), time.UTC)
	if err != nil {
		return nil
	}
	return &parsed
}
