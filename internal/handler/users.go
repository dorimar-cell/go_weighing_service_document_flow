package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

// UsersHandler — пользователи, весовые и настройки журнала.
type UsersHandler struct {
	service *service.UsersService
}

func NewUsersHandler(svc *service.UsersService) *UsersHandler {
	return &UsersHandler{service: svc}
}

// List: GET /spa/api/weighing/users?search&page&pageSize
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	result, err := h.service.List(
		r.Context(),
		r.URL.Query().Get("search"),
		optionalQueryID(r, "groupId"),
		parseQueryInt(r, "page", 1),
		parseQueryInt(r, "pageSize", 25),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить пользователей")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// WeightRooms: GET /spa/api/weighing/weight-rooms
func (h *UsersHandler) WeightRooms(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	items, err := h.service.WeightRooms(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить весовые")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// SetWeightRoom: PUT /spa/api/weighing/users/{id}/weight-room
func (h *UsersHandler) SetWeightRoom(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	id, err := parseNestedID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный пользователь")
		return
	}

	var payload dto.UpdateWeightRoomRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	if err := h.service.SetWeightRoom(r.Context(), id, payload.WeightRoomID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
			return
		}
		if validation, ok := err.(*service.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Update: PUT /spa/api/weighing/users/{id}
func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permSettings) {
		return
	}

	id, err := parseNestedID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный пользователь")
		return
	}

	var payload dto.UpdateUserRequest
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	if err := h.service.UpdateUser(r.Context(), id, payload); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
			return
		}
		if validation, ok := err.(*service.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FavoriteReports: GET /spa/api/weighing/reports/favorites
func (h *UsersHandler) FavoriteReports(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}
	codes, err := h.service.FavoriteReports(r.Context(), currentUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить избранное")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": codes})
}

// SetFavorite: PUT /spa/api/weighing/reports/{code}/favorite
func (h *UsersHandler) SetFavorite(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}
	code := chi.URLParam(r, "code")
	var payload struct {
		Favorite bool `json:"favorite"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	if err := h.service.SetFavoriteReport(r.Context(), currentUserID(r), code, payload.Favorite); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить избранное")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"code": code, "favorite": payload.Favorite})
}

// GetPreferences: GET /spa/api/weighing/journal/preferences
func (h *UsersHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	prefs, err := h.service.JournalPreferences(r.Context(), currentUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось загрузить настройки")
		return
	}
	writeJSON(w, http.StatusOK, prefs)
}

// SavePreferences: PUT /spa/api/weighing/journal/preferences
func (h *UsersHandler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	var payload dto.JournalPreferences
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	prefs, err := h.service.SaveJournalPreferences(r.Context(), currentUserID(r), payload)
	if err != nil {
		if validation, ok := err.(*service.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Не удалось сохранить настройки")
		return
	}
	writeJSON(w, http.StatusOK, prefs)
}
