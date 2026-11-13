package dto

import "time"

// ReportParamOption — вариант значения параметра отчёта.
type ReportParamOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ReportParam — параметр запуска отчёта. Набор типов намеренно узкий:
// отчёты этапа 6 фильтруются датой, статусом, строкой поиска и id рейса.
type ReportParam struct {
	Key      string              `json:"key"`
	Label    string              `json:"label"`
	Type     string              `json:"type"` // date | text | select | session
	Required bool                `json:"required"`
	Default  string              `json:"default,omitempty"`
	Options  []ReportParamOption `json:"options,omitempty"`
}

// ReportDefinition — карточка каталога отчётов.
type ReportDefinition struct {
	Code        string        `json:"code"`
	Title       string        `json:"title"`
	Group       string        `json:"group"`
	Description string        `json:"description"`
	Kind        string        `json:"kind"` // grid | document
	Formats     []string      `json:"formats"`
	Params      []ReportParam `json:"params"`
}

// ReportCatalogResponse — ответ каталога.
type ReportCatalogResponse struct {
	Items []ReportDefinition `json:"items"`
}

// ReportSettings — реквизиты организации для бланков.
type ReportSettings struct {
	OrganizationName string `json:"organizationName"`
	INN              string `json:"inn"`
	KPP              string `json:"kpp"`
	Address          string `json:"address"`
	OKPO             string `json:"okpo"`
	Chief            string `json:"chief"`
	ChiefPosition    string `json:"chiefPosition"`
	Accountant       string `json:"accountant"`
}

// AggregateRow — строка агрегатного отчёта (группировка + нетто).
type AggregateRow struct {
	Group    string  `json:"group"`
	SubGroup string  `json:"subGroup"`
	Sessions int     `json:"sessions"`
	NettoKg  float64 `json:"nettoKg"`
}

// SessionDetails — данные рейса, которых нет в журнальном DTO,
// но которые нужны бланкам (ИНН, адреса, цена груза, водительское).
type SessionDetails struct {
	Session SessionDTO

	DriverLicense string

	CounterpartyINN      string
	CounterpartyKPP      string
	CounterpartyAddress  string
	CounterpartyDirector string

	CarrierINN     string
	CarrierKPP     string
	CarrierAddress string

	CargoPricePerTonne *float64
	CargoSellingPrice  *float64

	OperatorName string
}

// TTNDocument — реквизиты ТТН (форма 1-Т), сохраняемые в ttn_document.
type TTNDocument struct {
	ID        int64 `json:"id"`
	SessionID int64 `json:"sessionId"`
	Saved     bool  `json:"saved"`

	TTNNumber   string `json:"ttnNumber"`
	Dogovor     string `json:"dogovor"`
	Putevka     string `json:"putevka"`
	Prilozhenie string `json:"prilozhenie"`

	CustomerOrganization  string `json:"customerOrganization"`
	ShipperOrganization   string `json:"shipperOrganization"`
	ConsigneeOrganization string `json:"consigneeOrganization"`
	PayerOrganization     string `json:"payerOrganization"`
	CarrierOrganization   string `json:"carrierOrganization"`

	TrailerNumber string  `json:"trailerNumber"`
	TypeOfTrailer string  `json:"typeOfTrailer"`
	LoadPoint     string  `json:"loadPoint"`
	UploadPoint   string  `json:"uploadPoint"`
	Distance      float64 `json:"distance"`

	RouteDistance float64 `json:"routeDistance"`
	CargoCapacity float64 `json:"cargoCapacity"`

	CargoReleaseName    string  `json:"cargoReleaseName"`
	CargoReleaseWorkAs  string  `json:"cargoReleaseWorkAs"`
	AccountantName      string  `json:"accountantName"`
	CargoPickUpName     string  `json:"cargoPickUpName"`
	CargoPickUpWorkAs   string  `json:"cargoPickUpWorkAs"`
	ConsigneeName       string  `json:"consigneeName"`
	ConsigneeWorkAs     string  `json:"consigneeWorkAs"`
	AdditionalInfo      string  `json:"additionalInfo"`
	PriceCargo          float64 `json:"priceCargo"`
	PriceTransportation float64 `json:"priceTransportation"`
	Forwarder           bool    `json:"forwarder"`

	UpdatedAt *time.Time `json:"updatedAt"`
}

// Torg12Document — реквизиты товарной накладной ТОРГ-12.
type Torg12Document struct {
	ID        int64 `json:"id"`
	SessionID int64 `json:"sessionId"`
	Saved     bool  `json:"saved"`

	Name         string  `json:"name"`
	Date         string  `json:"date"` // YYYY-MM-DD
	ConsigneeOrg string  `json:"consigneeOrgName"`
	ProviderOrg  string  `json:"providerOrgName"`
	StructName   string  `json:"structName"`
	Foundation   string  `json:"foundation"`
	NDS          float64 `json:"nds"`

	CargoReleaseName   string `json:"cargoReleaseName"`
	CargoReleaseWorkAs string `json:"cargoReleaseWorkAs"`
	AccountantName     string `json:"accountantName"`
	ShipperName        string `json:"shipperName"`
	ShipperWorkAs      string `json:"shipperWorkAs"`
	CargoPickUpName    string `json:"cargoPickUpName"`
	CargoPickUpWorkAs  string `json:"cargoPickUpWorkAs"`
	ConsigneeName      string `json:"consigneeName"`
	ConsigneeWorkAs    string `json:"consigneeWorkAs"`

	ProxyName    string `json:"proxyName"`
	ProxyDate    string `json:"proxyDate"`
	ProxyGivenBy string `json:"proxyGivenBy"`

	UpdatedAt *time.Time `json:"updatedAt"`
}
