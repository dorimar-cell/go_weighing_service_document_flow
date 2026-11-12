package registry

// FieldKind — тип поля справочника на уровне API и PostgreSQL.
type FieldKind string

const (
	KindText    FieldKind = "text"
	KindInt     FieldKind = "int"
	KindNumeric FieldKind = "numeric"
	KindBool    FieldKind = "bool"
)

// Field описывает одно поле справочника.
//
// Key — имя в JSON-контракте (camelCase), Column — имя колонки в PostgreSQL.
// Имена таблиц и колонок берутся только отсюда, поэтому динамический SQL
// не зависит от пользовательского ввода.
type Field struct {
	Key      string
	Column   string
	Kind     FieldKind
	Required bool
	Search   bool
}

// Dictionary — описание справочника: код API, таблица и поля.
type Dictionary struct {
	Code        string
	Title       string
	Group       string
	Description string
	LegacyTable string
	AccessTag   string
	Table       string
	OrderBy     string
	Fields      []Field
}

var dictionaries = []Dictionary{
	{
		Code: "transport", Title: "Транспорт", Group: "Транспорт",
		Description: "Госномера транспортных средств, чёрный список.",
		LegacyTable: "dir_transport", AccessTag: "D011",
		Table: "dir_transport", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "blackList", Column: "black_list", Kind: KindBool},
		},
	},
	{
		Code: "trailer", Title: "Прицепы", Group: "Транспорт",
		Description: "Номера прицепов и признак недопуска.",
		LegacyTable: "dir_trailer", AccessTag: "D013",
		Table: "dir_trailer", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "blackList", Column: "black_list", Kind: KindBool},
		},
	},
	{
		Code: "typeOfTransport", Title: "Виды ТС", Group: "Транспорт",
		Description: "Классификатор транспорта (тягач, самосвал, фургон…).",
		LegacyTable: "dir_typeOfTransport", AccessTag: "D012",
		Table: "dir_type_of_transport", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "description", Column: "description", Kind: KindText},
		},
	},
	{
		Code: "typeOfTrailer", Title: "Виды прицепов", Group: "Транспорт",
		Description: "Классификатор прицепов и полуприцепов.",
		LegacyTable: "dir_typeOfTrailer", AccessTag: "D013",
		Table: "dir_type_of_trailer", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
		},
	},
	{
		Code: "driver", Title: "Водители", Group: "Транспорт",
		Description: "ФИО, документы, контакты, чёрный список.",
		LegacyTable: "dir_driver", AccessTag: "D004",
		Table: "dir_driver", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "phone", Column: "phone", Kind: KindText, Search: true},
			{Key: "inn", Column: "inn", Kind: KindText, Search: true},
			{Key: "passportSeries", Column: "passport_series", Kind: KindText},
			{Key: "passportNumber", Column: "passport_number", Kind: KindText},
			{Key: "licenseSeries", Column: "license_series", Kind: KindText},
			{Key: "licenseNumber", Column: "license_number", Kind: KindText},
			{Key: "blackList", Column: "black_list", Kind: KindBool},
		},
	},
	{
		Code: "cargo", Title: "Грузы", Group: "Грузы и склады",
		Description: "Номенклатура грузов, цены, габариты, опасность.",
		LegacyTable: "dir_cargo", AccessTag: "D001",
		Table: "dir_cargo", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "pricePerTonne", Column: "price_per_tonne", Kind: KindNumeric},
			{Key: "sellingPrice", Column: "selling_price", Kind: KindNumeric},
			{Key: "isDanger", Column: "is_danger", Kind: KindBool},
		},
	},
	{
		Code: "typeOfCargo", Title: "Виды груза", Group: "Грузы и склады",
		Description: "Подклассификатор грузов (род груза).",
		LegacyTable: "dir_typeOfCargo", AccessTag: "D002",
		Table: "dir_type_of_cargo", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "cargoId", Column: "cargo_id", Kind: KindInt},
		},
	},
	{
		Code: "storage", Title: "Склады", Group: "Грузы и склады",
		Description: "Склады, ответственные лица, остаток.",
		LegacyTable: "dir_storage", AccessTag: "D005",
		Table: "dir_storage", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "company", Column: "company", Kind: KindText, Search: true},
			{Key: "phone", Column: "phone", Kind: KindText},
			{Key: "address", Column: "address", Kind: KindText},
		},
	},
	{
		Code: "placement", Title: "Места размещения", Group: "Грузы и склады",
		Description: "Зоны/места хранения внутри склада.",
		LegacyTable: "dir_placement", AccessTag: "D009",
		Table: "dir_placement", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
		},
	},
	{
		Code: "fund", Title: "Фонды", Group: "Грузы и склады",
		Description: "Фонды/источники финансирования, используемые в рейсе.",
		LegacyTable: "dir_fund", AccessTag: "D018",
		Table: "dir_fund", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "description", Column: "description", Kind: KindText},
		},
	},
	{
		Code: "carrier", Title: "Перевозчики", Group: "Контрагенты",
		Description: "Транспортные организации и реквизиты.",
		LegacyTable: "dir_carrier", AccessTag: "D003",
		Table: "dir_carrier", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "inn", Column: "inn", Kind: KindText, Search: true},
			{Key: "kpp", Column: "kpp", Kind: KindText},
			{Key: "phone", Column: "phone", Kind: KindText},
			{Key: "address", Column: "address", Kind: KindText},
		},
	},
	{
		Code: "counterparty", Title: "Контрагенты", Group: "Контрагенты",
		Description: "Юридические лица: реквизиты, руководитель, бухгалтер.",
		LegacyTable: "dir_counterparty", AccessTag: "D014",
		Table: "dir_counterparty", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "inn", Column: "inn", Kind: KindText, Search: true},
			{Key: "kpp", Column: "kpp", Kind: KindText},
			{Key: "phone", Column: "phone", Kind: KindText},
			{Key: "director", Column: "director", Kind: KindText},
			{Key: "address", Column: "address", Kind: KindText},
		},
	},
	{
		Code: "provider", Title: "Поставщики", Group: "Контрагенты",
		Description: "Поставщики груза и их реквизиты.",
		LegacyTable: "dir_provider", AccessTag: "D027",
		Table: "dir_provider", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "inn", Column: "inn", Kind: KindText, Search: true},
			{Key: "phone", Column: "phone", Kind: KindText},
			{Key: "address", Column: "address", Kind: KindText},
		},
	},
	{
		Code: "operationType", Title: "Типы операций", Group: "Система",
		Description: "Приход / расход / перемещение и системные признаки.",
		LegacyTable: "dir_typeOfOperation", AccessTag: "D019",
		Table: "dir_type_of_operation", OrderBy: "id",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "direction", Column: "direction", Kind: KindInt},
			{Key: "main", Column: "main", Kind: KindBool},
			{Key: "systemRecord", Column: "system_record", Kind: KindBool},
		},
	},
	{
		Code: "plantation", Title: "Плантации", Group: "КФХ",
		Description: "Поля/плантации для учёта КФХ.",
		LegacyTable: "dir_kfkh_plantation", AccessTag: "D017",
		Table: "dir_kfkh_plantation", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "description", Column: "description", Kind: KindText},
		},
	},
	{
		Code: "combine", Title: "Комбайны", Group: "КФХ",
		Description: "Техника КФХ: марка, производительность.",
		LegacyTable: "dir_kfkh_combine", AccessTag: "D015",
		Table: "dir_kfkh_combine", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "make", Column: "make", Kind: KindText, Search: true},
			{Key: "scope", Column: "scope", Kind: KindNumeric},
		},
	},
	{
		Code: "combiner", Title: "Комбайнёры", Group: "КФХ",
		Description: "Водители-комбайнёры и талоны.",
		LegacyTable: "dir_kfkh_combiner", AccessTag: "D016",
		Table: "dir_kfkh_combiner", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "talon", Column: "talon", Kind: KindText},
		},
	},
	{
		Code: "weightRoom", Title: "Весовые", Group: "Пользователи",
		Description: "Посты (весовые): привязка пользователей и настроек.",
		LegacyTable: "dir_WeightRoom", AccessTag: "S001-TP003",
		Table: "dir_weight_room", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "description", Column: "description", Kind: KindText},
		},
	},
	{
		Code: "userGroup", Title: "Группы пользователей", Group: "Пользователи",
		Description: "Группы для матрицы прав.",
		LegacyTable: "dir_UserGroups", AccessTag: "S001-TP004",
		Table: "dir_user_group", OrderBy: "lower(name)",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
			{Key: "systemRecord", Column: "system_record", Kind: KindBool},
		},
	},
	{
		Code: "passType", Title: "Типы пропусков", Group: "Система",
		Description: "Системный справочник типов пропусков.",
		LegacyTable: "system_dir_typeOfPass", AccessTag: "D020",
		Table: "system_dir_type_of_pass", OrderBy: "id",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
		},
	},
	{
		Code: "identifierState", Title: "Статусы пропусков", Group: "Система",
		Description: "Системный справочник состояний идентификаторов.",
		LegacyTable: "system_dir_identifierState", AccessTag: "D020",
		Table: "system_dir_identifier_state", OrderBy: "id",
		Fields: []Field{
			{Key: "name", Column: "name", Kind: KindText, Required: true, Search: true},
		},
	},
}

var byCode = func() map[string]Dictionary {
	m := make(map[string]Dictionary, len(dictionaries))
	for _, d := range dictionaries {
		m[d.Code] = d
	}
	return m
}()

// Get возвращает справочник по коду API.
func Get(code string) (Dictionary, bool) {
	d, ok := byCode[code]
	return d, ok
}

// All возвращает все справочники в порядке каталога.
func All() []Dictionary {
	return dictionaries
}
