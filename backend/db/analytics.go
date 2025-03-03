package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AnalyticsRepository provides specialized methods for click analytics
type AnalyticsRepository struct {
	Repository[ClickAnalytics]
}

// RecordClick logs a click event for a URL
func (r *AnalyticsRepository) RecordClick(urlID uint, ipAddr, userAgent, referrer string) error {
	analytics := ClickAnalytics{
		URLID:     urlID,
		ClickedAt: time.Now(),
		IPAddress: ipAddr,
		UserAgent: userAgent,
		Referrer:  referrer,
		// Additional fields like Browser, OS, Device would be extracted in a production system
	}
	
	return r.Create(&analytics)
}

// GetClicksForURL retrieves analytics for a specific URL with pagination
func (r *AnalyticsRepository) GetClicksForURL(urlID uint, since time.Time, page, pageSize int) ([]ClickAnalytics, *Pagination, error) {
	options := []QueryOption{
		WithOrder("clicked_at", SortDescending),
		func(db *gorm.DB) *gorm.DB {
			db = db.Where("url_id = ?", urlID)
			
			if !since.IsZero() {
				db = db.Where("clicked_at >= ?", since)
			}
			
			return db
		},
	}
	
	return r.Paginate(page, pageSize, options...)
}

// GetAggregatedStats returns analytics grouped by different dimensions
func (r *AnalyticsRepository) GetAggregatedStats(urlID uint, since time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get browser stats
	var browserStats []struct {
		Browser string
		Count   int
	}
	
	query := GetDB().Model(&ClickAnalytics{}).
		Select("browser, count(*) as count").
		Where("url_id = ?", urlID)
		
	if !since.IsZero() {
		query = query.Where("clicked_at >= ?", since)
	}
	
	err := query.Group("browser").
		Order("count DESC").
		Find(&browserStats).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats: %w", err)
	}
	
	stats["browsers"] = browserStats
	
	// Get daily clicks
	var dailyClicks []struct {
		Date  string
		Count int
	}
	
	// SQLite date formatting
	dateFormat := "date(clicked_at)"
	if GetDB().Dialector.Name() == "postgres" {
		dateFormat = "TO_CHAR(clicked_at, 'YYYY-MM-DD')" 
	} else if GetDB().Dialector.Name() == "mysql" {
		dateFormat = "DATE_FORMAT(clicked_at, '%Y-%m-%d')"
	}
	
	dailyQuery := GetDB().Model(&ClickAnalytics{}).
		Select(dateFormat + " as date, count(*) as count").
		Where("url_id = ?", urlID)
	
	if !since.IsZero() {
		dailyQuery = dailyQuery.Where("clicked_at >= ?", since)
	}
	
	err = dailyQuery.Group("date").
		Order("date").
		Find(&dailyClicks).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get daily clicks: %w", err)
	}
	
	stats["dailyClicks"] = dailyClicks
	
	// Get device stats
	var deviceStats []struct {
		Device string
		Count  int
	}
	
	deviceQuery := GetDB().Model(&ClickAnalytics{}).
		Select("device, count(*) as count").
		Where("url_id = ?", urlID)
	
	if !since.IsZero() {
		deviceQuery = deviceQuery.Where("clicked_at >= ?", since)
	}
	
	err = deviceQuery.Group("device").
		Order("count DESC").
		Find(&deviceStats).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}
	
	stats["devices"] = deviceStats
	
	// Get country stats
	var countryStats []struct {
		CountryCode string
		Count       int
	}
	
	countryQuery := GetDB().Model(&ClickAnalytics{}).
		Select("country_code, count(*) as count").
		Where("url_id = ? AND country_code IS NOT NULL AND country_code != ''", urlID)
	
	if !since.IsZero() {
		countryQuery = countryQuery.Where("clicked_at >= ?", since)
	}
	
	err = countryQuery.Group("country_code").
		Order("count DESC").
		Find(&countryStats).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}
	
	stats["countries"] = countryStats
	
	return stats, nil
}
