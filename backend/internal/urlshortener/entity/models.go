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
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ShortCode    string         `json:"short_code" gorm:"column:short_code;size:20;uniqueIndex;not null"`
	OriginalURL  string         `json:"original_url" gorm:"column:original_url;type:text;not null"`
	UserID       *uuid.UUID     `json:"user_id,omitempty" gorm:"column:user_id;type:uuid"`
	Title        string         `json:"title,omitempty" gorm:"column:title;size:255"`
	Description  string         `json:"description,omitempty" gorm:"column:description;type:text"`
	IsActive     bool           `json:"is_active" gorm:"column:is_active;default:true"`
	IsPrivate    bool           `json:"is_private" gorm:"column:is_private;default:false"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;size:255"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty" gorm:"column:expires_at"`
	MaxClicks    *int           `json:"max_clicks,omitempty" gorm:"column:max_clicks"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Associations
	Clicks []URLClick `json:"clicks,omitempty" gorm:"foreignKey:ShortURLID"`
}

// TableName specifies the table name for GORM
func (ShortenedURL) TableName() string {
	return "shortened_urls"
}

// URLClick represents a click analytics record
type URLClick struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ShortURLID  uuid.UUID `json:"short_url_id" gorm:"column:short_url_id;type:uuid;not null"`
	IPAddress   net.IP    `json:"ip_address,omitempty" gorm:"column:ip_address;type:inet"`
	UserAgent   string    `json:"user_agent,omitempty" gorm:"column:user_agent;type:text"`
	Referer     string    `json:"referer,omitempty" gorm:"column:referer;type:text"`
	CountryCode string    `json:"country_code,omitempty" gorm:"column:country_code;size:2"`
	City        string    `json:"city,omitempty" gorm:"column:city;size:100"`
	DeviceType  string    `json:"device_type,omitempty" gorm:"column:device_type;size:50"`
	Browser     string    `json:"browser,omitempty" gorm:"column:browser;size:50"`
	OS          string    `json:"os,omitempty" gorm:"column:os;size:50"`
	IsBot       bool      `json:"is_bot" gorm:"column:is_bot;default:false"`
	ClickedAt   time.Time `json:"clicked_at" gorm:"column:clicked_at;autoCreateTime"`

	// Associations
	URL *ShortenedURL `json:"url,omitempty" gorm:"foreignKey:ShortURLID;references:ID"`
}

// TableName specifies the table name for GORM
func (URLClick) TableName() string {
	return "url_clicks"
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
	TotalClicks     int64            `json:"total_clicks"`
	UniqueClicks    int64            `json:"unique_clicks"`
	ClicksByCountry map[string]int64 `json:"clicks_by_country"`
	ClicksByDate    map[string]int64 `json:"clicks_by_date"`
	RecentClicks    []URLClick       `json:"recent_clicks"`
	TopReferers     map[string]int64 `json:"top_referers"`
}

// UserURLsResponse represents a user's URLs with pagination
type UserURLsResponse struct {
	URLs       []URLWithStats     `json:"urls"`
	Pagination PaginationResponse `json:"pagination"`
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
