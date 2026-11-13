package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// ReportJobRepository — очередь заданий и кэш готовых отчётов.
type ReportJobRepository struct {
	pool *pgxpool.Pool
}

func NewReportJobRepository(pool *pgxpool.Pool) *ReportJobRepository {
	return &ReportJobRepository{pool: pool}
}

// CachedReport — сохранённый результат отчёта.
type CachedReport struct {
	FileName    string
	ContentType string
	Data        []byte
}

func (r *ReportJobRepository) CacheGet(ctx context.Context, key string) (*CachedReport, error) {
	var item CachedReport
	err := r.pool.QueryRow(ctx, `
		SELECT file_name, content_type, data
		FROM report_cache
		WHERE cache_key = $1 AND expires_at > now()`, key,
	).Scan(&item.FileName, &item.ContentType, &item.Data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("report cache get: %w", err)
	}
	return &item, nil
}

func (r *ReportJobRepository) CachePut(
	ctx context.Context,
	key string,
	code string,
	format string,
	params map[string][]string,
	fileName string,
	contentType string,
	data []byte,
	ttl time.Duration,
) error {
	raw, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal report cache params: %w", err)
	}
	if ttl <= 0 {
		return nil
	}
	if _, err := r.pool.Exec(ctx, `
		INSERT INTO report_cache
			(cache_key, code, format, params, file_name, content_type, data, created_at, expires_at)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6, $7, now(), now() + $8::interval)
		ON CONFLICT (cache_key) DO UPDATE SET
			file_name = EXCLUDED.file_name,
			content_type = EXCLUDED.content_type,
			data = EXCLUDED.data,
			created_at = now(),
			expires_at = EXCLUDED.expires_at`,
		key, code, format, string(raw), fileName, contentType, data,
		fmt.Sprintf("%d seconds", int(ttl.Seconds())),
	); err != nil {
		return fmt.Errorf("report cache put: %w", err)
	}
	// Лёгкая уборка просроченных записей без отдельного планировщика.
	_, _ = r.pool.Exec(ctx, `DELETE FROM report_cache WHERE expires_at < now()`)
	return nil
}

// CreateJob создаёт задание со статусом pending.
func (r *ReportJobRepository) CreateJob(
	ctx context.Context,
	code string,
	format string,
	params map[string][]string,
	userID int64,
) (*dto.ReportJob, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal report job params: %w", err)
	}

	var id int64
	if err := r.pool.QueryRow(ctx, `
		INSERT INTO report_jobs (code, format, params, status, created_by_user_id, created_at)
		VALUES ($1, $2, $3::jsonb, 'pending', $4, now())
		RETURNING id`,
		code, format, string(raw), userID,
	).Scan(&id); err != nil {
		return nil, fmt.Errorf("create report job: %w", err)
	}
	return r.Job(ctx, id)
}

func (r *ReportJobRepository) MarkJobRunning(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE report_jobs SET status = 'running', started_at = now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("mark report job running: %w", err)
	}
	return nil
}

func (r *ReportJobRepository) CompleteJob(
	ctx context.Context,
	id int64,
	fileName string,
	contentType string,
	data []byte,
) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE report_jobs
		SET status = 'done', file_name = $2, content_type = $3, data = $4,
		    error = NULL, finished_at = now()
		WHERE id = $1`, id, fileName, contentType, data)
	if err != nil {
		return fmt.Errorf("complete report job: %w", err)
	}
	return nil
}

func (r *ReportJobRepository) FailJob(ctx context.Context, id int64, message string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE report_jobs
		SET status = 'error', error = $2, finished_at = now()
		WHERE id = $1`, id, message)
	if err != nil {
		return fmt.Errorf("fail report job: %w", err)
	}
	return nil
}

func (r *ReportJobRepository) Job(ctx context.Context, id int64) (*dto.ReportJob, error) {
	var job dto.ReportJob
	err := r.pool.QueryRow(ctx, `
		SELECT id, code, format, status, COALESCE(error, ''),
		       COALESCE(file_name, ''), COALESCE(content_type, ''),
		       created_at, started_at, finished_at
		FROM report_jobs WHERE id = $1`, id,
	).Scan(
		&job.ID, &job.Code, &job.Format, &job.Status, &job.Error,
		&job.FileName, &job.ContentType, &job.CreatedAt, &job.StartedAt, &job.FinishedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load report job: %w", err)
	}
	return &job, nil
}

func (r *ReportJobRepository) JobFile(ctx context.Context, id int64) (*CachedReport, error) {
	var item CachedReport
	err := r.pool.QueryRow(ctx, `
		SELECT file_name, content_type, data
		FROM report_jobs
		WHERE id = $1 AND status = 'done' AND data IS NOT NULL`, id,
	).Scan(&item.FileName, &item.ContentType, &item.Data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load report job file: %w", err)
	}
	return &item, nil
}

func (r *ReportJobRepository) Jobs(ctx context.Context, limit int) ([]dto.ReportJob, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, code, format, status, COALESCE(error, ''),
		       COALESCE(file_name, ''), COALESCE(content_type, ''),
		       created_at, started_at, finished_at
		FROM report_jobs
		ORDER BY created_at DESC, id DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list report jobs: %w", err)
	}
	defer rows.Close()

	items := make([]dto.ReportJob, 0, limit)
	for rows.Next() {
		var job dto.ReportJob
		if err := rows.Scan(
			&job.ID, &job.Code, &job.Format, &job.Status, &job.Error,
			&job.FileName, &job.ContentType, &job.CreatedAt, &job.StartedAt, &job.FinishedAt,
		); err != nil {
			return nil, fmt.Errorf("scan report job: %w", err)
		}
		items = append(items, job)
	}
	return items, rows.Err()
}
