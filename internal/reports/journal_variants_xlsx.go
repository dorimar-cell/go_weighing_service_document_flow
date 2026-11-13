package reports

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"go_weighing_service/internal/dto"
)

type xlsxStyles struct {
	title  int
	meta   int
	header int
	cell   int
	number int
	total  int
}

func buildXLSXStyles(file *excelize.File) xlsxStyles {
	title, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	meta, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Color: "#444444"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	header, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 9},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E8EEF7"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: thinBorder(),
	})
	cell, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    thinBorder(),
	})
	numberFormat := "#,##0.000"
	number, _ := file.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 9},
		CustomNumFmt: &numberFormat,
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		Border:       thinBorder(),
	})
	total, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 9},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border:    thinBorder(),
	})

	return xlsxStyles{title: title, meta: meta, header: header, cell: cell, number: number, total: total}
}

func writeXLSXHeader(
	file *excelize.File,
	sheet string,
	title string,
	meta string,
	headers []string,
) (xlsxStyles, int, string) {
	styles := buildXLSXStyles(file)
	lastColumn, _ := excelize.ColumnNumberToName(len(headers))

	_ = file.MergeCell(sheet, "A1", lastColumn+"1")
	_ = file.SetCellValue(sheet, "A1", title)
	_ = file.SetCellStyle(sheet, "A1", lastColumn+"1", styles.title)

	_ = file.MergeCell(sheet, "A2", lastColumn+"2")
	_ = file.SetCellValue(sheet, "A2", meta)
	_ = file.SetCellStyle(sheet, "A2", lastColumn+"2", styles.meta)

	headerRow := 4
	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, headerRow)
		_ = file.SetCellValue(sheet, cell, header)
	}
	_ = file.SetCellStyle(sheet, "A4", lastColumn+"4", styles.header)
	_ = file.SetRowHeight(sheet, headerRow, 32)

	return styles, headerRow, lastColumn
}

func newReportFile(sheet string) (*excelize.File, error) {
	file := excelize.NewFile()
	if _, err := file.NewSheet(sheet); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("new sheet: %w", err)
	}
	_ = file.DeleteSheet("Sheet1")
	return file, nil
}

// renderFirstWeightXLSX — печатный отчёт по транспорту на первом взвешивании.
func renderFirstWeightXLSX(items []dto.SessionDTO, filter journalFilter) ([]byte, error) {
	file, err := newReportFile("Первое взвешивание")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	sheet := "Первое взвешивание"
	headers := []string{
		"№", "Дата", "Время", "Транспорт", "Прицеп", "Водитель", "Груз",
		"Контрагент", "Пропуск", "Вес 1-го, т",
	}
	styles, headerRow, lastColumn := writeXLSXHeader(
		file, sheet, "Транспорт, прошедший первое взвешивание", journalMetaText(filter), headers,
	)

	totalWeight := 0.0
	for index, item := range items {
		row := headerRow + 1 + index
		totalWeight += item.FirstWeightKg
		values := []any{
			item.ID,
			item.FirstMeasuredAt.Format("02.01.2006"),
			item.FirstMeasuredAt.Format("15:04"),
			item.TransportNumber,
			item.TrailerNumber,
			item.DriverName,
			item.CargoName,
			item.CounterpartyName,
			item.PassNumber,
			item.FirstWeightKg / 1000,
		}
		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = file.SetCellValue(sheet, cell, value)
			_ = file.SetCellStyle(sheet, cell, cell, styles.cell)
		}
		weightCell, _ := excelize.CoordinatesToCellName(len(values), row)
		_ = file.SetCellStyle(sheet, weightCell, weightCell, styles.number)
	}

	totalRow := headerRow + len(items) + 1
	firstCell, _ := excelize.CoordinatesToCellName(1, totalRow)
	lastCell, _ := excelize.CoordinatesToCellName(len(headers), totalRow)
	_ = file.MergeCell(sheet, firstCell, lastCell)
	_ = file.SetCellValue(sheet, firstCell, fmt.Sprintf(
		"Рейсов: %d. Итого вес 1-го взвешивания, т: %s",
		len(items), formatNumber(totalWeight/1000, 3),
	))
	_ = file.SetCellStyle(sheet, firstCell, lastCell, styles.total)

	widths := []float64{6, 12, 8, 18, 14, 22, 20, 24, 16, 14}
	for index, width := range widths {
		column, _ := excelize.ColumnNumberToName(index + 1)
		_ = file.SetColWidth(sheet, column, column, width)
	}
	_ = lastColumn

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write first weight xlsx: %w", err)
	}
	return buffer.Bytes(), nil
}

// renderTransportJournalXLSX — журнал по номеру ТС в Excel.
func renderTransportJournalXLSX(items []dto.SessionDTO, filter journalFilter) ([]byte, error) {
	file, err := newReportFile("Журнал по ТС")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	sheet := "Журнал по ТС"
	headers := []string{
		"Транспорт", "Дата", "№ рейса", "Прицеп", "Водитель", "Груз", "Склад", "Нетто, т",
	}
	styles, headerRow, _ := writeXLSXHeader(
		file, sheet, "Журнал взвешиваний по номеру ТС", journalMetaText(filter), headers,
	)

	totalNetto := 0.0
	for index, item := range items {
		row := headerRow + 1 + index
		netto := 0.0
		if item.NetWeightKg != nil {
			netto = *item.NetWeightKg
		}
		totalNetto += netto
		values := []any{
			item.TransportNumber,
			item.FirstMeasuredAt.Format("02.01.2006"),
			item.ID,
			item.TrailerNumber,
			item.DriverName,
			item.CargoName,
			item.StorageName,
			netto / 1000,
		}
		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = file.SetCellValue(sheet, cell, value)
			_ = file.SetCellStyle(sheet, cell, cell, styles.cell)
		}
		nettoCell, _ := excelize.CoordinatesToCellName(len(values), row)
		_ = file.SetCellStyle(sheet, nettoCell, nettoCell, styles.number)
	}

	totalRow := headerRow + len(items) + 1
	firstCell, _ := excelize.CoordinatesToCellName(1, totalRow)
	lastCell, _ := excelize.CoordinatesToCellName(len(headers), totalRow)
	_ = file.MergeCell(sheet, firstCell, lastCell)
	_ = file.SetCellValue(sheet, firstCell, fmt.Sprintf(
		"Всего рейсов: %d. Итого нетто, т: %s", len(items), formatNumber(totalNetto/1000, 3),
	))
	_ = file.SetCellStyle(sheet, firstCell, lastCell, styles.total)

	widths := []float64{18, 12, 10, 14, 22, 20, 18, 14}
	for index, width := range widths {
		column, _ := excelize.ColumnNumberToName(index + 1)
		_ = file.SetColWidth(sheet, column, column, width)
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write transport journal xlsx: %w", err)
	}
	return buffer.Bytes(), nil
}
