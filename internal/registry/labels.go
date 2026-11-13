package registry

// fieldLabels — русские заголовки полей для выгрузок.
var fieldLabels = map[string]string{
	"name":           "Наименование",
	"description":    "Описание",
	"blackList":      "Чёрный список",
	"company":        "Организация",
	"phone":          "Телефон",
	"address":        "Адрес",
	"inn":            "ИНН",
	"kpp":            "КПП",
	"director":       "Руководитель",
	"cargoId":        "Груз (id)",
	"isDanger":       "Опасный груз",
	"pricePerTonne":  "Цена за тонну",
	"sellingPrice":   "Цена продажи",
	"direction":      "Направление",
	"main":           "Основная",
	"systemRecord":   "Системная запись",
	"make":           "Марка",
	"talon":          "Талон",
	"scope":          "Объём",
	"passportSeries": "Серия паспорта",
	"passportNumber": "Номер паспорта",
	"licenseSeries":  "Серия прав",
	"licenseNumber":  "Номер прав",
}

// FieldLabel возвращает русский заголовок поля или сам ключ.
func FieldLabel(key string) string {
	if label, ok := fieldLabels[key]; ok {
		return label
	}
	return key
}
