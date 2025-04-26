package storage

import (
	"database/sql"
	"fmt"

	"url-sentinel/internal/model"

	"github.com/google/uuid"
)

// URLRepository implements storage of URLs
type URLRepository struct {
	db *sql.DB
}

// creates a new repository URL
func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

// insert new URL to database
func (r *URLRepository) SaveURL(u *model.URL) error {
	const op = "storage.SaveURL"
	query := `INSERT INTO urls (id, address, check_interval, created_at) VALUES ($1, $2, $3, $4)`

	stmt, err := r.db.Prepare(query)

	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	_, err = stmt.Exec(
		u.ID,
		u.Address,
		u.CheckInterval.String(),
		u.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// return URL on ID
func (r *URLRepository) GetURLByID(id uuid.UUID) (*model.URL, error) {
	const op = "storage.GetURLByID"
	query := `SELECT id, address,
		(EXTRACT(EPOCH FROM check_interval) * 1000000000)::BIGINT AS check_interval,
	 	created_at FROM urls WHERE id = $1`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u := new(model.URL)
	err = stmt.QueryRow(id).Scan(
		&u.ID,
		&u.Address,
		&u.CheckInterval,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

// return all URL from database
func (r *URLRepository) ListOfURLs() ([]model.URL, error) {
	const op = "storage.ListOfURLs"
	query := `SELECT id, address,
		(EXTRACT(EPOCH FROM check_interval) * 1000000000)::BIGINT AS check_interval,
		created_at FROM urls ORDER BY created_at ASC`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.URL
	for rows.Next() {
		u := model.URL{}
		if err := rows.Scan(
			&u.ID,
			&u.Address,
			&u.CheckInterval,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, rows.Err()
}

// delete URL on ID
func (r *URLRepository) Delete(id uuid.UUID) error {
	const op = "storage.Delete"
	query := `DELETE FROM urls WHERE id = $1`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}
