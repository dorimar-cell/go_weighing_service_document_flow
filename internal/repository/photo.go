package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

type PhotoRepository struct {
	pool *pgxpool.Pool
}

func NewPhotoRepository(pool *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{pool: pool}
}

func (r *PhotoRepository) Add(
	ctx context.Context,
	sessionID int64,
	fileName string,
	firstWeight bool,
	userID int64,
) (*dto.PhotoDTO, error) {
	var photo dto.PhotoDTO
	photo.FileName = fileName
	photo.FirstWeight = firstWeight

	err := r.pool.QueryRow(ctx, `
		INSERT INTO weighting_photo
			(id_weighing, file_name, created_by_user_id, created_date, first_weight, id_settings_video_surveillance)
		VALUES ($1, $2, $3, now(), $4, 1)
		RETURNING id, created_date`,
		sessionID, fileName, userID, firstWeight,
	).Scan(&photo.ID, &photo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert photo: %w", err)
	}

	return &photo, nil
}

func (r *PhotoRepository) List(ctx context.Context, sessionID int64) ([]dto.PhotoDTO, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, file_name, first_weight, created_date
		FROM weighting_photo
		WHERE id_weighing = $1
		ORDER BY first_weight, created_date`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}
	defer rows.Close()

	items := make([]dto.PhotoDTO, 0, 4)
	for rows.Next() {
		var item dto.PhotoDTO
		if err := rows.Scan(&item.ID, &item.FileName, &item.FirstWeight, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PhotoRepository) Get(ctx context.Context, sessionID int64, photoID int64) (*dto.PhotoDTO, error) {
	var item dto.PhotoDTO
	err := r.pool.QueryRow(ctx, `
		SELECT id, file_name, first_weight, created_date
		FROM weighting_photo
		WHERE id = $1 AND id_weighing = $2`,
		photoID, sessionID,
	).Scan(&item.ID, &item.FileName, &item.FirstWeight, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get photo: %w", err)
	}

	return &item, nil
}
