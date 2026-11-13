package service

import (
	"context"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

// RBACService — матрица прав групп весовой.
type RBACService struct {
	repo *repository.RBACRepository
}

func NewRBACService(repo *repository.RBACRepository) *RBACService {
	return &RBACService{repo: repo}
}

func (s *RBACService) Matrix(ctx context.Context) (*dto.RBACResponse, error) {
	return s.repo.Matrix(ctx)
}

func (s *RBACService) SaveGroup(
	ctx context.Context,
	groupID int64,
	codes []string,
	userID int64,
) (*dto.RBACResponse, error) {
	if groupID <= 0 {
		return nil, &ValidationError{Message: "Некорректная группа"}
	}
	if err := s.repo.SaveGroupPermissions(ctx, groupID, codes, userID); err != nil {
		return nil, err
	}
	return s.repo.Matrix(ctx)
}
