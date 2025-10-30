package visitor

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VisitorSession represents an anonymous visitor session
type VisitorSession struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionHash   string    `gorm:"type:varchar(64);not null;index"`
	CountryCode   string    `gorm:"type:varchar(2)"`
	Region        string    `gorm:"type:varchar(100)"`
	City          string    `gorm:"type:varchar(100)"`
	Timezone      string    `gorm:"type:varchar(50)"`
	Language      string    `gorm:"type:varchar(10)"`
	DeviceType    string    `gorm:"type:varchar(50);check:device_type IN ('desktop','mobile','tablet','other')"`
	BrowserFamily string    `gorm:"type:varchar(50)"`
	OSFamily      string    `gorm:"type:varchar(50)"`
	IsBot         bool      `gorm:"default:false"`
	CreatedAt     time.Time
	LastSeenAt    time.Time
	ExpiresAt     time.Time
}

// PageView represents a page view event
type PageView struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Path             string    `gorm:"type:varchar(255);not null;index"`
	ReferrerDomain   string    `gorm:"type:varchar(255)"`
	DurationSeconds  int
	CreatedAt        time.Time
}

// PrivacyConsent represents user consent for data collection
type PrivacyConsent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionHash string    `gorm:"type:varchar(64);not null;index"`
	ConsentType string    `gorm:"type:varchar(50);not null;check:consent_type IN ('analytics','functional','marketing')"`
	Granted     bool      `gorm:"not null"`
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// VisitorMetrics represents aggregated metrics
type VisitorMetrics struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	MetricDate         time.Time `gorm:"type:date;not null"`
	Hour               *int      `gorm:"check:hour >= 0 AND hour <= 23"`
	UniqueVisitors     int       `gorm:"default:0"`
	TotalPageViews     int       `gorm:"default:0"`
	AvgSessionDuration int       `gorm:"default:0"`
	BounceRate         float64   `gorm:"type:decimal(5,2)"`
	NewVisitors        int       `gorm:"default:0"`
	ReturningVisitors  int       `gorm:"default:0"`
	CreatedAt          time.Time
}

// VisitorDailySummary represents daily aggregated data
type VisitorDailySummary struct {
	ID                  uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SummaryDate         time.Time              `gorm:"type:date;not null;unique"`
	UniqueVisitors      int                    `gorm:"default:0"`
	TotalSessions       int                    `gorm:"default:0"`
	TotalPageViews      int                    `gorm:"default:0"`
	AvgPagesPerSession  float64                `gorm:"type:decimal(5,2)"`
	AvgSessionDuration  int
	TopCountries        map[string]int         `gorm:"type:jsonb;default:'{}'"`
	TopPages            map[string]int         `gorm:"type:jsonb;default:'{}'"`
	DeviceBreakdown     map[string]int         `gorm:"type:jsonb;default:'{}'"`
	BrowserBreakdown    map[string]int         `gorm:"type:jsonb;default:'{}'"`
	CreatedAt           time.Time
}

// VisitorRealtime represents real-time visitor tracking
type VisitorRealtime struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionHash  string    `gorm:"type:varchar(64);not null;index"`
	LastActivity time.Time `gorm:"index"`
	CurrentPage  string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
}

// VisitorLocation represents aggregated location data
type VisitorLocation struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Date         time.Time `gorm:"type:date;not null"`
	CountryCode  string    `gorm:"type:varchar(2);not null"`
	VisitorCount int       `gorm:"default:0"`
	CreatedAt    time.Time
}

// TableName overrides for GORM
func (VisitorSession) TableName() string      { return "visitor_sessions" }
func (PageView) TableName() string            { return "page_views" }
func (PrivacyConsent) TableName() string      { return "privacy_consents" }
func (VisitorMetrics) TableName() string      { return "visitor_metrics" }
func (VisitorDailySummary) TableName() string { return "visitor_daily_summary" }
func (VisitorRealtime) TableName() string     { return "visitor_realtime" }
func (VisitorLocation) TableName() string     { return "visitor_locations" }

// VisitorStats represents comprehensive visitor statistics
type VisitorStats struct {
	Today           MetricsSummary   `json:"today"`
	Last7Days       MetricsSummary   `json:"last7Days"`
	Last30Days      MetricsSummary   `json:"last30Days"`
	AllTime         MetricsSummary   `json:"allTime"`
	RealTimeCount   int              `json:"realTimeCount"`
	TrendComparison TrendData        `json:"trends"`
}

// MetricsSummary represents metrics for a time period
type MetricsSummary struct {
	UniqueVisitors     int     `json:"uniqueVisitors"`
	TotalPageViews     int     `json:"totalPageViews"`
	AvgSessionDuration int     `json:"avgSessionDuration"`
	BounceRate         float64 `json:"bounceRate"`
	NewVisitors        int     `json:"newVisitors"`
	ReturningVisitors  int     `json:"returningVisitors"`
}

// TrendData represents trend comparisons
type TrendData struct {
	TodayVsYesterday      string `json:"todayVsYesterday"`
	ThisWeekVsLastWeek    string `json:"thisWeekVsLastWeek"`
	ThisMonthVsLastMonth  string `json:"thisMonthVsLastMonth"`
}

// TimelineData represents time series data
type TimelineData struct {
	Timestamp      time.Time `json:"timestamp"`
	Visitors       int       `json:"visitors"`
	PageViews      int       `json:"pageViews"`
	AvgSessionTime int       `json:"avgSessionTime"`
}

// LocationData represents visitor location distribution
type LocationData struct {
	CountryCode  string `json:"countryCode"`
	CountryName  string `json:"countryName"`
	VisitorCount int    `json:"visitorCount"`
	Percentage   float64 `json:"percentage"`
}

// DeviceBrowserData represents device and browser statistics
type DeviceBrowserData struct {
	Devices  map[string]int `json:"devices"`
	Browsers map[string]int `json:"browsers"`
	OS       map[string]int `json:"os"`
}

// ConsentStatus represents the current consent status
type ConsentStatus struct {
	Analytics  bool `json:"analytics"`
	Functional bool `json:"functional"`
	Marketing  bool `json:"marketing"`
}

// BeforeCreate hooks for UUID generation
func (v *VisitorSession) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.ExpiresAt.IsZero() {
		v.ExpiresAt = time.Now().Add(24 * time.Hour)
	}
	return nil
}

func (p *PageView) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (c *PrivacyConsent) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExpiresAt.IsZero() {
		c.ExpiresAt = time.Now().Add(365 * 24 * time.Hour) // 1 year
	}
	return nil
}