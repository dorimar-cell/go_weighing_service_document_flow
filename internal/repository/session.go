package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

var ErrConflict = errors.New("conflict")

const sessionFrom = `
FROM weighing w
JOIN weighing_info wi ON wi.id = w.id_weighting_info
JOIN com_dir_transport_driver_trailer c ON c.id = wi.id_com_transport_driver_trailer
JOIN dir_transport t ON t.id = c.id_dir_transport
LEFT JOIN dir_trailer tr ON tr.id = c.id_dir_trailer
LEFT JOIN dir_driver d ON d.id = c.id_dir_driver
LEFT JOIN reg_weight rw ON rw.id = w.id_reg_weight
LEFT JOIN dir_invoice inv ON inv.id = wi.id_dir_invoice
LEFT JOIN dir_carrier cr ON cr.id = inv.id_dir_carrier
LEFT JOIN dir_counterparty cp ON cp.id = inv.id_dir_counterparty
LEFT JOIN dir_type_of_operation ot ON ot.id = inv.id_dir_type_of_operation
LEFT JOIN LATERAL (
    SELECT cg.id AS cargo_id,
           cg.name AS cargo_name,
           st.id AS storage_id,
           st.name AS storage_name
    FROM com_dir_invoice_dir_storage cids
    LEFT JOIN dir_cargo cg ON cg.id = cids.id_dir_cargo
    LEFT JOIN dir_storage st ON st.id = cids.id_dir_storage
    WHERE cids.id_dir_invoice = inv.id AND cids.deleted = false
    ORDER BY cids.id
    LIMIT 1
) line ON true`

const sessionSelect = `
SELECT w.id,
       t.name,
       COALESCE(tr.name, ''),
       COALESCE(d.name, ''),
       w.first_weight,
       w.second_weight,
       w.date_time_first_weight,
       w.date_time_second_weight,
       w.flag_weight,
       w.deleted,
       w.created_date,
       COALESCE(line.cargo_id, 0),
       COALESCE(line.cargo_name, ''),
       COALESCE(cp.id, 0),
       COALESCE(cp.name, ''),
       COALESCE(cr.id, 0),
       COALESCE(cr.name, ''),
       COALESCE(line.storage_id, 0),
       COALESCE(line.storage_name, ''),
       COALESCE(ot.id, 0),
       COALESCE(ot.name, ''),
       COALESCE(rw.id, 0),
       COALESCE(rw.name, '')` + sessionFrom

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

