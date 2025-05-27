package repository

import (
	"context"

	"url-sentinel/internal/domain/entity"

	"github.com/google/uuid"
)

// CheckRepository defines the interface for check result persistence operations
type CheckRepository interface {
	// Create saves a new check result to the repository
	Create(ctx context.Context, check *entity.Check) error

	// ListByURLID retrieves all check results for a specific URL
	ListByURLID(ctx context.Context, urlID uuid.UUID) ([]*entity.Check, error)

	// GetLatestByURLID retrieves the most recent check for a URL
	GetLatestByURLID(ctx context.Context, urlID uuid.UUID) (*entity.Check, error)
}
