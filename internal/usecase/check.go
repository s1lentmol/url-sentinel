package usecase

import (
	"context"
	"fmt"

	"url-sentinel/internal/domain/entity"
	"url-sentinel/internal/domain/repository"

	"github.com/google/uuid"
)

// CheckUseCase handles business logic for check operations
type CheckUseCase struct {
	checkRepo repository.CheckRepository
}

// NewCheckUseCase creates a new check use case
func NewCheckUseCase(checkRepo repository.CheckRepository) *CheckUseCase {
	return &CheckUseCase{
		checkRepo: checkRepo,
	}
}

// GetCheckHistory retrieves all checks for a URL
func (uc *CheckUseCase) GetCheckHistory(ctx context.Context, urlID uuid.UUID) ([]*entity.Check, error) {
	checks, err := uc.checkRepo.ListByURLID(ctx, urlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get check history: %w", err)
	}

	return checks, nil
}

// GetLatestCheck retrieves the most recent check for a URL
func (uc *CheckUseCase) GetLatestCheck(ctx context.Context, urlID uuid.UUID) (*entity.Check, error) {
	check, err := uc.checkRepo.GetLatestByURLID(ctx, urlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest check: %w", err)
	}

	return check, nil
}
