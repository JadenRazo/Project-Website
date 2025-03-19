package postgres

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"

	"jadenrazo.dev/internal/urlshortener/domain"
)

// URLRepository implements both domain.URLRepository and domain.StatsRepository
type URLRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new postgres repository for the URL shortener
func NewURLRepository(config interface{}) (*URLRepository, error) {
	// For demonstration purposes, we'll use a mock implementation
	// In a real application, we would connect to the database
	return &URLRepository{
		db: nil, // This would be a real DB connection in production
	}, nil
}

// Close closes the database connection
func (r *URLRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// GetByShortCode returns a URL by its short code
func (r *URLRepository) GetByShortCode(shortCode string) (*domain.URL, error) {
	// Mock implementation for demonstration
	if shortCode == "demo" {
		return &domain.URL{
			ID:          1,
			ShortCode:   "demo",
			OriginalURL: "https://example.com",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
			IsCustom:    true,
			IsActive:    true,
		}, nil
	}

	// Return nil for other short codes for this demo
	return nil, errors.New("URL not found")
}

// GetByOriginalURL returns a URL by its original URL
func (r *URLRepository) GetByOriginalURL(originalURL string) (*domain.URL, error) {
	// Mock implementation
	if originalURL == "https://example.com" {
		return &domain.URL{
			ID:          1,
			ShortCode:   "demo",
			OriginalURL: "https://example.com",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
			IsCustom:    true,
			IsActive:    true,
		}, nil
	}

	return nil, errors.New("URL not found")
}

// Save saves a URL to the database
func (r *URLRepository) Save(url *domain.URL) error {
	// Mock implementation
	url.ID = 2 // Simulate auto-increment ID
	return nil
}

// Update updates an existing URL in the database
func (r *URLRepository) Update(url *domain.URL) error {
	// Mock implementation
	return nil
}

// Delete deletes a URL from the database
func (r *URLRepository) Delete(shortCode string) error {
	// Mock implementation
	return nil
}

// IsShortCodeAvailable checks if a short code is available for use
func (r *URLRepository) IsShortCodeAvailable(shortCode string) (bool, error) {
	// Mock implementation
	if shortCode == "demo" {
		return false, nil
	}
	return true, nil
}

// GetAllURLs returns all URLs with pagination
func (r *URLRepository) GetAllURLs(page, limit int) ([]domain.URL, error) {
	// Mock implementation
	urls := []domain.URL{
		{
			ID:          1,
			ShortCode:   "demo",
			OriginalURL: "https://example.com",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
			IsCustom:    true,
			IsActive:    true,
		},
	}
	return urls, nil
}

// GetTotalURLCount returns the total count of URLs
func (r *URLRepository) GetTotalURLCount() (int, error) {
	// Mock implementation
	return 1, nil
}

// SaveClick saves a click to the database
func (r *URLRepository) SaveClick(click *domain.Click) error {
	// Mock implementation
	click.ID = 1 // Simulate auto-increment ID
	return nil
}

// GetTotalClicks returns the total number of clicks for a URL
func (r *URLRepository) GetTotalClicks(shortCode string) (int64, error) {
	// Mock implementation
	return 42, nil
}

// GetUniqueClicks returns the number of unique clicks for a URL
func (r *URLRepository) GetUniqueClicks(shortCode string) (int64, error) {
	// Mock implementation
	return 30, nil
}

// GetTopReferrers returns the top referrers for a URL
func (r *URLRepository) GetTopReferrers(shortCode string, limit int) ([]domain.Referrer, error) {
	// Mock implementation
	referrers := []domain.Referrer{
		{Referrer: "google.com", Count: 20},
		{Referrer: "twitter.com", Count: 15},
		{Referrer: "facebook.com", Count: 7},
	}
	return referrers, nil
}

// GetTopCountries returns the top countries for a URL
func (r *URLRepository) GetTopCountries(shortCode string, limit int) ([]domain.Country, error) {
	// Mock implementation
	countries := []domain.Country{
		{Country: "United States", Count: 25},
		{Country: "Canada", Count: 10},
		{Country: "United Kingdom", Count: 7},
	}
	return countries, nil
}

// GetTopBrowsers returns the top browsers for a URL
func (r *URLRepository) GetTopBrowsers(shortCode string, limit int) ([]domain.Browser, error) {
	// Mock implementation
	browsers := []domain.Browser{
		{Browser: "Chrome", Count: 30},
		{Browser: "Firefox", Count: 8},
		{Browser: "Safari", Count: 4},
	}
	return browsers, nil
}

// GetClicksOverTime returns clicks over time aggregated by day
func (r *URLRepository) GetClicksOverTime(shortCode string, filter domain.StatsFilter) ([]domain.ClicksOverTime, error) {
	// Mock implementation
	now := time.Now()
	clicksOverTime := []domain.ClicksOverTime{
		{Date: now.AddDate(0, 0, -2), Count: 15},
		{Date: now.AddDate(0, 0, -1), Count: 12},
		{Date: now, Count: 15},
	}
	return clicksOverTime, nil
}