// Create создаёт первое взвешивание. operationId делает операцию идемпотентной.
func (r *SessionRepository) Create(
	ctx context.Context,
	req dto.CreateSessionRequest,
	userID int64,
) (*dto.SessionDTO, error) {
	if req.Measurement != nil {
		existing, err := r.findByOperation(ctx, req.Measurement.OperationID)
		if err != nil {
			return nil, err
		}
		if existing != 0 {
			return r.Get(ctx, existing)
		}
	}

	weight := 0.0
	measuredAt := time.Now()
	manual := false
	if req.ManualWeightKg != nil {
		weight = *req.ManualWeightKg
		manual = true
	} else {
		weight = req.Measurement.WeightKg
		measuredAt = req.Measurement.MeasuredAt
	}

	var handleFirst any
	if manual {
		handleFirst = weight
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create session: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	transportID, err := findOrCreateNamed(ctx, tx, "dir_transport", req.TransportNumber)
	if err != nil {
		return nil, err
	}

	var trailerID *int64
	if strings.TrimSpace(req.TrailerNumber) != "" {
		id, err := findOrCreateNamed(ctx, tx, "dir_trailer", req.TrailerNumber)
		if err != nil {
			return nil, err
		}
		trailerID = &id
	}

	var driverID *int64
	if strings.TrimSpace(req.DriverName) != "" {
		id, err := findOrCreateNamed(ctx, tx, "dir_driver", req.DriverName)
		if err != nil {
			return nil, err
		}
		driverID = &id
	}

	var linkID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO com_dir_transport_driver_trailer
			(id_dir_driver, id_dir_transport, id_dir_trailer, deleted, created_by_user_id, created_date)
		VALUES ($1, $2, $3, false, $4, now())
		RETURNING id`,
		driverID, transportID, trailerID, userID,
	).Scan(&linkID); err != nil {
		return nil, fmt.Errorf("insert transport link: %w", err)
	}

	invoiceID, err := createInvoice(ctx, tx, req, userID)
	if err != nil {
		return nil, err
	}

	var infoID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO weighing_info (id_com_transport_driver_trailer, id_dir_invoice, deleted, created_by_user_id, created_date)
		VALUES ($1, $2, false, $3, now())
		RETURNING id`,
		linkID, invoiceID, userID,
	).Scan(&infoID); err != nil {
		return nil, fmt.Errorf("insert weighing_info: %w", err)
	}

	var weighingID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO weighing
			(id_weighting_info, first_weight, handle_first_weight, flag_weight, deleted,
			 created_by_user_id, created_date, id_first_weight_user, date_time_first_weight, ref,
			 first_weight_1c, is_flow_weight, has_edited, second_weight_1c, id_reg_weight)
		VALUES ($1, $2, $3, false, false, $6, now(), $6, $4, gen_random_uuid(),
		        false, false, false, false, $5)
		RETURNING id`,
		infoID, weight, handleFirst, measuredAt, req.PassID, userID,
	).Scan(&weighingID); err != nil {
		return nil, fmt.Errorf("insert weighing: %w", err)
	}

	if req.Measurement != nil {
		if err := insertMeasurement(ctx, tx, weighingID, 1, *req.Measurement); err != nil {
			return nil, err
		}
	}

	if err := insertSessionAudit(ctx, tx, weighingID, "created", userID, map[string]any{
		"firstWeightKg":   weight,
		"transportNumber": req.TransportNumber,
		"manual":          manual,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create session: %w", err)
	}

	return r.Get(ctx, weighingID)
}

// Close закрывает вторым взвешиванием.
func (r *SessionRepository) Close(
	ctx context.Context,
	id int64,
	measurement dto.MeasurementDTO,
	userID int64,
) (*dto.SessionDTO, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin close session: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var firstWeight float64
	var flagWeight, deleted bool
	var passID sql.NullInt64
	if err := tx.QueryRow(ctx,
		`SELECT first_weight, flag_weight, deleted, id_reg_weight FROM weighing WHERE id = $1 FOR UPDATE`, id,
	).Scan(&firstWeight, &flagWeight, &deleted, &passID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("lock weighing: %w", err)
	}

	if deleted {
		return nil, fmt.Errorf("%w: взвешивание отменено", ErrConflict)
	}
	if flagWeight {
		return nil, fmt.Errorf("%w: взвешивание уже завершено", ErrConflict)
	}

	existing, err := findByOperationTx(ctx, tx, measurement.OperationID)
	if err != nil {
		return nil, err
	}
	if existing != 0 {
		if existing != id {
			return nil, fmt.Errorf("%w: измерение уже привязано к другому рейсу", ErrConflict)
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit idempotent close: %w", err)
		}
		return r.Get(ctx, id)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE weighing
		SET second_weight = $2,
		    date_time_second_weight = $3,
		    flag_weight = true,
		    id_second_weight_user = $4,
		    modified_by_user_id = $4,
		    modified_date = now()
		WHERE id = $1`,
		id, measurement.WeightKg, measurement.MeasuredAt, userID,
	); err != nil {
		return nil, fmt.Errorf("update weighing: %w", err)
	}

	if err := insertMeasurement(ctx, tx, id, 2, measurement); err != nil {
		return nil, err
	}

	if err := insertSessionAudit(ctx, tx, id, "closed", userID, map[string]any{
		"secondWeightKg": measurement.WeightKg,
	}); err != nil {
		return nil, err
	}

	if passID.Valid {
		if _, err := tx.Exec(ctx,
			`UPDATE reg_weight SET pass_enabled = false WHERE id = $1 AND id_system_dir_type_of_pass = 3`,
			passID.Int64,
		); err != nil {
			return nil, fmt.Errorf("close one-time pass: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit close session: %w", err)
	}

	return r.Get(ctx, id)
}

