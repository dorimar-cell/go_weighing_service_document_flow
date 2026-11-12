package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/repository"
)

type SessionService struct {
	repo     *repository.SessionRepository
	passes   *repository.PassRepository
	verifier *MeasurementVerifier
}

func NewSessionService(
	repo *repository.SessionRepository,
	passes *repository.PassRepository,
	verifier *MeasurementVerifier,
) *SessionService {
	return &SessionService{repo: repo, passes: passes, verifier: verifier}
}

func (s *SessionService) Create(ctx context.Context, request dto.CreateSessionRequest) (*dto.SessionDTO, error) {
	request.TransportNumber = strings.TrimSpace(request.TransportNumber)
	request.TrailerNumber = strings.TrimSpace(request.TrailerNumber)
	request.DriverName = strings.TrimSpace(request.DriverName)

	if request.PassID != nil {
		pass, err := s.passes.Get(ctx, *request.PassID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &ValidationError{Message: "Пропуск не найден"}
			}
			return nil, err
		}
		if !pass.PassEnabled {
			return nil, &ValidationError{Message: "Пропуск неактивен"}
		}

		if request.TransportNumber == "" {
			request.TransportNumber = pass.TransportNumber
		}
		if request.TrailerNumber == "" {
			request.TrailerNumber = pass.TrailerNumber
		}
		if request.DriverName == "" {
			request.DriverName = pass.DriverName
		}
		if request.CargoID == nil {
			request.CargoID = pass.CargoID
		}
		if request.CounterpartyID == nil {
			request.CounterpartyID = pass.CounterpartyID
		}
		if request.CarrierID == nil {
			request.CarrierID = pass.CarrierID
		}
	}

	if request.TransportNumber == "" {
		return nil, &ValidationError{Message: "Укажите номер транспортного средства"}
	}
	if err := s.verifier.Verify(request.Measurement); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, request)
}

// CreateFlow — поточный режим: рейс сразу закрывается введённой тарой.
func (s *SessionService) CreateFlow(
	ctx context.Context,
	request dto.FlowSessionRequest,
) (*dto.SessionDTO, error) {
	if request.TareKg <= 0 {
		return nil, &ValidationError{Message: "Укажите тару больше нуля"}
	}
	if request.TareKg >= request.Measurement.WeightKg {
		return nil, &ValidationError{Message: "Тара должна быть меньше брутто"}
	}

	session, err := s.Create(ctx, request.CreateSessionRequest)
	if err != nil {
		return nil, err
	}

	return s.repo.CloseWithManualTare(ctx, session.ID, request.TareKg, time.Now())
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
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
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
	return s.repo.List(ctx, status, strings.TrimSpace(search), dateFrom, dateTo, page, pageSize)
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

func (s *SessionService) Measurements(ctx context.Context, id int64) ([]dto.MeasurementRecord, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Measurements(ctx, id)
}
