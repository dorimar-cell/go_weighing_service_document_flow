package dto

import (
	"encoding/json"
	"time"
)

// CameraDTO — камера весовой.
type CameraDTO struct {
	ID             int64      `json:"id"`
	Enable         bool       `json:"enable"`
	Name           string     `json:"name"`
	IP             string     `json:"ip"`
	Port           int64      `json:"port"`
	Login          string     `json:"login"`
	WaitTime       int64      `json:"waitTime"`
	TryTimes       int64      `json:"tryTimes"`
	StreamTypeMini int64      `json:"streamTypeMini"`
	StreamTypeFull int64      `json:"streamTypeFull"`
	PlatformNumber int64      `json:"platformNumber"`
	WeightRoomName string     `json:"weightRoomName"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// SaveCameraRequest — правка камеры.
type SaveCameraRequest struct {
	Enable         bool   `json:"enable"`
	Name           string `json:"name"`
	IP             string `json:"ip"`
	Port           int64  `json:"port"`
	Login          string `json:"login"`
	WaitTime       int64  `json:"waitTime"`
	TryTimes       int64  `json:"tryTimes"`
	StreamTypeMini int64  `json:"streamTypeMini"`
	StreamTypeFull int64  `json:"streamTypeFull"`
	PlatformNumber int64  `json:"platformNumber"`
}

// LicenseDTO — лицензия/модуль.
type LicenseDTO struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	ModuleCode string     `json:"moduleCode"`
	ValidUntil *string    `json:"validUntil"`
	Seats      int64      `json:"seats"`
	Active     bool       `json:"active"`
	Notes      string     `json:"notes"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

// LicenseRequest — создание/правка лицензии.
type LicenseRequest struct {
	Name       string  `json:"name"`
	ModuleCode string  `json:"moduleCode"`
	ValidUntil *string `json:"validUntil"`
	Seats      int64   `json:"seats"`
	Active     bool    `json:"active"`
	Notes      string  `json:"notes"`
}

// SettingsAuditItem — запись аудита настроек.
type SettingsAuditItem struct {
	ID        int64           `json:"id"`
	Scope     string          `json:"scope"`
	EntityID  int64           `json:"entityId"`
	Action    string          `json:"action"`
	UserID    int64           `json:"userId"`
	UserName  string          `json:"userName"`
	Version   int             `json:"version"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"createdAt"`
}

// SettingsAuditResponse — страница аудита настроек.
type SettingsAuditResponse struct {
	Items    []SettingsAuditItem `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	Total    int                 `json:"total"`
}