// CloseWithManualTare закрывает рейс в поточном режиме: тара вводится вручную.
func (r *SessionRepository) CloseWithManualTare(
	ctx context.Context,
	id int64,
	tareKg float64,
	measuredAt time.Time,
	userID int64,
) (*dto.SessionDTO, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin flow close: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var flagWeight, deleted bool
	var passID sql.NullInt64
	if err := tx.QueryRow(ctx,
		`SELECT flag_weight, deleted, id_reg_weight FROM weighing WHERE id = $1 FOR UPDATE`, id,
	).Scan(&flagWeight, &deleted, &passID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("lock weighing: %w", err)
	}

	if deleted {
		return nil, fmt.Errorf("%w: взвешивание отменено", ErrConflict)
	}
	if flagWeight {
		return nil, fmt.Errorf("%w: взвешивание уже завершено", ErrConflict)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE weighing
		SET second_weight = $2,
		    handle_second_weight = $2,
		    date_time_second_weight = $3,
		    flag_weight = true,
		    id_second_weight_user = $4,
		    modified_by_user_id = $4,
		    modified_date = now()
		WHERE id = $1`,
		id, tareKg, measuredAt, userID,
	); err != nil {
		return nil, fmt.Errorf("update flow weighing: %w", err)
	}

	if passID.Valid {
		if _, err := tx.Exec(ctx,
			`UPDATE reg_weight SET pass_enabled = false WHERE id = $1 AND id_system_dir_type_of_pass = 3`,
			passID.Int64,
		); err != nil {
			return nil, fmt.Errorf("close one-time pass: %w", err)
		}
	}

	if err := insertSessionAudit(ctx, tx, id, "flow_closed", userID, map[string]any{
		"tareKg": tareKg,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit flow close: %w", err)
	}

	return r.Get(ctx, id)
}

func (r *SessionRepository) List(
	ctx context.Context,
	status string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
	page int,
	pageSize int,
) (*dto.SessionListResponse, error) {
	where, args := buildSessionWhere(status, search, dateFrom, dateTo)

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) "+sessionFrom+" "+where, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count sessions: %w", err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	query := sessionSelect + " " + where +
		fmt.Sprintf(" ORDER BY w.created_date DESC LIMIT $%d OFFSET $%d", limitArg, offsetArg)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	items := make([]dto.SessionDTO, 0, pageSize)
	for rows.Next() {
		item, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}

	return &dto.SessionListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// Measurements возвращает историю измерений рейса.
func (r *SessionRepository) Measurements(
	ctx context.Context,
	id int64,
) ([]dto.MeasurementRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT operation_id, kind, weight_kg::float8, stable, terminal, platform, measured_at, checksum
		FROM weighing_measurement
		WHERE id_weighing = $1
		ORDER BY kind, measured_at`, id)
	if err != nil {
		return nil, fmt.Errorf("list measurements: %w", err)
	}
	defer rows.Close()

	items := make([]dto.MeasurementRecord, 0, 2)
	for rows.Next() {
		var item dto.MeasurementRecord
		if err := rows.Scan(
			&item.OperationID,
			&item.Kind,
			&item.WeightKg,
			&item.Stable,
			&item.Terminal,
			&item.Platform,
			&item.MeasuredAt,
			&item.Checksum,
		); err != nil {
			return nil, fmt.Errorf("scan measurement: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate measurements: %w", err)
	}

	return items, nil
}

// Export возвращает все рейсы по фильтру без пагинации (для CSV).
func (r *SessionRepository) Export(
	ctx context.Context,
	status string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.SessionDTO, error) {
	where, args := buildSessionWhere(status, search, dateFrom, dateTo)
	query := sessionSelect + " " + where + " ORDER BY w.created_date DESC LIMIT 10000"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("export sessions: %w", err)
	}
	defer rows.Close()

	items := make([]dto.SessionDTO, 0, 256)
	for rows.Next() {
		item, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate export sessions: %w", err)
	}

	return items, nil
}

