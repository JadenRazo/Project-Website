package metrics

import (
	"context"
	"fmt"
	"net/http"
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
	config    Config
	providers []Provider
}

// NewManager creates a new metrics manager
func NewManager(config Config) (*Manager, error) {
	if !config.Enabled {
		return &Manager{config: config}, nil
	}

	manager := &Manager{
		config:    config,
		providers: make([]Provider, 0),
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

// Shutdown stops any running servers
func (m *Manager) Shutdown(ctx context.Context) error {
	if !m.config.Enabled {
		return nil
	}

	var lastErr error
	for _, provider := range m.providers {
		if err := provider.Shutdown(ctx); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
