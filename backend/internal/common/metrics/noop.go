package metrics

import (
	"context"
	"net/http"
)

// NoOpProvider is a metrics provider that doesn't collect any metrics
type NoOpProvider struct{}

// NewNoOpProvider creates a new NoOpProvider
func NewNoOpProvider() *NoOpProvider {
	return &NoOpProvider{}
}

// Middleware returns a pass-through middleware function that doesn't collect any metrics
func (p *NoOpProvider) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}

// RegisterHandlers is a no-op implementation that doesn't register any handlers
func (p *NoOpProvider) RegisterHandlers(mux *http.ServeMux) {
	// No handlers to register
}

// SetAppInfo is a no-op implementation that doesn't collect application metrics
func (p *NoOpProvider) SetAppInfo(version, commit, buildDate string) {
	// No metrics to set
}

// SetUserCount is a no-op implementation that doesn't track user counts
func (p *NoOpProvider) SetUserCount(count int) {
	// No metrics to set
}

// SetURLCount is a no-op implementation that doesn't track URL counts
func (p *NoOpProvider) SetURLCount(count int) {
	// No metrics to set
}

// SetMessageCount is a no-op implementation that doesn't track message counts
func (p *NoOpProvider) SetMessageCount(count int) {
	// No metrics to set
}

// Shutdown is a no-op implementation that doesn't need to clean up any resources
func (p *NoOpProvider) Shutdown(ctx context.Context) error {
	return nil
}
