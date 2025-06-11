package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"gorm.io/gorm"
)

// DefaultConfig returns the default metrics configuration
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Grafana: DefaultGrafanaConfig(),
	}
}

// Manager manages application metrics
type Manager struct {
	config         Config
	providers      []Provider
	latencyTracker *LatencyTracker
}

// NewManager creates a new metrics manager
func NewManager(config Config) (*Manager, error) {
	if !config.Enabled {
		return &Manager{config: config}, nil
	}

	manager := &Manager{
		config:         config,
		providers:      make([]Provider, 0),
		latencyTracker: NewLatencyTracker(10080), // Store 1 week of 1-minute intervals
	}

	// Initialize Grafana provider if enabled
	if config.Grafana.Enabled {
		grafanaProvider, err := NewGrafanaProvider(config.Grafana)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Grafana metrics: %w", err)
		}
		manager.providers = append(manager.providers, grafanaProvider)
	}

	return manager, nil
}

// NewManagerWithDB creates a new metrics manager with database support
func NewManagerWithDB(config Config, db *gorm.DB) (*Manager, error) {
	if !config.Enabled {
		return &Manager{config: config}, nil
	}

	manager := &Manager{
		config:         config,
		providers:      make([]Provider, 0),
		latencyTracker: NewLatencyTrackerWithDB(1440, db), // Store 1 day in memory + unlimited in DB
	}

	// Initialize Grafana provider if enabled
	if config.Grafana.Enabled {
		grafanaProvider, err := NewGrafanaProvider(config.Grafana)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Grafana metrics: %w", err)
		}
		manager.providers = append(manager.providers, grafanaProvider)
	}

	return manager, nil
}

// Middleware returns an HTTP middleware for tracking metrics
func (m *Manager) Middleware(next http.Handler) http.Handler {
	if !m.config.Enabled || len(m.providers) == 0 {
		return next
	}

	// Apply all provider middlewares
	handler := next
	for i := len(m.providers) - 1; i >= 0; i-- {
		handler = m.providers[i].Middleware()(handler)
	}

	return handler
}

// RegisterHandlers registers the metrics endpoints
func (m *Manager) RegisterHandlers(mux *http.ServeMux) {
	if !m.config.Enabled {
		return
	}

	for _, provider := range m.providers {
		provider.RegisterHandlers(mux)
	}
}

// SetAppInfo sets application information
func (m *Manager) SetAppInfo(version, goVersion, buildDate string) {
	if !m.config.Enabled {
		return
	}

	for _, provider := range m.providers {
		provider.SetAppInfo(version, goVersion, buildDate)
	}
}

// SetUserCount sets the current user count
func (m *Manager) SetUserCount(count int) {
	if !m.config.Enabled {
		return
	}

	for _, provider := range m.providers {
		provider.SetUserCount(count)
	}
}

// SetURLCount sets the current URL count
func (m *Manager) SetURLCount(count int) {
	if !m.config.Enabled {
		return
	}

	for _, provider := range m.providers {
		provider.SetURLCount(count)
	}
}

// SetMessageCount sets the current message count
func (m *Manager) SetMessageCount(count int) {
	if !m.config.Enabled {
		return
	}

	for _, provider := range m.providers {
		provider.SetMessageCount(count)
	}
}

// RecordLatency records a latency metric
func (m *Manager) RecordLatency(latency float64, endpoint string) {
	if !m.config.Enabled || m.latencyTracker == nil {
		return
	}
	m.latencyTracker.AddMetric(latency, endpoint)
}

// GetLatencyMetrics returns latency metrics for the specified period
func (m *Manager) GetLatencyMetrics(period TimePeriod) []LatencyMetric {
	if !m.config.Enabled || m.latencyTracker == nil {
		return []LatencyMetric{}
	}
	return m.latencyTracker.GetAggregatedMetrics(period)
}

// GetLatencyStats returns latency statistics for the specified period
func (m *Manager) GetLatencyStats(period TimePeriod) LatencyStats {
	if !m.config.Enabled || m.latencyTracker == nil {
		return LatencyStats{Period: string(period)}
	}
	return m.latencyTracker.GetLatencyStats(period)
}

// HasSufficientLatencyData checks if there's enough data for the specified period
func (m *Manager) HasSufficientLatencyData(period TimePeriod) bool {
	if !m.config.Enabled || m.latencyTracker == nil {
		return false
	}
	return m.latencyTracker.HasSufficientData(period)
}

// GetLatestLatency returns the most recent latency measurement
func (m *Manager) GetLatestLatency() *LatencyMetric {
	if !m.config.Enabled || m.latencyTracker == nil {
		return nil
	}
	return m.latencyTracker.GetLatestMetric()
}

// Shutdown stops any running servers
func (m *Manager) Shutdown(ctx context.Context) error {
	if !m.config.Enabled {
		return nil
	}

	// Force flush any remaining metrics to database
	if m.latencyTracker != nil {
		m.latencyTracker.ForceFlush()
	}

	var lastErr error
	for _, provider := range m.providers {
		if err := provider.Shutdown(ctx); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// StartPeriodicCleanup starts a background goroutine to clean up old metrics
func (m *Manager) StartPeriodicCleanup() {
	if !m.config.Enabled || m.latencyTracker == nil {
		return
	}
	
	go func() {
		// Clean up old metrics every 24 hours
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			m.latencyTracker.Cleanup()
		}
	}()
}
