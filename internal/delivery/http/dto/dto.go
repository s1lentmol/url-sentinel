package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateURLRequest represents the request to create a new URL
type CreateURLRequest struct {
	Address       string `json:"address"`
	CheckInterval string `json:"check_interval"` // e.g. "30s", "1m", "5m"
}

// URLResponse represents a URL in API responses
type URLResponse struct {
	ID            uuid.UUID `json:"id"`
	Address       string    `json:"address"`
	CheckInterval string    `json:"check_interval"`
	CreatedAt     time.Time `json:"created_at"`
}

// CheckResponse represents a check result in API responses
type CheckResponse struct {
	ID        uuid.UUID `json:"id"`
	URLID     uuid.UUID `json:"url_id"`
	Status    bool      `json:"status"`
	Code      int       `json:"code"`
	Duration  string    `json:"duration"` // e.g. "123ms"
	CheckedAt time.Time `json:"checked_at"`
}

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Error string `json:"error"`
}
