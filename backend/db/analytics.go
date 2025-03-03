package db

import (
	"time"

	"gorm.io/gorm"
)

// ClickAnalytics tracks detailed information about URL clicks
type ClickAnalytics struct {
	Model
	URLID       uint   `json:"urlId"`
	URL         *URL   `gorm:"foreignKey:URLID" json:"-"`
	ClickedAt   time.Time `gorm:"index" json:"clickedAt"`
	IPAddress   string `json:"ipAddress,omitempty"`
	UserAgent   string `json:"userAgent,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	City        string `json:"city,omitempty"`
	Device      string `json:"device,omitempty"` // mobile, desktop, tablet, etc.
	Browser     string `json:"browser,omitempty"`
	OS          string `json:"os,omitempty"`
}

// AnalyticsRepository handles analytics data operations
type AnalyticsRepository struct {
	*BaseRepository
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// RecordClick logs a click event for a URL
func (r *AnalyticsRepository) RecordClick(urlID uint, ipAddr, userAgent, referrer string) error {
	analytics := ClickAnalytics{
		URLID:     urlID,
		ClickedAt: time.Now(),
		IPAddress: ipAddr,
		UserAgent: userAgent,
		Referrer:  referrer,
		// Extract browser/OS/device from user agent in a real implementation
	}
	
	return r.Create(&analytics)
}

// GetClicksForURL retrieves analytics for a specific URL
func (r *AnalyticsRepository) GetClicksForURL(urlID uint, since time.Time, page, pageSize int) ([]ClickAnalytics, int64, error) {
	var analytics []ClickAnalytics
	var count int64
	
	query := r.DB.Model(&ClickAnalytics{}).Where("url_id = ?", urlID)
	
	if !since.IsZero() {
		query = query.Where("clicked_at >= ?", since)
	}
	
	// Get total count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	offset := (page - 1) * pageSize
	err = query.Order("clicked_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&analytics).Error
		
	return analytics, count, err
}

// GetAggregatedStats returns analytics grouped by different dimensions
func (r *AnalyticsRepository) GetAggregatedStats(urlID uint, since time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get browser stats
	var browserStats []struct {
		Browser string
		Count   int
	}
	
	err := r.DB.Model(&ClickAnalytics{}).
		Select("browser, count(*) as count").
		Where("url_id = ? AND clicked_at >= ?", urlID, since).
		Group("browser").
		Order("count DESC").
		Find(&browserStats).Error
		
	if err != nil {
		return nil, err
	}
	
	stats["browsers"] = browserStats
	
	// Get daily clicks
	var dailyClicks []struct {
		Date  string
		Count int
	}
	
	// SQLite date formatting
	dateFormat := "date(clicked_at)"
	if r.DB.Dialector.Name() == "postgres" {
		dateFormat = "TO_CHAR(clicked_at, 'YYYY-MM-DD')" 
	} else if r.DB.Dialector.Name() == "mysql" {
		dateFormat = "DATE_FORMAT(clicked_at, '%Y-%m-%d')"
	}
	
	err = r.DB.Model(&ClickAnalytics{}).
		Select(dateFormat + " as date, count(*) as count").
		Where("url_id = ? AND clicked_at >= ?", urlID, since).
		Group("date").
		Order("date").
		Find(&dailyClicks).Error
		
	if err != nil {
		return nil, err
	}
	
	stats["dailyClicks"] = dailyClicks
	
	// Add other stats as needed (countries, devices, etc.)
	
	return stats, nil
}
