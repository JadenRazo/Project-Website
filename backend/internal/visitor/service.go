package visitor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
)

// Service handles visitor tracking operations
type Service struct {
	*core.BaseService
	db      *gorm.DB
	cache   *cache.SecureCache
	metrics *metrics.Manager
	config  Config
	hub     *Hub
}

// Config holds visitor service configuration
type Config struct {
	EnableTracking     bool
	SessionTimeout     time.Duration
	RealtimeTimeout    time.Duration
	MaxPageViewsPerSession int
	EnableBotDetection bool
	PrivacyMode        string // "strict", "balanced", "minimal"
}

// NewService creates a new visitor service
func NewService(db *gorm.DB, cache *cache.SecureCache, metrics *metrics.Manager, config Config) *Service {
	hub := NewHub()
	go hub.Run()
	return &Service{
		BaseService: core.NewBaseService("visitor"),
		db:          db,
		cache:       cache,
		metrics:     metrics,
		config:      config,
		hub:         hub,
	}
}

// TrackPageView tracks a page view for the current session
func (s *Service) TrackPageView(ctx context.Context, r *http.Request, path string) error {
	if !s.config.EnableTracking {
		return nil
	}

	// Get or create session
	session, isNew, err := s.getOrCreateSession(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to get/create session: %w", err)
	}

	// Check consent
	consent, err := s.GetConsentStatus(ctx, session.SessionHash)
	if err != nil {
		// If no consent record, assume minimal tracking
		consent = &ConsentStatus{Analytics: false}
	}

	if !consent.Analytics && s.config.PrivacyMode == "strict" {
		return nil // Don't track without consent in strict mode
	}

	// Create page view record
	pageView := &PageView{
		SessionID:      session.ID,
		Path:           path,
		ReferrerDomain: s.extractReferrerDomain(r.Header.Get("Referer")),
		CreatedAt:      time.Now(),
	}

	if err := s.db.Create(pageView).Error; err != nil {
		return fmt.Errorf("failed to create page view: %w", err)
	}

	// Update session last seen
	session.LastSeenAt = time.Now()
	if err := s.db.Save(session).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	// Update real-time tracking
	if err := s.updateRealtimeTracking(ctx, session.SessionHash, path); err != nil {
		// Log but don't fail the request
		fmt.Printf("Failed to update realtime tracking: %v\n", err)
	}

	// Broadcast page view event
	go func() {
		realtimeCount := s.GetRealTimeCount(ctx)
		message, err := json.Marshal(map[string]interface{}{
			"type":      "pageview",
			"path":      path,
			"realtime":  realtimeCount,
		})
		if err == nil {
			s.hub.broadcast <- message
		}
	}()

	// Record metrics
	if s.metrics != nil {
		s.metrics.RecordLatency(0, "/visitor/track")
		if isNew {
			s.recordNewVisitorMetrics()
		}
	}

	return nil
}

// getOrCreateSession gets existing session or creates new one
func (s *Service) getOrCreateSession(ctx context.Context, r *http.Request) (*VisitorSession, bool, error) {
	sessionHash := s.generateSessionHash(r)

	// Try to get from cache first if cache is available
	cacheKey := fmt.Sprintf("visitor:session:%s", sessionHash)
	var session VisitorSession

	if s.cache != nil {
		if err := s.cache.GetDecrypted(ctx, cacheKey, &session); err == nil {
			return &session, false, nil
		}
	}

	// Check database
	if err := s.db.Where("session_hash = ? AND expires_at > ?", sessionHash, time.Now()).First(&session).Error; err == nil {
		// Cache for quick access if cache is available
		if s.cache != nil {
			s.cache.SetEncrypted(ctx, cacheKey, session, 5*time.Minute)
		}
		return &session, false, nil
	}

	// Create new session
	ua := user_agent.New(r.UserAgent())
	browserName, _ := ua.Browser()
	
	session = VisitorSession{
		SessionHash:   sessionHash,
		Language:      s.extractLanguage(r.Header.Get("Accept-Language")),
		DeviceType:    s.getDeviceType(ua, r.UserAgent()),
		BrowserFamily: browserName,
		OSFamily:      ua.OS(),
		IsBot:         s.config.EnableBotDetection && ua.Bot(),
		CreatedAt:     time.Now(),
		LastSeenAt:    time.Now(),
		ExpiresAt:     time.Now().Add(s.config.SessionTimeout),
	}

	// Get location from IP (privacy-compliant)
	if location := s.getLocationFromIP(r); location != nil {
		session.CountryCode = location.CountryCode
		session.Timezone = location.Timezone
		
		// Only include detailed location with consent
		if consent, _ := s.GetConsentStatus(ctx, sessionHash); consent != nil && consent.Analytics {
			session.Region = location.Region
			session.City = location.City
		}
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, false, fmt.Errorf("failed to create session: %w", err)
	}

	// Cache the new session if cache is available
	if s.cache != nil {
		s.cache.SetEncrypted(ctx, cacheKey, session, 5*time.Minute)
	}

	return &session, true, nil
}

