package reports

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"go_weighing_service/internal/dto"
)

// renderJournalXLSX собирает печатную книгу журнала: шапка, период, таблица, итог.
func renderJournalXLSX(items []dto.SessionDTO, filter journalFilter) ([]byte, error) {
	file := excelize.NewFile()
	defer func() { _ = file.Close() }()

	sheet := "Журнал"
	if _, err := file.NewSheet(sheet); err != nil {
		return nil, fmt.Errorf("new sheet: %w", err)
	}
	_ = file.DeleteSheet("Sheet1")

	titleStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	metaStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Color: "#444444"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
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
	weightFormat := "#,##0.000"
	weightStyle, _ := file.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 9},
		CustomNumFmt: &weightFormat,
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		Border:       thinBorder(),
	})
	totalStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 9},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border:    thinBorder(),
	})

	headers := []string{
		"№", "Статус", "Транспорт", "Прицеп", "Водитель", "Груз", "Контрагент",
		"Перевозчик", "Склад", "Тип операции", "1-е взвешивание", "Вес 1, т",
		"2-е взвешивание", "Вес 2, т", "Брутто, т", "Тара, т", "Нетто, т",
	}
	lastColumn, _ := excelize.ColumnNumberToName(len(headers))

	_ = file.MergeCell(sheet, "A1", lastColumn+"1")
	_ = file.SetCellValue(sheet, "A1", "Журнал взвешиваний")
	_ = file.SetCellStyle(sheet, "A1", lastColumn+"1", titleStyle)
	_ = file.SetRowHeight(sheet, 1, 26)

	_ = file.MergeCell(sheet, "A2", lastColumn+"2")
	_ = file.SetCellValue(sheet, "A2", journalMetaText(filter))
	_ = file.SetCellStyle(sheet, "A2", lastColumn+"2", metaStyle)

	headerRow := 4
	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, headerRow)
		_ = file.SetCellValue(sheet, cell, header)
	}
	_ = file.SetCellStyle(sheet, "A4", lastColumn+"4", headerStyle)
	_ = file.SetRowHeight(sheet, headerRow, 34)

	var totalNetto float64
	for index, item := range items {
		row := headerRow + 1 + index
		brutto, tare := bruttoTare(&item)
		if item.NetWeightKg != nil {
			totalNetto += *item.NetWeightKg
		}

		values := []any{
			item.ID,
			statusLabel(item.Status),
			item.TransportNumber,
			item.TrailerNumber,
			item.DriverName,
			item.CargoName,
			item.CounterpartyName,
			item.CarrierName,
			item.StorageName,
			item.OperationTypeName,
			item.FirstMeasuredAt.Format("02.01.2006 15:04"),
			item.FirstWeightKg / 1000,
			secondWeightTime(item.SecondMeasuredAt),
			weightOrNil(item.SecondWeightKg),
			brutto / 1000,
			tare / 1000,
			weightOrNil(item.NetWeightKg),
		}
		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = file.SetCellValue(sheet, cell, value)
			_ = file.SetCellStyle(sheet, cell, cell, cellStyle)
		}
		for _, column := range []int{12, 14, 15, 16, 17} {
			cell, _ := excelize.CoordinatesToCellName(column, row)
			_ = file.SetCellStyle(sheet, cell, cell, weightStyle)
		}
	}

	lastDataRow := headerRow + len(items)
	totalRow := lastDataRow + 1
	firstTotalCell, _ := excelize.CoordinatesToCellName(1, totalRow)
	lastTotalCell, _ := excelize.CoordinatesToCellName(len(headers), totalRow)
	_ = file.MergeCell(sheet, firstTotalCell, lastTotalCell)
	_ = file.SetCellValue(sheet, firstTotalCell,
		fmt.Sprintf("Всего строк: %d. Итого нетто, т: %s", len(items), formatNumber(totalNetto/1000, 3)))
	_ = file.SetCellStyle(sheet, firstTotalCell, lastTotalCell, totalStyle)

	widths := map[string]float64{
		"A": 6, "B": 10, "C": 16, "D": 13, "E": 22, "F": 20, "G": 22,
		"H": 18, "I": 16, "J": 15, "K": 18, "L": 11,
		"M": 18, "N": 11, "O": 11, "P": 11, "Q": 11,
	}
	for column, width := range widths {
		_ = file.SetColWidth(sheet, column, column, width)
	}

	_ = file.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      headerRow,
		TopLeftCell: fmt.Sprintf("A%d", headerRow+1),
		ActivePane:  "bottomLeft",
	})

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write journal xlsx: %w", err)
	}

	return buffer.Bytes(), nil
}

func journalMetaText(filter journalFilter) string {
	from := filter.DateFrom
	to := filter.DateTo
	if from == "" {
		from = "начало"
	}
	if to == "" {
		to = "по настоящее время"
	}

	text := ""
	if filter.Organization != "" {
		text = filter.Organization + ". "
	}
	text += "Период: " + from + " — " + to
	text += ". Статус: " + filter.StatusLabel
	if filter.Search != "" {
		text += ". Поиск: " + filter.Search
	}
	text += ". Сформирован: " + time.Now().Format("02.01.2006 15:04")

	return text
}

func thinBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "#999999", Style: 1},
		{Type: "right", Color: "#999999", Style: 1},
		{Type: "top", Color: "#999999", Style: 1},
		{Type: "bottom", Color: "#999999", Style: 1},
	}
}

func secondWeightTime(value *time.Time) any {
	if value == nil {
		return ""
	}
	return value.Format("02.01.2006 15:04")
}

func weightOrNil(value *float64) any {
	if value == nil {
		return nil
	}
	return *value / 1000
}
