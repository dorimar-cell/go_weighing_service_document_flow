package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// SettingsExtraRepository — камеры, лицензии и аудит настроек.
type SettingsExtraRepository struct {
	pool *pgxpool.Pool
}

func NewSettingsExtraRepository(pool *pgxpool.Pool) *SettingsExtraRepository {
	return &SettingsExtraRepository{pool: pool}
}

func (r *SettingsExtraRepository) Cameras(ctx context.Context) ([]dto.CameraDTO, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.enable, COALESCE(c.name, ''), COALESCE(c.ip, ''), c.port,
		       COALESCE(c.login, ''), c.wait_time, c.try_times,
		       c.stream_type_mini_window, c.stream_type_full_window,
		       c.platform_number, COALESCE(wr.name, ''), c.updated_date
		FROM settings_video_surveillance c
		LEFT JOIN dir_weight_room wr ON wr.id = c.id_dir_weight_room
		ORDER BY c.id`)
	if err != nil {
		return nil, fmt.Errorf("list cameras: %w", err)
	}
	defer rows.Close()

	items := make([]dto.CameraDTO, 0, 8)
	for rows.Next() {
		var item dto.CameraDTO
		if err := rows.Scan(
			&item.ID, &item.Enable, &item.Name, &item.IP, &item.Port, &item.Login,
			&item.WaitTime, &item.TryTimes, &item.StreamTypeMini, &item.StreamTypeFull,
			&item.PlatformNumber, &item.WeightRoomName, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan camera: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SettingsExtraRepository) SaveCamera(
	ctx context.Context,
	id int64,
	req dto.SaveCameraRequest,
	userID int64,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE settings_video_surveillance
		SET enable = $2, name = $3, ip = $4, port = $5, login = $6,
		    wait_time = $7, try_times = $8, stream_type_mini_window = $9,
		    stream_type_full_window = $10, platform_number = $11,
		    updated_date = now(), updated_by_user_id = $12
		WHERE id = $1`,
		id, req.Enable, req.Name, req.IP, req.Port, req.Login,
		req.WaitTime, req.TryTimes, req.StreamTypeMini, req.StreamTypeFull,
		req.PlatformNumber, userID,
	)
	if err != nil {
		return fmt.Errorf("save camera: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SettingsExtraRepository) Licenses(ctx context.Context) ([]dto.LicenseDTO, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, module_code, to_char(valid_until, 'YYYY-MM-DD'),
		       seats, active, COALESCE(notes, ''), COALESCE(modified_date, created_date)
		FROM weighing_license
		ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list licenses: %w", err)
	}
	defer rows.Close()

	items := make([]dto.LicenseDTO, 0, 16)
	for rows.Next() {
		var item dto.LicenseDTO
		var validUntil *string
		if err := rows.Scan(
			&item.ID, &item.Name, &item.ModuleCode, &validUntil,
			&item.Seats, &item.Active, &item.Notes, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan license: %w", err)
		}
		item.ValidUntil = validUntil
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SettingsExtraRepository) CreateLicense(
	ctx context.Context,
	req dto.LicenseRequest,
	userID int64,
) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO weighing_license
			(name, module_code, valid_until, seats, active, notes, created_by_user_id, created_date)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7, now())`,
		req.Name, req.ModuleCode, req.ValidUntil, req.Seats, req.Active, req.Notes, userID,
	)
	if err != nil {
		return fmt.Errorf("create license: %w", err)
	}
	return nil
}

func (r *SettingsExtraRepository) UpdateLicense(
	ctx context.Context,
	id int64,
	req dto.LicenseRequest,
	userID int64,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE weighing_license
		SET name = $2, module_code = $3, valid_until = $4::date, seats = $5,
		    active = $6, notes = $7, modified_by_user_id = $8, modified_date = now()
		WHERE id = $1`,
		id, req.Name, req.ModuleCode, req.ValidUntil, req.Seats, req.Active, req.Notes, userID,
	)
	if err != nil {
		return fmt.Errorf("update license: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SettingsExtraRepository) DeleteLicense(ctx context.Context, id int64) error {
	if _, err := r.pool.Exec(ctx, `DELETE FROM weighing_license WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete license: %w", err)
	}
	return nil
}

func (r *SettingsExtraRepository) SettingsAudit(
	ctx context.Context,
	page int,
	pageSize int,
) (*dto.SettingsAuditResponse, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM settings_audit`).Scan(&total); err != nil {
		return nil, fmt.Errorf("count settings audit: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, scope, COALESCE(entity_id, 0), action, user_id, COALESCE(user_name, ''),
		       COALESCE(payload::text, '{}'), created_at
		FROM settings_audit
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, fmt.Errorf("list settings audit: %w", err)
	}
	defer rows.Close()

	items := make([]dto.SettingsAuditItem, 0, pageSize)
	for rows.Next() {
		var item dto.SettingsAuditItem
		var payload string
		if err := rows.Scan(
			&item.ID, &item.Scope, &item.EntityID, &item.Action, &item.UserID,
			&item.UserName, &payload, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan settings audit: %w", err)
		}
		item.Payload = json.RawMessage(payload)
		items = append(items, item)
	}

	return &dto.SettingsAuditResponse{
		Items: items, Page: page, PageSize: pageSize, Total: total,
	}, rows.Err()
}

// LogSettings пишет аудит изменения настроек.
func (r *SettingsExtraRepository) LogSettings(
	ctx context.Context,
	scope string,
	entityID int64,
	action string,
	userID int64,
	payload any,
) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal settings audit: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO settings_audit (scope, entity_id, action, user_id, user_name, payload)
		VALUES ($1, $2, $3, $4, (SELECT name FROM dir_users WHERE id = $4), $5::jsonb)`,
		scope, entityID, action, userID, string(data),
	)
	if err != nil {
		return fmt.Errorf("insert settings audit: %w", err)
	}
	return nil
}
