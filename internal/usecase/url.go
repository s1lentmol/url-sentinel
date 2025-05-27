package usecase

import (
	"context"
	"fmt"
	"time"

	"url-sentinel/internal/domain/entity"
	"url-sentinel/internal/domain/repository"

	"github.com/google/uuid"
)

// Monitor defines the interface for URL monitoring
type Monitor interface {
	AddURL(ctx context.Context, url *entity.URL)
	RemoveURL(urlID string)
}

// URLUseCase handles business logic for URL operations
type URLUseCase struct {
	urlRepo repository.URLRepository
	monitor Monitor
}

// NewURLUseCase creates a new URL use case
func NewURLUseCase(urlRepo repository.URLRepository, monitor Monitor) *URLUseCase {
	return &URLUseCase{
		urlRepo: urlRepo,
		monitor: monitor,
	}
}

// CreateURL creates a new URL with validation and starts monitoring
func (uc *URLUseCase) CreateURL(ctx context.Context, address string, interval time.Duration) (*entity.URL, error) {
	// Check if URL already exists
	exists, err := uc.urlRepo.ExistsByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to check url existence: %w", err)
	}
	if exists {
		return nil, repository.ErrURLAddressExists
	}

	// Create new URL entity
	url, err := entity.NewURL(address, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to create url entity: %w", err)
	}

	// Save to repository
	if err := uc.urlRepo.Create(ctx, url); err != nil {
		return nil, fmt.Errorf("failed to save url: %w", err)
	}

	// Add to monitoring if monitor is available
	if uc.monitor != nil {
		uc.monitor.AddURL(ctx, url)
	}

	return url, nil
}

// GetURLByID retrieves a URL by its ID
func (uc *URLUseCase) GetURLByID(ctx context.Context, id uuid.UUID) (*entity.URL, error) {
	url, err := uc.urlRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	return url, nil
}

// ListURLs retrieves all URLs
func (uc *URLUseCase) ListURLs(ctx context.Context) ([]*entity.URL, error) {
	urls, err := uc.urlRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list urls: %w", err)
	}

	return urls, nil
}

// DeleteURL deletes a URL by its ID and stops monitoring
func (uc *URLUseCase) DeleteURL(ctx context.Context, id uuid.UUID) error {
	// Remove from monitoring first
	if uc.monitor != nil {
		uc.monitor.RemoveURL(id.String())
	}

	// Delete from repository
	if err := uc.urlRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}

	return nil
}
