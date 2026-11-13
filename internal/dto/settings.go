package dto

// AppSettings — бизнес-настройки полей рейса: обязательность и видимость.
type AppSettings struct {
	RequireDriver        bool `json:"requireDriver"`
	RequireTrailer       bool `json:"requireTrailer"`
	RequireCargo         bool `json:"requireCargo"`
	RequireCounterparty  bool `json:"requireCounterparty"`
	RequireCarrier       bool `json:"requireCarrier"`
	RequireStorage       bool `json:"requireStorage"`
	RequireOperationType bool `json:"requireOperationType"`

	VisibleDriver        bool `json:"visibleDriver"`
	VisibleTrailer       bool `json:"visibleTrailer"`
	VisibleCargo         bool `json:"visibleCargo"`
	VisibleCounterparty  bool `json:"visibleCounterparty"`
	VisibleCarrier       bool `json:"visibleCarrier"`
	VisibleStorage       bool `json:"visibleStorage"`
	VisibleOperationType bool `json:"visibleOperationType"`
	VisiblePass          bool `json:"visiblePass"`
}
