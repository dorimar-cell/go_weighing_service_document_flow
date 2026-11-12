package dto

import "time"

// MeasurementDTO — подписанное измерение от Station Agent.
type MeasurementDTO struct {
	OperationID string    `json:"operationId"`
	StationID   string    `json:"stationId"`
	Platform    int       `json:"platform"`
	Terminal    string    `json:"terminal"`
	WeightKg    float64   `json:"weightKg"`
	Stable      bool      `json:"stable"`
	MeasuredAt  time.Time `json:"measuredAt"`
	Sequence    int64     `json:"sequence"`
	Checksum    string    `json:"checksum"`
}

// CreateSessionRequest — первое взвешивание.
type CreateSessionRequest struct {
	Measurement     MeasurementDTO `json:"measurement"`
	TransportNumber string         `json:"transportNumber"`
	TrailerNumber   string         `json:"trailerNumber"`
	DriverName      string         `json:"driverName"`
	CargoID         *int64         `json:"cargoId"`
	CounterpartyID  *int64         `json:"counterpartyId"`
	CarrierID       *int64         `json:"carrierId"`
	StorageID       *int64         `json:"storageId"`
	OperationTypeID *int64         `json:"operationTypeId"`
}

// CloseSessionRequest — второе взвешивание.
type CloseSessionRequest struct {
	Measurement MeasurementDTO `json:"measurement"`
}

// SessionDTO — карточка рейса для пульта и журнала.
type SessionDTO struct {
	ID               int64      `json:"id"`
	Status           string     `json:"status"`
	TransportNumber  string     `json:"transportNumber"`
	TrailerNumber    string     `json:"trailerNumber"`
	DriverName       string     `json:"driverName"`
	FirstWeightKg    float64    `json:"firstWeightKg"`
	SecondWeightKg   *float64   `json:"secondWeightKg"`
	NetWeightKg      *float64   `json:"netWeightKg"`
	FirstMeasuredAt  time.Time  `json:"firstMeasuredAt"`
	SecondMeasuredAt *time.Time `json:"secondMeasuredAt"`
	CreatedAt        time.Time  `json:"createdAt"`

	CargoID           *int64 `json:"cargoId"`
	CargoName         string `json:"cargoName"`
	CounterpartyID    *int64 `json:"counterpartyId"`
	CounterpartyName  string `json:"counterpartyName"`
	CarrierID         *int64 `json:"carrierId"`
	CarrierName       string `json:"carrierName"`
	StorageID         *int64 `json:"storageId"`
	StorageName       string `json:"storageName"`
	OperationTypeID   *int64 `json:"operationTypeId"`
	OperationTypeName string `json:"operationTypeName"`
}

// SessionListResponse — постраничный журнал.
type SessionListResponse struct {
	Items    []SessionDTO `json:"items"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	Total    int          `json:"total"`
}
