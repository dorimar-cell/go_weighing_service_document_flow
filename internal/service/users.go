package service

import (
	"context"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

// UsersService — пользователи, весовые и настройки журнала.
type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) List(
	ctx context.Context,
	search string,
	page int,
	pageSize int,
) (*dto.UsersResponse, error) {
	page, pageSize = normalizePaging(page, pageSize)
	return s.repo.List(ctx, search, page, pageSize)
}

func (s *UsersService) WeightRooms(ctx context.Context) ([]dto.WeightRoomDTO, error) {
	return s.repo.WeightRooms(ctx)
}

func (s *UsersService) SetWeightRoom(ctx context.Context, userID, weightRoomID int64) error {
	if userID <= 0 || weightRoomID <= 0 {
		return &ValidationError{Message: "Некорректный пользователь или весовая"}
	}
	return s.repo.SetWeightRoom(ctx, userID, weightRoomID)
}

func (s *UsersService) JournalPreferences(
	ctx context.Context,
	userID int64,
) (*dto.JournalPreferences, error) {
	return s.repo.JournalPreferences(ctx, userID)
}

func (s *UsersService) SaveJournalPreferences(
	ctx context.Context,
	userID int64,
	prefs dto.JournalPreferences,
) (*dto.JournalPreferences, error) {
	if len(prefs.Columns) > 64 {
		return nil, &ValidationError{Message: "Слишком много колонок"}
	}
	if err := s.repo.SaveJournalPreferences(ctx, userID, prefs); err != nil {
		return nil, err
	}
	return s.repo.JournalPreferences(ctx, userID)
}
