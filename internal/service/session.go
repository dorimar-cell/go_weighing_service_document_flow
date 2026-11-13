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
	settings *repository.SettingsRepository
	verifier *MeasurementVerifier
}

func NewSessionService(
	repo *repository.SessionRepository,
	passes *repository.PassRepository,
	settings *repository.SettingsRepository,
	verifier *MeasurementVerifier,
) *SessionService {
	return &SessionService{repo: repo, passes: passes, settings: settings, verifier: verifier}
}

func (s *SessionService) Create(
	ctx context.Context,
	request dto.CreateSessionRequest,
	userID int64,
) (*dto.SessionDTO, error) {
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
	if request.ManualWeightKg == nil && request.Measurement == nil {
		return nil, &ValidationError{Message: "Нет измерения веса: нужен подписанный снимок или ручной вес"}
	}
	if request.ManualWeightKg != nil {
		if *request.ManualWeightKg < 0 {
			return nil, &ValidationError{Message: "Ручной вес не может быть отрицательным"}
		}
	} else if err := s.verifier.Verify(*request.Measurement); err != nil {
		return nil, err
	}
	if err := s.checkBlacklist(ctx, request.TransportNumber, request.TrailerNumber, request.DriverName); err != nil {
		return nil, err
	}
	if err := s.validateRequiredFields(
		ctx,
		request.TransportNumber,
		request.DriverName,
		request.TrailerNumber,
		request.CargoID,
		request.CounterpartyID,
		request.CarrierID,
		request.StorageID,
		request.OperationTypeID,
	); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, request, userID)
}

// CreateFlow — поточный режим: рейс сразу закрывается введённой тарой.
func (s *SessionService) CreateFlow(
	ctx context.Context,
	request dto.FlowSessionRequest,
	userID int64,
) (*dto.SessionDTO, error) {
	if request.TareKg <= 0 {
		return nil, &ValidationError{Message: "Укажите тару больше нуля"}
	}

	gross := 0.0
	switch {
	case request.ManualWeightKg != nil:
		gross = *request.ManualWeightKg
	case request.Measurement != nil:
		gross = request.Measurement.WeightKg
	default:
		return nil, &ValidationError{Message: "Нет измерения веса: нужен подписанный снимок или ручной вес"}
	}
	if request.TareKg >= gross {
		return nil, &ValidationError{Message: "Тара должна быть меньше брутто"}
	}

	session, err := s.Create(ctx, request.CreateSessionRequest, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.CloseWithManualTare(ctx, session.ID, request.TareKg, time.Now(), userID)
}

func (s *SessionService) Close(
	ctx context.Context,
	id int64,
	measurement dto.MeasurementDTO,
	userID int64,
) (*dto.SessionDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	if err := s.verifier.Verify(measurement); err != nil {
		return nil, err
	}

	return s.repo.Close(ctx, id, measurement, userID)
}

func (s *SessionService) List(
	ctx context.Context,
	status string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
	sortBy string,
	direction string,
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
	return s.repo.List(ctx, status, strings.TrimSpace(search), dateFrom, dateTo,
		sortBy, direction, page, pageSize)
}

func (s *SessionService) Get(ctx context.Context, id int64) (*dto.SessionDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Get(ctx, id)
}

func (s *SessionService) Cancel(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Cancel(ctx, id, userID)
}

// AuditExport возвращает все записи журнала аудита для CSV.
func (s *SessionService) AuditExport(
	ctx context.Context,
	action string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.AuditJournalItem, error) {
	return s.repo.AuditExport(ctx, action, search, dateFrom, dateTo)
}

// Bulk выполняет массовую отмену или восстановление рейсов.
func (s *SessionService) Bulk(
	ctx context.Context,
	action string,
	ids []int64,
	userID int64,
) (*dto.BulkSessionResponse, error) {
	if len(ids) == 0 {
		return nil, &ValidationError{Message: "Не выбрано ни одного рейса"}
	}
	if len(ids) > 200 {
		return nil, &ValidationError{Message: "Слишком много рейсов (макс. 200)"}
	}

	result := &dto.BulkSessionResponse{Failed: make([]dto.BulkSessionFailure, 0)}
	for _, id := range ids {
		var err error
		switch action {
		case "cancel":
			err = s.Cancel(ctx, id, userID)
		case "restore":
			err = s.Restore(ctx, id, userID)
		default:
			return nil, &ValidationError{Message: "Неизвестное действие: " + action}
		}
		if err != nil {
			result.Failed = append(result.Failed, dto.BulkSessionFailure{ID: id, Message: err.Error()})
			continue
		}
		result.Updated++
	}

	return result, nil
}

// Restore возвращает отменённый рейс.
func (s *SessionService) Restore(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Restore(ctx, id, userID)
}

// AuditJournal возвращает журнал аудита по всем рейсам.
func (s *SessionService) AuditJournal(
	ctx context.Context,
	action string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
	page int,
	pageSize int,
) (*dto.AuditJournalResponse, error) {
	action = strings.TrimSpace(action)
	switch action {
	case "", "created", "closed", "flow_closed", "cancelled", "restored", "edited", "document_saved":
	default:
		return nil, &ValidationError{Message: fmt.Sprintf("Неизвестное действие: %s", action)}
	}

	page, pageSize = normalizePaging(page, pageSize)
	return s.repo.AuditJournal(ctx, action, strings.TrimSpace(search), dateFrom, dateTo, page, pageSize)
}

func (s *SessionService) Measurements(ctx context.Context, id int64) ([]dto.MeasurementRecord, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Measurements(ctx, id)
}

// Update правит рейс: реквизиты, веса и время взвешиваний.
func (s *SessionService) Update(
	ctx context.Context,
	id int64,
	request dto.UpdateSessionRequest,
	userID int64,
) (*dto.SessionDTO, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}

	request.TransportNumber = strings.TrimSpace(request.TransportNumber)
	request.TrailerNumber = strings.TrimSpace(request.TrailerNumber)
	request.DriverName = strings.TrimSpace(request.DriverName)

	if request.TransportNumber == "" {
		return nil, &ValidationError{Message: "Укажите номер транспортного средства"}
	}
	if err := s.checkBlacklist(ctx, request.TransportNumber, request.TrailerNumber, request.DriverName); err != nil {
		return nil, err
	}
	if request.FirstWeightKg < 0 {
		return nil, &ValidationError{Message: "Вес первого взвешивания не может быть отрицательным"}
	}
	if request.FirstMeasuredAt.IsZero() {
		return nil, &ValidationError{Message: "Укажите дату и время первого взвешивания"}
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.Status == "cancelled" {
		return nil, &ValidationError{Message: "Отменённый рейс нельзя редактировать"}
	}

	if existing.Status == "closed" {
		if request.SecondWeightKg == nil {
			return nil, &ValidationError{Message: "Укажите вес второго взвешивания"}
		}
		if *request.SecondWeightKg < 0 {
			return nil, &ValidationError{Message: "Вес второго взвешивания не может быть отрицательным"}
		}
		if request.SecondMeasuredAt == nil {
			return nil, &ValidationError{Message: "Укажите дату и время второго взвешивания"}
		}
	} else {
		request.SecondWeightKg = nil
		request.SecondMeasuredAt = nil
	}

	if err := s.validateRequiredFields(
		ctx,
		request.TransportNumber,
		request.DriverName,
		request.TrailerNumber,
		request.CargoID,
		request.CounterpartyID,
		request.CarrierID,
		request.StorageID,
		request.OperationTypeID,
	); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, request, userID)
}

