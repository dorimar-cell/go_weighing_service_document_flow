package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

const maxPhotoBytes = 8 << 20 // 8 MiB

type PhotoService struct {
	repo *repository.PhotoRepository
	dir  string
}

func NewPhotoService(repo *repository.PhotoRepository, dir string) *PhotoService {
	return &PhotoService{repo: repo, dir: dir}
}

func (s *PhotoService) Save(
	ctx context.Context,
	sessionID int64,
	userID int64,
	firstWeight bool,
	contentType string,
	encoded string,
) (*dto.PhotoDTO, error) {
	if sessionID <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	if strings.TrimSpace(encoded) == "" {
		return nil, &ValidationError{Message: "Пустой файл изображения"}
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, &ValidationError{Message: "Некорректный base64 изображения"}
	}
	if len(data) == 0 || len(data) > maxPhotoBytes {
		return nil, &ValidationError{Message: "Недопустимый размер изображения"}
	}

	extension := ".jpg"
	switch strings.ToLower(contentType) {
	case "image/png":
		extension = ".png"
	case "image/jpeg", "image/jpg", "":
		extension = ".jpg"
	default:
		return nil, &ValidationError{Message: "Поддерживаются только JPEG и PNG"}
	}

	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return nil, fmt.Errorf("generate photo name: %w", err)
	}

	sessionDir := strconv.FormatInt(sessionID, 10)
	fullDir := filepath.Join(s.dir, sessionDir)
	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		return nil, fmt.Errorf("create photo dir: %w", err)
	}

	fileName := hex.EncodeToString(random) + extension
	fullPath := filepath.Join(fullDir, fileName)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return nil, fmt.Errorf("write photo: %w", err)
	}

	relative := sessionDir + "/" + fileName
	return s.repo.Add(ctx, sessionID, relative, firstWeight, userID)
}

func (s *PhotoService) List(ctx context.Context, sessionID int64) ([]dto.PhotoDTO, error) {
	if sessionID <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.List(ctx, sessionID)
}

func (s *PhotoService) FilePath(ctx context.Context, sessionID, photoID int64) (string, error) {
	photo, err := s.repo.Get(ctx, sessionID, photoID)
	if err != nil {
		return "", err
	}

	full := filepath.Join(s.dir, filepath.FromSlash(photo.FileName))
	base := filepath.Clean(s.dir) + string(os.PathSeparator)
	if !strings.HasPrefix(filepath.Clean(full), base) {
		return "", repository.ErrNotFound
	}

	return full, nil
}
