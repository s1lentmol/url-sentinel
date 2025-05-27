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
)

type checkRepository struct {
	db *sql.DB
}

// NewCheckRepository creates a new PostgreSQL check repository
func NewCheckRepository(db *sql.DB) repository.CheckRepository {
	return &checkRepository{db: db}
}

func (r *checkRepository) Create(ctx context.Context, check *entity.Check) error {
	query := `
		INSERT INTO checks (id, url_id, status, code, duration, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		check.ID,
		check.URLID,
		check.Status,
		check.Code,
		check.Duration,
		check.CheckedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create check: %w", err)
	}

	return nil
}

func (r *checkRepository) ListByURLID(ctx context.Context, urlID uuid.UUID) ([]*entity.Check, error) {
	query := `
		SELECT id, url_id, status, code,
			EXTRACT(EPOCH FROM duration)::BIGINT * 1000000000 AS duration_ns,
			checked_at
		FROM checks
		WHERE url_id = $1
		ORDER BY checked_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, urlID)
	if err != nil {
		return nil, fmt.Errorf("failed to list checks: %w", err)
	}
	defer rows.Close()

	var checks []*entity.Check
	for rows.Next() {
		var check entity.Check
		var durationNs int64

		if err := rows.Scan(
			&check.ID,
			&check.URLID,
			&check.Status,
			&check.Code,
			&durationNs,
			&check.CheckedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan check: %w", err)
		}

		check.Duration = time.Duration(durationNs)
		checks = append(checks, &check)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return checks, nil
}

func (r *checkRepository) GetLatestByURLID(ctx context.Context, urlID uuid.UUID) (*entity.Check, error) {
	query := `
		SELECT id, url_id, status, code,
			EXTRACT(EPOCH FROM duration)::BIGINT * 1000000000 AS duration_ns,
			checked_at
		FROM checks
		WHERE url_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`

	var check entity.Check
	var durationNs int64

	err := r.db.QueryRowContext(ctx, query, urlID).Scan(
		&check.ID,
		&check.URLID,
		&check.Status,
		&check.Code,
		&durationNs,
		&check.CheckedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No checks yet
		}
		return nil, fmt.Errorf("failed to get latest check: %w", err)
	}

	check.Duration = time.Duration(durationNs)

	return &check, nil
}
