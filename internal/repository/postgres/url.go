package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"url-sentinel/internal/domain/entity"
	"url-sentinel/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type urlRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new PostgreSQL URL repository
func NewURLRepository(db *sql.DB) repository.URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(ctx context.Context, url *entity.URL) error {
	query := `
		INSERT INTO urls (id, address, check_interval, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		url.ID,
		url.Address,
		url.CheckInterval,
		url.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return repository.ErrURLAddressExists
		}
		return fmt.Errorf("failed to create url: %w", err)
	}

	return nil
}

func (r *urlRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.URL, error) {
	query := `
		SELECT id, address,
			EXTRACT(EPOCH FROM check_interval)::BIGINT * 1000000000 AS check_interval_ns,
			created_at
		FROM urls
		WHERE id = $1
	`

	var url entity.URL
	var intervalNs int64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&url.ID,
		&url.Address,
		&intervalNs,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to get url by id: %w", err)
	}

	url.CheckInterval = time.Duration(intervalNs)

	return &url, nil
}

func (r *urlRepository) List(ctx context.Context) ([]*entity.URL, error) {
	query := `
		SELECT id, address,
			EXTRACT(EPOCH FROM check_interval)::BIGINT * 1000000000 AS check_interval_ns,
			created_at
		FROM urls
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list urls: %w", err)
	}
	defer rows.Close()

	var urls []*entity.URL
	for rows.Next() {
		var url entity.URL
		var intervalNs int64

		if err := rows.Scan(
			&url.ID,
			&url.Address,
			&intervalNs,
			&url.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan url: %w", err)
		}

		url.CheckInterval = time.Duration(intervalNs)
		urls = append(urls, &url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return urls, nil
}

func (r *urlRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM urls WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}

func (r *urlRepository) ExistsByAddress(ctx context.Context, address string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE address = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, address).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check url existence: %w", err)
	}

	return exists, nil
}
