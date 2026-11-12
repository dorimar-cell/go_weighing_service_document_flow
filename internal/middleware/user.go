package middleware

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userContextKey string

const currentUserKey userContextKey = "weighing_user"

// User — пользователь весовой, связанный с пользователем портала.
type User struct {
	ID          int64
	PortalID    int64
	Username    string
	GroupID     int64
	IsRoot      bool
	Permissions map[string]bool
}

func (u User) Can(permission string) bool {
	return u.IsRoot || u.Permissions[permission]
}

func GetUser(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(currentUserKey).(User)
	return user, ok
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

// UserStore находит или создаёт пользователя весовой по id портала.
type UserStore struct {
	pool  *pgxpool.Pool
	cache sync.Map
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

// DevUser — сервисный администратор для AUTH_MODE=dev.
func (s *UserStore) DevUser(ctx context.Context) (User, error) {
	return s.load(ctx, 1)
}

// Resolve находит пользователя по portal_user_id или создаёт нового оператора.
func (s *UserStore) Resolve(ctx context.Context, portalID int64, username string) (User, error) {
	if cached, ok := s.cache.Load(portalID); ok {
		return cached.(User), nil
	}

	user, err := s.findByPortalID(ctx, portalID)
	if errors.Is(err, pgx.ErrNoRows) {
		user, err = s.create(ctx, portalID, username)
	}
	if err != nil {
		return User{}, err
	}

	s.cache.Store(portalID, user)
	return user, nil
}

func (s *UserStore) load(ctx context.Context, id int64) (User, error) {
	var user User
	var groupID int64
	var username *string

	err := s.pool.QueryRow(ctx, `
		SELECT id, COALESCE(id_dir_user_groups, 0), portal_username
		FROM dir_users
		WHERE id = $1 AND deleted = false`, id,
	).Scan(&user.ID, &groupID, &username)
	if err != nil {
		return User{}, fmt.Errorf("load user %d: %w", id, err)
	}

	if username != nil {
		user.Username = *username
	}
	user.PortalID = id
	user.GroupID = groupID
	user.IsRoot = groupID == 0 || groupID == 1

	permissions, err := s.loadPermissions(ctx, groupID)
	if err != nil {
		return User{}, err
	}
	user.Permissions = permissions

	s.cache.Store(id, user)
	return user, nil
}

func (s *UserStore) findByPortalID(ctx context.Context, portalID int64) (User, error) {
	var user User
	var groupID int64
	var username *string

	err := s.pool.QueryRow(ctx, `
		SELECT id, COALESCE(id_dir_user_groups, 0), COALESCE(portal_username, name)
		FROM dir_users
		WHERE portal_user_id = $1 AND deleted = false`, portalID,
	).Scan(&user.ID, &groupID, &username)
	if err != nil {
		return User{}, err
	}

	if username != nil {
		user.Username = *username
	}
	user.PortalID = portalID
	user.GroupID = groupID
	user.IsRoot = groupID == 0 || groupID == 1

	permissions, err := s.loadPermissions(ctx, groupID)
	if err != nil {
		return User{}, err
	}
	user.Permissions = permissions

	return user, nil
}

func (s *UserStore) create(ctx context.Context, portalID int64, username string) (User, error) {
	if username == "" {
		username = fmt.Sprintf("portal-%d", portalID)
	}

	var userID int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO dir_users
			(name, portal_user_id, portal_username, deleted,
			 id_dir_weight_room, id_dir_user_groups,
			 created_by_user_id, created_date, ref)
		VALUES ($1, $2, $1, false, 1, 2, 1, now(), gen_random_uuid())
		RETURNING id`,
		username, portalID,
	).Scan(&userID)
	if err != nil {
		return User{}, fmt.Errorf("create portal user: %w", err)
	}

	return s.load(ctx, userID)
}

func (s *UserStore) loadPermissions(ctx context.Context, groupID int64) (map[string]bool, error) {
	permissions := make(map[string]bool)
	if groupID == 0 || groupID == 1 {
		return permissions, nil
	}

	rows, err := s.pool.Query(ctx, `
		SELECT p.control_tag
		FROM com_dir_user_groups_access_permissions ugap
		JOIN system_dir_access_permissions p ON p.id = ugap.id_system_dir_access_permission
		WHERE ugap.id_dir_user_group = $1`, groupID)
	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		permissions[tag] = true
	}

	return permissions, rows.Err()
}
