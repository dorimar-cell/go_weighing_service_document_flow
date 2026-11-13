package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

const passFrom = `
FROM reg_weight rw
JOIN dir_transport t ON t.id = rw.id_dir_transport
LEFT JOIN dir_trailer tr ON tr.id = rw.id_dir_trailer
LEFT JOIN dir_driver d ON d.id = rw.id_dir_driver
LEFT JOIN dir_cargo cg ON cg.id = rw.id_dir_cargo
LEFT JOIN dir_counterparty cp ON cp.id = rw.id_dir_counterparty
LEFT JOIN dir_carrier cr ON cr.id = rw.id_dir_carrier
JOIN system_dir_type_of_pass pt ON pt.id = rw.id_system_dir_type_of_pass
JOIN system_dir_identifier_state st ON st.id = rw.id_system_dir_identifier_state`

const passSelect = `
SELECT rw.id,
       rw.name,
       t.name,
       COALESCE(tr.name, ''),
       COALESCE(d.name, ''),
       COALESCE(cg.id, 0),
       COALESCE(cg.name, ''),
       COALESCE(cp.id, 0),
       COALESCE(cp.name, ''),
       COALESCE(cr.id, 0),
       COALESCE(cr.name, ''),
       rw.id_system_dir_type_of_pass,
       pt.name,
       rw.id_system_dir_identifier_state,
       st.name,
       rw.pass_enabled,
       rw.date_pass_expiry,
       rw.created_date` + passFrom

type PassRepository struct {
	pool *pgxpool.Pool
}

func NewPassRepository(pool *pgxpool.Pool) *PassRepository {
	return &PassRepository{pool: pool}
}

func (r *PassRepository) Create(
	ctx context.Context,
	request dto.CreatePassRequest,
	userID int64,
) (*dto.PassDTO, error) {
	passTypeID := request.PassTypeID
	if passTypeID == 0 {
		passTypeID = 1
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create pass: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	transportID, err := findOrCreateNamed(ctx, tx, "dir_transport", request.TransportNumber)
	if err != nil {
		return nil, err
	}

	var trailerID *int64
	if strings.TrimSpace(request.TrailerNumber) != "" {
		id, err := findOrCreateNamed(ctx, tx, "dir_trailer", request.TrailerNumber)
		if err != nil {
			return nil, err
		}
		trailerID = &id
	}

	var driverID *int64
	if strings.TrimSpace(request.DriverName) != "" {
		id, err := findOrCreateNamed(ctx, tx, "dir_driver", request.DriverName)
		if err != nil {
			return nil, err
		}
		driverID = &id
	}

	var passID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO reg_weight
			(name, pass_enabled, auto_save_tara,
			 id_system_dir_type_of_regweight, id_system_dir_identifier_state,
			 id_system_dir_type_of_pass, id_dir_transport, id_dir_trailer, id_dir_driver,
			 id_dir_cargo, id_dir_counterparty, id_dir_carrier,
			 date_pass_expiry, deleted, created_by_user_id, created_date)
		VALUES ($1, true, false, 1, 1, $2, $3, $4, $5, $6, $7, $8, $9, false, $10, now())
		RETURNING id`,
		request.Number,
		passTypeID,
		transportID,
		trailerID,
		driverID,
		request.CargoID,
		request.CounterpartyID,
		request.CarrierID,
		request.ValidUntil,
		userID,
	).Scan(&passID); err != nil {
		return nil, fmt.Errorf("insert pass: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create pass: %w", err)
	}

	return r.Get(ctx, passID)
}

func (r *PassRepository) List(
	ctx context.Context,
	status string,
	search string,
	page int,
	pageSize int,
) (*dto.PassListResponse, error) {
	where := []string{"rw.deleted = false"}
	args := make([]any, 0, 3)

	if status == "active" {
		where = append(where, "rw.pass_enabled = true")
	}

	if term := strings.TrimSpace(search); term != "" {
		args = append(args, "%"+term+"%")
		where = append(where, fmt.Sprintf("(rw.name ILIKE $%d OR t.name ILIKE $%d)", len(args), len(args)))
	}

	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) "+passFrom+whereSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count passes: %w", err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	query := passSelect + whereSQL +
		fmt.Sprintf(" ORDER BY rw.created_date DESC LIMIT $%d OFFSET $%d", limitArg, offsetArg)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list passes: %w", err)
	}
	defer rows.Close()

	items := make([]dto.PassDTO, 0, pageSize)
	for rows.Next() {
		item, err := scanPass(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate passes: %w", err)
	}

	return &dto.PassListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (r *PassRepository) Get(ctx context.Context, id int64) (*dto.PassDTO, error) {
	row := r.pool.QueryRow(ctx, passSelect+" WHERE rw.id = $1", id)
	item, err := scanPass(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (r *PassRepository) SetEnabled(
	ctx context.Context,
	id int64,
	enabled bool,
	userID int64,
) (*dto.PassDTO, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE reg_weight
		SET pass_enabled = $2,
		    modified_by_user_id = $3,
		    modified_date = now()
		WHERE id = $1 AND deleted = false`,
		id, enabled, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("set pass enabled: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	return r.Get(ctx, id)
}

func scanPass(row scanner) (*dto.PassDTO, error) {
	var item dto.PassDTO
	var cargoID, counterpartyID, carrierID int64
	var validUntil sql.NullTime

	if err := row.Scan(
		&item.ID,
		&item.Number,
		&item.TransportNumber,
		&item.TrailerNumber,
		&item.DriverName,
		&cargoID,
		&item.CargoName,
		&counterpartyID,
		&item.CounterpartyName,
		&carrierID,
		&item.CarrierName,
		&item.PassTypeID,
		&item.PassTypeName,
		&item.StateID,
		&item.StateName,
		&item.PassEnabled,
		&validUntil,
		&item.CreatedAt,
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
	if validUntil.Valid {
		value := validUntil.Time
		item.ValidUntil = &value
	}

	return &item, nil
}
