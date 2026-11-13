package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// RBACRepository — группы пользователей и доменные права весовой.
type RBACRepository struct {
	pool *pgxpool.Pool
}

func NewRBACRepository(pool *pgxpool.Pool) *RBACRepository {
	return &RBACRepository{pool: pool}
}

// Matrix возвращает права W* и группы с их назначением.
func (r *RBACRepository) Matrix(ctx context.Context) (*dto.RBACResponse, error) {
	permissionRows, err := r.pool.Query(ctx, `
		SELECT control_tag, name
		FROM system_dir_access_permissions
		WHERE control_tag LIKE 'W%'
		ORDER BY control_tag`)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	permissions := make([]dto.Permission, 0, 16)
	for permissionRows.Next() {
		var item dto.Permission
		if err := permissionRows.Scan(&item.Code, &item.Name); err != nil {
			permissionRows.Close()
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, item)
	}
	permissionRows.Close()
	if err := permissionRows.Err(); err != nil {
		return nil, err
	}

	groupRows, err := r.pool.Query(ctx, `
		SELECT g.id, g.name, COALESCE(g.system_record, false),
		       COALESCE(array_agg(p.control_tag) FILTER (WHERE p.control_tag IS NOT NULL), '{}')
		FROM dir_user_group g
		LEFT JOIN com_dir_user_groups_access_permissions ugap
		       ON ugap.id_dir_user_group = g.id
		LEFT JOIN system_dir_access_permissions p
		       ON p.id = ugap.id_system_dir_access_permission AND p.control_tag LIKE 'W%'
		GROUP BY g.id, g.name, g.system_record
		ORDER BY g.id`)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	defer groupRows.Close()

	groups := make([]dto.GroupPermissions, 0, 8)
	for groupRows.Next() {
		var item dto.GroupPermissions
		if err := groupRows.Scan(&item.ID, &item.Name, &item.System, &item.Permissions); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		groups = append(groups, item)
	}

	return &dto.RBACResponse{Permissions: permissions, Groups: groups}, groupRows.Err()
}

// SaveGroupPermissions заменяет набор прав группы.
func (r *RBACRepository) SaveGroupPermissions(
	ctx context.Context,
	groupID int64,
	codes []string,
	userID int64,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin rbac save: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var exists bool
	if err := tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM dir_user_group WHERE id = $1)`, groupID,
	).Scan(&exists); err != nil {
		return fmt.Errorf("check group: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM com_dir_user_groups_access_permissions ugap
		USING system_dir_access_permissions p
		WHERE ugap.id_dir_user_group = $1
		  AND p.id = ugap.id_system_dir_access_permission
		  AND p.control_tag LIKE 'W%'`, groupID,
	); err != nil {
		return fmt.Errorf("clear group permissions: %w", err)
	}

	if len(codes) > 0 {
		if _, err := tx.Exec(ctx, `
			INSERT INTO com_dir_user_groups_access_permissions
				(id_dir_user_group, id_system_dir_access_permission)
			SELECT $1, id FROM system_dir_access_permissions
			WHERE control_tag = ANY($2) AND control_tag LIKE 'W%'`,
			groupID, codes,
		); err != nil {
			return fmt.Errorf("insert group permissions: %w", err)
		}
	}

	_ = userID
	return tx.Commit(ctx)
}