// generateSessionHash generates a privacy-safe session hash
func (s *Service) generateSessionHash(r *http.Request) string {
	// Use non-PII data for session identification
	ua := r.UserAgent()
	acceptLang := r.Header.Get("Accept-Language")
	acceptEnc := r.Header.Get("Accept-Encoding")
	
	// Create hash from browser characteristics (no IP)
	data := fmt.Sprintf("%s|%s|%s|%d", ua, acceptLang, acceptEnc, time.Now().Unix()/3600)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// getLocationFromIP gets location from IP without storing the IP
func (s *Service) getLocationFromIP(r *http.Request) *LocationInfo {
	// Extract IP
	ip := s.extractIP(r)
	if ip == "" {
		return nil
	}

	// Use IP geolocation service (implement your preferred service)
	// This is a placeholder - integrate with MaxMind, IP2Location, etc.
	// Important: Do NOT store the IP address
	location := s.geolocateIP(ip)
	
	return location
}

// extractIP extracts client IP from request
func (s *Service) extractIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	
	return host
}

// LocationInfo represents geolocation data
type LocationInfo struct {
	CountryCode string
	Region      string
	City        string
	Timezone    string
}

// geolocateIP performs IP geolocation
func (s *Service) geolocateIP(ip string) *LocationInfo {
	// Use ip-api.com for geolocation
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var geoData struct {
		CountryCode string `json:"countryCode"`
		RegionName  string `json:"regionName"`
		City        string `json:"city"`
		Timezone    string `json:"timezone"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return nil
	}

	if geoData.Status != "success" {
		return nil
	}

	return &LocationInfo{
		CountryCode: geoData.CountryCode,
		Region:      geoData.RegionName,
		City:        geoData.City,
		Timezone:    geoData.Timezone,
	}
}

// extractLanguage extracts primary language from Accept-Language header
func (s *Service) extractLanguage(acceptLang string) string {
	if acceptLang == "" {
		return "en"
	}
	
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		lang := strings.TrimSpace(parts[0])
		if idx := strings.Index(lang, "-"); idx > 0 {
			return lang[:idx]
		}
		return lang
	}
	
	return "en"
}

// getDeviceType determines device type from user agent
func (s *Service) getDeviceType(ua *user_agent.UserAgent, userAgentString string) string {
	if ua.Mobile() {
		return "mobile"
	}
	// Simple tablet detection
	if strings.Contains(strings.ToLower(userAgentString), "tablet") || 
	   strings.Contains(strings.ToLower(userAgentString), "ipad") {
		return "tablet"
	}
	if ua.Bot() {
		return "other"
	}
	return "desktop"
}

// extractReferrerDomain extracts domain from referrer URL
func (s *Service) extractReferrerDomain(referrer string) string {
	if referrer == "" {
		return ""
	}
	
	// Parse URL to extract domain
	parts := strings.Split(referrer, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	
	return ""
}

// updateRealtimeTracking updates real-time visitor tracking
func (s *Service) updateRealtimeTracking(ctx context.Context, sessionHash, currentPage string) error {
	realtime := &VisitorRealtime{
		SessionHash:  sessionHash,
		LastActivity: time.Now(),
		CurrentPage:  currentPage,
		CreatedAt:    time.Now(),
	}
	
	// Upsert realtime tracking
	return s.db.Where("session_hash = ?", sessionHash).
		Assign(VisitorRealtime{
			LastActivity: time.Now(),
			CurrentPage:  currentPage,
		}).
		FirstOrCreate(realtime).Error
}

// recordNewVisitorMetrics records metrics for new visitors
func (s *Service) recordNewVisitorMetrics() {
	if s.metrics != nil {
		// Update Prometheus metrics safely
		if metrics.CurrentVisitors != nil {
			metrics.CurrentVisitors.Inc()
		}
		if metrics.VisitorsByPeriod != nil {
			metrics.VisitorsByPeriod.WithLabelValues("today").Inc()
		}
	}
}

// GetVisitorStats returns comprehensive visitor statistics
func (s *Service) GetVisitorStats(ctx context.Context) (*VisitorStats, error) {
	// Try cache first if available
	cacheKey := "visitor:stats:complete"
	var stats VisitorStats

	if s.cache != nil {
		if err := s.cache.GetDecrypted(ctx, cacheKey, &stats); err == nil {
			return &stats, nil
		}
	}
	
	// Calculate stats
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	stats = VisitorStats{
		Today:         s.getMetricsForPeriod(ctx, todayStart, now),
		Last7Days:     s.getMetricsForPeriod(ctx, now.AddDate(0, 0, -7), now),
		Last30Days:    s.getMetricsForPeriod(ctx, now.AddDate(0, 0, -30), now),
		AllTime:       s.getAllTimeMetrics(ctx),
		RealTimeCount: s.GetRealTimeCount(ctx),
		TrendComparison: s.calculateTrends(ctx),
	}

	// Cache results if cache is available
	if s.cache != nil {
		s.cache.SetEncrypted(ctx, cacheKey, stats, 1*time.Minute)
	}

	return &stats, nil
}

// getMetricsForPeriod calculates metrics for a specific time period
func (s *Service) getMetricsForPeriod(ctx context.Context, start, end time.Time) MetricsSummary {
	var summary MetricsSummary
	
	// Get unique visitors
	var uniqueVisitors int64
	s.db.Model(&VisitorSession{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&uniqueVisitors)
	summary.UniqueVisitors = int(uniqueVisitors)
	
	// Get total page views
	var totalPageViews int64
	s.db.Model(&PageView{}).
		Joins("JOIN visitor_sessions ON visitor_sessions.id = page_views.session_id").
		Where("page_views.created_at BETWEEN ? AND ?", start, end).
		Count(&totalPageViews)
	summary.TotalPageViews = int(totalPageViews)
	
	// Calculate average session duration
	var avgDuration float64
	s.db.Model(&VisitorSession{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Select("AVG(EXTRACT(EPOCH FROM (last_seen_at - created_at)))").
		Scan(&avgDuration)
	summary.AvgSessionDuration = int(avgDuration)
	
	// Calculate bounce rate (sessions with only 1 page view AND short duration)
	var totalSessions, bouncedSessions int64
	s.db.Model(&VisitorSession{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Where("is_bot = ?", false).
		Count(&totalSessions)

	s.db.Raw(`
		SELECT COUNT(*) FROM visitor_sessions vs
		WHERE vs.created_at BETWEEN ? AND ?
		AND vs.is_bot = false
		AND (
			SELECT COUNT(*) FROM page_views pv WHERE pv.session_id = vs.id
		) <= 1
		AND EXTRACT(EPOCH FROM (vs.last_seen_at - vs.created_at)) < 10
	`, start, end).Scan(&bouncedSessions)

	if totalSessions > 0 {
		summary.BounceRate = float64(bouncedSessions) / float64(totalSessions) * 100
	} else {
		summary.BounceRate = 0
	}
	
	// Get new vs returning visitors
	summary.NewVisitors = summary.UniqueVisitors // Simplified for now
	
	return summary
}

// getAllTimeMetrics gets all-time metrics
func (s *Service) getAllTimeMetrics(ctx context.Context) MetricsSummary {
	var summary MetricsSummary
	
	// Use aggregated data for better performance
	var uniqueVisitors int64
	s.db.Model(&VisitorSession{}).Count(&uniqueVisitors)
	summary.UniqueVisitors = int(uniqueVisitors)

	var totalPageViews int64
	s.db.Model(&PageView{}).Count(&totalPageViews)
	summary.TotalPageViews = int(totalPageViews)
	
	return summary
}

// GetRealTimeCount gets current active visitor count
func (s *Service) GetRealTimeCount(ctx context.Context) int {
	var count int64
	s.db.Model(&VisitorRealtime{}).
		Where("last_activity > ?", time.Now().Add(-5*time.Minute)).
		Count(&count)
	
	return int(count)
}

// calculateTrends calculates trend percentages
func (s *Service) calculateTrends(ctx context.Context) TrendData {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// Today vs Yesterday
	todayVisitors := s.getVisitorCount(ctx, todayStart, now)
	yesterdayVisitors := s.getVisitorCount(ctx, yesterdayStart, todayStart)

	// This week vs last week
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))
	lastWeekStart := weekStart.AddDate(0, 0, -7)
	thisWeekVisitors := s.getVisitorCount(ctx, weekStart, now)
	lastWeekVisitors := s.getVisitorCount(ctx, lastWeekStart, weekStart)

	// This month vs last month
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	thisMonthVisitors := s.getVisitorCount(ctx, monthStart, now)
	lastMonthVisitors := s.getVisitorCount(ctx, lastMonthStart, monthStart)

	trends := TrendData{
		TodayVsYesterday:     s.calculateTrendPercentage(todayVisitors, yesterdayVisitors),
		ThisWeekVsLastWeek:   s.calculateTrendPercentage(thisWeekVisitors, lastWeekVisitors),
		ThisMonthVsLastMonth: s.calculateTrendPercentage(thisMonthVisitors, lastMonthVisitors),
	}

	return trends
}

// getVisitorCount gets visitor count for a period
func (s *Service) getVisitorCount(ctx context.Context, start, end time.Time) int {
	var count int64
	s.db.Model(&VisitorSession{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count)
	return int(count)
}

// calculateTrendPercentage calculates percentage change
func (s *Service) calculateTrendPercentage(current, previous int) string {
	if previous == 0 {
		if current > 0 {
			return "+100%"
		}
		return "0%"
	}
	
	change := float64(current-previous) / float64(previous) * 100
	if change >= 0 {
		return fmt.Sprintf("+%.1f%%", change)
	}
	return fmt.Sprintf("%.1f%%", change)
}

// RecordConsent records user privacy consent
func (s *Service) RecordConsent(ctx context.Context, sessionHash string, consentType string, granted bool) error {
	consent := &PrivacyConsent{
		SessionHash: sessionHash,
		ConsentType: consentType,
		Granted:     granted,
		CreatedAt:   time.Now(),
	}
	
	return s.db.Create(consent).Error
}

// GetConsentStatus gets current consent status for a session
func (s *Service) GetConsentStatus(ctx context.Context, sessionHash string) (*ConsentStatus, error) {
	var consents []PrivacyConsent
	
	if err := s.db.Where("session_hash = ? AND expires_at > ?", sessionHash, time.Now()).
		Find(&consents).Error; err != nil {
		return nil, err
	}
	
	status := &ConsentStatus{}
	for _, consent := range consents {
		switch consent.ConsentType {
		case "analytics":
			status.Analytics = consent.Granted
		case "functional":
			status.Functional = consent.Granted
		case "marketing":
			status.Marketing = consent.Granted
		}
	}
	
	return status, nil
}

// GetTimelineData returns visitor timeline data
func (s *Service) GetTimelineData(ctx context.Context, period string, interval string) ([]TimelineData, error) {
	var data []TimelineData

	// Try cache first if available
	cacheKey := fmt.Sprintf("visitor:timeline:%s:%s", period, interval)
	if s.cache != nil {
		if err := s.cache.GetDecrypted(ctx, cacheKey, &data); err == nil {
			return data, nil
		}
	}

	// Parse period to determine date range
	now := time.Now()
	var startDate time.Time
	switch period {
	case "1d":
		startDate = now.AddDate(0, 0, -1)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "1y":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, 0, -7) // Default to 7 days
	}

	// Query data based on interval
	query := `
		SELECT
			DATE_TRUNC($1, pv.created_at) as timestamp,
			COUNT(DISTINCT vs.session_hash) as visitors,
			COUNT(pv.id) as page_views,
			AVG(EXTRACT(EPOCH FROM (vs.last_seen_at - vs.created_at))) as avg_session_time
		FROM visitor_sessions vs
		LEFT JOIN page_views pv ON vs.id = pv.session_id
		WHERE vs.created_at >= $2 AND vs.created_at <= $3
		GROUP BY timestamp
		ORDER BY timestamp ASC
	`

	intervalUnit := "hour"
	if period == "30d" || period == "1y" {
		intervalUnit = "day"
	}

	rows, err := s.db.Raw(query, intervalUnit, startDate, now).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query timeline data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var td TimelineData
		if err := rows.Scan(&td.Timestamp, &td.Visitors, &td.PageViews, &td.AvgSessionTime); err != nil {
			continue
		}
		data = append(data, td)
	}

	// Cache results if cache is available
	if s.cache != nil {
		s.cache.SetEncrypted(ctx, cacheKey, data, 5*time.Minute)
	}

	return data, nil
}

// GetLocationDistribution returns visitor distribution by location
func (s *Service) GetLocationDistribution(ctx context.Context, period string) ([]LocationData, error) {
	var locations []LocationData

	// Parse period to determine date range
	now := time.Now()
	var startDate time.Time
	switch period {
	case "1d":
		startDate = now.AddDate(0, 0, -1)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "1y":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, 0, -30) // Default to 30 days
	}

	// Query location distribution
	query := `
		SELECT
			country_code,
			COUNT(*) as visitor_count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at <= ?
			AND country_code IS NOT NULL
			AND country_code != ''
		GROUP BY country_code
		ORDER BY visitor_count DESC
		LIMIT 20
	`

	rows, err := s.db.Raw(query, startDate, now).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query location distribution: %w", err)
	}
	defer rows.Close()

	var totalVisitors int
	tempLocations := []struct {
		CountryCode  string
		VisitorCount int
	}{}

	// First pass to get data and total
	for rows.Next() {
		var loc struct {
			CountryCode  string
			VisitorCount int
		}
		if err := rows.Scan(&loc.CountryCode, &loc.VisitorCount); err != nil {
			continue
		}
		tempLocations = append(tempLocations, loc)
		totalVisitors += loc.VisitorCount
	}

	// Second pass to calculate percentages and country names
	countryNames := map[string]string{
		"US": "United States",
		"GB": "United Kingdom",
		"CA": "Canada",
		"AU": "Australia",
		"DE": "Germany",
		"FR": "France",
		"JP": "Japan",
		"CN": "China",
		"IN": "India",
		"BR": "Brazil",
		"MX": "Mexico",
		"ES": "Spain",
		"IT": "Italy",
		"NL": "Netherlands",
		"SE": "Sweden",
		"NO": "Norway",
		"DK": "Denmark",
		"FI": "Finland",
		"PL": "Poland",
		"RU": "Russia",
		// Add more as needed
	}

	for _, loc := range tempLocations {
		countryName := countryNames[loc.CountryCode]
		if countryName == "" {
			countryName = loc.CountryCode // Fallback to code if name not found
		}

		percentage := 0.0
		if totalVisitors > 0 {
			percentage = (float64(loc.VisitorCount) / float64(totalVisitors)) * 100
		}

		locations = append(locations, LocationData{
			CountryCode:  loc.CountryCode,
			CountryName:  countryName,
			VisitorCount: loc.VisitorCount,
			Percentage:   percentage,
		})
	}

	return locations, nil
}

// GetDeviceBrowserStats returns device and browser statistics
func (s *Service) GetDeviceBrowserStats(ctx context.Context, period string) (*DeviceBrowserData, error) {
	stats := &DeviceBrowserData{
		Devices:  make(map[string]int),
		Browsers: make(map[string]int),
		OS:       make(map[string]int),
	}

	// Parse period to determine date range
	now := time.Now()
	var startDate time.Time
	switch period {
	case "1d":
		startDate = now.AddDate(0, 0, -1)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "1y":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, 0, -30) // Default to 30 days
	}

	// Query device types
	deviceQuery := `
		SELECT device_type, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at <= ?
			AND device_type IS NOT NULL
		GROUP BY device_type
	`
	rows, err := s.db.Raw(deviceQuery, startDate, now).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var deviceType string
			var count int
			if err := rows.Scan(&deviceType, &count); err == nil {
				stats.Devices[deviceType] = count
			}
		}
	}

	// Query browsers
	browserQuery := `
		SELECT browser_family, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at <= ?
			AND browser_family IS NOT NULL
			AND browser_family != ''
		GROUP BY browser_family
		ORDER BY count DESC
		LIMIT 10
	`
	rows, err = s.db.Raw(browserQuery, startDate, now).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var browser string
			var count int
			if err := rows.Scan(&browser, &count); err == nil {
				stats.Browsers[browser] = count
			}
		}
	}

	// Query OS
	osQuery := `
		SELECT os_family, COUNT(*) as count
		FROM visitor_sessions
		WHERE created_at >= ? AND created_at <= ?
			AND os_family IS NOT NULL
			AND os_family != ''
		GROUP BY os_family
		ORDER BY count DESC
		LIMIT 10
	`
	rows, err = s.db.Raw(osQuery, startDate, now).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var os string
			var count int
			if err := rows.Scan(&os, &count); err == nil {
				stats.OS[os] = count
			}
		}
	}

	return stats, nil
}

// GetHub returns the WebSocket hub for external registration
func (s *Service) GetHub() *Hub {
	return s.hub
}