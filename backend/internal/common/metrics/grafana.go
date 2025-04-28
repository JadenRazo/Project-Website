package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultGrafanaConfig returns the default Grafana configuration
func DefaultGrafanaConfig() GrafanaConfig {
	return GrafanaConfig{
		Enabled:       true,
		Endpoint:      "/metrics",
		MetricsPrefix: "app_",
		Namespace:     "api",
		Subsystem:     "http",
		Address:       ":9090",
		DefaultLabels: map[string]string{
			"service": "backend",
		},
	}
}

// GrafanaProvider implements the Provider interface for Grafana
type GrafanaProvider struct {
	Config GrafanaConfig

	// Prometheus registry
	registry *prometheus.Registry

	// HTTP metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestsDuration *prometheus.HistogramVec

	// Application metrics
	appInfo       *prometheus.GaugeVec
	usersTotal    prometheus.Gauge
	urlsTotal     prometheus.Gauge
	messagesTotal prometheus.Gauge
}

// NewGrafanaProvider creates a new Grafana metrics provider
func NewGrafanaProvider(config GrafanaConfig) (*GrafanaProvider, error) {
	registry := prometheus.NewRegistry()

	// Create a new factory with our registry
	factory := promauto.With(registry)

	// Prefix for all metrics
	prefix := config.MetricsPrefix

	// HTTP metrics
	httpRequestsTotal := factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: prefix + "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestsDuration := factory.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prefix + "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// Application info
	appInfo := factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "info",
			Help: "Application information",
		},
		[]string{"version", "go_version", "build_date"},
	)

	// Business metrics
	usersTotal := factory.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "users_total",
			Help: "Total number of users",
		},
	)

	urlsTotal := factory.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "urls_total",
			Help: "Total number of URLs",
		},
	)

	messagesTotal := factory.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "messages_total",
			Help: "Total number of messages",
		},
	)

	provider := &GrafanaProvider{
		Config:               config,
		registry:             registry,
		httpRequestsTotal:    httpRequestsTotal,
		httpRequestsDuration: httpRequestsDuration,
		appInfo:              appInfo,
		usersTotal:           usersTotal,
		urlsTotal:            urlsTotal,
		messagesTotal:        messagesTotal,
	}

	return provider, nil
}

// Middleware returns an HTTP middleware for tracking metrics
func (g *GrafanaProvider) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture the status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler
			next.ServeHTTP(ww, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(ww.statusCode)

			g.httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
			g.httpRequestsDuration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration)
		})
	}
}

// RegisterHandlers registers metrics endpoints
func (g *GrafanaProvider) RegisterHandlers(mux *http.ServeMux) {
	mux.Handle(g.Config.Endpoint, promhttp.HandlerFor(g.registry, promhttp.HandlerOpts{}))
}

// SetAppInfo sets application information
func (g *GrafanaProvider) SetAppInfo(version, goVersion, buildDate string) {
	g.appInfo.WithLabelValues(version, goVersion, buildDate).Set(1)
}

// SetUserCount sets the current user count
func (g *GrafanaProvider) SetUserCount(count int) {
	g.usersTotal.Set(float64(count))
}

// SetURLCount sets the current URL count
func (g *GrafanaProvider) SetURLCount(count int) {
	g.urlsTotal.Set(float64(count))
}

// SetMessageCount sets the current message count
func (g *GrafanaProvider) SetMessageCount(count int) {
	g.messagesTotal.Set(float64(count))
}

// Shutdown stops any running servers
func (g *GrafanaProvider) Shutdown(ctx context.Context) error {
	// No active services to shut down
	return nil
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if it hasn't been set yet
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
