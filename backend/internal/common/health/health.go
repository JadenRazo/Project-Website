package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
)

// Status represents the health of a component
type Status string

const (
	// StatusUp indicates the component is healthy
	StatusUp Status = "UP"
	// StatusDown indicates the component is unhealthy
	StatusDown Status = "DOWN"
	// StatusDegraded indicates the component has reduced functionality
	StatusDegraded Status = "DEGRADED"
)

// Component represents a service component to health check
type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// HealthCheck represents the overall health of the application
type HealthCheck struct {
	Status     Status      `json:"status"`
	Components []Component `json:"components,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Version    string      `json:"version"`
}

// Checker is the interface for health check implementations
type Checker interface {
	Name() string
	Check(ctx context.Context) (Status, string)
}

// Health manages application health checking
type Health struct {
	checkers []Checker
	version  string
	mu       sync.RWMutex
}

// NewHealth creates a new Health instance
func NewHealth(version string) *Health {
	return &Health{
		checkers: make([]Checker, 0),
		version:  version,
	}
}

// AddChecker adds a component health checker
func (h *Health) AddChecker(checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

// CheckHealth performs a health check on all components
func (h *Health) CheckHealth(ctx context.Context) HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := HealthCheck{
		Status:     StatusUp,
		Components: make([]Component, 0, len(h.checkers)),
		Timestamp:  time.Now(),
		Version:    h.version,
	}

	// Check each component in parallel
	var wg sync.WaitGroup
	componentsCh := make(chan Component, len(h.checkers))

	for _, checker := range h.checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()
			status, detail := c.Check(ctx)
			componentsCh <- Component{
				Name:   c.Name(),
				Status: status,
				Detail: detail,
			}
		}(checker)
	}

	// Wait for all checks to complete
	go func() {
		wg.Wait()
		close(componentsCh)
	}()

	// Collect results
	for component := range componentsCh {
		result.Components = append(result.Components, component)
		if component.Status == StatusDown {
			result.Status = StatusDown
		} else if component.Status == StatusDegraded && result.Status != StatusDown {
			result.Status = StatusDegraded
		}
	}

	return result
}

// HTTPHandler creates an HTTP handler for health checks
func (h *Health) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		health := h.CheckHealth(ctx)

		w.Header().Set("Content-Type", "application/json")
		if health.Status != StatusUp {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(health)
	})
}

// LivenessHandler creates an HTTP handler for liveness probes
func (h *Health) LivenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
	})
}

// ReadinessHandler creates an HTTP handler for readiness probes
func (h *Health) ReadinessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		health := h.CheckHealth(ctx)

		w.Header().Set("Content-Type", "application/json")
		if health.Status != StatusUp {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(health)
	})
}

// DatabaseChecker checks database connectivity
type DatabaseChecker struct {
	db database.Database
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(db database.Database) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Name returns the checker name
func (c *DatabaseChecker) Name() string {
	return "database"
}

// Check performs the health check
func (c *DatabaseChecker) Check(ctx context.Context) (Status, string) {
	start := time.Now()
	err := c.db.Ping(ctx)
	duration := time.Since(start)

	if err != nil {
		return StatusDown, fmt.Sprintf("Database connection failed: %v", err)
	}

	if duration > 500*time.Millisecond {
		return StatusDegraded, fmt.Sprintf("Database response time degraded: %v", duration)
	}

	return StatusUp, fmt.Sprintf("Response time: %v", duration)
}

// CacheChecker checks cache connectivity
type CacheChecker struct {
	cache cache.Cache
}

// NewCacheChecker creates a new cache health checker
func NewCacheChecker(cache cache.Cache) *CacheChecker {
	return &CacheChecker{cache: cache}
}

// Name returns the checker name
func (c *CacheChecker) Name() string {
	return "cache"
}

// Check performs the health check
func (c *CacheChecker) Check(ctx context.Context) (Status, string) {
	start := time.Now()
	key := "health-check-probe"
	value := time.Now().String()

	err := c.cache.Set(ctx, key, value, 5*time.Second)
	if err != nil {
		return StatusDown, fmt.Sprintf("Cache write failed: %v", err)
	}

	_, err = c.cache.Get(ctx, key)
	duration := time.Since(start)

	if err != nil {
		return StatusDown, fmt.Sprintf("Cache read failed: %v", err)
	}

	if duration > 200*time.Millisecond {
		return StatusDegraded, fmt.Sprintf("Cache response time degraded: %v", duration)
	}

	return StatusUp, fmt.Sprintf("Response time: %v", duration)
}

// DiskChecker checks disk space
type DiskChecker struct {
	path string
}

// NewDiskChecker creates a new disk health checker
func NewDiskChecker(path string) *DiskChecker {
	return &DiskChecker{path: path}
}

// Name returns the checker name
func (c *DiskChecker) Name() string {
	return "disk"
}

// Check performs the health check
func (c *DiskChecker) Check(ctx context.Context) (Status, string) {
	// Example implementation - in a real system, use os.Stat or similar
	// This is a simplified example
	freeDiskPercentage := 75.0 // Normally would calculate this

	if freeDiskPercentage < 10 {
		return StatusDown, fmt.Sprintf("Critical disk space: %.2f%% free", freeDiskPercentage)
	}

	if freeDiskPercentage < 20 {
		return StatusDegraded, fmt.Sprintf("Low disk space: %.2f%% free", freeDiskPercentage)
	}

	return StatusUp, fmt.Sprintf("%.2f%% free disk space", freeDiskPercentage)
}

// MariaDBChecker checks MariaDB connectivity
type MariaDBChecker struct {
	pinger func(ctx context.Context) error
}

// NewMariaDBChecker creates a new MariaDB health checker
func NewMariaDBChecker(pinger func(ctx context.Context) error) *MariaDBChecker {
	return &MariaDBChecker{pinger: pinger}
}

// Name returns the checker name
func (c *MariaDBChecker) Name() string {
	return "mariadb"
}

// Check performs the health check
func (c *MariaDBChecker) Check(ctx context.Context) (Status, string) {
	start := time.Now()
	err := c.pinger(ctx)
	duration := time.Since(start)

	if err != nil {
		return StatusDown, fmt.Sprintf("MariaDB connection failed: %v", err)
	}

	if duration > 500*time.Millisecond {
		return StatusDegraded, fmt.Sprintf("MariaDB response time degraded: %v", duration)
	}

	return StatusUp, fmt.Sprintf("Response time: %v", duration)
}

// SetupHealthChecks initializes all health checks
func SetupHealthChecks(cfg *config.Config, db database.Database, cacheClient cache.Cache) *Health {
	health := NewHealth(cfg.App.Version)

	// Add checkers
	if db != nil {
		health.AddChecker(NewDatabaseChecker(db))
	}

	if cacheClient != nil {
		health.AddChecker(NewCacheChecker(cacheClient))
	}

	health.AddChecker(NewDiskChecker("/"))

	return health
}
