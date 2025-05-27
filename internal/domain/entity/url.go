package entity

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidURLFormat       = errors.New("invalid URL format")
	ErrInvalidCheckInterval   = errors.New("check interval must be positive")
	ErrURLIDRequired          = errors.New("url id is required")
)

// URL represents a monitored web address with its configuration
type URL struct {
	ID            uuid.UUID
	Address       string
	CheckInterval time.Duration
	CreatedAt     time.Time
}

// NewURL creates a new URL entity with validation
func NewURL(address string, interval time.Duration) (*URL, error) {
	if _, err := url.ParseRequestURI(address); err != nil {
		return nil, ErrInvalidURLFormat
	}
	if interval <= 0 {
		return nil, ErrInvalidCheckInterval
	}

	return &URL{
		ID:            uuid.New(),
		Address:       address,
		CheckInterval: interval,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

// Validate checks the correctness of the URL entity
func (u *URL) Validate() error {
	if u.ID == uuid.Nil {
		return ErrURLIDRequired
	}
	if _, err := url.ParseRequestURI(u.Address); err != nil {
		return ErrInvalidURLFormat
	}
	if u.CheckInterval <= 0 {
		return ErrInvalidCheckInterval
	}
	return nil
}
