package entity

import (
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShortenedURL represents a shortened URL in the database
type ShortenedURL struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ShortCode   string         `json:"short_code" gorm:"size:10;uniqueIndex;not null"`
	OriginalURL string         `json:"original_url" gorm:"type:text;not null"`
	UserID      *uuid.UUID     `json:"user_id,omitempty" gorm:"type:uuid"`
	Title       string         `json:"title,omitempty" gorm:"size:255"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Associations
	Clicks []URLClick `json:"clicks,omitempty" gorm:"foreignKey:URLID"`
}

// URLClick represents a click analytics record
type URLClick struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	URLID     uuid.UUID `json:"url_id" gorm:"type:uuid;not null"`
	IPAddress net.IP    `json:"ip_address,omitempty" gorm:"type:inet"`
	UserAgent string    `json:"user_agent,omitempty" gorm:"type:text"`
	Referer   string    `json:"referer,omitempty" gorm:"type:text"`
	Country   string    `json:"country,omitempty" gorm:"size:2"`
	City      string    `json:"city,omitempty" gorm:"size:100"`
	ClickedAt time.Time `json:"clicked_at" gorm:"autoCreateTime"`

	// Associations
	URL *ShortenedURL `json:"url,omitempty" gorm:"foreignKey:URLID"`
}

// Request/Response DTOs

// ShortenURLRequest represents the request to shorten a URL
type ShortenURLRequest struct {
	URL         string `json:"url" binding:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" binding:"omitempty,min=4,max=10,alphanum"`
	Title       string `json:"title,omitempty" binding:"omitempty,max=255"`
	Description string `json:"description,omitempty" binding:"omitempty,max=1000"`
	ExpiresAt   string `json:"expires_at,omitempty" binding:"omitempty"`
}

// ShortenURLResponse represents the response when shortening a URL
type ShortenURLResponse struct {
	ID          uuid.UUID  `json:"id"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// URLStatsResponse represents URL analytics data
type URLStatsResponse struct {
	*ShortenedURL
	TotalClicks     int64                `json:"total_clicks"`
	UniqueClicks    int64                `json:"unique_clicks"`
	ClicksByCountry map[string]int64     `json:"clicks_by_country"`
	ClicksByDate    map[string]int64     `json:"clicks_by_date"`
	RecentClicks    []URLClick          `json:"recent_clicks"`
	TopReferers     map[string]int64     `json:"top_referers"`
}

// UserURLsResponse represents a user's URLs with pagination
type UserURLsResponse struct {
	URLs       []URLWithStats      `json:"urls"`
	Pagination PaginationResponse  `json:"pagination"`
}

// URLWithStats represents a URL with basic statistics
type URLWithStats struct {
	*ShortenedURL
	TotalClicks  int64 `json:"total_clicks"`
	UniqueClicks int64 `json:"unique_clicks"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Validation helpers

// IsValidShortCode checks if a short code is valid
func IsValidShortCode(code string) bool {
	if len(code) < 4 || len(code) > 10 {
		return false
	}
	
	// Only allow alphanumeric characters
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	
	return true
}

// SanitizeURL cleans and validates a URL
func SanitizeURL(url string) string {
	url = strings.TrimSpace(url)
	
	// Add https:// if no protocol is specified
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	
	return url
}

// IsURLExpired checks if a URL has expired
func (u *ShortenedURL) IsURLExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

// IsAccessible checks if a URL can be accessed (active and not expired)
func (u *ShortenedURL) IsAccessible() bool {
	return u.IsActive && !u.IsURLExpired()
}