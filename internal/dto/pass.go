package dto

import "time"

// CreatePassRequest — регистрация пропуска.
type CreatePassRequest struct {
	Number          string     `json:"number"`
	TransportNumber string     `json:"transportNumber"`
	TrailerNumber   string     `json:"trailerNumber"`
	DriverName      string     `json:"driverName"`
	CargoID         *int64     `json:"cargoId"`
	CounterpartyID  *int64     `json:"counterpartyId"`
	CarrierID       *int64     `json:"carrierId"`
	PassTypeID      int64      `json:"passTypeId"`
	ValidUntil      *time.Time `json:"validUntil"`
}

// PassDTO — пропуск для пульта и журнала.
type PassDTO struct {
	ID               int64      `json:"id"`
	Number           string     `json:"number"`
	TransportNumber  string     `json:"transportNumber"`
	TrailerNumber    string     `json:"trailerNumber"`
	DriverName       string     `json:"driverName"`
	CargoID          *int64     `json:"cargoId"`
	CargoName        string     `json:"cargoName"`
	CounterpartyID   *int64     `json:"counterpartyId"`
	CounterpartyName string     `json:"counterpartyName"`
	CarrierID        *int64     `json:"carrierId"`
	CarrierName      string     `json:"carrierName"`
	PassTypeID       int64      `json:"passTypeId"`
	PassTypeName     string     `json:"passTypeName"`
	StateID          int64      `json:"stateId"`
	StateName        string     `json:"stateName"`
	PassEnabled      bool       `json:"passEnabled"`
	ValidUntil       *time.Time `json:"validUntil"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type PassListResponse struct {
	Items    []PassDTO `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total"`
}
