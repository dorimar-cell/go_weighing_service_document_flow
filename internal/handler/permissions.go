package handler

import (
	"net/http"

	"go_weighing_service/internal/middleware"
)

const (
	permFirstWeight  = "W001"
	permSecondWeight = "W002"
	permFlowWeight   = "W003"
	permCancelWeight = "W004"
	permJournal      = "W005"
	permExport       = "W006"
	permPasses       = "W007"
)

func requirePermission(w http.ResponseWriter, r *http.Request, permission string) bool {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Пользователь не определён")
		return false
	}

	if !user.Can(permission) {
		writeError(w, http.StatusForbidden, "forbidden", "Недостаточно прав для этого действия")
		return false
	}

	return true
}

func currentUserID(r *http.Request) int64 {
	user, ok := middleware.GetUser(r.Context())
	if !ok || user.ID == 0 {
		return 1
	}
	return user.ID
}
