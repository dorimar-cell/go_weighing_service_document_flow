package service

import (
	"context"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

type PassService struct {
	repo *repository.PassRepository
}

func NewPassService(repo *repository.PassRepository) *PassService {
	return &PassService{repo: repo}
}

func (s *PassService) Create(
	ctx context.Context,
	request dto.CreatePassRequest,
	userID int64,
) (*dto.PassDTO, error) {
	request.Number = strings.TrimSpace(request.Number)
	request.TransportNumber = strings.TrimSpace(request.TransportNumber)
	request.TrailerNumber = strings.TrimSpace(request.TrailerNumber)
	request.DriverName = strings.TrimSpace(request.DriverName)

	if request.Number == "" {
		return nil, &ValidationError{Message: "Укажите номер пропуска"}
	}
	if request.TransportNumber == "" {
		return nil, &ValidationError{Message: "Укажите номер транспортного средства"}
	}

	return s.repo.Create(ctx, request, userID)
}

func (s *PassService) List(
	ctx context.Context,
	status string,
	search string,
	page int,
	pageSize int,
) (*dto.PassListResponse, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		status = "active"
	}
	if status != "active" && status != "all" {
		return nil, &ValidationError{Message: "Неизвестный статус пропусков"}
	}

	page, pageSize = normalizePaging(page, pageSize)
	return s.repo.List(ctx, status, search, page, pageSize)
}

func (s *PassService) Get(ctx context.Context, id int64) (*dto.PassDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор пропуска"}
	}
	return s.repo.Get(ctx, id)
}

func (s *PassService) SetEnabled(
	ctx context.Context,
	id int64,
	enabled bool,
	userID int64,
) (*dto.PassDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор пропуска"}
	}
	return s.repo.SetEnabled(ctx, id, enabled, userID)
}
