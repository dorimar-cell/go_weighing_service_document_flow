package service

import (
	"context"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

// SettingsExtraService — камеры, лицензии и аудит настроек.
type SettingsExtraService struct {
	repo *repository.SettingsExtraRepository
}

func NewSettingsExtraService(repo *repository.SettingsExtraRepository) *SettingsExtraService {
	return &SettingsExtraService{repo: repo}
}

func (s *SettingsExtraService) Cameras(ctx context.Context) ([]dto.CameraDTO, error) {
	return s.repo.Cameras(ctx)
}

func (s *SettingsExtraService) SaveCamera(
	ctx context.Context,
	id int64,
	req dto.SaveCameraRequest,
	userID int64,
) ([]dto.CameraDTO, error) {
	if req.Port <= 0 {
		return nil, &ValidationError{Message: "Укажите порт камеры"}
	}
	req.Name = strings.TrimSpace(req.Name)
	req.IP = strings.TrimSpace(req.IP)
	if err := s.repo.SaveCamera(ctx, id, req, userID); err != nil {
		return nil, err
	}
	_ = s.repo.LogSettings(ctx, "camera", id, "updated", userID, req)
	return s.repo.Cameras(ctx)
}

func (s *SettingsExtraService) Licenses(ctx context.Context) ([]dto.LicenseDTO, error) {
	return s.repo.Licenses(ctx)
}

func (s *SettingsExtraService) CreateLicense(
	ctx context.Context,
	req dto.LicenseRequest,
	userID int64,
) ([]dto.LicenseDTO, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.ModuleCode) == "" {
		return nil, &ValidationError{Message: "Укажите название и код модуля"}
	}
	if req.Seats <= 0 {
		req.Seats = 1
	}
	if err := s.repo.CreateLicense(ctx, req, userID); err != nil {
		return nil, err
	}
	_ = s.repo.LogSettings(ctx, "license", 0, "created", userID, req)
	return s.repo.Licenses(ctx)
}

func (s *SettingsExtraService) UpdateLicense(
	ctx context.Context,
	id int64,
	req dto.LicenseRequest,
	userID int64,
) ([]dto.LicenseDTO, error) {
	if err := s.repo.UpdateLicense(ctx, id, req, userID); err != nil {
		return nil, err
	}
	_ = s.repo.LogSettings(ctx, "license", id, "updated", userID, req)
	return s.repo.Licenses(ctx)
}

func (s *SettingsExtraService) DeleteLicense(
	ctx context.Context,
	id int64,
	userID int64,
) ([]dto.LicenseDTO, error) {
	if err := s.repo.DeleteLicense(ctx, id); err != nil {
		return nil, err
	}
	_ = s.repo.LogSettings(ctx, "license", id, "deleted", userID, nil)
	return s.repo.Licenses(ctx)
}

func (s *SettingsExtraService) Audit(
	ctx context.Context,
	page int,
	pageSize int,
) (*dto.SettingsAuditResponse, error) {
	page, pageSize = normalizePaging(page, pageSize)
	return s.repo.SettingsAudit(ctx, page, pageSize)
}

func (s *SettingsExtraService) Log(
	ctx context.Context,
	scope string,
	entityID int64,
	action string,
	userID int64,
	payload any,
) {
	_ = s.repo.LogSettings(ctx, scope, entityID, action, userID, payload)
}
