package reports

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"go_weighing_service/internal/dto"
)

type aggregateFilter struct {
	Title        string
	Subtitle     string
	ShowSubGroup bool
	DateFrom     string
	DateTo       string
	Organization string
}

type aggregateGroupTotal struct {
	NettoKg  float64
	Sessions int
}

func aggregateMeta(filter aggregateFilter) string {
	if filter.DateFrom == "" && filter.DateTo == "" {
		return "Период: весь журнал"
	}
	from := filter.DateFrom
	to := filter.DateTo
	if from == "" {
		from = "начало"
	} else {
		from = prettyDate(from)
	}
	if to == "" {
		to = "по настоящее время"
	} else {
		to = prettyDate(to)
	}
	return "Период: " + from + " — " + to
}

func renderAggregate(rows []dto.AggregateRow, filter aggregateFilter) []byte {
	var body strings.Builder

	body.WriteString("<h1>" + esc(filter.Title) + "</h1>")
	if filter.Subtitle != "" {
		body.WriteString("<p class=\"subtitle\">" + esc(filter.Subtitle) + "</p>")
	}

	body.WriteString("<div class=\"meta\">")
	if filter.Organization != "" {
		body.WriteString("<span><b>" + esc(filter.Organization) + "</b></span>")
	}
	body.WriteString("<span>" + esc(aggregateMeta(filter)) + "</span>")
	body.WriteString("<span>Сформирован: " + nowInLocation().Format("02.01.2006 15:04") + "</span>")
	body.WriteString("</div>")

	body.WriteString("<table><thead><tr>")
	body.WriteString("<th>№</th><th>" + esc(aggregateGroupLabel(filter)) + "</th>")
	if filter.ShowSubGroup {
		body.WriteString("<th>Груз / операция</th>")
	}
	body.WriteString("<th>Рейсов</th><th>Нетто, т</th><th>Среднее нетто, т</th></tr></thead><tbody>")

	var (
		totalNetto    float64
		totalSessions int
		currentGroup  string
		group         aggregateGroupTotal
		groupRows     int
		rowNumber     int
	)

	flushGroup := func() {
		if currentGroup == "" || groupRows < 2 {
			group = aggregateGroupTotal{}
			groupRows = 0
			return
		}
		body.WriteString("<tr class=\"total\"><td colspan=\"" + fmt.Sprintf("%d", aggregateColumnCount(filter)) +
			"\" class=\"num\">Итого по группе «" + esc(currentGroup) + "»: рейсов " +
			fmt.Sprintf("%d", group.Sessions) + ", нетто " + formatNumber(group.NettoKg/1000, 3) + " т</td></tr>")
		group = aggregateGroupTotal{}
		groupRows = 0
	}

	for _, item := range rows {
		if item.Group != currentGroup {
			flushGroup()
			currentGroup = item.Group
		}

		rowNumber++
		groupRows++
		group.NettoKg += item.NettoKg
		group.Sessions += item.Sessions
		totalNetto += item.NettoKg
		totalSessions += item.Sessions

		average := 0.0
		if item.Sessions > 0 {
			average = item.NettoKg / 1000 / float64(item.Sessions)
		}

		body.WriteString("<tr>")
		body.WriteString("<td class=\"center\">" + fmt.Sprintf("%d", rowNumber) + "</td>")
		body.WriteString("<td>" + esc(item.Group) + "</td>")
		if filter.ShowSubGroup {
			body.WriteString("<td>" + esc(item.SubGroup) + "</td>")
		}
		body.WriteString("<td class=\"num\">" + fmt.Sprintf("%d", item.Sessions) + "</td>")
		body.WriteString("<td class=\"num\">" + formatNumber(item.NettoKg/1000, 3) + "</td>")
		body.WriteString("<td class=\"num\">" + formatNumber(average, 3) + "</td>")
		body.WriteString("</tr>")
	}
	flushGroup()

	if len(rows) == 0 {
		body.WriteString("<tr><td class=\"center muted\" colspan=\"" +
			fmt.Sprintf("%d", aggregateColumnCount(filter)) +
			"\">За выбранный период закрытых рейсов нет</td></tr>")
	} else {
		body.WriteString("<tr class=\"total\"><td colspan=\"" +
			fmt.Sprintf("%d", aggregateColumnCount(filter)) +
			"\" class=\"num\">Всего: рейсов " + fmt.Sprintf("%d", totalSessions) +
			", нетто " + formatNumber(totalNetto/1000, 3) + " т</td></tr>")
	}
	body.WriteString("</tbody></table>")

	body.WriteString("<p class=\"note\">Учитываются только закрытые рейсы (два взвешивания или поточный режим). " +
		"Нетто = |первое − второе взвешивание|.</p>")

	return wrapHTML(filter.Title, false, body.String())
}

