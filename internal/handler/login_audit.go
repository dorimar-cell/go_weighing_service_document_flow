package handler

import (
	"net/http"

	"go_weighing_service/internal/repository"
)

// LoginAuditHandler — журнал входов в портал.
type LoginAuditHandler struct {
	repo *repository.LoginAuditRepository
}

func NewLoginAuditHandler(repo *repository.LoginAuditRepository) *LoginAuditHandler {
	return &LoginAuditHandler{repo: repo}
}

// List: GET /spa/api/weighing/login-audit?result&search&dateFrom&dateTo&page&pageSize
func (h *LoginAuditHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	dateFrom, err := parseDateParam(r.URL.Query().Get("dateFrom"), false)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}
	dateTo, err := parseDateParam(r.URL.Query().Get("dateTo"), true)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_date", err.Error())
		return
	}

	response, err := h.repo.List(r.Context(), repository.LoginAuditFilter{
		Result:   r.URL.Query().Get("result"),
		Search:   r.URL.Query().Get("search"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Page:     parseQueryInt(r, "page", 1),
		PageSize: parseQueryInt(r, "pageSize", 25),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить журнал входов")
		return
	}
	writeJSON(w, http.StatusOK, response)
}
