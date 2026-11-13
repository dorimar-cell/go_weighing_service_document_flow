package reports

import (
	"html"
	"strconv"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
)

// reportLocation — часовой пояс печатных бланков. Задаётся WEIGHING_TZ.
var reportLocation = time.UTC

// SetTimezone устанавливает часовой пояс для дат и времени в отчётах.
func SetTimezone(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	location, err := time.LoadLocation(name)
	if err != nil {
		return err
	}

	reportLocation = location
	return nil
}

func inLocation(value time.Time) time.Time {
	return value.In(reportLocation)
}

func nowInLocation() time.Time {
	return time.Now().In(reportLocation)
}

// wrapHTML оборачивает готовое тело бланка в самодостаточный HTML-документ
// с печатными стилями A4. Тело собирается функциями пакета и уже экранировано.
func wrapHTML(title string, landscape bool, body string) []byte {
	orientation := "portrait"
	if landscape {
		orientation = "landscape"
	}

	var builder strings.Builder
	builder.WriteString("<!DOCTYPE html><html lang=\"ru\"><head><meta charset=\"utf-8\">")
	builder.WriteString("<title>" + esc(title) + "</title><style>")
	builder.WriteString(`
@page { size: A4 ` + orientation + `; margin: 10mm; }
* { box-sizing: border-box; }
body { margin: 0; padding: 0; color: #000; background: #fff;
       font-family: "Times New Roman", "Liberation Serif", Arial, sans-serif; font-size: 10pt; }
h1 { margin: 0 0 2mm; font-size: 15pt; text-align: center; text-transform: uppercase; letter-spacing: .04em; }
h2 { margin: 4mm 0 2mm; font-size: 11pt; }
.subtitle { margin: 0 0 3mm; text-align: center; font-size: 10pt; }
.meta { margin: 0 0 3mm; font-size: 9pt; line-height: 1.5; }
.meta span { margin-right: 6mm; }
table { width: 100%; border-collapse: collapse; }
th, td { border: 1px solid #000; padding: 1.2mm 1.5mm; vertical-align: top; font-size: 9pt; }
th { text-align: center; font-weight: 600; background: #f2f2f2; }
thead { display: table-header-group; }
tr { page-break-inside: avoid; }
.num { text-align: right; white-space: nowrap; }
.center { text-align: center; }
.muted { color: #444; }
.total td { font-weight: 600; background: #f7f7f7; }
.pairs { width: 100%; border-collapse: collapse; margin-bottom: 3mm; }
.pairs td { border: 1px solid #000; }
.pairs .label { width: 32%; color: #333; }
.signs { display: grid; grid-template-columns: 1fr 1fr; gap: 6mm 10mm; margin-top: 6mm; font-size: 9pt; }
.sign-line { border-bottom: 1px solid #000; display: inline-block; min-width: 45mm; }
.note { margin-top: 5mm; font-size: 8pt; color: #555; }
.photos { display: grid; grid-template-columns: 1fr 1fr; gap: 4mm; margin-top: 3mm; }
.photos figure { margin: 0; }
.photos img { width: 100%; height: auto; border: 1px solid #000; }
.photos figcaption { font-size: 8pt; text-align: center; margin-top: 1mm; }
.words { margin-top: 3mm; font-size: 9pt; }
.talon-page { page-break-after: always; }
.talon-page:last-child { page-break-after: auto; }
`)
	builder.WriteString("</style></head><body>")
	builder.WriteString(body)
	builder.WriteString("</body></html>")

	return []byte(builder.String())
}

func esc(value string) string {
	return html.EscapeString(value)
}

func orDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "—"
	}
	return esc(value)
}

// formatNumber форматирует число с разделителем разрядов и запятой.
func formatNumber(value float64, decimals int) string {
	negative := value < 0
	if negative {
		value = -value
	}

	raw := strconv.FormatFloat(value, 'f', decimals, 64)
	integer := raw
	fraction := ""
	if dot := strings.Index(raw, "."); dot >= 0 {
		integer = raw[:dot]
		fraction = raw[dot+1:]
	}

	var grouped strings.Builder
	for index, digit := range integer {
		if index > 0 && (len(integer)-index)%3 == 0 {
			grouped.WriteString("\u00a0")
		}
		grouped.WriteRune(digit)
	}

	result := grouped.String()
	if fraction != "" {
		result += "," + fraction
	}
	if negative && strings.Trim(result, "0,\u00a0") != "" {
		result = "−" + result
	}

	return result
}

// tonnes переводит килограммы в тонны для печати.
func tonnes(kilograms float64) string {
	return formatNumber(kilograms/1000, 3)
}

func tonnesPtr(kilograms *float64) string {
	if kilograms == nil {
		return "—"
	}
	return tonnes(*kilograms)
}

func dateOnly(value time.Time) string {
	if value.IsZero() {
		return "—"
	}
	return inLocation(value).Format("02.01.2006")
}

func dateOnlyPtr(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "—"
	}
	return inLocation(*value).Format("02.01.2006")
}

func timeOnly(value time.Time) string {
	if value.IsZero() {
		return "—"
	}
	return inLocation(value).Format("15:04")
}

func timeOnlyPtr(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "—"
	}
	return inLocation(*value).Format("15:04")
}

func dateTimePtr(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "—"
	}
	return inLocation(*value).Format("02.01.2006 15:04")
}

func dateOnlyNow() string {
	return nowInLocation().Format("02.01.2006")
}

// prettyDate переводит дату из формы (ГГГГ-ММ-ДД) в печатный вид.
func prettyDate(value string) string {
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(value), time.UTC)
	if err != nil {
		return value
	}
	return parsed.Format("02.01.2006")
}

func int64Text(value int64) string {
	return strconv.FormatInt(value, 10)
}

func weightOrZero(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func nettoValue(session *dto.SessionDTO) float64 {
	if session.NetWeightKg == nil {
		return 0
	}
	return *session.NetWeightKg
}

// bruttoTare возвращает брутто и тару по легаси-правилу: большее — брутто.
func bruttoTare(session *dto.SessionDTO) (brutto float64, tare float64) {
	if session.SecondWeightKg == nil {
		return session.FirstWeightKg, 0
	}
	second := *session.SecondWeightKg
	if session.FirstWeightKg >= second {
		return session.FirstWeightKg, second
	}
	return second, session.FirstWeightKg
}
