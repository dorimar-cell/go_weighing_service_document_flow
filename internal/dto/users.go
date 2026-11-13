package dto

// UserDTO — пользователь весовой.
type UserDTO struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	PortalUsername string `json:"portalUsername"`
	GroupName      string `json:"groupName"`
	WeightRoomID   int64  `json:"weightRoomId"`
	WeightRoomName string `json:"weightRoomName"`
	Deleted        bool   `json:"deleted"`
}

// UsersResponse — список пользователей.
type UsersResponse struct {
	Items    []UserDTO `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total"`
}

// WeightRoomDTO — весовая (пост).
type WeightRoomDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UpdateWeightRoomRequest — привязка пользователя к весовой.
type UpdateWeightRoomRequest struct {
	WeightRoomID int64 `json:"weightRoomId"`
}

// UpdateUserRequest — группа и активность пользователя.
type UpdateUserRequest struct {
	GroupID *int64 `json:"groupId"`
	Deleted *bool  `json:"deleted"`
}

// JournalPreferences — серверные настройки журнала пользователя.
type JournalPreferences struct {
	Filters   map[string]any        `json:"filters"`
	Columns   []string              `json:"columns"`
	Presets   []JournalFilterPreset `json:"presets"`
	UpdatedAt *string               `json:"updatedAt,omitempty"`
}

// JournalFilterPreset — именованный набор фильтров журнала.
type JournalFilterPreset struct {
	Name    string         `json:"name"`
	Filters map[string]any `json:"filters"`
}

// BulkSessionRequest — массовая отмена/восстановление рейсов.
type BulkSessionRequest struct {
	Action string  `json:"action"`
	IDs    []int64 `json:"ids"`
}

// BulkSessionFailure — ошибка по конкретному рейсу.
type BulkSessionFailure struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// BulkSessionResponse — результат массовой операции.
type BulkSessionResponse struct {
	Updated int                  `json:"updated"`
	Failed  []BulkSessionFailure `json:"failed"`
}
