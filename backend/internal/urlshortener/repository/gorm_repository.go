// internal/urlshortener/repository/gorm_repository.go
package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/core/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/patrickmn/go-cache"
)

// GormRepository implements URLRepository using Gorm
type GormRepository struct {
	db       *gorm.DB
	cache    *cache.Cache
	baseRepo *repository.BaseRepository[domain.ShortURL]
}

// NewGormRepository creates a new URL shortener repository
func NewGormRepository(db *gorm.DB) *GormRepository {
	// Initialize cache with 5 minute default expiration and 10 minute cleanup interval
	baseRepo := repository.NewBaseRepository(db, domain.ShortURL{})
	return &GormRepository{
		db:    db,
		cache: cache.New(5*time.Minute, 10*time.Minute),
		baseRepo: baseRepo.WithPreloads(
			func(db *gorm.DB) *gorm.DB { return db.Preload("Creator", "id, username, email") },
		),
	}
}

// Create stores a new short URL
func (r *GormRepository) Create(ctx context.Context, shortURL *domain.ShortURL) error {
	return r.baseRepo.Create(ctx, shortURL)
}

// FindByShortCode finds a URL by its short code
func (r *GormRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	if shortCode == "" {
		return nil, repository.ErrInvalidInput
	}

	// Check cache first for better performance
	if cachedURL, found := r.cache.Get(shortCode); found {
		return cachedURL.(*domain.ShortURL), nil
	}

	var shortURL domain.ShortURL
	if err := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&shortURL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Store in cache for future requests
	r.cache.Set(shortCode, &shortURL, cache.DefaultExpiration)

	return &shortURL, nil
}

// FindByUserID finds all URLs created by a specific user with pagination
func (r *GormRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]domain.ShortURL, int64, error) {
	if userID == 0 {
		return nil, 0, repository.ErrInvalidInput
	}

	var shortURLs []domain.ShortURL
	var total int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.ShortURL{}).Where("creator_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).Where("creator_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shortURLs).Error; err != nil {
		return nil, 0, err
	}

	return shortURLs, total, nil
}

// UpdateStats updates the click count and last accessed time for a URL
func (r *GormRepository) UpdateStats(ctx context.Context, shortURL *domain.ShortURL) error {
	if shortURL == nil || shortURL.ID == 0 {
		return repository.ErrInvalidInput
	}

	// Update in database
	err := r.db.WithContext(ctx).Model(shortURL).
		Updates(map[string]interface{}{
			"click_count":      shortURL.ClickCount,
			"last_accessed_at": shortURL.LastAccessedAt,
		}).Error

	if err != nil {
		return err
	}

	// Update in cache
	r.cache.Set(shortURL.ShortCode, shortURL, cache.DefaultExpiration)

	return nil
}

// Delete removes a shortened URL
func (r *GormRepository) Delete(ctx context.Context, shortCode string, userID uint) error {
	if shortCode == "" || userID == 0 {
		return repository.ErrInvalidInput
	}

	result := r.db.WithContext(ctx).Where("short_code = ? AND creator_id = ?", shortCode, userID).
		Delete(&domain.ShortURL{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.ErrUnauthorized
	}

	// Remove from cache
	r.cache.Delete(shortCode)

	return nil
}

// IsShortCodeUnique checks if a short code is already in use
func (r *GormRepository) IsShortCodeUnique(ctx context.Context, shortCode string) (bool, error) {
	if shortCode == "" {
		return false, repository.ErrInvalidInput
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ShortURL{}).Where("short_code = ?", shortCode).Count(&count).Error
	return count == 0, err
}

// AddClick records a new click for analytics
func (r *GormRepository) AddClick(ctx context.Context, click *domain.URLClick) error {
	if click == nil || click.ShortURLID == 0 {
		return repository.ErrInvalidInput
	}

	return r.db.WithContext(ctx).Create(click).Error
}

// GetClicks retrieves analytics data for a shortened URL
func (r *GormRepository) GetClicks(ctx context.Context, shortCode string, from, to time.Time) ([]domain.URLClick, error) {
	if shortCode == "" {
		return nil, repository.ErrInvalidInput
	}

	var shortURL domain.ShortURL
	if err := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&shortURL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	var clicks []domain.URLClick
	query := r.db.WithContext(ctx).Where("short_url_id = ?", shortURL.ID)

	// Add time constraints if provided
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}

	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Order("created_at DESC").Find(&clicks).Error; err != nil {
		return nil, err
	}

	return clicks, nil
}
