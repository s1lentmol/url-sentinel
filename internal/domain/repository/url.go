package repository

import (
	"context"
	"errors"

	"url-sentinel/internal/domain/entity"

	"github.com/google/uuid"
)

var (
	ErrURLNotFound      = errors.New("url not found")
	ErrURLAlreadyExists = errors.New("url already exists")
	ErrURLAddressExists = errors.New("url with this address already exists")
)

// URLRepository defines the interface for URL persistence operations
type URLRepository interface {
	// Create saves a new URL to the repository
	Create(ctx context.Context, url *entity.URL) error

	// GetByID retrieves a URL by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.URL, error)

	// List retrieves all URLs from the repository
	List(ctx context.Context) ([]*entity.URL, error)

	// Delete removes a URL by its ID
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByAddress checks if a URL with the given address already exists
	ExistsByAddress(ctx context.Context, address string) (bool, error)
}
