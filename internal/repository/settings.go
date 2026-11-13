package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// SettingsRepository — бизнес-настройки полей рейса.
type SettingsRepository struct {
	pool *pgxpool.Pool
}

func NewSettingsRepository(pool *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{pool: pool}
}

// App возвращает настройки приложения; при отсутствии строки — значения по умолчанию.
func (r *SettingsRepository) App(ctx context.Context) (*dto.AppSettings, error) {
	var settings dto.AppSettings

	err := r.pool.QueryRow(ctx, `
		SELECT require_driver, require_trailer, require_cargo, require_counterparty,
		       require_carrier, require_storage, need_operation_type,
		       visible_driver, visible_trailer, visible_cargo, visible_counterparty,
		       visible_carrier, visible_storage, visible_operation_type, visible_pass
		FROM settings_app
		WHERE id = 1`,
	).Scan(
		&settings.RequireDriver, &settings.RequireTrailer, &settings.RequireCargo,
		&settings.RequireCounterparty, &settings.RequireCarrier, &settings.RequireStorage,
		&settings.RequireOperationType,
		&settings.VisibleDriver, &settings.VisibleTrailer, &settings.VisibleCargo,
		&settings.VisibleCounterparty, &settings.VisibleCarrier, &settings.VisibleStorage,
		&settings.VisibleOperationType, &settings.VisiblePass,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.AppSettings{
				VisibleDriver: true, VisibleTrailer: true, VisibleCargo: true,
				VisibleCounterparty: true, VisibleCarrier: true, VisibleStorage: true,
				VisibleOperationType: true, VisiblePass: true,
			}, nil
		}
		return nil, fmt.Errorf("load app settings: %w", err)
	}

	return &settings, nil
}

// SaveApp сохраняет настройки приложения.
func (r *SettingsRepository) SaveApp(
	ctx context.Context,
	settings dto.AppSettings,
	userID int64,
) error {
	if _, err := r.pool.Exec(ctx, `
		INSERT INTO settings_app (
			id, ts_number_mask_use, ts_number_mask, need_operation_type, need_fund,
			need_required_fields, placement_visible, ves_doc_visible, axle_load_visible,
			need_cargo_type, visible_cargo_type, need_placement, visible_fund,
			visible_placement, conditional_zero, use_com_object, use_region,
			union_recognize_skud,
			require_driver, require_trailer, require_cargo, require_counterparty,
			require_carrier, require_storage,
			visible_driver, visible_trailer, visible_cargo, visible_counterparty,
			visible_carrier, visible_storage, visible_operation_type, visible_pass,
			modified_by_user_id, modified_date
		) VALUES (
			1, false, '', $1, false,
			true, true, true, false,
			false, true, false, true,
			true, 0, false, false,
			false,
			$2, $3, $4, $5,
			$6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15,
			$16, now()
		)
		ON CONFLICT (id) DO UPDATE SET
			need_operation_type = EXCLUDED.need_operation_type,
			require_driver = EXCLUDED.require_driver,
			require_trailer = EXCLUDED.require_trailer,
			require_cargo = EXCLUDED.require_cargo,
			require_counterparty = EXCLUDED.require_counterparty,
			require_carrier = EXCLUDED.require_carrier,
			require_storage = EXCLUDED.require_storage,
			visible_driver = EXCLUDED.visible_driver,
			visible_trailer = EXCLUDED.visible_trailer,
			visible_cargo = EXCLUDED.visible_cargo,
			visible_counterparty = EXCLUDED.visible_counterparty,
			visible_carrier = EXCLUDED.visible_carrier,
			visible_storage = EXCLUDED.visible_storage,
			visible_operation_type = EXCLUDED.visible_operation_type,
			visible_pass = EXCLUDED.visible_pass,
			modified_by_user_id = EXCLUDED.modified_by_user_id,
			modified_date = now()`,
		settings.RequireOperationType,
		settings.RequireDriver, settings.RequireTrailer, settings.RequireCargo,
		settings.RequireCounterparty, settings.RequireCarrier, settings.RequireStorage,
		settings.VisibleDriver, settings.VisibleTrailer, settings.VisibleCargo,
		settings.VisibleCounterparty, settings.VisibleCarrier, settings.VisibleStorage,
		settings.VisibleOperationType, settings.VisiblePass,
		userID,
	); err != nil {
		return fmt.Errorf("save app settings: %w", err)
	}

	return nil
}
