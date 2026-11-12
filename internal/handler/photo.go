package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"go_weighing_service/internal/service"
)

type PhotoHandler struct {
	service *service.PhotoService
}

type uploadPhotoRequest struct {
	Base64      string `json:"base64"`
	ContentType string `json:"contentType"`
	FirstWeight bool   `json:"firstWeight"`
}

func NewPhotoHandler(svc *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{service: svc}
}

// Upload: POST /spa/api/weighing/sessions/{id}/photos
func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	var request uploadPhotoRequest
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", "Некорректный JSON изображения")
		return
	}

	permission := permSecondWeight
	if request.FirstWeight {
		permission = permFirstWeight
	}
	if !requirePermission(w, r, permission) {
		return
	}

	photo, err := h.service.Save(
		r.Context(),
		id,
		currentUserID(r),
		request.FirstWeight,
		request.ContentType,
		request.Base64,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	photo.URL = photoURL(id, photo.ID)
	writeJSON(w, http.StatusCreated, photo)
}

// List: GET /spa/api/weighing/sessions/{id}/photos
func (h *PhotoHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	items, err := h.service.List(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	for index := range items {
		items[index].URL = photoURL(id, items[index].ID)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// File: GET /spa/api/weighing/sessions/{id}/photos/{photoId}/file
func (h *PhotoHandler) File(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(w, r, permJournal) {
		return
	}

	sessionID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор рейса")
		return
	}

	photoID, err := parseNestedID(r, "photoId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "Некорректный идентификатор фото")
		return
	}

	fullPath, err := h.service.FilePath(r.Context(), sessionID, photoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Cache-Control", "private, max-age=3600")
	http.ServeFile(w, r, fullPath)
}

func (h *PhotoHandler) handleError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, "validation_error", validation.Message)
		return
	}

	slog.Error("photo handler", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
}

func photoURL(sessionID, photoID int64) string {
	return fmt.Sprintf("/spa/api/weighing/sessions/%d/photos/%d/file", sessionID, photoID)
}
