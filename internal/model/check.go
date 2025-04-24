package model

import (
	"time"

	"github.com/google/uuid"
)

// Check present the result of a single URL check
// ID is generated automatically, URLID is associated with the URL entity
// Status = true if code < 300, Code — HTTP response code,
// Duration — duration of the request, CheckedAt — check time

// Check the structure of the check results
type Check struct {
	ID        uuid.UUID     `db:"id"`
	URLID     uuid.UUID     `db:"url_id"`
	Status    bool          `db:"status"`
	Code      int           `db:"code"`
	Duration  time.Duration `db:"duration"`
	CheckedAt time.Time     `db:"checked_at"`
}

// NewCheck creates new record checks, generates an ID and records the time
func NewCheck(urlID uuid.UUID, status bool, code int, duration time.Duration) *Check {
	return &Check{
		ID:        uuid.New(),
		URLID:     urlID,
		Status:    status,
		Code:      code,
		Duration:  duration,
		CheckedAt: time.Now(),
	}
}
