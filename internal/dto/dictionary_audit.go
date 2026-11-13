package dto

import (
	"encoding/json"
	"time"
)

// DictionaryAuditItem — запись аудита справочника.
type DictionaryAuditItem struct {
	ID        int64           `json:"id"`
	Code      string          `json:"code"`
	RecordID  int64           `json:"recordId"`
	Action    string          `json:"action"`
	UserID    int64           `json:"userId"`
	UserName  string          `json:"userName"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"createdAt"`
}

// DictionaryAuditResponse — постраничный аудит справочников.
type DictionaryAuditResponse struct {
	Items    []DictionaryAuditItem `json:"items"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
	Total    int                   `json:"total"`
}
