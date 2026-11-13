package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"go_weighing_service/internal/dto"
)

// insertSessionAudit пишет запись истории изменений рейса в текущей транзакции.
func insertSessionAudit(
	ctx context.Context,
	tx pgx.Tx,
	weighingID int64,
	action string,
	userID int64,
	payload any,
) error {
	if payload == nil {
		payload = map[string]any{}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal audit payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO weighing_audit (id_weighing, action, user_id, user_name, payload)
		VALUES ($1, $2, $3, (SELECT name FROM dir_users WHERE id = $3), $4::jsonb)`,
		weighingID, action, userID, string(data),
	); err != nil {
		return fmt.Errorf("insert session audit: %w", err)
	}

	return nil
}

// Audit возвращает историю изменений рейса (новые записи первыми).
func (r *SessionRepository) Audit(ctx context.Context, id int64) ([]dto.SessionAuditRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, action, user_id, COALESCE(user_name, ''), payload, created_at
		FROM weighing_audit
		WHERE id_weighing = $1
		ORDER BY created_at DESC, id DESC`, id)
	if err != nil {
		return nil, fmt.Errorf("list session audit: %w", err)
	}
	defer rows.Close()

	items := make([]dto.SessionAuditRecord, 0, 8)
	for rows.Next() {
		var item dto.SessionAuditRecord
		var payload []byte
		if err := rows.Scan(
			&item.ID, &item.Action, &item.UserID, &item.UserName, &payload, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan session audit: %w", err)
		}
		item.Payload = json.RawMessage(payload)
		items = append(items, item)
	}

	return items, rows.Err()
}

// AuditJournal возвращает постраничный журнал аудита по всем рейсам.
func (r *SessionRepository) AuditJournal(
	ctx context.Context,
	action string,
	search string,
	dateFrom *time.Time,
	dateTo *time.Time,
	page int,
	pageSize int,
) (*dto.AuditJournalResponse, error) {
	from := `
		FROM weighing_audit a
		LEFT JOIN weighing w ON w.id = a.id_weighing
		LEFT JOIN weighing_info wi ON wi.id = w.id_weighting_info
		LEFT JOIN com_dir_transport_driver_trailer c ON c.id = wi.id_com_transport_driver_trailer
		LEFT JOIN dir_transport t ON t.id = c.id_dir_transport`

	where := []string{"1 = 1"}
	args := make([]any, 0, 4)

	if action != "" {
		args = append(args, action)
		where = append(where, fmt.Sprintf("a.action = $%d", len(args)))
	}
	if term := strings.TrimSpace(search); term != "" {
		args = append(args, "%"+term+"%")
		position := len(args)
		where = append(where, fmt.Sprintf(
			"(COALESCE(t.name, '') ILIKE $%d OR COALESCE(a.user_name, '') ILIKE $%d)",
			position, position,
		))
	}
	if dateFrom != nil {
		args = append(args, *dateFrom)
		where = append(where, fmt.Sprintf("a.created_at >= $%d", len(args)))
	}
	if dateTo != nil {
		args = append(args, *dateTo)
		where = append(where, fmt.Sprintf("a.created_at < $%d", len(args)))
	}

	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) "+from+whereSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count audit journal: %w", err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	query := `
		SELECT a.id, a.id_weighing, a.action, a.user_id, COALESCE(a.user_name, ''),
		       COALESCE(t.name, ''), a.payload, a.created_at ` + from + whereSQL +
		fmt.Sprintf(" ORDER BY a.created_at DESC, a.id DESC LIMIT $%d OFFSET $%d", limitArg, offsetArg)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit journal: %w", err)
	}
	defer rows.Close()

	items := make([]dto.AuditJournalItem, 0, pageSize)
	for rows.Next() {
		var item dto.AuditJournalItem
		var payload []byte
		if err := rows.Scan(
			&item.ID, &item.SessionID, &item.Action, &item.UserID, &item.UserName,
			&item.TransportNumber, &payload, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan audit journal: %w", err)
		}
		item.Payload = json.RawMessage(payload)
		items = append(items, item)
	}

	return &dto.AuditJournalResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, rows.Err()
}

