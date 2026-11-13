package reports

import (
	"fmt"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
)

type journalFilter struct {
	DateFrom     string
	DateTo       string
	StatusLabel  string
	Search       string
	Organization string
}

func journalMeta(filter journalFilter) string {
	var period strings.Builder
	period.WriteString("<div class=\"meta\">")
	if filter.Organization != "" {
		period.WriteString("<span><b>" + esc(filter.Organization) + "</b></span>")
	}
	from := filter.DateFrom
	to := filter.DateTo
	if from == "" {
		from = "начало"
	}
	if to == "" {
		to = "по настоящее время"
	}
	if filter.DateFrom == "" && filter.DateTo == "" {
		period.WriteString("<span>Период: весь журнал</span>")
	} else {
		period.WriteString("<span>Период: " + esc(from) + " — " + esc(to) + "</span>")
	}
	period.WriteString("<span>Статус: " + esc(filter.StatusLabel) + "</span>")
	if filter.Search != "" {
		period.WriteString("<span>Поиск: " + esc(filter.Search) + "</span>")
	}
	period.WriteString("<span>Сформирован: " + time.Now().Format("02.01.2006 15:04") + "</span>")
	period.WriteString("</div>")

	return period.String()
}

func renderJournal(items []dto.SessionDTO, filter journalFilter) []byte {
	var body strings.Builder
	body.WriteString("<h1>Журнал взвешиваний</h1>")
	body.WriteString("<p class=\"subtitle\">Весовая: рейсы за выбранный период</p>")
	body.WriteString(journalMeta(filter))

	body.WriteString("<table><thead><tr>")
	headers := []string{
		"№", "Статус", "Транспорт", "Прицеп", "Водитель", "Груз", "Контрагент",
		"Перевозчик", "Склад", "Тип операции", "1-е взвешивание", "Вес 1, т",
		"2-е взвешивание", "Вес 2, т", "Брутто, т", "Тара, т", "Нетто, т",
	}
	for _, header := range headers {
		body.WriteString("<th>" + esc(header) + "</th>")
	}
	body.WriteString("</tr></thead><tbody>")

	var totalNetto float64
	for _, item := range items {
		brutto, tare := bruttoTare(&item)
		if item.NetWeightKg != nil {
			totalNetto += *item.NetWeightKg
		}

		body.WriteString("<tr>")
		cells := []struct {
			value string
			class string
		}{
			{fmt.Sprintf("%d", item.ID), "center"},
			{statusLabel(item.Status), "center"},
			{esc(item.TransportNumber), ""},
			{orDash(item.TrailerNumber), ""},
			{orDash(item.DriverName), ""},
			{orDash(item.CargoName), ""},
			{orDash(item.CounterpartyName), ""},
			{orDash(item.CarrierName), ""},
			{orDash(item.StorageName), ""},
			{orDash(item.OperationTypeName), ""},
			{dateOnly(item.FirstMeasuredAt) + " " + timeOnly(item.FirstMeasuredAt), "center"},
			{formatNumber(item.FirstWeightKg/1000, 3), "num"},
			{dateOnlyPtr(item.SecondMeasuredAt) + " " + timeOnlyPtr(item.SecondMeasuredAt), "center"},
			{tonnesPtr(item.SecondWeightKg), "num"},
			{formatNumber(brutto/1000, 3), "num"},
			{formatNumber(tare/1000, 3), "num"},
			{tonnesPtr(item.NetWeightKg), "num"},
		}
		for _, cell := range cells {
			if cell.class == "" {
				body.WriteString("<td>" + cell.value + "</td>")
				continue
			}
			body.WriteString("<td class=\"" + cell.class + "\">" + cell.value + "</td>")
		}
		body.WriteString("</tr>")
	}

	if len(items) == 0 {
		body.WriteString("<tr><td class=\"center muted\" colspan=\"17\">За выбранный период рейсов нет</td></tr>")
	} else {
		body.WriteString("<tr class=\"total\"><td colspan=\"16\" class=\"num\">Итого нетто, т</td><td class=\"num\">" +
			formatNumber(totalNetto/1000, 3) + "</td></tr>")
	}
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"note\">Брутто — большее из двух взвешиваний, тара — меньшее. " +
		"Веса приведены в тоннах с тремя знаками после запятой. Всего строк: " +
		fmt.Sprintf("%d", len(items)) + ".</p>")

	return wrapHTML("Журнал взвешиваний", true, body.String())
}

func statusLabel(status string) string {
	switch status {
	case "open":
		return "Открыт"
	case "closed":
		return "Закрыт"
	case "cancelled":
		return "Отменён"
	default:
		return status
	}
}
