package service

import (
	"context"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

// SettingsService — бизнес-настройки полей рейса.
type SettingsService struct {
	repo *repository.SettingsRepository
}

func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) App(ctx context.Context) (*dto.AppSettings, error) {
	return s.repo.App(ctx)
}

func (s *SettingsService) SaveApp(
	ctx context.Context,
	settings dto.AppSettings,
	userID int64,
) (*dto.AppSettings, error) {
	// Обязательное поле не может быть скрыто: сервер не даст сохранить рейс.
	if settings.RequireDriver {
		settings.VisibleDriver = true
	}
	if settings.RequireTrailer {
		settings.VisibleTrailer = true
	}
	if settings.RequireCargo {
		settings.VisibleCargo = true
	}
	if settings.RequireCounterparty {
		settings.VisibleCounterparty = true
	}
	if settings.RequireCarrier {
		settings.VisibleCarrier = true
	}
	if settings.RequireStorage {
		settings.VisibleStorage = true
	}
	if settings.RequireOperationType {
		settings.VisibleOperationType = true
	}

	if err := s.repo.SaveApp(ctx, settings, userID); err != nil {
		return nil, err
	}
	return s.repo.App(ctx)
}
