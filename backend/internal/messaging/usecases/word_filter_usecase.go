package usecases

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

type WordFilterUseCase struct {
	wordFilterRepo domain.WordFilterRepository
}

func NewWordFilterUseCase(wordFilterRepo domain.WordFilterRepository) *WordFilterUseCase {
	return &WordFilterUseCase{
		wordFilterRepo: wordFilterRepo,
	}
}

func (uc *WordFilterUseCase) CreateFilter(ctx context.Context, filter *domain.WordFilter) error {
	// Validate the filter
	if filter.Word == "" {
		return domain.NewValidationError("word is required")
	}
	if filter.Scope == "" {
		return domain.NewValidationError("scope is required")
	}
	if filter.Scope != "private" && filter.Scope != "public" && filter.Scope != "all" {
		return domain.NewValidationError("invalid scope")
	}

	// Create the filter
	return uc.wordFilterRepo.CreateFilter(ctx, filter)
}

func (uc *WordFilterUseCase) UpdateFilter(ctx context.Context, filter *domain.WordFilter) error {
	// Validate the filter
	if filter.Word == "" {
		return domain.NewValidationError("word is required")
	}
	if filter.Scope == "" {
		return domain.NewValidationError("scope is required")
	}
	if filter.Scope != "private" && filter.Scope != "public" && filter.Scope != "all" {
		return domain.NewValidationError("invalid scope")
	}

	// Update the filter
	return uc.wordFilterRepo.UpdateFilter(ctx, filter)
}

func (uc *WordFilterUseCase) DeleteFilter(ctx context.Context, id uuid.UUID) error {
	return uc.wordFilterRepo.DeleteFilter(ctx, id)
}

func (uc *WordFilterUseCase) GetFilter(ctx context.Context, id uuid.UUID) (*domain.WordFilter, error) {
	return uc.wordFilterRepo.GetFilter(ctx, id)
}

func (uc *WordFilterUseCase) ListFilters(ctx context.Context, serverID uuid.UUID) ([]*domain.WordFilter, error) {
	return uc.wordFilterRepo.ListFilters(ctx, serverID)
}

func (uc *WordFilterUseCase) GetFiltersByScope(ctx context.Context, serverID uuid.UUID, scope domain.WordFilterScope) ([]*domain.WordFilter, error) {
	return uc.wordFilterRepo.GetFiltersByScope(ctx, serverID, scope)
}

func (uc *WordFilterUseCase) CheckMessage(ctx context.Context, serverID uuid.UUID, message string, scope domain.WordFilterScope) (bool, []*domain.WordFilter, error) {
	return uc.wordFilterRepo.CheckMessage(ctx, serverID, message, scope)
}
