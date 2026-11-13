package dto

// Permission — доменное право весовой (W001–W011).
type Permission struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GroupPermissions — группа пользователей и её права.
type GroupPermissions struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	System      bool     `json:"system"`
	Permissions []string `json:"permissions"`
}

// RBACResponse — матрица прав.
type RBACResponse struct {
	Permissions []Permission       `json:"permissions"`
	Groups      []GroupPermissions `json:"groups"`
}

// SaveGroupPermissionsRequest — новый набор прав группы.
type SaveGroupPermissionsRequest struct {
	Permissions []string `json:"permissions"`
}
