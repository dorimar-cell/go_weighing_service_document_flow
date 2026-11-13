package reports

import (
	"fmt"
	"strings"
)

var (
	unitsMale = []string{"", "один", "два", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять"}
	unitsFem  = []string{"", "одна", "две", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять"}
	teens     = []string{"десять", "одиннадцать", "двенадцать", "тринадцать", "четырнадцать",
		"пятнадцать", "шестнадцать", "семнадцать", "восемнадцать", "девятнадцать"}
	tensWords = []string{"", "", "двадцать", "тридцать", "сорок", "пятьдесят",
		"шестьдесят", "семьдесят", "восемьдесят", "девяносто"}
	hundredsWords = []string{"", "сто", "двести", "триста", "четыреста",
		"пятьсот", "шестьсот", "семьсот", "восемьсот", "девятьсот"}
)

type scaleWord struct {
	female bool
	one    string
	few    string
	many   string
}

var scales = []scaleWord{
	{false, "", "", ""},
	{true, "тысяча", "тысячи", "тысяч"},
	{false, "миллион", "миллиона", "миллионов"},
	{false, "миллиард", "миллиарда", "миллиардов"},
	{false, "триллион", "триллиона", "триллионов"},
}

func pluralRu(value int64, one, few, many string) string {
	value = value % 100
	if value >= 11 && value <= 14 {
		return many
	}
	switch value % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return few
	default:
		return many
	}
}

func tripletWords(value int, female bool) string {
	parts := make([]string, 0, 3)

	if value >= 100 {
		parts = append(parts, hundredsWords[value/100])
		value %= 100
	}
	if value >= 10 && value < 20 {
		parts = append(parts, teens[value-10])
		return strings.Join(parts, " ")
	}
	if value >= 20 {
		parts = append(parts, tensWords[value/10])
		value %= 10
	}
	if value > 0 {
		if female {
			parts = append(parts, unitsFem[value])
		} else {
			parts = append(parts, unitsMale[value])
		}
	}

	return strings.Join(parts, " ")
}

// intToWordsRu записывает целое число словами (до триллионов).
func intToWordsRu(value int64, female bool) string {
	if value == 0 {
		return "ноль"
	}
	if value < 0 {
		return "минус " + intToWordsRu(-value, female)
	}

	groups := make([]int64, 0, len(scales))
	for value > 0 {
		groups = append(groups, value%1000)
		value /= 1000
	}

	parts := make([]string, 0, 8)
	for index := len(groups) - 1; index >= 0; index-- {
		current := groups[index]
		if current == 0 {
			continue
		}

		isFemaleGroup := scales[index].female
		if index == 0 {
			isFemaleGroup = female
		}
		words := tripletWords(int(current), isFemaleGroup)
		if words != "" {
			parts = append(parts, words)
		}

		if index > 0 {
			scale := scales[index]
			parts = append(parts, pluralRu(current, scale.one, scale.few, scale.many))
		}
	}

	return strings.Join(parts, " ")
}

// rublesWords — сумма прописью для ТОРГ-12.
func rublesWords(kopecks int64) string {
	if kopecks < 0 {
		kopecks = -kopecks
	}

	rubles := kopecks / 100
	kopeckPart := kopecks % 100

	text := fmt.Sprintf("%s %s %02d %s",
		intToWordsRu(rubles, false),
		pluralRu(rubles, "рубль", "рубля", "рублей"),
		kopeckPart,
		pluralRu(kopeckPart, "копейка", "копейки", "копеек"),
	)

	return capitalize(text)
}

// kilogramsWords — масса прописью для ТТН.
func kilogramsWords(kilograms int64) string {
	if kilograms < 0 {
		kilograms = -kilograms
	}
	return capitalize(fmt.Sprintf("%s %s",
		intToWordsRu(kilograms, false),
		pluralRu(kilograms, "килограмм", "килограмма", "килограммов"),
	))
}

func capitalize(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	if runes[0] >= 'а' && runes[0] <= 'я' {
		runes[0] -= 'а' - 'А'
	}
	return string(runes)
}
