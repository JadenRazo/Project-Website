// internal/urlshortener/service/service.go
package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// URLShortenerService defines the interface for URL shortener operations
type URLShortenerService interface {
	// CreateShortURL creates a new shortened URL
	CreateShortURL(ctx context.Context, originalURL string, userID uint, expiry *time.Time) (*domain.ShortURL, error)

	// GetOriginalURL retrieves the original URL from a short code and tracks the visit
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)

	// ListUserURLs lists all URLs created by a user
	ListUserURLs(ctx context.Context, userID uint, page, limit int) ([]domain.ShortURL, int64, error)

	// DeleteURL deletes a shortened URL by short code
	DeleteURL(ctx context.Context, shortCode string, userID uint) error

	// GetURLStats gets statistics for a shortened URL
	GetURLStats(ctx context.Context, shortCode string, userID uint) (*URLStats, error)
}

// URLStats represents usage statistics for a shortened URL
type URLStats struct {
	ShortCode    string
	OriginalURL  string
	ClickCount   int64
	CreatedAt    time.Time
	ExpiresAt    *time.Time
	LastAccessed *time.Time
}
