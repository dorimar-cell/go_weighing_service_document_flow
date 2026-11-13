package reports

import (
	"fmt"
	"strings"
)

// code128Patterns — ширины полос Code128 (6 элементов, у stop — 7).
var code128Patterns = []string{
	"212222", "222122", "222221", "121223", "121322", "131222", "122213", "122312",
	"132212", "221213", "221312", "231212", "112232", "122132", "122231", "113222",
	"123122", "123221", "223211", "221132", "221231", "213212", "223112", "312131",
	"311222", "321122", "321221", "312212", "322112", "322211", "212123", "212321",
	"232121", "111323", "131123", "131321", "112313", "132113", "132311", "211313",
	"231113", "231311", "112133", "112331", "132131", "113123", "113321", "133121",
	"313121", "211331", "231131", "213113", "213311", "213131", "311123", "311321",
	"331121", "312113", "312311", "332111", "314111", "221411", "431111", "111224",
	"111422", "121124", "121421", "141122", "141221", "112214", "112412", "122114",
	"122411", "142112", "142211", "241211", "221114", "413111", "241112", "134111",
	"111242", "121142", "121241", "114212", "124112", "124211", "411212", "421112",
	"421211", "212141", "214121", "412121", "111143", "111341", "131141", "114113",
	"114311", "411113", "411311", "113141", "114131", "311141", "411131", "211412",
	"211214", "211232", "2331112",
}

// BarcodeSVG кодирует ASCII-строку в Code128B и возвращает inline SVG.
// Внешние сервисы и шрифты не используются: только прямоугольники.
func BarcodeSVG(value string, module int, height int) string {
	if module <= 0 {
		module = 2
	}
	if height <= 0 {
		height = 56
	}

	clean := make([]rune, 0, len(value))
	for _, char := range value {
		if char < 32 || char > 126 {
			continue
		}
		clean = append(clean, char)
	}
	if len(clean) == 0 {
		clean = []rune{'0'}
	}

	codes := make([]int, 0, len(clean)+3)
	codes = append(codes, 104) // Start B
	for _, char := range clean {
		codes = append(codes, int(char)-32)
	}

	checksum := 104
	for index, code := range codes[1:] {
		checksum += code * (index + 1)
	}
	codes = append(codes, checksum%103, 106) // checksum + stop

	totalModules := 20 // quiet zones
	for _, code := range codes {
		totalModules += patternWidth(code128Patterns[code])
	}

	width := totalModules * module
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="Штрих-код %s">`,
		width, height, width, height, esc(value),
	))
	builder.WriteString(fmt.Sprintf(
		`<rect x="0" y="0" width="%d" height="%d" fill="#ffffff"/>`, width, height,
	))

	x := 10 * module
	for _, code := range codes {
		pattern := code128Patterns[code]
		bar := true
		for _, digit := range pattern {
			run := int(digit-'0') * module
			if bar {
				builder.WriteString(fmt.Sprintf(
					`<rect x="%d" y="0" width="%d" height="%d" fill="#000000"/>`,
					x, run, height,
				))
			}
			x += run
			bar = !bar
		}
	}

	builder.WriteString("</svg>")
	return builder.String()
}

func patternWidth(pattern string) int {
	width := 0
	for _, digit := range pattern {
		width += int(digit - '0')
	}
	return width
}
