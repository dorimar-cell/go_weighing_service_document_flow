package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"go_weighing_service/internal/dto"
)

const profileSelect = `
SELECT scope, draft, applied, version, updated_at,
       COALESCE(updated_by_name, ''), applied_at, COALESCE(applied_by_name, '')
FROM weighing_settings_profile`

type profileScanner interface {
	Scan(dest ...any) error
}

func scanProfile(row profileScanner) (*dto.SettingsProfile, error) {
	var item dto.SettingsProfile
	var draft, applied []byte
	var appliedAt *time.Time

	if err := row.Scan(
		&item.Scope, &draft, &applied, &item.Version, &item.UpdatedAt,
		&item.UpdatedByName, &appliedAt, &item.AppliedByName,
	); err != nil {
		return nil, err
	}

	item.Draft = make(map[string]any)
	item.Applied = make(map[string]any)
	if err := json.Unmarshal(draft, &item.Draft); err != nil {
		return nil, fmt.Errorf("decode settings draft: %w", err)
	}
	if err := json.Unmarshal(applied, &item.Applied); err != nil {
		return nil, fmt.Errorf("decode applied settings: %w", err)
	}
	item.AppliedAt = appliedAt
	return &item, nil
}

// Profile возвращает профиль настроек раздела или nil, если он ещё не сохранялся.
func (r *SettingsExtraRepository) Profile(ctx context.Context, scope string) (*dto.SettingsProfile, error) {
	row := r.pool.QueryRow(ctx, profileSelect+` WHERE scope = $1`, scope)
	item, err := scanProfile(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("load settings profile %s: %w", scope, err)
	}
	return item, nil
}

// Profiles возвращает все сохранённые профили настроек.
func (r *SettingsExtraRepository) Profiles(ctx context.Context) ([]dto.SettingsProfile, error) {
	rows, err := r.pool.Query(ctx, profileSelect+` ORDER BY scope`)
	if err != nil {
		return nil, fmt.Errorf("list settings profiles: %w", err)
	}
	defer rows.Close()

	items := make([]dto.SettingsProfile, 0, 8)
	for rows.Next() {
		item, err := scanProfile(rows)
		if err != nil {
			return nil, fmt.Errorf("scan settings profile: %w", err)
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

// SaveDraft сохраняет черновик настроек раздела.
func (r *SettingsExtraRepository) SaveDraft(
	ctx context.Context,
	scope string,
	values map[string]any,
	userID int64,
) error {
	payload, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("marshal settings draft: %w", err)
	}

	if _, err := r.pool.Exec(ctx, `
		INSERT INTO weighing_settings_profile
			(scope, draft, applied, version, updated_at, updated_by_user_id, updated_by_name)
		VALUES ($1, $2::jsonb, '{}'::jsonb, 0, now(), $3, (SELECT name FROM dir_users WHERE id = $3))
		ON CONFLICT (scope) DO UPDATE SET
			draft = EXCLUDED.draft,
			updated_at = now(),
			updated_by_user_id = EXCLUDED.updated_by_user_id,
			updated_by_name = EXCLUDED.updated_by_name`,
		scope, string(payload), userID,
	); err != nil {
		return fmt.Errorf("save settings draft %s: %w", scope, err)
	}

	return nil
}

// Apply переносит черновик в applied и поднимает версию.
func (r *SettingsExtraRepository) Apply(
	ctx context.Context,
	scope string,
	userID int64,
) (*dto.SettingsProfile, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE weighing_settings_profile
		SET applied = draft,
		    version = version + 1,
		    applied_at = now(),
		    applied_by_user_id = $2,
		    applied_by_name = (SELECT name FROM dir_users WHERE id = $2),
		    updated_at = now(),
		    updated_by_user_id = $2,
		    updated_by_name = (SELECT name FROM dir_users WHERE id = $2)
		WHERE scope = $1
		RETURNING scope, draft, applied, version, updated_at,
		          COALESCE(updated_by_name, ''), applied_at, COALESCE(applied_by_name, '')`,
		scope, userID,
	)
	item, err := scanProfile(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("apply settings profile %s: %w", scope, err)
	}
	return item, nil
}

// AppliedProfiles возвращает применённые настройки выбранных разделов.
func (r *SettingsExtraRepository) AppliedProfiles(
	ctx context.Context,
	scopes []string,
) (map[string]dto.AgentSettingsScope, *time.Time, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT scope, applied, version, applied_at
		FROM weighing_settings_profile
		WHERE scope = ANY($1) AND applied <> '{}'::jsonb
		ORDER BY scope`, scopes)
	if err != nil {
		return nil, nil, fmt.Errorf("load applied settings: %w", err)
	}
	defer rows.Close()

	result := make(map[string]dto.AgentSettingsScope, len(scopes))
	var updatedAt *time.Time
	for rows.Next() {
		var scope string
		var raw []byte
		var version int
		var appliedAt *time.Time
		if err := rows.Scan(&scope, &raw, &version, &appliedAt); err != nil {
			return nil, nil, fmt.Errorf("scan applied settings: %w", err)
		}
		values := make(map[string]any)
		if err := json.Unmarshal(raw, &values); err != nil {
			return nil, nil, fmt.Errorf("decode applied settings %s: %w", scope, err)
		}
		result[scope] = dto.AgentSettingsScope{Version: version, Values: values}
		if appliedAt != nil && (updatedAt == nil || appliedAt.After(*updatedAt)) {
			updatedAt = appliedAt
		}
	}
	return result, updatedAt, rows.Err()
}
