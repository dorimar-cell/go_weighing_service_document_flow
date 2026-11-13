package dto

import (
	"encoding/json"
	"time"
)

// UpdateSessionRequest — правка завершённого (или открытого) рейса.
type UpdateSessionRequest struct {
	TransportNumber string `json:"transportNumber"`
	TrailerNumber   string `json:"trailerNumber"`
	DriverName      string `json:"driverName"`
	CargoID         *int64 `json:"cargoId"`
	CounterpartyID  *int64 `json:"counterpartyId"`
	CarrierID       *int64 `json:"carrierId"`
	StorageID       *int64 `json:"storageId"`
	OperationTypeID *int64 `json:"operationTypeId"`

	FirstWeightKg    float64    `json:"firstWeightKg"`
	SecondWeightKg   *float64   `json:"secondWeightKg"`
	FirstMeasuredAt  time.Time  `json:"firstMeasuredAt"`
	SecondMeasuredAt *time.Time `json:"secondMeasuredAt"`
}

// SessionAuditRecord — запись истории изменений рейса.
type SessionAuditRecord struct {
	ID        int64           `json:"id"`
	Action    string          `json:"action"`
	UserID    int64           `json:"userId"`
	UserName  string          `json:"userName"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"createdAt"`
}

// AuditJournalItem — запись журнала аудита с контекстом рейса.
type AuditJournalItem struct {
	ID              int64           `json:"id"`
	SessionID       int64           `json:"sessionId"`
	Action          string          `json:"action"`
	UserID          int64           `json:"userId"`
	UserName        string          `json:"userName"`
	TransportNumber string          `json:"transportNumber"`
	Payload         json.RawMessage `json:"payload"`
	CreatedAt       time.Time       `json:"createdAt"`
}

// AuditJournalResponse — постраничный журнал аудита.
type AuditJournalResponse struct {
	Items    []AuditJournalItem `json:"items"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
	Total    int                `json:"total"`
}