func (r *SessionRepository) Get(ctx context.Context, id int64) (*dto.SessionDTO, error) {
	row := r.pool.QueryRow(ctx, sessionSelect+" WHERE w.id = $1", id)
	item, err := scanSession(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (r *SessionRepository) Cancel(ctx context.Context, id int64, userID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin cancel session: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var infoID int64
	err = tx.QueryRow(ctx, `
		UPDATE weighing
		SET deleted = true, modified_by_user_id = $2, modified_date = now()
		WHERE id = $1 AND flag_weight = false AND deleted = false
		RETURNING id_weighting_info`, id, userID,
	).Scan(&infoID)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("cancel weighing: %w", err)
		}

		var flagWeight, deleted bool
		if err := tx.QueryRow(ctx,
			`SELECT flag_weight, deleted FROM weighing WHERE id = $1`, id,
		).Scan(&flagWeight, &deleted); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}

		if deleted {
			return fmt.Errorf("%w: взвешивание уже отменено", ErrConflict)
		}
		return fmt.Errorf("%w: завершённое взвешивание нельзя отменить", ErrConflict)
	}

	if _, err := tx.Exec(ctx,
		`UPDATE weighing_info SET deleted = true WHERE id = $1`, infoID,
	); err != nil {
		return fmt.Errorf("cancel weighing_info: %w", err)
	}

	if err := insertSessionAudit(ctx, tx, id, "cancelled", userID, nil); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SessionRepository) findByOperation(ctx context.Context, operationID string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`SELECT id_weighing FROM weighing_measurement WHERE operation_id = $1`, operationID,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("find measurement: %w", err)
	}
	return id, nil
}

func findByOperationTx(ctx context.Context, tx pgx.Tx, operationID string) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx,
		`SELECT id_weighing FROM weighing_measurement WHERE operation_id = $1`, operationID,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("find measurement: %w", err)
	}
	return id, nil
}

func insertMeasurement(
	ctx context.Context,
	tx pgx.Tx,
	weighingID int64,
	kind int,
	measurement dto.MeasurementDTO,
) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO weighing_measurement
			(id_weighing, kind, operation_id, weight_kg, stable, terminal, platform, measured_at, checksum)
		VALUES ($1, $2, $3, $4, true, $5, $6, $7, $8)`,
		weighingID,
		kind,
		measurement.OperationID,
		measurement.WeightKg,
		measurement.Terminal,
		measurement.Platform,
		measurement.MeasuredAt,
		measurement.Checksum,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: измерение уже сохранено", ErrConflict)
		}
		return fmt.Errorf("insert measurement: %w", err)
	}

	return nil
}

// createInvoice сохраняет груз/контрагента/перевозчика/склад/тип операции.
// Если ничего не выбрано, накладная не создаётся.
func createInvoice(
	ctx context.Context,
	tx pgx.Tx,
	req dto.CreateSessionRequest,
	userID int64,
) (*int64, error) {
	if req.CargoID == nil && req.StorageID == nil && req.CounterpartyID == nil &&
		req.CarrierID == nil && req.OperationTypeID == nil {
		return nil, nil
	}

	var invoiceID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO dir_invoice
			(name, id_dir_carrier, id_dir_counterparty, id_dir_type_of_operation, date, time, deleted, created_by_user_id, created_date)
		VALUES ('', $1, $2, $3, current_date, current_time, false, $4, now())
		RETURNING id`,
		req.CarrierID, req.CounterpartyID, req.OperationTypeID, userID,
	).Scan(&invoiceID); err != nil {
		return nil, fmt.Errorf("insert invoice: %w", err)
	}

	if req.CargoID != nil || req.StorageID != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO com_dir_invoice_dir_storage
				(id_dir_cargo, id_dir_invoice, id_dir_storage, deleted, created_by_user_id, created_date)
			VALUES ($1, $2, $3, false, $4, now())`,
			req.CargoID, invoiceID, req.StorageID, userID,
		); err != nil {
			return nil, fmt.Errorf("insert invoice line: %w", err)
		}
	}

	return &invoiceID, nil
}

func findOrCreateNamed(ctx context.Context, tx pgx.Tx, table, name string) (int64, error) {
	var id int64
	selectQuery := fmt.Sprintf(
		`SELECT id FROM %s WHERE lower(name) = lower($1) AND deleted_at IS NULL LIMIT 1`,
		table,
	)

	err := tx.QueryRow(ctx, selectQuery, strings.TrimSpace(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("find %s: %w", table, err)
	}

	insertQuery := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1) RETURNING id`, table)
	if err := tx.QueryRow(ctx, insertQuery, strings.TrimSpace(name)).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert %s: %w", table, err)
	}

	return id, nil
}

