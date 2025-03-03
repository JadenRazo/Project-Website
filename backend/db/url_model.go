package db

import (
	"time"
)

// URL represents a shortened URL
type URL struct {
	Model
	ShortCode    string    `gorm:"uniqueIndex;not null" json:"shortCode"`
	OriginalURL  string    `gorm:"not null" json:"originalUrl"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	Clicks       uint       `gorm:"default:0" json:"clicks"`
	UserID       uint       `json:"userId,omitempty"`
	User         *User      `gorm:"foreignKey:UserID" json:"-"`
	IsPrivate    bool       `gorm:"default:false" json:"isPrivate"`
	CustomDomain string    `json:"customDomain,omitempty"`
	Tags         string    `json:"tags,omitempty"`
	Analytics    []ClickAnalytics `gorm:"foreignKey:URLID" json:"-"`
}

// URLRepository defines operations for URL data access
type URLRepository struct {
	*BaseRepository
}

// NewURLRepository creates a new URL repository
func NewURLRepository() *URLRepository {
	return &URLRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByShortCode retrieves a URL by short code
func (r *URLRepository) FindByShortCode(shortCode string) (*URL, error) {
	var url URL
	err := r.DB.Where("short_code = ?", shortCode).First(&url).Error
	return &url, err
}

// IncrementClicks increases the click count for a URL
func (r *URLRepository) IncrementClicks(id uint) error {
	return r.DB.Model(&URL{}).Where("id = ?", id).
		UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error
}

// FindByUser retrieves all URLs for a specific user
func (r *URLRepository) FindByUser(userID uint) ([]URL, error) {
	var urls []URL
	err := r.DB.Where("user_id = ?", userID).Find(&urls).Error
	return urls, err
}

// SearchURLs searches for URLs with pagination and filters
func (r *URLRepository) SearchURLs(userID uint, query string, tag string, 
	page, pageSize int, sortBy, sortOrder string) ([]URL, int64, error) {
	var urls []URL
	var count int64

	db := r.DB.Model(&URL{}).Where("user_id = ?", userID)

	// Apply filters
	if query != "" {
		db = db.Where("original_url LIKE ? OR short_code LIKE ?", 
			"%"+query+"%", "%"+query+"%")
	}
	
	if tag != "" {
		db = db.Where("tags LIKE ?", "%"+tag+"%")
	}

	// Get total count
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	offset := (page - 1) * pageSize
	
	// Validate and apply sort
	validSortColumns := map[string]bool{
		"created_at": true, "clicks": true, "expires_at": true,
	}
	
	if !validSortColumns[sortBy] {
		sortBy = "created_at"
	}
	
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	
	err = db.Order(sortBy + " " + sortOrder).
		Offset(offset).Limit(pageSize).
		Find(&urls).Error

	return urls, count, err
}
