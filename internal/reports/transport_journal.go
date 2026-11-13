package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// renderFirstWeight — отчёт «Транспорт, прошедший первое взвешивание».
func renderFirstWeight(items []dto.SessionDTO, filter journalFilter) []byte {
	var body strings.Builder

	body.WriteString("<h1>Транспорт, прошедший первое взвешивание</h1>")
	body.WriteString("<p class=\"subtitle\">Открытые рейсы: ожидают второго взвешивания</p>")
	body.WriteString(journalMeta(filter))

	body.WriteString("<table><thead><tr>")
	headers := []string{
		"№", "Дата", "Время", "Транспорт", "Прицеп", "Водитель", "Груз",
		"Контрагент", "Пропуск", "Вес 1-го, т",
	}
	for _, header := range headers {
		body.WriteString("<th>" + esc(header) + "</th>")
	}
	body.WriteString("</tr></thead><tbody>")

	totalWeight := 0.0
	for _, item := range items {
		totalWeight += item.FirstWeightKg
		body.WriteString("<tr>")
		body.WriteString("<td class=\"center\">" + fmt.Sprintf("%d", item.ID) + "</td>")
		body.WriteString("<td class=\"center\">" + dateOnly(item.FirstMeasuredAt) + "</td>")
		body.WriteString("<td class=\"center\">" + timeOnly(item.FirstMeasuredAt) + "</td>")
		body.WriteString("<td>" + esc(item.TransportNumber) + "</td>")
		body.WriteString("<td>" + orDash(item.TrailerNumber) + "</td>")
		body.WriteString("<td>" + orDash(item.DriverName) + "</td>")
		body.WriteString("<td>" + orDash(item.CargoName) + "</td>")
		body.WriteString("<td>" + orDash(item.CounterpartyName) + "</td>")
		body.WriteString("<td>" + orDash(item.PassNumber) + "</td>")
		body.WriteString("<td class=\"num\">" + formatNumber(item.FirstWeightKg/1000, 3) + "</td>")
		body.WriteString("</tr>")
	}

	if len(items) == 0 {
		body.WriteString("<tr><td class=\"center muted\" colspan=\"10\">Открытых рейсов нет</td></tr>")
	} else {
		body.WriteString("<tr class=\"total\"><td colspan=\"9\" class=\"num\">Рейсов: " +
			fmt.Sprintf("%d", len(items)) + ". Итого вес 1-го взвешивания, т</td>" +
			"<td class=\"num\">" + formatNumber(totalWeight/1000, 3) + "</td></tr>")
	}
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"note\">В отчёт попадают рейсы без второго взвешивания и без отмены. " +
		"Веса приведены в тоннах с тремя знаками.</p>")

	return wrapHTML("Транспорт, прошедший первое взвешивание", true, body.String())
}

type transportGroup struct {
	Transport string
	Sessions  int
	NettoKg   float64
	Date      string
}

// renderTransportJournal — журнал взвешиваний с группировкой по номеру ТС.
func renderTransportJournal(items []dto.SessionDTO, filter journalFilter) []byte {
	var body strings.Builder

	body.WriteString("<h1>Журнал взвешиваний по номеру ТС</h1>")
	body.WriteString("<p class=\"subtitle\">Закрытые рейсы с группировкой по транспорту и дате</p>")
	body.WriteString(journalMeta(filter))

	body.WriteString("<table><thead><tr>")
	headers := []string{
		"Транспорт", "Дата", "№ рейса", "Прицеп", "Водитель", "Груз", "Склад", "Нетто, т",
	}
	for _, header := range headers {
		body.WriteString("<th>" + esc(header) + "</th>")
	}
	body.WriteString("</tr></thead><tbody>")

	var (
		totalNetto    float64
		totalSessions int
		current       transportGroup
		groupRows     int
	)

	flushGroup := func() {
		if current.Transport == "" || groupRows < 2 {
			current = transportGroup{}
			groupRows = 0
			return
		}
		body.WriteString("<tr class=\"total\"><td colspan=\"7\" class=\"num\">Итого по ТС " +
			esc(current.Transport) + ": рейсов " + fmt.Sprintf("%d", current.Sessions) +
			", нетто " + formatNumber(current.NettoKg/1000, 3) + " т</td><td class=\"num\">" +
			formatNumber(current.NettoKg/1000, 3) + "</td></tr>")
		current = transportGroup{}
		groupRows = 0
	}

	for _, item := range items {
		if item.TransportNumber != current.Transport {
			flushGroup()
			current.Transport = item.TransportNumber
		}

		groupRows++
		current.Sessions++
		netto := 0.0
		if item.NetWeightKg != nil {
			netto = *item.NetWeightKg
		}
		current.NettoKg += netto
		totalNetto += netto
		totalSessions++

		body.WriteString("<tr>")
		body.WriteString("<td>" + esc(item.TransportNumber) + "</td>")
		body.WriteString("<td class=\"center\">" + dateOnly(item.FirstMeasuredAt) + "</td>")
		body.WriteString("<td class=\"center\">" + fmt.Sprintf("%d", item.ID) + "</td>")
		body.WriteString("<td>" + orDash(item.TrailerNumber) + "</td>")
		body.WriteString("<td>" + orDash(item.DriverName) + "</td>")
		body.WriteString("<td>" + orDash(item.CargoName) + "</td>")
		body.WriteString("<td>" + orDash(item.StorageName) + "</td>")
		body.WriteString("<td class=\"num\">" + tonnesPtr(item.NetWeightKg) + "</td>")
		body.WriteString("</tr>")
	}
	flushGroup()

	if len(items) == 0 {
		body.WriteString("<tr><td class=\"center muted\" colspan=\"8\">За выбранный период рейсов нет</td></tr>")
	} else {
		body.WriteString("<tr class=\"total\"><td colspan=\"7\" class=\"num\">Всего: рейсов " +
			fmt.Sprintf("%d", totalSessions) + ", нетто " + formatNumber(totalNetto/1000, 3) +
			" т</td><td class=\"num\">" + formatNumber(totalNetto/1000, 3) + "</td></tr>")
	}
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"note\">Учитываются закрытые рейсы. Подытог по ТС выводится, " +
		"если у транспорта больше одного рейса за период.</p>")

	return wrapHTML("Журнал взвешиваний по номеру ТС", true, body.String())
}
