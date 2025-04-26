package storage

import (
	"database/sql"
	"fmt"

	"url-sentinel/internal/model"

	"github.com/google/uuid"
)

// CheckRepository implements storage of check results
type CheckRepository struct {
	db *sql.DB
}

// NewCheckRepository creates a new check repository
func NewCheckRepository(db *sql.DB) *CheckRepository {
	return &CheckRepository{db: db}
}

// Create new test results
func (r *CheckRepository) SaveCheck(c *model.Check) error {
	const op = "storage.SaveCheck"
	query := `INSERT INTO checks (id, url_id, status, code, duration, checked_at) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	_, err = stmt.Exec(
		c.ID,
		c.URLID,
		c.Status,
		c.Code,
		c.Duration.String(),
		c.CheckedAt,
	)
	
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// ListByURL returns all checks for a given URL
func (r *CheckRepository) ListOfChecksByURL(urlID uuid.UUID) ([]model.Check, error) {
	const op = "storage.ListOfChecksByURL"
	query := `SELECT id, url_id, status, code,
		(EXTRACT(EPOCH FROM duration) * 1000000000)::BIGINT AS duration,
	 	checked_at FROM checks WHERE url_id = $1 ORDER BY checked_at DESC`
			
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	rows, err := stmt.Query(urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []model.Check
	for rows.Next() {
		c := model.Check{}
		if err := rows.Scan(
			&c.ID,
			&c.URLID,
			&c.Status,
			&c.Code,
			&c.Duration,
			&c.CheckedAt,
		); err != nil {
			return nil, err
		}
		checks = append(checks, c)
	}
	return checks, rows.Err()
}
