package reports

import "go_weighing_service/internal/dto"

// Группы каталога.
const (
	GroupJournals  = "Журналы и реестры"
	GroupDocuments = "Документы по рейсу"
)

var catalog = []dto.ReportDefinition{
	{
		Code:        "journal",
		Title:       "Журнал взвешиваний",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Рейсы за период: транспорт, прицеп, водитель, груз, контрагент, перевозчик, склад, веса и нетто. Печать и выгрузка XLSX.",
		Formats:     []string{"html", "xlsx"},
		Params: []dto.ReportParam{
			{Key: "dateFrom", Label: "С даты", Type: "date"},
			{Key: "dateTo", Label: "По дату", Type: "date"},
			{
				Key: "status", Label: "Статус", Type: "select", Default: "all",
				Options: []dto.ReportParamOption{
					{Value: "all", Label: "Все"},
					{Value: "closed", Label: "Закрытые"},
					{Value: "open", Label: "Открытые"},
					{Value: "cancelled", Label: "Отменённые"},
				},
			},
			{Key: "search", Label: "Номер ТС или пропуска", Type: "text"},
		},
	},
	{
		Code:        "ts-card",
		Title:       "Карточка ТС",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Реквизиты рейса, оба взвешивания, брутто/тара/нетто, фото с камер и подписи.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
	{
		Code:        "ttn",
		Title:       "ТТН (форма 1-Т)",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Товарно-транспортная накладная: стороны, автомобиль, груз, цена, подписи. Реквизиты сохраняются на рейсе.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
	{
		Code:        "torg12",
		Title:       "ТОРГ-12",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Товарная накладная: одна строка груза, НДС, суммы и сумма прописью.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
}

func sessionParams() []dto.ReportParam {
	return []dto.ReportParam{
		{Key: "sessionId", Label: "Рейс", Type: "session", Required: true},
	}
}

// Catalog возвращает каталог отчётов, доступных в веб-версии.
func Catalog() []dto.ReportDefinition {
	items := make([]dto.ReportDefinition, len(catalog))
	copy(items, catalog)
	return items
}

// Definition ищет отчёт по коду.
func Definition(code string) (dto.ReportDefinition, bool) {
	for _, item := range catalog {
		if item.Code == code {
			return item, true
		}
	}
	return dto.ReportDefinition{}, false
}
