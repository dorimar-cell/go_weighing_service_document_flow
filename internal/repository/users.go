package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// UsersRepository — пользователи портала и их привязка к весовым.
type UsersRepository struct {
	pool *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{pool: pool}
}

func (r *UsersRepository) List(
	ctx context.Context,
	search string,
	page int,
	pageSize int,
) (*dto.UsersResponse, error) {
	where := []string{"1 = 1"}
	args := make([]any, 0, 2)

	if term := strings.TrimSpace(search); term != "" {
		args = append(args, "%"+term+"%")
		where = append(where, fmt.Sprintf(
			"(lower(u.name) LIKE lower($%d) OR lower(COALESCE(u.portal_username, '')) LIKE lower($%d))",
			len(args), len(args),
		))
	}
	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.pool.QueryRow(ctx,
		"SELECT count(*) FROM dir_users u "+whereSQL, args...,
	).Scan(&total); err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.name, COALESCE(u.portal_username, ''),
		       COALESCE(g.name, ''), COALESCE(wr.id, 0), COALESCE(wr.name, ''), u.deleted
		FROM dir_users u
		LEFT JOIN dir_user_group g ON g.id = u.id_dir_user_groups
		LEFT JOIN dir_weight_room wr ON wr.id = u.id_dir_weight_room
		`+whereSQL+fmt.Sprintf(" ORDER BY u.name LIMIT $%d OFFSET $%d", limitArg, offsetArg),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	items := make([]dto.UserDTO, 0, pageSize)
	for rows.Next() {
		var item dto.UserDTO
		if err := rows.Scan(
			&item.ID, &item.Name, &item.PortalUsername, &item.GroupName,
			&item.WeightRoomID, &item.WeightRoomName, &item.Deleted,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		items = append(items, item)
	}

	return &dto.UsersResponse{Items: items, Page: page, PageSize: pageSize, Total: total}, rows.Err()
}

func (r *UsersRepository) WeightRooms(ctx context.Context) ([]dto.WeightRoomDTO, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name FROM dir_weight_room WHERE deleted_at IS NULL ORDER BY lower(name)`)
	if err != nil {
		return nil, fmt.Errorf("list weight rooms: %w", err)
	}
	defer rows.Close()

	items := make([]dto.WeightRoomDTO, 0, 8)
	for rows.Next() {
		var item dto.WeightRoomDTO
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, fmt.Errorf("scan weight room: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *UsersRepository) SetWeightRoom(ctx context.Context, userID, weightRoomID int64) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE dir_users
		SET id_dir_weight_room = $2, modified_date = now()
		WHERE id = $1`, userID, weightRoomID)
	if err != nil {
		return fmt.Errorf("set user weight room: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// JournalPreferences возвращает сохранённые фильтры и колонки пользователя.
func (r *UsersRepository) JournalPreferences(
	ctx context.Context,
	userID int64,
) (*dto.JournalPreferences, error) {
	var filters []byte
	var columns []byte
	var updatedAt *string

	err := r.pool.QueryRow(ctx, `
		SELECT filters, columns, to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS')
		FROM user_journal_preferences
		WHERE user_id = $1`, userID,
	).Scan(&filters, &columns, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &dto.JournalPreferences{Filters: map[string]any{}, Columns: []string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load journal preferences: %w", err)
	}

	prefs := &dto.JournalPreferences{}
	_ = json.Unmarshal(filters, &prefs.Filters)
	_ = json.Unmarshal(columns, &prefs.Columns)
	prefs.UpdatedAt = updatedAt
	if prefs.Filters == nil {
		prefs.Filters = map[string]any{}
	}
	if prefs.Columns == nil {
		prefs.Columns = []string{}
	}
	return prefs, nil
}

// SaveJournalPreferences сохраняет фильтры и колонки пользователя.
func (r *UsersRepository) SaveJournalPreferences(
	ctx context.Context,
	userID int64,
	prefs dto.JournalPreferences,
) error {
	filters, err := json.Marshal(prefs.Filters)
	if err != nil {
		return fmt.Errorf("marshal filters: %w", err)
	}
	columns, err := json.Marshal(prefs.Columns)
	if err != nil {
		return fmt.Errorf("marshal columns: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO user_journal_preferences (user_id, filters, columns, updated_at)
		VALUES ($1, $2::jsonb, $3::jsonb, now())
		ON CONFLICT (user_id) DO UPDATE SET
			filters = EXCLUDED.filters,
			columns = EXCLUDED.columns,
			updated_at = now()`,
		userID, string(filters), string(columns),
	)
	if err != nil {
		return fmt.Errorf("save journal preferences: %w", err)
	}
	return nil
}