func aggregateColumnCount(filter aggregateFilter) int {
	if filter.ShowSubGroup {
		return 6
	}
	return 5
}

func aggregateGroupLabel(filter aggregateFilter) string {
	switch {
	case strings.Contains(filter.Title, "контрагент"):
		return "Контрагент"
	case strings.Contains(filter.Title, "перевозчик"):
		return "Перевозчик"
	case strings.Contains(filter.Title, "водител"):
		return "Водитель"
	case strings.Contains(filter.Title, "склад"):
		return "Склад"
	case strings.Contains(filter.Title, "груз"):
		return "Груз"
	default:
		return "Группа"
	}
}

// renderAggregateXLSX собирает агрегатный отчёт книгой Excel.
func renderAggregateXLSX(rows []dto.AggregateRow, filter aggregateFilter) ([]byte, error) {
	file := excelize.NewFile()
	defer func() { _ = file.Close() }()

	sheet := "Отчёт"
	if _, err := file.NewSheet(sheet); err != nil {
		return nil, fmt.Errorf("new sheet: %w", err)
	}
	_ = file.DeleteSheet("Sheet1")

	titleStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	metaStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Color: "#444444"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	headerStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 9},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E8EEF7"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: thinBorder(),
	})
	cellStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    thinBorder(),
	})
	numberStyle, _ := file.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 9},
		CustomNumFmt: ptrString("#,##0.000"),
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		Border:       thinBorder(),
	})
	countStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border:    thinBorder(),
	})
	totalStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 9},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border:    thinBorder(),
	})

	columns := []string{"№", aggregateGroupLabel(filter)}
	if filter.ShowSubGroup {
		columns = append(columns, "Груз / операция")
	}
	columns = append(columns, "Рейсов", "Нетто, т", "Среднее нетто, т")
	lastColumn, _ := excelize.ColumnNumberToName(len(columns))

	_ = file.MergeCell(sheet, "A1", lastColumn+"1")
	_ = file.SetCellValue(sheet, "A1", filter.Title)
	_ = file.SetCellStyle(sheet, "A1", lastColumn+"1", titleStyle)

	_ = file.MergeCell(sheet, "A2", lastColumn+"2")
	_ = file.SetCellValue(sheet, "A2", aggregateMeta(filter))
	_ = file.SetCellStyle(sheet, "A2", lastColumn+"2", metaStyle)

	headerRow := 4
	for index, column := range columns {
		cell, _ := excelize.CoordinatesToCellName(index+1, headerRow)
		_ = file.SetCellValue(sheet, cell, column)
	}
	_ = file.SetCellStyle(sheet, "A4", lastColumn+"4", headerStyle)

	for index, item := range rows {
		row := headerRow + 1 + index
		values := []any{index + 1, item.Group}
		if filter.ShowSubGroup {
			values = append(values, item.SubGroup)
		}
		values = append(values, item.Sessions, item.NettoKg/1000, averageNetto(item))

		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = file.SetCellValue(sheet, cell, value)
			_ = file.SetCellStyle(sheet, cell, cell, cellStyle)
		}
		for _, column := range []int{len(values) - 2, len(values) - 1, len(values)} {
			cell, _ := excelize.CoordinatesToCellName(column, row)
			_ = file.SetCellStyle(sheet, cell, cell, numberStyle)
		}
		countCell, _ := excelize.CoordinatesToCellName(len(values)-2, row)
		_ = file.SetCellStyle(sheet, countCell, countCell, countStyle)
	}

	totalRow := headerRow + len(rows) + 1
	firstCell, _ := excelize.CoordinatesToCellName(1, totalRow)
	lastCell, _ := excelize.CoordinatesToCellName(len(columns), totalRow)
	_ = file.MergeCell(sheet, firstCell, lastCell)

	totalNetto := 0.0
	totalSessions := 0
	for _, item := range rows {
		totalNetto += item.NettoKg
		totalSessions += item.Sessions
	}
	_ = file.SetCellValue(sheet, firstCell, fmt.Sprintf(
		"Всего: рейсов %d, нетто %s т", totalSessions, formatNumber(totalNetto/1000, 3),
	))
	_ = file.SetCellStyle(sheet, firstCell, lastCell, totalStyle)

	for index, width := range []float64{6, 34, 26, 10, 14, 18} {
		if index >= len(columns) {
			break
		}
		column, _ := excelize.ColumnNumberToName(index + 1)
		_ = file.SetColWidth(sheet, column, column, width)
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write aggregate xlsx: %w", err)
	}

	return buffer.Bytes(), nil
}

func averageNetto(item dto.AggregateRow) float64 {
	if item.Sessions == 0 {
		return 0
	}
	return item.NettoKg / 1000 / float64(item.Sessions)
}

func ptrString(value string) *string {
	return &value
}