// Restore возвращает отменённый рейс.
func (r *SessionRepository) Restore(ctx context.Context, id int64, userID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin restore session: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var infoID int64
	var deleted bool
	if err := tx.QueryRow(ctx,
		`SELECT id_weighting_info, deleted FROM weighing WHERE id = $1 FOR UPDATE`, id,
	).Scan(&infoID, &deleted); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("lock session for restore: %w", err)
	}

	if !deleted {
		return fmt.Errorf("%w: рейс не отменён", ErrConflict)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE weighing
		SET deleted = false, modified_by_user_id = $2, modified_date = now()
		WHERE id = $1`, id, userID,
	); err != nil {
		return fmt.Errorf("restore weighing: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE weighing_info
		SET deleted = false, modified_by_user_id = $2, modified_date = now()
		WHERE id = $1`, infoID, userID,
	); err != nil {
		return fmt.Errorf("restore weighing_info: %w", err)
	}

	if err := insertSessionAudit(ctx, tx, id, "restored", userID, nil); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Update правит рейс: транспорт, реквизиты, веса и время взвешиваний.
func (r *SessionRepository) Update(
	ctx context.Context,
	id int64,
	req dto.UpdateSessionRequest,
	userID int64,
) (*dto.SessionDTO, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin update session: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		oldFirstWeight, oldSecondWeight                       sql.NullFloat64
		oldFirstAt, oldSecondAt                               sql.NullTime
		flagWeight, deleted                                   bool
		infoID, linkID, invoiceID                             int64
		transportID, oldTrailerID, oldDriverID                int64
		oldCounterpartyID, oldCarrierID, oldOperationTypeID   int64
		lineID, oldCargoID, oldStorageID                      int64
		oldTransportName, oldTrailerName, oldDriverName       string
		oldCounterpartyName, oldCarrierName, oldOperationName string
		oldCargoName, oldStorageName                          string
	)

	err = tx.QueryRow(ctx, `
		SELECT w.first_weight::float8,
		       w.second_weight::float8,
		       w.date_time_first_weight,
		       w.date_time_second_weight,
		       w.flag_weight,
		       w.deleted,
		       w.id_weighting_info,
		       wi.id_com_transport_driver_trailer,
		       COALESCE(wi.id_dir_invoice, 0),
		       c.id_dir_transport,
		       COALESCE(c.id_dir_trailer, 0),
		       COALESCE(c.id_dir_driver, 0),
		       COALESCE(inv.id_dir_counterparty, 0),
		       COALESCE(inv.id_dir_carrier, 0),
		       COALESCE(inv.id_dir_type_of_operation, 0),
		       COALESCE(line.id, 0),
		       COALESCE(line.id_dir_cargo, 0),
		       COALESCE(line.id_dir_storage, 0),
		       COALESCE(t.name, ''),
		       COALESCE(tr.name, ''),
		       COALESCE(d.name, ''),
		       COALESCE(cp.name, ''),
		       COALESCE(cr.name, ''),
		       COALESCE(ot.name, ''),
		       COALESCE(cg.name, ''),
		       COALESCE(st.name, '')
		FROM weighing w
		JOIN weighing_info wi ON wi.id = w.id_weighting_info
		JOIN com_dir_transport_driver_trailer c ON c.id = wi.id_com_transport_driver_trailer
		LEFT JOIN dir_transport t ON t.id = c.id_dir_transport
		LEFT JOIN dir_trailer tr ON tr.id = c.id_dir_trailer
		LEFT JOIN dir_driver d ON d.id = c.id_dir_driver
		LEFT JOIN dir_invoice inv ON inv.id = wi.id_dir_invoice
		LEFT JOIN dir_counterparty cp ON cp.id = inv.id_dir_counterparty
		LEFT JOIN dir_carrier cr ON cr.id = inv.id_dir_carrier
		LEFT JOIN dir_type_of_operation ot ON ot.id = inv.id_dir_type_of_operation
		LEFT JOIN LATERAL (
		    SELECT cids.id, cids.id_dir_cargo, cids.id_dir_storage
		    FROM com_dir_invoice_dir_storage cids
		    WHERE cids.id_dir_invoice = inv.id AND cids.deleted = false
		    ORDER BY cids.id
		    LIMIT 1
		) line ON true
		LEFT JOIN dir_cargo cg ON cg.id = line.id_dir_cargo
		LEFT JOIN dir_storage st ON st.id = line.id_dir_storage
		WHERE w.id = $1
		FOR UPDATE OF w`, id,
	).Scan(
		&oldFirstWeight, &oldSecondWeight, &oldFirstAt, &oldSecondAt,
		&flagWeight, &deleted, &infoID, &linkID, &invoiceID,
		&transportID, &oldTrailerID, &oldDriverID,
		&oldCounterpartyID, &oldCarrierID, &oldOperationTypeID,
		&lineID, &oldCargoID, &oldStorageID,
		&oldTransportName, &oldTrailerName, &oldDriverName,
		&oldCounterpartyName, &oldCarrierName, &oldOperationName,
		&oldCargoName, &oldStorageName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lock session for update: %w", err)
	}

	if deleted {
		return nil, fmt.Errorf("%w: отменённый рейс нельзя редактировать", ErrConflict)
	}

	newTransportID, err := findOrCreateNamed(ctx, tx, "dir_transport", req.TransportNumber)
	if err != nil {
		return nil, err
	}

	var newTrailerID, newDriverID *int64
	if strings.TrimSpace(req.TrailerNumber) != "" {
		value, err := findOrCreateNamed(ctx, tx, "dir_trailer", req.TrailerNumber)
		if err != nil {
			return nil, err
		}
		newTrailerID = &value
	}
	if strings.TrimSpace(req.DriverName) != "" {
		value, err := findOrCreateNamed(ctx, tx, "dir_driver", req.DriverName)
		if err != nil {
			return nil, err
		}
		newDriverID = &value
	}

	if _, err := tx.Exec(ctx, `
		UPDATE com_dir_transport_driver_trailer
		SET id_dir_transport = $2, id_dir_trailer = $3, id_dir_driver = $4,
		    modified_by_user_id = $5, modified_date = now()
		WHERE id = $1`,
		linkID, newTransportID, newTrailerID, newDriverID, userID,
	); err != nil {
		return nil, fmt.Errorf("update transport link: %w", err)
	}

	var newInvoiceID int64
	if invoiceID != 0 {
		newInvoiceID = invoiceID
		if _, err := tx.Exec(ctx, `
			UPDATE dir_invoice
			SET id_dir_counterparty = $2, id_dir_carrier = $3, id_dir_type_of_operation = $4,
			    modified_by_user_id = $5, modified_date = now()
			WHERE id = $1`,
			invoiceID, req.CounterpartyID, req.CarrierID, req.OperationTypeID, userID,
		); err != nil {
			return nil, fmt.Errorf("update invoice: %w", err)
		}
	} else if req.CounterpartyID != nil || req.CarrierID != nil || req.OperationTypeID != nil {
		if err := tx.QueryRow(ctx, `
			INSERT INTO dir_invoice
				(name, id_dir_carrier, id_dir_counterparty, id_dir_type_of_operation,
				 date, time, deleted, created_by_user_id, created_date)
			VALUES ('', $1, $2, $3, current_date, current_time, false, $4, now())
			RETURNING id`,
			req.CarrierID, req.CounterpartyID, req.OperationTypeID, userID,
		).Scan(&newInvoiceID); err != nil {
			return nil, fmt.Errorf("insert invoice: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE weighing_info SET id_dir_invoice = $2 WHERE id = $1`,
			infoID, newInvoiceID,
		); err != nil {
			return nil, fmt.Errorf("link invoice: %w", err)
		}
	}

	if lineID != 0 {
		if _, err := tx.Exec(ctx, `
			UPDATE com_dir_invoice_dir_storage
			SET id_dir_cargo = $2, id_dir_storage = $3,
			    modified_by_user_id = $4, modified_date = now()
			WHERE id = $1`,
			lineID, req.CargoID, req.StorageID, userID,
		); err != nil {
			return nil, fmt.Errorf("update invoice line: %w", err)
		}
	} else if newInvoiceID != 0 && (req.CargoID != nil || req.StorageID != nil) {
		if _, err := tx.Exec(ctx, `
			INSERT INTO com_dir_invoice_dir_storage
				(id_dir_cargo, id_dir_invoice, id_dir_storage, deleted, created_by_user_id, created_date)
			VALUES ($1, $2, $3, false, $4, now())`,
			req.CargoID, newInvoiceID, req.StorageID, userID,
		); err != nil {
			return nil, fmt.Errorf("insert invoice line: %w", err)
		}
	}

	var secondWeight *float64
	var secondMeasuredAt *time.Time
	if flagWeight {
		secondWeight = req.SecondWeightKg
		secondMeasuredAt = req.SecondMeasuredAt
	}

	if _, err := tx.Exec(ctx, `
		UPDATE weighing
		SET first_weight = $2,
		    second_weight = $3,
		    date_time_first_weight = $4,
		    date_time_second_weight = $5,
		    has_edited = true,
		    modified_by_user_id = $6,
		    modified_date = now()
		WHERE id = $1`,
		id, req.FirstWeightKg, secondWeight, req.FirstMeasuredAt, secondMeasuredAt, userID,
	); err != nil {
		return nil, fmt.Errorf("update weighing: %w", err)
	}

	changes := make(map[string]sessionChange)
	addStringChange(changes, "transport", oldTransportName, req.TransportNumber)
	addStringChange(changes, "trailer", oldTrailerName, req.TrailerNumber)
	addStringChange(changes, "driver", oldDriverName, req.DriverName)

	lookup := func(table string, value int64) string {
		if value == 0 {
			return ""
		}
		name, _ := lookupName(ctx, tx, table, value)
		return name
	}

	addIDChange(changes, "cargo", oldCargoID, req.CargoID, lookup("dir_cargo", idValue(req.CargoID)), oldCargoName)
	addIDChange(changes, "counterparty", oldCounterpartyID, req.CounterpartyID,
		lookup("dir_counterparty", idValue(req.CounterpartyID)), oldCounterpartyName)
	addIDChange(changes, "carrier", oldCarrierID, req.CarrierID,
		lookup("dir_carrier", idValue(req.CarrierID)), oldCarrierName)
	addIDChange(changes, "storage", oldStorageID, req.StorageID,
		lookup("dir_storage", idValue(req.StorageID)), oldStorageName)
	addIDChange(changes, "operationType", oldOperationTypeID, req.OperationTypeID,
		lookup("dir_type_of_operation", idValue(req.OperationTypeID)), oldOperationName)

	if oldFirstWeight.Valid && oldFirstWeight.Float64 != req.FirstWeightKg {
		changes["firstWeightKg"] = sessionChange{Old: oldFirstWeight.Float64, New: req.FirstWeightKg}
	}

	oldSecond := any(nil)
	if oldSecondWeight.Valid {
		oldSecond = oldSecondWeight.Float64
	}
	newSecond := any(nil)
	if secondWeight != nil {
		newSecond = *secondWeight
	}
	if fmt.Sprint(oldSecond) != fmt.Sprint(newSecond) {
		changes["secondWeightKg"] = sessionChange{Old: oldSecond, New: newSecond}
	}

	oldFirstAText := ""
	if oldFirstAt.Valid {
		oldFirstAText = oldFirstAt.Time.Format("2006-01-02T15:04")
	}
	newFirstAText := req.FirstMeasuredAt.Format("2006-01-02T15:04")
	if oldFirstAText != newFirstAText {
		changes["firstMeasuredAt"] = sessionChange{Old: oldFirstAText, New: newFirstAText}
	}

	oldSecondAText := ""
	if oldSecondAt.Valid {
		oldSecondAText = oldSecondAt.Time.Format("2006-01-02T15:04")
	}
	newSecondAText := ""
	if secondMeasuredAt != nil {
		newSecondAText = secondMeasuredAt.Format("2006-01-02T15:04")
	}
	if oldSecondAText != newSecondAText {
		changes["secondMeasuredAt"] = sessionChange{Old: oldSecondAText, New: newSecondAText}
	}

	if len(changes) > 0 {
		if err := insertSessionAudit(ctx, tx, id, "edited", userID, map[string]any{
			"changes": changes,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit update session: %w", err)
	}

	return r.Get(ctx, id)
}

type sessionChange struct {
	Old any `json:"old"`
	New any `json:"new"`
}

func addStringChange(changes map[string]sessionChange, field, oldValue, newValue string) {
	if strings.TrimSpace(oldValue) != strings.TrimSpace(newValue) {
		changes[field] = sessionChange{Old: oldValue, New: newValue}
	}
}

func addIDChange(
	changes map[string]sessionChange,
	field string,
	oldID int64,
	newID *int64,
	newName string,
	oldName string,
) {
	newValue := int64(0)
	if newID != nil {
		newValue = *newID
	}
	if oldID == newValue {
		return
	}
	changes[field] = sessionChange{Old: oldName, New: newName}
}

func idValue(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

// lookupName читает имя из whitelist-таблицы справочника.
func lookupName(ctx context.Context, tx pgx.Tx, table string, id int64) (string, error) {
	switch table {
	case "dir_cargo", "dir_counterparty", "dir_carrier", "dir_storage", "dir_type_of_operation",
		"dir_transport", "dir_trailer", "dir_driver":
	default:
		return "", fmt.Errorf("lookupName: unsupported table %s", table)
	}

	var name string
	if err := tx.QueryRow(ctx,
		fmt.Sprintf(`SELECT name FROM %s WHERE id = $1`, table), id,
	).Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return name, nil
}
