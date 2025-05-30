package status

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/jadenrazo/project-website/backend/internal/common/metrics"
)

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"` // operational, degraded, down
	Latency     int64     `json:"latency_ms"`
	LastChecked time.Time `json:"last_checked"`
	Error       string    `json:"error,omitempty"`
	Uptime      float64   `json:"uptime_percentage"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// SystemStatus represents the overall system status
type SystemStatus struct {
	Status      string           `json:"status"` // operational, partial_outage, major_outage
	Services    []ServiceStatus  `json:"services"`
	LastUpdated time.Time        `json:"last_updated"`
	Incidents   []Incident       `json:"incidents,omitempty"`
}

// Incident represents a service incident
type Incident struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Service     string    `json:"service"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // investigating, identified, monitoring, resolved
	Severity    string    `json:"severity"` // minor, major, critical
	StartedAt   time.Time `json:"started_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StatusHistory tracks service status over time
type StatusHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// Service handles status monitoring
type Service struct {
	db              *gorm.DB
	mu              sync.RWMutex
	serviceStatuses map[string]*ServiceStatus
	checkInterval   time.Duration
	checks          map[string]HealthCheck
	metricsManager  *metrics.Manager
	metricsHandler  *metrics.SystemMetricsHandler
}

// HealthCheck is a function that checks the health of a service
type HealthCheck func(ctx context.Context) error

// NewService creates a new status service
func NewService(db *gorm.DB, checkInterval time.Duration, metricsManager *metrics.Manager) *Service {
	service := &Service{
		db:              db,
		serviceStatuses: make(map[string]*ServiceStatus),
		checkInterval:   checkInterval,
		checks:          make(map[string]HealthCheck),
		metricsManager:  metricsManager,
		metricsHandler:  metrics.NewSystemMetricsHandler(metricsManager),
	}

	// Auto-migrate status tables
	db.AutoMigrate(&Incident{}, &StatusHistory{})

	// Register default health checks
	service.registerDefaultChecks()

	// Start monitoring
	go service.startMonitoring()

	// Start periodic latency collection
	if metricsManager != nil {
		metricsManager.StartPeriodicLatencyCollector(1*time.Minute, "http://localhost:8080/api/status")
	}

	return service
}

// RegisterHealthCheck registers a health check for a service
func (s *Service) RegisterHealthCheck(name string, check HealthCheck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[name] = check
}

// registerDefaultChecks registers default health checks
func (s *Service) registerDefaultChecks() {
	// Database check
	s.RegisterHealthCheck("database", func(ctx context.Context) error {
		db, err := s.db.DB()
		if err != nil {
			return err
		}
		return db.PingContext(ctx)
	})

	// API check (self-ping)
	s.RegisterHealthCheck("api", func(ctx context.Context) error {
		// This is a simple check that the API service itself is running
		return nil
	})

	// Code Stats check
	s.RegisterHealthCheck("code_stats", func(ctx context.Context) error {
		var count int64
		return s.db.WithContext(ctx).Model(&StatusHistory{}).Count(&count).Error
	})
}

// startMonitoring starts the background monitoring
func (s *Service) startMonitoring() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	// Initial check
	s.performChecks()

	for range ticker.C {
		s.performChecks()
	}
}

// performChecks performs all registered health checks
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
				errorMsg = "Health check timed out"
			} else {
				status = "degraded"
				errorMsg = err.Error()
			}
		}

		s.updateServiceStatus(name, status, latency, errorMsg)
		
		// Record to history
		s.recordHistory(name, status, latency, errorMsg)
		
		// Record latency metrics for successful checks
		if s.metricsManager != nil && status == "operational" {
			s.metricsManager.RecordLatency(float64(latency), name)
		}
	}
}

// updateServiceStatus updates the status of a service
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

	// Calculate uptime (last 24 hours)
	s.serviceStatuses[name].Uptime = s.calculateUptime(name)
}

// calculateUptime calculates the uptime percentage for the last 24 hours
func (s *Service) calculateUptime(service string) float64 {
	var totalChecks, successfulChecks int64
	
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	
	s.db.Model(&StatusHistory{}).
		Where("service = ? AND checked_at > ?", service, twentyFourHoursAgo).
		Count(&totalChecks)
	
	s.db.Model(&StatusHistory{}).
		Where("service = ? AND checked_at > ? AND status = ?", service, twentyFourHoursAgo, "operational").
		Count(&successfulChecks)
	
	if totalChecks == 0 {
		return 100.0
	}
	
	return float64(successfulChecks) / float64(totalChecks) * 100.0
}

// recordHistory records a status check to history
func (s *Service) recordHistory(service, status string, latency int64, errorMsg string) {
	history := &StatusHistory{
		Service:   service,
		Status:    status,
		Latency:   latency,
		Error:     errorMsg,
		CheckedAt: time.Now(),
	}
	s.db.Create(history)
}

// GetSystemStatus returns the overall system status
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

	// Determine overall status
	overallStatus := "operational"
	if operationalCount == 0 {
		overallStatus = "major_outage"
	} else if operationalCount < len(services) {
		overallStatus = "partial_outage"
	}

	// Get active incidents
	var incidents []Incident
	s.db.Where("resolved_at IS NULL").Order("started_at DESC").Find(&incidents)

	return &SystemStatus{
		Status:      overallStatus,
		Services:    services,
		LastUpdated: time.Now(),
		Incidents:   incidents,
	}
}

// GetServiceStatus returns the status of a specific service
func (s *Service) GetServiceStatus(name string) (*ServiceStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.serviceStatuses[name]
	if !exists {
		return nil, false
	}

	return status, true
}

// GetStatusHistory returns the status history for a service
func (s *Service) GetStatusHistory(service string, duration time.Duration) ([]StatusHistory, error) {
	var history []StatusHistory
	since := time.Now().Add(-duration)
	
	err := s.db.Where("service = ? AND checked_at > ?", service, since).
		Order("checked_at DESC").
		Find(&history).Error
		
	return history, err
}

// CreateIncident creates a new incident
func (s *Service) CreateIncident(incident *Incident) error {
	incident.StartedAt = time.Now()
	incident.CreatedAt = time.Now()
	incident.UpdatedAt = time.Now()
	return s.db.Create(incident).Error
}

// UpdateIncident updates an existing incident
func (s *Service) UpdateIncident(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.db.Model(&Incident{}).Where("id = ?", id).Updates(updates).Error
}

// ResolveIncident resolves an incident
func (s *Service) ResolveIncident(id uint) error {
	now := time.Now()
	return s.db.Model(&Incident{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "resolved",
		"resolved_at": &now,
		"updated_at":  now,
	}).Error
}

// RegisterRoutes registers HTTP routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", s.handleGetSystemStatus)
	router.GET("/services/:name", s.handleGetServiceStatus)
	router.GET("/services/:name/history", s.handleGetServiceHistory)
	router.GET("/incidents", s.handleGetIncidents)
	router.POST("/incidents", s.handleCreateIncident)
	router.PUT("/incidents/:id", s.handleUpdateIncident)
	router.POST("/incidents/:id/resolve", s.handleResolveIncident)
	
	// Metrics routes
	router.GET("/metrics/day", s.handleMetricsDay)
	router.GET("/metrics/week", s.handleMetricsWeek)
	router.GET("/metrics/month", s.handleMetricsMonth)
	router.GET("/metrics/latest", s.handleLatestMetrics)
}

// HTTP handlers
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
		c.JSON(500, gin.H{"error": "Failed to get history"})
		return
	}
	
	c.JSON(200, history)
}

func (s *Service) handleGetIncidents(c *gin.Context) {
	var incidents []Incident
	query := s.db.Order("started_at DESC")
	
	if active := c.Query("active"); active == "true" {
		query = query.Where("resolved_at IS NULL")
	}
	
	if err := query.Find(&incidents).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get incidents"})
		return
	}
	
	c.JSON(200, incidents)
}

func (s *Service) handleCreateIncident(c *gin.Context) {
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	if err := s.CreateIncident(&incident); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create incident"})
		return
	}
	
	c.JSON(201, incident)
}

func (s *Service) handleUpdateIncident(c *gin.Context) {
	id := c.Param("id")
	
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var incidentID uint
	if _, err := fmt.Sscanf(id, "%d", &incidentID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid incident ID"})
		return
	}
	
	if err := s.UpdateIncident(incidentID, updates); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update incident"})
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
		c.JSON(500, gin.H{"error": "Failed to resolve incident"})
		return
	}
	
	c.JSON(200, gin.H{"message": "Incident resolved"})
}

// Metrics handlers
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