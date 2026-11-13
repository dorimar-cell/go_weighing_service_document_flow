package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// LoginAuditRepository — аудит входов в портал.
type LoginAuditRepository struct {
	pool *pgxpool.Pool
}

func NewLoginAuditRepository(pool *pgxpool.Pool) *LoginAuditRepository {
	return &LoginAuditRepository{pool: pool}
}

// LoginAuditFilter — фильтры журнала входов.
type LoginAuditFilter struct {
	Result   string
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PageSize int
}

func (r *LoginAuditRepository) RecordLogin(
	ctx context.Context,
	userID int64,
	username string,
	ip string,
	userAgent string,
	result string,
	reason string,
) error {
	var userIDArg any
	if userID > 0 {
		userIDArg = userID
	}
	if _, err := r.pool.Exec(ctx, `
		INSERT INTO login_audit (user_id, username, ip, user_agent, result, reason)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userIDArg, truncate(username, 150), truncate(ip, 64),
		truncate(userAgent, 300), truncate(result, 20), truncate(reason, 80),
	); err != nil {
		return fmt.Errorf("insert login audit: %w", err)
	}
	return nil
}

func (r *LoginAuditRepository) List(
	ctx context.Context,
	filter LoginAuditFilter,
) (*dto.LoginAuditResponse, error) {
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 25
	}

	where := make([]string, 0, 4)
	args := make([]any, 0, 6)

	if filter.Result == "success" || filter.Result == "failed" {
		args = append(args, filter.Result)
		where = append(where, fmt.Sprintf("result = $%d", len(args)))
	}
	if search := strings.TrimSpace(filter.Search); search != "" {
		args = append(args, "%"+search+"%")
		where = append(where, fmt.Sprintf(
			"(username ILIKE $%d OR ip ILIKE $%d)", len(args), len(args),
		))
	}
	if filter.DateFrom != nil {
		args = append(args, *filter.DateFrom)
		where = append(where, fmt.Sprintf("created_at >= $%d", len(args)))
	}
	if filter.DateTo != nil {
		args = append(args, *filter.DateTo)
		where = append(where, fmt.Sprintf("created_at < $%d", len(args)))
	}

	condition := ""
	if len(where) > 0 {
		condition = " WHERE " + strings.Join(where, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM login_audit`+condition, args...).
		Scan(&total); err != nil {
		return nil, fmt.Errorf("count login audit: %w", err)
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(user_id, 0), COALESCE(username, ''), COALESCE(ip, ''),
		       COALESCE(user_agent, ''), result, COALESCE(reason, ''), created_at
		FROM login_audit`+condition+fmt.Sprintf(`
		ORDER BY created_at DESC, id DESC
		LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, fmt.Errorf("list login audit: %w", err)
	}
	defer rows.Close()

	items := make([]dto.LoginAuditItem, 0, pageSize)
	for rows.Next() {
		var item dto.LoginAuditItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Username, &item.IP,
			&item.UserAgent, &item.Result, &item.Reason, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan login audit: %w", err)
		}
		items = append(items, item)
	}

	return &dto.LoginAuditResponse{
		Items: items, Page: page, PageSize: pageSize, Total: total,
	}, rows.Err()
}

func truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
