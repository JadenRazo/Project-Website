package status

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceStatus struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Latency     int64                  `json:"latency_ms"`
	LastChecked time.Time              `json:"last_checked"`
	Error       string                 `json:"error,omitempty"`
	Uptime      float64                `json:"uptime_percentage"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type SystemStatus struct {
	Status      string          `json:"status"`
	Services    []ServiceStatus `json:"services"`
	LastUpdated time.Time       `json:"last_updated"`
	Incidents   []Incident      `json:"incidents,omitempty"`
}

type Incident struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Service     string     `json:"service"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	StartedAt   time.Time  `json:"started_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type StatusHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

type Service struct {
	db              *StatusDatabase
	originalDB      *gorm.DB
	mu              sync.RWMutex
	serviceStatuses map[string]*ServiceStatus
	checkInterval   time.Duration
	checks          map[string]HealthCheck
	metricsManager  *metrics.Manager
	metricsHandler  *metrics.SystemMetricsHandler
}

type HealthCheck func(ctx context.Context) error

func NewService(db *gorm.DB, checkInterval time.Duration, metricsManager *metrics.Manager) *Service {
	statusDB := NewStatusDatabase(nil, db)

	service := &Service{
		db:              statusDB,
		originalDB:      db,
		serviceStatuses: make(map[string]*ServiceStatus),
		checkInterval:   checkInterval,
		checks:          make(map[string]HealthCheck),
		metricsManager:  metricsManager,
		metricsHandler:  metrics.NewSystemMetricsHandler(metricsManager),
	}

	service.initializeDatabase()

	service.registerDefaultChecks()

	go service.startMonitoring()

	if metricsManager != nil {
		go service.startPeriodicLatencyCollection()

		metricsManager.StartPeriodicCleanup()
	}

	return service
}

func (s *Service) initializeDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.AutoMigrate(&Incident{}, &StatusHistory{})
	})

	if err != nil {
		fmt.Printf("Warning: Failed to initialize database tables: %v\n", err)
	}
}

func (s *Service) RegisterHealthCheck(name string, check HealthCheck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[name] = check
}

func (s *Service) registerDefaultChecks() {
	s.RegisterHealthCheck("database", func(ctx context.Context) error {
		return s.db.Ping(ctx)
	})

	s.RegisterHealthCheck("api", func(ctx context.Context) error {
		return nil
	})

	s.RegisterHealthCheck("code_stats", func(ctx context.Context) error {
		var count int64
		return s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
			return db.WithContext(ctx).Model(&StatusHistory{}).Count(&count).Error
		})
	})
}

func (s *Service) startMonitoring() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	s.performChecks()

	for range ticker.C {
		s.performChecks()
	}
}

func (s *Service) performChecks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for name, check := range s.checks {
		start := time.Now()
		err := check(ctx)
		latency := time.Since(start).Milliseconds()

		status := "operational"
		errorMsg := ""
		if err != nil {
			if err == context.DeadlineExceeded {
				status = "down"
				errorMsg = SanitizeErrorMessage(err)
			} else {
				status = "degraded"
				errorMsg = SanitizeErrorMessage(err)
			}
		}

		s.updateServiceStatus(name, status, latency, errorMsg)
		s.recordHistory(name, status, latency, errorMsg)

		if s.metricsManager != nil && status == "operational" {
			s.metricsManager.RecordLatency(float64(latency), name)
		}
	}
}

func (s *Service) updateServiceStatus(name, status string, latency int64, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.serviceStatuses[name]; !exists {
		s.serviceStatuses[name] = &ServiceStatus{
			Name: name,
		}
	}

	s.serviceStatuses[name].Status = status
	s.serviceStatuses[name].Latency = latency
	s.serviceStatuses[name].LastChecked = time.Now()
	s.serviceStatuses[name].Error = errorMsg

	s.serviceStatuses[name].Uptime = s.calculateUptime(name)
}

func (s *Service) calculateUptime(service string) float64 {
	if !s.db.IsHealthy() {
		return 100.0
	}

	var totalChecks, successfulChecks int64
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		if err := db.Model(&StatusHistory{}).
			Where("service = ? AND checked_at > ?", service, twentyFourHoursAgo).
			Count(&totalChecks).Error; err != nil {
			return err
		}

		return db.Model(&StatusHistory{}).
			Where("service = ? AND checked_at > ? AND status = ?", service, twentyFourHoursAgo, "operational").
			Count(&successfulChecks).Error
	})

	if err != nil {
		return 100.0
	}

	if totalChecks == 0 {
		return 100.0
	}

	return float64(successfulChecks) / float64(totalChecks) * 100.0
}

func (s *Service) recordHistory(service, status string, latency int64, errorMsg string) {
	history := &StatusHistory{
		Service:   service,
		Status:    status,
		Latency:   latency,
		Error:     errorMsg,
		CheckedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.Create(history).Error
	})

	if err != nil {
		fmt.Printf("Warning: Failed to record status history: %v\n", err)
	}
}

func (s *Service) GetSystemStatus() *SystemStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	services := make([]ServiceStatus, 0, len(s.serviceStatuses))
	operationalCount := 0

	for _, status := range s.serviceStatuses {
		services = append(services, *status)
		if status.Status == "operational" {
			operationalCount++
		}
	}

	overallStatus := "operational"
	if operationalCount == 0 {
		overallStatus = "major_outage"
	} else if operationalCount < len(services) {
		overallStatus = "partial_outage"
	}

	var incidents []Incident
	if s.db.IsHealthy() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
			return db.Where("resolved_at IS NULL").Order("started_at DESC").Find(&incidents).Error
		})

		if err != nil {
			incidents = []Incident{
				{
					Service:     "database",
					Title:       "Database connectivity issues",
					Description: "Unable to retrieve incident data",
					Status:      "investigating",
					Severity:    "major",
					StartedAt:   time.Now(),
				},
			}
		}
	}

	return &SystemStatus{
		Status:      overallStatus,
		Services:    services,
		LastUpdated: time.Now(),
		Incidents:   incidents,
	}
}

func (s *Service) GetServiceStatus(name string) (*ServiceStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.serviceStatuses[name]
	if !exists {
		return nil, false
	}

	return status, true
}

func (s *Service) GetStatusHistory(service string, duration time.Duration) ([]StatusHistory, error) {
	var history []StatusHistory
	since := time.Now().Add(-duration)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.Where("service = ? AND checked_at > ?", service, since).
			Order("checked_at DESC").
			Find(&history).Error
	})

	return history, err
}

func (s *Service) CreateIncident(incident *Incident) error {
	incident.StartedAt = time.Now()
	incident.CreatedAt = time.Now()
	incident.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.Create(incident).Error
	})
}

func (s *Service) UpdateIncident(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.Model(&Incident{}).Where("id = ?", id).Updates(updates).Error
	})
}

func (s *Service) ResolveIncident(id uint) error {
	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		return db.Model(&Incident{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":      "resolved",
			"resolved_at": &now,
			"updated_at":  now,
		}).Error
	})
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", s.handleGetSystemStatus)
	router.GET("/services/:name", s.handleGetServiceStatus)
	router.GET("/services/:name/history", s.handleGetServiceHistory)
	router.GET("/incidents", s.handleGetIncidents)
	router.POST("/incidents", s.handleCreateIncident)
	router.PUT("/incidents/:id", s.handleUpdateIncident)
	router.POST("/incidents/:id/resolve", s.handleResolveIncident)

	router.GET("/metrics/day", s.handleMetricsDay)
	router.GET("/metrics/week", s.handleMetricsWeek)
	router.GET("/metrics/month", s.handleMetricsMonth)
	router.GET("/metrics/latest", s.handleLatestMetrics)
}

func (s *Service) handleGetSystemStatus(c *gin.Context) {
	c.JSON(200, s.GetSystemStatus())
}

func (s *Service) handleGetServiceStatus(c *gin.Context) {
	name := c.Param("name")
	status, exists := s.GetServiceStatus(name)
	if !exists {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(200, status)
}

func (s *Service) handleGetServiceHistory(c *gin.Context) {
	name := c.Param("name")
	durationStr := c.DefaultQuery("duration", "24h")

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid duration"})
		return
	}

	history, err := s.GetStatusHistory(name, duration)
	if err != nil {
		c.JSON(500, gin.H{"error": "Service history temporarily unavailable"})
		return
	}

	c.JSON(200, history)
}

func (s *Service) handleGetIncidents(c *gin.Context) {
	var incidents []Incident

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.db.ExecuteWithRetry(ctx, func(db *gorm.DB) error {
		query := db.Order("started_at DESC")

		if active := c.Query("active"); active == "true" {
			query = query.Where("resolved_at IS NULL")
		}

		return query.Find(&incidents).Error
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Incident data temporarily unavailable"})
		return
	}

	c.JSON(200, incidents)
}

func (s *Service) handleCreateIncident(c *gin.Context) {
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(400, gin.H{"error": "Invalid incident data provided"})
		return
	}

	if err := s.CreateIncident(&incident); err != nil {
		c.JSON(500, gin.H{"error": "Unable to create incident at this time"})
		return
	}

	c.JSON(201, incident)
}

func (s *Service) handleUpdateIncident(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": "Invalid update data provided"})
		return
	}

	var incidentID uint
	if _, err := fmt.Sscanf(id, "%d", &incidentID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid incident ID"})
		return
	}

	if err := s.UpdateIncident(incidentID, updates); err != nil {
		c.JSON(500, gin.H{"error": "Unable to update incident at this time"})
		return
	}

	c.JSON(200, gin.H{"message": "Incident updated"})
}

func (s *Service) handleResolveIncident(c *gin.Context) {
	id := c.Param("id")

	var incidentID uint
	if _, err := fmt.Sscanf(id, "%d", &incidentID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid incident ID"})
		return
	}

	if err := s.ResolveIncident(incidentID); err != nil {
		c.JSON(500, gin.H{"error": "Unable to resolve incident at this time"})
		return
	}

	c.JSON(200, gin.H{"message": "Incident resolved"})
}

func (s *Service) handleMetricsDay(c *gin.Context) {
	if s.metricsHandler == nil {
		c.JSON(503, gin.H{"error": "Metrics not available"})
		return
	}
	s.metricsHandler.HandleMetricsDay(c.Writer, c.Request)
}

func (s *Service) handleMetricsWeek(c *gin.Context) {
	if s.metricsHandler == nil {
		c.JSON(503, gin.H{"error": "Metrics not available"})
		return
	}
	s.metricsHandler.HandleMetricsWeek(c.Writer, c.Request)
}

func (s *Service) handleMetricsMonth(c *gin.Context) {
	if s.metricsHandler == nil {
		c.JSON(503, gin.H{"error": "Metrics not available"})
		return
	}
	s.metricsHandler.HandleMetricsMonth(c.Writer, c.Request)
}

func (s *Service) handleLatestMetrics(c *gin.Context) {
	if s.metricsHandler == nil {
		c.JSON(503, gin.H{"error": "Metrics not available"})
		return
	}
	s.metricsHandler.HandleLatestMetrics(c.Writer, c.Request)
}

func (s *Service) startPeriodicLatencyCollection() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.collectLatencyMetrics()
	}
}

func (s *Service) collectLatencyMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()

	err := s.performAPIHealthCheck(ctx)
	latency := time.Since(start).Seconds() * 1000 // Convert to milliseconds

	if err == nil && s.metricsManager != nil {
		s.metricsManager.RecordLatency(latency, "api_health_check")
	}
}

func (s *Service) performAPIHealthCheck(ctx context.Context) error {
	return s.db.Ping(ctx)
}
