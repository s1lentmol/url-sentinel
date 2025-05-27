package entity

import (
	"time"

	"github.com/google/uuid"
)

// Check represents the result of a single URL health check
type Check struct {
	ID        uuid.UUID
	URLID     uuid.UUID
	Status    bool
	Code      int
	Duration  time.Duration
	CheckedAt time.Time
}

// NewCheck creates a new check result entity
func NewCheck(urlID uuid.UUID, status bool, code int, duration time.Duration) *Check {
	return &Check{
		ID:        uuid.New(),
		URLID:     urlID,
		Status:    status,
		Code:      code,
		Duration:  duration,
		CheckedAt: time.Now().UTC(),
	}
}
