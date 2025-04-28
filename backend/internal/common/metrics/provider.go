package metrics

import (
	"context"
	"net/http"
)

// Provider defines the interface for metrics collection
type Provider interface {
	// Middleware returns an HTTP middleware function that collects request metrics
	Middleware() func(http.Handler) http.Handler

	// RegisterHandlers registers any HTTP handlers needed by the metrics provider
	RegisterHandlers(mux *http.ServeMux)

	// SetAppInfo sets application information for metrics
	SetAppInfo(version, commit, buildDate string)

	// SetUserCount updates the total user count metric
	SetUserCount(count int)

	// SetURLCount updates the total URL count metric
	SetURLCount(count int)

	// SetMessageCount updates the total message count metric
	SetMessageCount(count int)

	// Shutdown gracefully shuts down the metrics provider
	Shutdown(ctx context.Context) error
}

// NewProvider creates a new metrics provider based on the given configuration
func NewProvider(config *Config) Provider {
	if config != nil && config.Grafana.Enabled {
		provider, err := NewGrafanaProvider(config.Grafana)
		if err == nil {
			return provider
		}
		// If Grafana provider creation fails, fallback to NoOp provider
	}
	return NewNoOpProvider()
}