func buildSessionWhere(
	status string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
) (string, []any) {
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)

	switch status {
	case "open":
		where = append(where, "w.flag_weight = false AND w.deleted = false")
	case "closed":
		where = append(where, "w.flag_weight = true AND w.deleted = false")
	case "cancelled":
		where = append(where, "w.deleted = true")
	case "all", "":
		// без фильтра
	default:
		where = append(where, "1 = 0")
	}

	if term := strings.TrimSpace(search); term != "" {
		args = append(args, "%"+term+"%")
		position := len(args)
		where = append(where, fmt.Sprintf(
			"(t.name ILIKE $%d OR COALESCE(tr.name, '') ILIKE $%d OR COALESCE(rw.name, '') ILIKE $%d)",
			position, position, position,
		))
	}

	if dateFrom != nil {
		args = append(args, *dateFrom)
		where = append(where, fmt.Sprintf("w.date_time_first_weight >= $%d", len(args)))
	}

	if dateTo != nil {
		args = append(args, *dateTo)
		where = append(where, fmt.Sprintf("w.date_time_first_weight < $%d", len(args)))
	}

	return " WHERE " + strings.Join(where, " AND "), args
}

type scanner interface {
	Scan(dest ...any) error
}

// aggregateSpec — SQL-выражения группировки для агрегатных отчётов.
type aggregateSpec struct {
	groupExpr string
	subExpr   string
}

var aggregateSpecs = map[string]aggregateSpec{
	"cargo": {
		groupExpr: "COALESCE(line.cargo_name, 'Без груза')",
		subExpr:   "''::text",
	},
	"counterparty": {
		groupExpr: "COALESCE(cp.name, 'Без контрагента')",
		subExpr:   "COALESCE(line.cargo_name, 'Без груза')",
	},
	"carrier": {
		groupExpr: "COALESCE(cr.name, 'Без перевозчика')",
		subExpr:   "COALESCE(line.cargo_name, 'Без груза')",
	},
	"driver": {
		groupExpr: "COALESCE(d.name, 'Без водителя')",
		subExpr:   "COALESCE(line.cargo_name, 'Без груза')",
	},
	"storage": {
		groupExpr: "COALESCE(line.storage_name, 'Без склада')",
		subExpr:   "COALESCE(ot.name, 'Без типа операции')",
	},
	"operation": {
		groupExpr: "COALESCE(ot.name, 'Без типа операции')",
		subExpr:   "''::text",
	},
}

