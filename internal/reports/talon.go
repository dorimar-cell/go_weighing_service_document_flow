package reports

import (
	"strings"

	"go_weighing_service/internal/dto"
)

type talonVariant struct {
	Title    string
	Subtitle string
	Note     string
}

// talonVariants — варианты отвесного талона. Легаси содержит 22 бланка
// ПромСорт-Тула; здесь сведены к шести печатным шаблонам по назначению.
var talonVariants = map[string]talonVariant{
	"base": {
		Title:    "Талон взвешивания",
		Subtitle: "Отвесной талон",
		Note:     "Предназначен для внутреннего учёта и выдачи водителю.",
	},
	"vzrm": {
		Title:    "Отвесной талон ВЗРМ",
		Subtitle: "Взвешивание с подписями отпускающего и принимающего",
		Note:     "Бланк ВЗРМ: подписи отпускающего и принимающего обязательны.",
	},
	"tula-prihod": {
		Title:    "Отвесной талон · приход",
		Subtitle: "Приём груза (ПромСорт-Тула, шаблон прихода)",
		Note:     "Приходный вариант: поставщик, сорт и качество заполняются на весовой.",
	},
	"tula-prodazha": {
		Title:    "Отвесной талон · продажа",
		Subtitle: "Отгрузка груза (ПромСорт-Тула, шаблон продажи)",
		Note:     "Вариант продажи: получатель, цена и стоимость рассчитаны по нетто рейса.",
	},
	"tula-second": {
		Title:    "Отвесной талон · второе взвешивание",
		Subtitle: "Выезд с территории (ПромСорт-Тула, шаблон второго взвешивания)",
		Note:     "Второе взвешивание: основной вес — при выезде, первое указано справочно.",
	},
	"tula-defect": {
		Title:    "Отвесной талон · дефектный",
		Subtitle: "Приём с качественными показателями (ПромСорт-Тула, дефектный)",
		Note:     "Дефектный вариант: засорённость, влажность и скидка заполняются лабораторией.",
	},
}

// TalonVariantValid проверяет код варианта талона.
func TalonVariantValid(code string) bool {
	_, ok := talonVariants[code]
	return ok
}

func renderTalon(details *dto.SessionDetails, settings dto.ReportSettings) []byte {
	return renderTalonVariant(details, settings, "base")
}

// renderTalonVariant собирает отвесной талон в выбранном варианте.
func renderTalonVariant(details *dto.SessionDetails, settings dto.ReportSettings, variant string) []byte {
	meta, ok := talonVariants[variant]
	if !ok {
		meta = talonVariants["base"]
	}
	return wrapHTML(
		meta.Title+" №"+int64Text(details.Session.ID),
		false,
		talonVariantBody(details, settings, variant),
	)
}

