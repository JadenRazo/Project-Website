// internal/urlshortener/service/service_impl.go
package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener/repository"
)

// URLShortenerServiceImpl implements the URLShortenerService interface
type URLShortenerServiceImpl struct {
	repo         repository.URLRepository
	config       config.URLShortenerConfig
	allowedChars string
}

// NewURLShortenerService creates a new URL shortener service
func NewURLShortenerService(repo repository.URLRepository, cfg config.URLShortenerConfig) URLShortenerService {
	return &URLShortenerServiceImpl{
		repo:         repo,
		config:       cfg,
		allowedChars: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

// CreateShortURL creates a new shortened URL
func (s *URLShortenerServiceImpl) CreateShortURL(ctx context.Context, originalURL string, userID uint, expiry *time.Time) (*domain.ShortURL, error) {
	if originalURL == "" {
		return nil, errors.New("original URL cannot be empty")
	}

	// Generate a unique short code
	shortCode, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		return nil, err
	}

	// Create the short URL entity
	shortURL := &domain.ShortURL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatorID:   userID,
	}

	// Set expiry if provided
	if expiry != nil {
		shortURL.ExpiresAt = *expiry
	}

	// Save to repository
	if err := s.repo.Create(ctx, shortURL); err != nil {
		return nil, err
	}

	return shortURL, nil
}

// GetOriginalURL retrieves the original URL from a short code and tracks the visit
func (s *URLShortenerServiceImpl) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	shortURL, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// Check if the URL has expired
	if shortURL.IsExpired() {
		return "", errors.New("URL has expired")
	}

	// Increment the click count and update last accessed time
	shortURL.IncrementClicks()

	// Update stats asynchronously to not block the redirect
	go func(su *domain.ShortURL) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.repo.UpdateStats(bgCtx, su)
	}(shortURL)

	return shortURL.OriginalURL, nil
}

// ListUserURLs lists all URLs created by a user
func (s *URLShortenerServiceImpl) ListUserURLs(ctx context.Context, userID uint, page, limit int) ([]domain.ShortURL, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	return s.repo.FindByUserID(ctx, userID, page, limit)
}

// DeleteURL deletes a shortened URL by short code
func (s *URLShortenerServiceImpl) DeleteURL(ctx context.Context, shortCode string, userID uint) error {
	return s.repo.Delete(ctx, shortCode, userID)
}

// GetURLStats gets statistics for a shortened URL
func (s *URLShortenerServiceImpl) GetURLStats(ctx context.Context, shortCode string, userID uint) (*URLStats, error) {
	shortURL, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Only allow the creator to see stats
	if shortURL.CreatorID != userID {
		return nil, errors.New("you don't have permission to view these stats")
	}

	var expiresAt *time.Time
	if !shortURL.ExpiresAt.IsZero() {
		expiresAt = &shortURL.ExpiresAt
	}

	var lastAccessed *time.Time
	if !shortURL.LastAccessedAt.IsZero() {
		lastAccessed = &shortURL.LastAccessedAt
	}

	return &URLStats{
		ShortCode:    shortURL.ShortCode,
		OriginalURL:  shortURL.OriginalURL,
		ClickCount:   shortURL.ClickCount,
		CreatedAt:    shortURL.CreatedAt,
		ExpiresAt:    expiresAt,
		LastAccessed: lastAccessed,
	}, nil
}

// generateUniqueShortCode generates a random short code that isn't already in use
func (s *URLShortenerServiceImpl) generateUniqueShortCode(ctx context.Context) (string, error) {
	maxAttempts := 5
	codeLength := s.config.ShortCodeLen

	for i := 0; i < maxAttempts; i++ {
		code, err := s.generateRandomString(codeLength)
		if err != nil {
			return "", err
		}

		// Check if code is unique
		isUnique, err := s.repo.IsShortCodeUnique(ctx, code)
		if err != nil {
			return "", err
		}

		if isUnique {
			return code, nil
		}

		// If we found a collision, try again with a longer code
		codeLength++
	}

	return "", errors.New("failed to generate a unique short code")
}

// generateRandomString generates a random string of the given length
func (s *URLShortenerServiceImpl) generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charCount := big.NewInt(int64(len(s.allowedChars)))

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, charCount)
		if err != nil {
			return "", err
		}
		result[i] = s.allowedChars[index.Int64()]
	}

	return string(result), nil
}
