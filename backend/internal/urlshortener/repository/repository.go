package repository

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// URLRepository defines the interface for URL shortener data access
type URLRepository interface {
	// Create stores a new short URL
	Create(ctx context.Context, shortURL *domain.ShortURL) error

	// FindByShortCode finds a URL by its short code
	FindByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error)

	// FindByUserID finds all URLs created by a specific user with pagination
	FindByUserID(ctx context.Context, userID uint, page, limit int) ([]domain.ShortURL, int64, error)

	// UpdateStats updates the click count and last accessed time for a URL
	UpdateStats(ctx context.Context, shortURL *domain.ShortURL) error

	// Delete removes a shortened URL
	Delete(ctx context.Context, shortCode string, userID uint) error

	// IsShortCodeUnique checks if a short code is already in use
	IsShortCodeUnique(ctx context.Context, shortCode string) (bool, error)

	// AddClick records a new click for analytics
	AddClick(ctx context.Context, click *domain.URLClick) error

	// GetClicks retrieves analytics data for a shortened URL
	GetClicks(ctx context.Context, shortCode string, from, to time.Time) ([]domain.URLClick, error)
}
