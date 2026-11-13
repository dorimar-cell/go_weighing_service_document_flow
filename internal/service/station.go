package service

import (
	"context"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

// StationService — весовые и терминалы.
type StationService struct {
	repo *repository.StationRepository
}

func NewStationService(repo *repository.StationRepository) *StationService {
	return &StationService{repo: repo}
}

func (s *StationService) List(ctx context.Context) ([]dto.StationDTO, error) {
	return s.repo.List(ctx)
}

func (s *StationService) Save(
	ctx context.Context,
	id int64,
	req dto.SaveStationRequest,
	userID int64,
) ([]dto.StationDTO, error) {
	if req.MultiplierWeight <= 0 {
		return nil, &ValidationError{Message: "Множитель веса должен быть больше нуля"}
	}
	if req.PlatformNumber <= 0 {
		return nil, &ValidationError{Message: "Укажите номер платформы"}
	}

	req.TerminalModel = strings.TrimSpace(req.TerminalModel)
	req.InterfaceType = strings.TrimSpace(req.InterfaceType)
	req.ComName = strings.TrimSpace(req.ComName)
	req.IP = strings.TrimSpace(req.IP)
	req.Port = strings.TrimSpace(req.Port)
	req.SocketIP = strings.TrimSpace(req.SocketIP)
	req.SocketPort = strings.TrimSpace(req.SocketPort)

	if err := s.repo.Save(ctx, id, req, userID); err != nil {
		return nil, err
	}
	return s.repo.List(ctx)
}
