package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// renderCargoAct собирает акт приёма-передачи товара в трёх экземплярах.
func renderCargoAct(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	copies := []struct {
		Title string
		Role  string
	}{
		{Title: "охраны", Role: "охрана"},
		{Title: "водителя", Role: "водитель"},
		{Title: "весовщика", Role: "весовщик"},
	}

	var body strings.Builder
	body.WriteString("<h1>Акт приёма-передачи товара</h1>")
	body.WriteString("<p class=\"subtitle\">Рейс №" + esc(int64Text(details.Session.ID)) +
		" от " + esc(documentDate(&details.Session)) + "</p>")
	body.WriteString("<div class=\"meta\">")
	if settings.OrganizationName != "" {
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
	}
	body.WriteString("<span>Экземпляров: 3</span>")
	body.WriteString("<span>Сформирован: " + nowInLocation().Format("02.01.2006 15:04") + "</span>")
	body.WriteString("</div>")

	for _, copy := range copies {
		body.WriteString("<section class=\"talon-page\">")
		body.WriteString("<h2>Экземпляр " + esc(copy.Title) + "</h2>")

		body.WriteString(pairRows([][2]string{
			{"Транспортное средство", esc(details.Session.TransportNumber)},
			{"Прицеп", orDash(details.Session.TrailerNumber)},
			{"Водитель", orDash(details.Session.DriverName)},
			{"Перевозчик", carrierText(details)},
			{"Контрагент", counterpartyText(details)},
			{"Склад", orDash(details.Session.StorageName)},
			{"Груз (фракция)", orDash(details.Session.CargoName)},
		}))

		body.WriteString("<table><thead><tr><th>Показатель</th><th>Значение</th></tr></thead><tbody>")
		body.WriteString("<tr><td>Нетто, т</td><td class=\"num\">" +
			tonnesPtr(details.Session.NetWeightKg) + "</td></tr>")
		body.WriteString("<tr><td>Объём, м³</td><td class=\"num\">&nbsp;</td></tr>")
		body.WriteString("</tbody></table>")

		body.WriteString("<div class=\"signs\">")
		body.WriteString("<div>Сдал (водитель): <span class=\"sign-line\"></span> " +
			orDash(details.Session.DriverName) + "</div>")
		body.WriteString("<div>Принял (" + esc(copy.Role) + "): <span class=\"sign-line\"></span></div>")
		body.WriteString("</div>")
		body.WriteString("</section>")
	}

	body.WriteString("<p class=\"note\">Акт сформирован веб-системой весовой; объём и подписи " +
		"заполняются при передаче товара. Один экземпляр — охране, один — водителю, один — весовщику.</p>")

	return wrapHTML("Акт приёма-передачи товара", false, body.String())
}

// renderUploadingCertificate собирает свидетельство о разгрузке/погрузке Т/С.
func renderUploadingCertificate(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	session := &details.Session
	brutto, tare := bruttoTare(session)

	documentNumber := session.PassNumber
	if documentNumber == "" {
		documentNumber = fmt.Sprintf("%d", session.ID)
	}

	var body strings.Builder
	if settings.OrganizationName != "" {
		body.WriteString("<div class=\"meta\"><span><b>" + esc(settings.OrganizationName) + "</b></span>")
		if settings.Address != "" {
			body.WriteString("<span>" + esc(settings.Address) + "</span>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<h1>Свидетельство о разгрузке/погрузке Т/С</h1>")
	body.WriteString("<p class=\"subtitle\">№ " + esc(documentNumber) + " от " +
		esc(documentDate(session)) + "</p>")

	body.WriteString(pairRows([][2]string{
		{"Тип операции", orDash(session.OperationTypeName)},
		{"Номер регистрации", esc(documentNumber)},
		{"Транспортное средство", esc(session.TransportNumber)},
		{"Прицеп", orDash(session.TrailerNumber)},
		{"Водитель", orDash(session.DriverName)},
		{"Продукт (груз)", orDash(session.CargoName)},
		{"Контрагент", counterpartyText(details)},
		{"Перевозчик", carrierText(details)},
		{"Склад", orDash(session.StorageName)},
	}))

	body.WriteString("<h2>Параметры взвешивания</h2>")
	body.WriteString("<table><thead><tr><th>Показатель</th><th>Кг</th><th>Тонн</th></tr></thead><tbody>")
	body.WriteString("<tr><td>Брутто</td><td class=\"num\">" + formatNumber(brutto, 1) +
		"</td><td class=\"num\">" + formatNumber(brutto/1000, 3) + "</td></tr>")
	body.WriteString("<tr><td>Тара</td><td class=\"num\">" + formatNumber(tare, 1) +
		"</td><td class=\"num\">" + formatNumber(tare/1000, 3) + "</td></tr>")
	body.WriteString("<tr class=\"total\"><td>Нетто</td><td class=\"num\">" +
		formatNumber(nettoValue(session), 1) + "</td><td class=\"num\">" +
		tonnesPtr(session.NetWeightKg) + "</td></tr>")
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"words\">Время операции: " + esc(dateOnly(session.FirstMeasuredAt)) +
		" " + esc(timeOnly(session.FirstMeasuredAt)) + "</p>")

	body.WriteString("<div class=\"signs\">")
	body.WriteString("<div>Весовщик: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
	body.WriteString("<div>Водитель: <span class=\"sign-line\"></span> " + orDash(session.DriverName) + "</div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Свидетельство сформировано веб-системой весовой и подтверждает " +
		"факт разгрузки/погрузки транспортного средства на весовой.</p>")

	return wrapHTML("Свидетельство №"+documentNumber, false, body.String())
}
