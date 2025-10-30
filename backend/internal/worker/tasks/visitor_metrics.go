package tasks

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/visitor"
)

// VisitorMetricsTask handles visitor metrics aggregation
type VisitorMetricsTask struct {
	db *gorm.DB
}

// NewVisitorMetricsTask creates a new visitor metrics task
func NewVisitorMetricsTask(db *gorm.DB) *VisitorMetricsTask {
	return &VisitorMetricsTask{db: db}
}

// AggregateHourlyMetrics aggregates visitor metrics for the past hour
func (t *VisitorMetricsTask) AggregateHourlyMetrics(ctx context.Context) error {
	now := time.Now()
	hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	hourEnd := hourStart.Add(time.Hour)

	// Skip if we're not at the end of the hour yet
	if now.Before(hourEnd.Add(-5 * time.Minute)) {
		return nil
	}

	hour := hourStart.Hour()
	
	// Calculate unique visitors for this hour
	var uniqueVisitors int64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", hourStart, hourEnd).
		Count(&uniqueVisitors)

	// Calculate total page views
	var totalPageViews int64
	t.db.Model(&visitor.PageView{}).
		Joins("JOIN visitor_sessions ON visitor_sessions.id = page_views.session_id").
		Where("page_views.created_at >= ? AND page_views.created_at < ?", hourStart, hourEnd).
		Count(&totalPageViews)

	// Calculate average session duration
	var avgDuration float64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", hourStart, hourEnd).
		Select("AVG(EXTRACT(EPOCH FROM (last_seen_at - created_at)))").
		Scan(&avgDuration)

	// Calculate bounce rate
	var totalSessions, bouncedSessions int64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", hourStart, hourEnd).
		Count(&totalSessions)

	t.db.Raw(`
		SELECT COUNT(*) FROM visitor_sessions vs
		WHERE vs.created_at >= ? AND vs.created_at < ?
		AND (SELECT COUNT(*) FROM page_views pv WHERE pv.session_id = vs.id) = 1
	`, hourStart, hourEnd).Scan(&bouncedSessions)

	bounceRate := float64(0)
	if totalSessions > 0 {
		bounceRate = float64(bouncedSessions) / float64(totalSessions) * 100
	}

	// Count new vs returning visitors
	var newVisitors int64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", hourStart, hourEnd).
		Where("session_hash NOT IN (SELECT DISTINCT session_hash FROM visitor_sessions WHERE created_at < ?)", hourStart).
		Count(&newVisitors)

	returningVisitors := uniqueVisitors - newVisitors

	// Create or update metrics record
	metrics := visitor.VisitorMetrics{
		MetricDate:         hourStart,
		Hour:               &hour,
		UniqueVisitors:     int(uniqueVisitors),
		TotalPageViews:     int(totalPageViews),
		AvgSessionDuration: int(avgDuration),
		BounceRate:         bounceRate,
		NewVisitors:        int(newVisitors),
		ReturningVisitors:  int(returningVisitors),
		CreatedAt:          time.Now(),
	}

	// Upsert the metrics
	return t.db.Where("metric_date = ? AND hour = ?", hourStart, hour).
		Assign(metrics).
		FirstOrCreate(&metrics).Error
}

