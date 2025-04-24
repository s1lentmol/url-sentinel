package model

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// URL - monitored address
// ID genereted upon creation, Address - valid URL,
// CheckInterval — positive check interval
// CreatedAt records the time of addition
type URL struct {
	ID            uuid.UUID     `db:"id"`
	Address       string        `db:"address"`
	CheckInterval time.Duration `db:"check_interval"`
	CreatedAt     time.Time     `db:"created_at"`
}

// NewURL creates new entity URL with generating ID and time of creation,
// checking correctness of address and interval
func NewURL(address string, interval time.Duration) (*URL, error) {
	if _, err := url.ParseRequestURI(address); err != nil {
		return nil, errors.New("invalid URL format")
	}
	if interval <= 0 {
		return nil, errors.New("check interval must be positive")
	}
	return &URL{
		ID:            uuid.New(),
		Address:       address,
		CheckInterval: interval,
		CreatedAt:     time.Now(),
	}, nil
}

// Validate checks the correctness of the URL entity fields
func (u *URL) Validate() error {
	if u.ID == uuid.Nil {
		return errors.New("id must be set")
	}
	if _, err := url.ParseRequestURI(u.Address); err != nil {
		return errors.New("invalid URL format")
	}
	if u.CheckInterval <= 0 {
		return errors.New("check interval must be positive")
	}
	return nil
}