package dto

import "time"

// SettingsProfile — версионируемый профиль настроек раздела.
// Draft редактируется в UI, applied применяется к агенту/сервису.
type SettingsProfile struct {
	Scope         string         `json:"scope"`
	Draft         map[string]any `json:"draft"`
	Applied       map[string]any `json:"applied"`
	Version       int            `json:"version"`
	UpdatedAt     *time.Time     `json:"updatedAt"`
	UpdatedByName string         `json:"updatedByName"`
	AppliedAt     *time.Time     `json:"appliedAt"`
	AppliedByName string         `json:"appliedByName"`
}

// SettingsProfilesResponse — список профилей настроек.
type SettingsProfilesResponse struct {
	Items []SettingsProfile `json:"items"`
}

// SaveSettingsProfileRequest — сохранение черновика раздела.
type SaveSettingsProfileRequest struct {
	Values map[string]any `json:"values"`
}

// AgentSettingsScope — применённые настройки одного раздела для агента.
type AgentSettingsScope struct {
	Version int            `json:"version"`
	Values  map[string]any `json:"values"`
}

// AgentSettingsResponse — конфигурация агента станции (только applied).
type AgentSettingsResponse struct {
	Scopes    map[string]AgentSettingsScope `json:"scopes"`
	UpdatedAt *time.Time                    `json:"updatedAt"`
}

// LoginAuditItem — запись аудита входов.
type LoginAuditItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Result    string    `json:"result"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginAuditResponse — страница аудита входов.
type LoginAuditResponse struct {
	Items    []LoginAuditItem `json:"items"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int              `json:"total"`
}