// Aggregate считает закрытые рейсы по измерению: рейсов и нетто.
func (r *SessionRepository) Aggregate(
	ctx context.Context,
	dimension string,
	dateFrom *time.Time,
	dateTo *time.Time,
	filter dto.AggregateFilter,
) ([]dto.AggregateRow, error) {
	spec, ok := aggregateSpecs[dimension]
	if !ok {
		return nil, ErrNotFound
	}

	where := []string{"w.flag_weight = true", "w.deleted = false"}
	args := make([]any, 0, 6)

	if dateFrom != nil {
		args = append(args, *dateFrom)
		where = append(where, fmt.Sprintf("w.date_time_first_weight >= $%d", len(args)))
	}
	if dateTo != nil {
		args = append(args, *dateTo)
		where = append(where, fmt.Sprintf("w.date_time_first_weight < $%d", len(args)))
	}
	if len(filter.CargoIDs) > 0 {
		args = append(args, filter.CargoIDs)
		where = append(where, fmt.Sprintf("line.cargo_id = ANY($%d)", len(args)))
	}
	if len(filter.CounterpartyIDs) > 0 {
		args = append(args, filter.CounterpartyIDs)
		where = append(where, fmt.Sprintf("cp.id = ANY($%d)", len(args)))
	}
	if len(filter.CarrierIDs) > 0 {
		args = append(args, filter.CarrierIDs)
		where = append(where, fmt.Sprintf("cr.id = ANY($%d)", len(args)))
	}
	if len(filter.StorageIDs) > 0 {
		args = append(args, filter.StorageIDs)
		where = append(where, fmt.Sprintf("line.storage_id = ANY($%d)", len(args)))
	}

	query := fmt.Sprintf(`
		SELECT %s AS grp,
		       %s AS sub,
		       count(*) AS sessions,
		       COALESCE(sum(abs(w.first_weight - COALESCE(w.second_weight, w.first_weight))), 0)::float8 AS netto
		%s
		WHERE %s
		GROUP BY 1, 2
		ORDER BY 1, 2`,
		spec.groupExpr, spec.subExpr, sessionFrom, strings.Join(where, " AND "),
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("aggregate %s: %w", dimension, err)
	}
	defer rows.Close()

	items := make([]dto.AggregateRow, 0, 64)
	for rows.Next() {
		var item dto.AggregateRow
		if err := rows.Scan(&item.Group, &item.SubGroup, &item.Sessions, &item.NettoKg); err != nil {
			return nil, fmt.Errorf("scan aggregate %s: %w", dimension, err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func scanSession(row scanner) (*dto.SessionDTO, error) {
	var item dto.SessionDTO
	var secondWeight sql.NullFloat64
	var secondAt sql.NullTime
	var flagWeight, deleted bool
	var cargoID, counterpartyID, carrierID, storageID, operationTypeID, passID int64

	if err := row.Scan(
		&item.ID,
		&item.TransportNumber,
		&item.TrailerNumber,
		&item.DriverName,
		&item.FirstWeightKg,
		&secondWeight,
		&item.FirstMeasuredAt,
		&secondAt,
		&flagWeight,
		&deleted,
		&item.CreatedAt,
		&cargoID,
		&item.CargoName,
		&counterpartyID,
		&item.CounterpartyName,
		&carrierID,
		&item.CarrierName,
		&storageID,
		&item.StorageName,
		&operationTypeID,
		&item.OperationTypeName,
		&passID,
		&item.PassNumber,
	); err != nil {
		return nil, err
	}

	if cargoID != 0 {
		item.CargoID = &cargoID
	}
	if counterpartyID != 0 {
		item.CounterpartyID = &counterpartyID
	}
	if carrierID != 0 {
		item.CarrierID = &carrierID
	}
	if storageID != 0 {
		item.StorageID = &storageID
	}
	if operationTypeID != 0 {
		item.OperationTypeID = &operationTypeID
	}
	if passID != 0 {
		item.PassID = &passID
	}

	if secondWeight.Valid {
		value := secondWeight.Float64
		item.SecondWeightKg = &value
		net := math.Abs(item.FirstWeightKg - value)
		item.NetWeightKg = &net
	}
	if secondAt.Valid {
		value := secondAt.Time
		item.SecondMeasuredAt = &value
	}

	switch {
	case deleted:
		item.Status = "cancelled"
	case flagWeight:
		item.Status = "closed"
	default:
		item.Status = "open"
	}

	return &item, nil
}
