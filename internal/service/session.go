package service

import (
	"context"
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

type SessionService struct {
	repo     *repository.SessionRepository
	verifier *MeasurementVerifier
}

func NewSessionService(repo *repository.SessionRepository, verifier *MeasurementVerifier) *SessionService {
	return &SessionService{repo: repo, verifier: verifier}
}

func (s *SessionService) Create(ctx context.Context, request dto.CreateSessionRequest) (*dto.SessionDTO, error) {
	request.TransportNumber = strings.TrimSpace(request.TransportNumber)
	request.TrailerNumber = strings.TrimSpace(request.TrailerNumber)
	request.DriverName = strings.TrimSpace(request.DriverName)

	if request.TransportNumber == "" {
		return nil, &ValidationError{Message: "Укажите номер транспортного средства"}
	}
	if err := s.verifier.Verify(request.Measurement); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, request)
}

func (s *SessionService) Close(
	ctx context.Context,
	id int64,
	measurement dto.MeasurementDTO,
) (*dto.SessionDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	if err := s.verifier.Verify(measurement); err != nil {
		return nil, err
	}

	return s.repo.Close(ctx, id, measurement)
}

func (s *SessionService) List(
	ctx context.Context,
	status string,
	page int,
	pageSize int,
) (*dto.SessionListResponse, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "", "open", "closed", "cancelled", "all":
	default:
		return nil, &ValidationError{Message: fmt.Sprintf("Неизвестный статус: %s", status)}
	}
	if status == "" {
		status = "open"
	}

	page, pageSize = normalizePaging(page, pageSize)
	return s.repo.List(ctx, status, page, pageSize)
}

func (s *SessionService) Get(ctx context.Context, id int64) (*dto.SessionDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Get(ctx, id)
}

func (s *SessionService) Cancel(ctx context.Context, id int64) error {
	if id <= 0 {
		return &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Cancel(ctx, id)
}
