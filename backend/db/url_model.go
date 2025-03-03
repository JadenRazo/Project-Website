package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// URLRepository wraps the generic repository with URL-specific methods
type URLRepository struct {
	Repository[URL]
}

// FindByShortCode finds a URL by its short code
func (r *URLRepository) FindByShortCode(shortCode string) (*URL, error) {
	return r.FindBy("short_code", shortCode)
}

// IncrementClicks increases the click count for a URL
func (r *URLRepository) IncrementClicks(id uint) error {
	result := GetDB().Model(&URL{}).
		Where("id = ?", id).
		UpdateColumn("clicks", gorm.Expr("clicks + 1"))
	
	if result.Error != nil {
		return fmt.Errorf("failed to increment clicks: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// FindExpired finds all URLs that have expired
func (r *URLRepository) FindExpired() ([]URL, error) {
	return r.FindAllWhere("expires_at IS NOT NULL AND expires_at < ?", time.Now())
}

// FindByUser retrieves all URLs for a specific user with optional pagination
func (r *URLRepository) FindByUser(userID uint, page, pageSize int) ([]URL, *Pagination, error) {
	return r.Paginate(page, pageSize, WithOrder("created_at", SortDescending),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ?", userID)
		})
}

// SearchURLs provides advanced searching functionality for URLs
func (r *URLRepository) SearchURLs(userID uint, query string, tag string, page, pageSize int, sortBy, sortOrder string) ([]URL, *Pagination, error) {
	// Validate sort parameters
	validSortFields := map[string]bool{
		"created_at": true, "clicks": true, "expires_at": true,
	}
	
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	
	var direction SortDirection = SortDescending
	if sortOrder == "asc" {
		direction = SortAscending
	}
	
	// Build query options
	options := []QueryOption{
		WithOrder(sortBy, direction),
		func(db *gorm.DB) *gorm.DB {
			db = db.Where("user_id = ?", userID)
			
			if query != "" {
				db = db.Where("original_url LIKE ? OR short_code LIKE ?", 
					"%"+query+"%", "%"+query+"%")
			}
			
			if tag != "" {
				db = db.Where("tags LIKE ?", "%"+tag+"%")
			}
			
			return db
		},
	}
	
	return r.Paginate(page, pageSize, options...)
}