// Audit возвращает историю изменений рейса.
func (s *SessionService) Audit(ctx context.Context, id int64) ([]dto.SessionAuditRecord, error) {
	if id <= 0 {
		return nil, &ValidationError{Message: "Некорректный идентификатор рейса"}
	}
	return s.repo.Audit(ctx, id)
}

// checkBlacklist отклоняет рейс, если ТС/прицеп/водитель в чёрном списке.
func (s *SessionService) checkBlacklist(
	ctx context.Context,
	transport string,
	trailer string,
	driver string,
) error {
	kind, name, found, err := s.repo.BlacklistMatch(ctx, transport, trailer, driver)
	if err != nil {
		return err
	}
	if found {
		return &ValidationError{Message: fmt.Sprintf("%s «%s» в чёрном списке", kind, name)}
	}
	return nil
}

// validateRequiredFields проверяет обязательность полей по настройкам весовой.
func (s *SessionService) validateRequiredFields(
	ctx context.Context,
	transport string,
	driver string,
	trailer string,
	cargo *int64,
	counterparty *int64,
	carrier *int64,
	storage *int64,
	operationType *int64,
) error {
	if s.settings == nil {
		return nil
	}

	settings, err := s.settings.App(ctx)
	if err != nil {
		return err
	}

	missing := make([]string, 0, 8)
	if strings.TrimSpace(transport) == "" {
		missing = append(missing, "номер ТС")
	}
	if settings.RequireDriver && strings.TrimSpace(driver) == "" {
		missing = append(missing, "водитель")
	}
	if settings.RequireTrailer && strings.TrimSpace(trailer) == "" {
		missing = append(missing, "прицеп")
	}
	if settings.RequireCargo && cargo == nil {
		missing = append(missing, "груз")
	}
	if settings.RequireCounterparty && counterparty == nil {
		missing = append(missing, "контрагент")
	}
	if settings.RequireCarrier && carrier == nil {
		missing = append(missing, "перевозчик")
	}
	if settings.RequireStorage && storage == nil {
		missing = append(missing, "склад")
	}
	if settings.RequireOperationType && operationType == nil {
		missing = append(missing, "тип операции")
	}

	if len(missing) > 0 {
		return &ValidationError{Message: "Заполните обязательные поля: " + strings.Join(missing, ", ")}
	}

	return nil
}

func (s *SessionService) Export(
	ctx context.Context,
	status string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.SessionDTO, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		status = "all"
	}

	return s.repo.Export(ctx, status, strings.TrimSpace(search), dateFrom, dateTo)
}
