package reports

import "go_weighing_service/internal/dto"

// Группы каталога.
const (
	GroupJournals   = "Журналы и реестры"
	GroupAggregates = "Агрегаты"
	GroupDocuments  = "Документы по рейсу"
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
		Code:        "talons",
		Title:       "Талоны за период",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Пакетная печать талонов взвешивания по закрытым рейсам за период: один талон на страницу.",
		Formats:     []string{"html"},
		Params:      periodParams(),
	},
	{
		Code:        "first-weight",
		Title:       "Транспорт на первом взвешивании",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Открытые рейсы, ожидающие второго взвешивания: транспорт, груз, вес первого взвешивания.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
	},
	{
		Code:        "journal-by-transport",
		Title:       "Журнал по номеру ТС",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Закрытые рейсы с группировкой по номеру транспортного средства и дате.",
		Formats:     []string{"html", "xlsx"},
		Params: append(periodParams(), dto.ReportParam{
			Key: "search", Label: "Номер ТС или пропуска", Type: "text",
		}),
	},
	{
		Code:        "ttn-batch",
		Title:       "ТТН за период",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Пакетная печать ТТН по закрытым рейсам: по документу на страницу, с сохранёнными реквизитами.",
		Formats:     []string{"html"},
		Params:      periodParams(),
	},
	{
		Code:        "torg12-batch",
		Title:       "ТОРГ-12 за период",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Пакетная печать товарных накладных по закрытым рейсам: документ на страницу.",
		Formats:     []string{"html"},
		Params:      periodParams(),
	},
	{
		Code:        "invoice-batch",
		Title:       "Счета-фактуры за период",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Пакетная печать счетов-фактур по закрытым рейсам: документ на страницу.",
		Formats:     []string{"html"},
		Params:      periodParams(),
	},
	{
		Code:        "cargo-act-batch",
		Title:       "Акты приёма-передачи за период",
		Group:       GroupJournals,
		Kind:        "grid",
		Description: "Пакетная печать актов приёма-передачи по закрытым рейсам: три экземпляра на рейс.",
		Formats:     []string{"html"},
		Params:      periodParams(),
	},
	{
		Code:        "cargo-summary",
		Title:       "Отчёт по грузам",
		Group:       GroupAggregates,
		Kind:        "grid",
		Description: "Нетто и количество рейсов по каждому грузу за период.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
	},
	{
		Code:        "counterparty-summary",
		Title:       "Отчёт по контрагентам",
		Group:       GroupAggregates,
		Kind:        "grid",
		Description: "Нетто по контрагентам с разбивкой по грузам.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
	},
	{
		Code:        "carrier-summary",
		Title:       "Отчёт по перевозчикам",
		Group:       GroupAggregates,
		Kind:        "grid",
		Description: "Нетто по перевозчикам с разбивкой по грузам.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
	},
	{
		Code:        "driver-summary",
		Title:       "Отчёт по водителям",
		Group:       GroupAggregates,
		Kind:        "grid",
		Description: "Нетто по водителям с разбивкой по грузам.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
	},
	{
		Code:        "storage-summary",
		Title:       "Взвешивания по складам",
		Group:       GroupAggregates,
		Kind:        "grid",
		Description: "Нетто по складам с разбивкой по типу операции.",
		Formats:     []string{"html", "xlsx"},
		Params:      periodParams(),
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
		Code:        "session-pack",
		Title:       "Комплект документов",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Карточка ТС, талон, ТТН, ТОРГ-12 и счёт-фактура одним файлом с разбивкой по страницам.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
	{
		Code:        "talon",
		Title:       "Талон взвешивания",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Отвесной талон: базовый, ВЗРМ, приход, продажа, второе взвешивание, дефектный.",
		Formats:     []string{"html"},
		Params: append(sessionParams(), dto.ReportParam{
			Key: "variant", Label: "Вариант бланка", Type: "select", Default: "base",
			Options: []dto.ReportParamOption{
				{Value: "base", Label: "Базовый"},
				{Value: "vzrm", Label: "ВЗРМ"},
				{Value: "tula-prihod", Label: "Приход"},
				{Value: "tula-prodazha", Label: "Продажа"},
				{Value: "tula-second", Label: "Второе взвешивание"},
				{Value: "tula-defect", Label: "Дефектный"},
			},
		}),
	},
	{
		Code:        "cargo-act",
		Title:       "Акт приёма-передачи товара",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Три экземпляра: охраны, водителя и весовщика. Нетто, объём и подписи сторон.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
	{
		Code:        "uploading-certificate",
		Title:       "Свидетельство о разгрузке/погрузке",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Подтверждение факта разгрузки/погрузки Т/С: операция, ТС, продукт, параметры взвешивания.",
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
	{
		Code:        "invoice",
		Title:       "Счёт-фактура",
		Group:       GroupDocuments,
		Kind:        "document",
		Description: "Счёт-фактура: продавец, покупатель, груз, НДС, сумма прописью. Реквизиты сохраняются на рейсе.",
		Formats:     []string{"html"},
		Params:      sessionParams(),
	},
}

func sessionParams() []dto.ReportParam {
	return []dto.ReportParam{
		{Key: "sessionId", Label: "Рейс", Type: "session", Required: true},
	}
}

func periodParams() []dto.ReportParam {
	return []dto.ReportParam{
		{Key: "dateFrom", Label: "С даты", Type: "date"},
		{Key: "dateTo", Label: "По дату", Type: "date"},
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