// talonVariantBody — тело талона без HTML-обёртки.
func talonVariantBody(details *dto.SessionDetails, settings dto.ReportSettings, variant string) string {
	meta, ok := talonVariants[variant]
	if !ok {
		meta = talonVariants["base"]
		variant = "base"
	}

	session := &details.Session
	brutto, tare := bruttoTare(session)

	var body strings.Builder

	if settings.OrganizationName != "" {
		body.WriteString("<div class=\"meta\"><span><b>" + esc(settings.OrganizationName) + "</b></span>")
		if settings.INN != "" {
			body.WriteString("<span>ИНН: " + esc(settings.INN) + "</span>")
		}
		body.WriteString("</div>")
	}

	body.WriteString("<h1>" + esc(meta.Title) + "</h1>")
	body.WriteString("<p class=\"subtitle\">" + esc(meta.Subtitle) + " · № " +
		esc(int64Text(session.ID)) + " от " + esc(documentDate(session)) + "</p>")

	rows := [][2]string{
		{"Транспортное средство", esc(session.TransportNumber)},
		{"Прицеп", orDash(session.TrailerNumber)},
		{"Водитель", orDash(session.DriverName)},
		{"Пропуск", orDash(session.PassNumber)},
		{"Груз", orDash(session.CargoName)},
		{"Склад", orDash(session.StorageName)},
		{"Тип операции", orDash(session.OperationTypeName)},
	}

	switch variant {
	case "vzrm":
		rows = append(rows,
			[2]string{"Отпускающий (контрагент)", counterpartyText(details)},
			[2]string{"Перевозчик", carrierText(details)},
		)
	case "tula-prihod":
		rows = append(rows,
			[2]string{"Поставщик", counterpartyText(details)},
			[2]string{"Сорт", "&nbsp;"},
			[2]string{"Влажность, %", "&nbsp;"},
			[2]string{"Засорённость, %", "&nbsp;"},
		)
	case "tula-prodazha":
		rows = append(rows,
			[2]string{"Покупатель", counterpartyText(details)},
			[2]string{"Перевозчик", carrierText(details)},
		)
	case "tula-second":
		rows = append(rows,
			[2]string{"Перевозчик", carrierText(details)},
			[2]string{"Контрагент", counterpartyText(details)},
		)
	case "tula-defect":
		rows = append(rows,
			[2]string{"Поставщик", counterpartyText(details)},
			[2]string{"Влажность, %", "&nbsp;"},
			[2]string{"Засорённость, %", "&nbsp;"},
			[2]string{"Скидка, %", "&nbsp;"},
		)
	default:
		rows = append(rows,
			[2]string{"Контрагент", counterpartyText(details)},
			[2]string{"Перевозчик", carrierText(details)},
		)
	}

	body.WriteString(pairRows(rows))

	body.WriteString("<h2>Результат взвешивания</h2>")
	body.WriteString("<table><thead><tr><th>Показатель</th><th>Кг</th><th>Тонн</th></tr></thead><tbody>")
	if variant == "tula-second" {
		body.WriteString("<tr><td>Вес при выезде (второе взвешивание)</td><td class=\"num\">" +
			formatNumber(weightOrZero(session.SecondWeightKg), 1) + "</td><td class=\"num\">" +
			formatNumber(weightOrZero(session.SecondWeightKg)/1000, 3) + "</td></tr>")
		body.WriteString("<tr><td>Вес при въезде (первое, справочно)</td><td class=\"num\">" +
			formatNumber(session.FirstWeightKg, 1) + "</td><td class=\"num\">" +
			formatNumber(session.FirstWeightKg/1000, 3) + "</td></tr>")
	} else {
		body.WriteString("<tr><td>Брутто</td><td class=\"num\">" + formatNumber(brutto, 1) +
			"</td><td class=\"num\">" + formatNumber(brutto/1000, 3) + "</td></tr>")
		body.WriteString("<tr><td>Тара</td><td class=\"num\">" + formatNumber(tare, 1) +
			"</td><td class=\"num\">" + formatNumber(tare/1000, 3) + "</td></tr>")
	}
	body.WriteString("<tr class=\"total\"><td>Нетто</td><td class=\"num\">" +
		formatNumber(nettoValue(session), 1) + "</td><td class=\"num\">" +
		tonnesPtr(session.NetWeightKg) + "</td></tr>")
	body.WriteString("</tbody></table>")

	if variant == "tula-prodazha" {
		price := cargoPrice(details, 0)
		quantity := cargoQuantity(session)
		body.WriteString("<p class=\"words\">Цена, руб./т: " + formatNumber(price, 2) +
			". Стоимость: " + formatNumber(price*quantity, 2) + " руб.</p>")
	}

	body.WriteString("<div class=\"signs\">")
	switch variant {
	case "vzrm":
		body.WriteString("<div>Отпуск разрешил: <span class=\"sign-line\"></span></div>")
		body.WriteString("<div>Принял: <span class=\"sign-line\"></span></div>")
	case "tula-prihod":
		body.WriteString("<div>Весовщик: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
		body.WriteString("<div>Принял: <span class=\"sign-line\"></span></div>")
	case "tula-prodazha":
		body.WriteString("<div>Весовщик: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
		body.WriteString("<div>Получил: <span class=\"sign-line\"></span></div>")
	case "tula-defect":
		body.WriteString("<div>Весовщик: <span class=\"sign-line\"></span> " + orDash(details.OperatorName) + "</div>")
		body.WriteString("<div>Лаборант: <span class=\"sign-line\"></span></div>")
	default:
		body.WriteString("<div>Оператор весовой: <span class=\"sign-line\"></span> " +
			orDash(details.OperatorName) + "</div>")
	}
	body.WriteString("<div>Водитель: <span class=\"sign-line\"></span> " + orDash(session.DriverName) + "</div>")
	body.WriteString("</div>")

	body.WriteString("<p class=\"note\">Талон сформирован веб-системой весовой по данным рейса. " +
		esc(meta.Note) + "</p>")

	return body.String()
}