// GenerateDailySummary generates daily summary statistics
func (t *VisitorMetricsTask) GenerateDailySummary(ctx context.Context) error {
	// Run at midnight for the previous day
	now := time.Now()
	if now.Hour() != 0 || now.Minute() > 10 {
		return nil // Only run shortly after midnight
	}

	yesterday := now.AddDate(0, 0, -1)
	dayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	dayEnd := dayStart.AddDate(0, 0, 1)

	// Calculate daily metrics
	var uniqueVisitors int64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Count(&uniqueVisitors)

	var totalSessions int64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Count(&totalSessions)

	var totalPageViews int64
	t.db.Model(&visitor.PageView{}).
		Joins("JOIN visitor_sessions ON visitor_sessions.id = page_views.session_id").
		Where("visitor_sessions.created_at >= ? AND visitor_sessions.created_at < ?", dayStart, dayEnd).
		Count(&totalPageViews)

	avgPagesPerSession := float64(0)
	if totalSessions > 0 {
		avgPagesPerSession = float64(totalPageViews) / float64(totalSessions)
	}

	// Calculate average session duration
	var avgDuration float64
	t.db.Model(&visitor.VisitorSession{}).
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Select("AVG(EXTRACT(EPOCH FROM (last_seen_at - created_at)))").
		Scan(&avgDuration)

	// Get top countries
	topCountries := make(map[string]int)
	var countryStats []struct {
		CountryCode string
		Count       int
	}
	t.db.Raw(`
		SELECT country_code, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at < ?
		AND country_code IS NOT NULL
		GROUP BY country_code
		ORDER BY count DESC
		LIMIT 10
	`, dayStart, dayEnd).Scan(&countryStats)

	for _, stat := range countryStats {
		topCountries[stat.CountryCode] = stat.Count
	}

	// Get top pages
	topPages := make(map[string]int)
	var pageStats []struct {
		Path  string
		Count int
	}
	t.db.Raw(`
		SELECT path, COUNT(*) as count
		FROM page_views pv
		JOIN visitor_sessions vs ON vs.id = pv.session_id
		WHERE vs.created_at >= ? AND vs.created_at < ?
		GROUP BY path
		ORDER BY count DESC
		LIMIT 10
	`, dayStart, dayEnd).Scan(&pageStats)

	for _, stat := range pageStats {
		topPages[stat.Path] = stat.Count
	}

	// Get device breakdown
	deviceBreakdown := make(map[string]int)
	var deviceStats []struct {
		DeviceType string
		Count      int
	}
	t.db.Raw(`
		SELECT device_type, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at < ?
		AND device_type IS NOT NULL
		GROUP BY device_type
	`, dayStart, dayEnd).Scan(&deviceStats)

	for _, stat := range deviceStats {
		deviceBreakdown[stat.DeviceType] = stat.Count
	}

	// Get browser breakdown
	browserBreakdown := make(map[string]int)
	var browserStats []struct {
		BrowserFamily string
		Count         int
	}
	t.db.Raw(`
		SELECT browser_family, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at < ?
		AND browser_family IS NOT NULL
		GROUP BY browser_family
		ORDER BY count DESC
		LIMIT 10
	`, dayStart, dayEnd).Scan(&browserStats)

	for _, stat := range browserStats {
		browserBreakdown[stat.BrowserFamily] = stat.Count
	}

	// Create daily summary
	summary := visitor.VisitorDailySummary{
		SummaryDate:        dayStart,
		UniqueVisitors:     int(uniqueVisitors),
		TotalSessions:      int(totalSessions),
		TotalPageViews:     int(totalPageViews),
		AvgPagesPerSession: avgPagesPerSession,
		AvgSessionDuration: int(avgDuration),
		TopCountries:       topCountries,
		TopPages:           topPages,
		DeviceBreakdown:    deviceBreakdown,
		BrowserBreakdown:   browserBreakdown,
		CreatedAt:          time.Now(),
	}

	// Upsert the summary
	return t.db.Where("summary_date = ?", dayStart).
		Assign(summary).
		FirstOrCreate(&summary).Error
}

// CleanupExpiredSessions removes expired visitor sessions
func (t *VisitorMetricsTask) CleanupExpiredSessions(ctx context.Context) error {
	// Delete expired sessions
	result := t.db.Where("expires_at < ?", time.Now()).Delete(&visitor.VisitorSession{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", result.Error)
	}

	// Delete old real-time data
	result = t.db.Where("created_at < ?", time.Now().Add(-24*time.Hour)).Delete(&visitor.VisitorRealtime{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete old realtime data: %w", result.Error)
	}

	// Delete expired consents
	result = t.db.Where("expires_at < ?", time.Now()).Delete(&visitor.PrivacyConsent{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired consents: %w", result.Error)
	}

	return nil
}

// UpdateLocationAggregates updates visitor location aggregates
func (t *VisitorMetricsTask) UpdateLocationAggregates(ctx context.Context) error {
	// Aggregate location data for today
	today := time.Now().Truncate(24 * time.Hour)
	
	var locationStats []struct {
		CountryCode string
		Count       int
	}
	
	t.db.Raw(`
		SELECT country_code, COUNT(DISTINCT session_hash) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at < ?
		AND country_code IS NOT NULL
		GROUP BY country_code
	`, today, today.Add(24*time.Hour)).Scan(&locationStats)

	// Update location aggregates
	for _, stat := range locationStats {
		location := visitor.VisitorLocation{
			Date:         today,
			CountryCode:  stat.CountryCode,
			VisitorCount: stat.Count,
			CreatedAt:    time.Now(),
		}

		if err := t.db.Where("date = ? AND country_code = ?", today, stat.CountryCode).
			Assign(location).
			FirstOrCreate(&location).Error; err != nil {
			return fmt.Errorf("failed to update location aggregate for %s: %w", stat.CountryCode, err)
		}
	}

	return nil
}

// CleanupOldMetrics removes metrics older than retention period
func (t *VisitorMetricsTask) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	// Delete old hourly metrics
	if err := t.db.Where("metric_date < ?", cutoffDate).Delete(&visitor.VisitorMetrics{}).Error; err != nil {
		return fmt.Errorf("failed to delete old hourly metrics: %w", err)
	}

	// Delete old daily summaries (keep longer)
	if err := t.db.Where("summary_date < ?", cutoffDate.AddDate(0, -6, 0)).Delete(&visitor.VisitorDailySummary{}).Error; err != nil {
		return fmt.Errorf("failed to delete old daily summaries: %w", err)
	}

	// Delete old location data
	if err := t.db.Where("date < ?", cutoffDate).Delete(&visitor.VisitorLocation{}).Error; err != nil {
		return fmt.Errorf("failed to delete old location data: %w", err)
	}

	return nil
}