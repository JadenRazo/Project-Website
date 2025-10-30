package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestDuration tracks HTTP request durations
	RequestDuration *prometheus.HistogramVec

	// RequestsTotal tracks total HTTP requests
	RequestsTotal *prometheus.CounterVec

	// ErrorsTotal tracks total errors
	ErrorsTotal *prometheus.CounterVec

	// ActiveRequests tracks currently active requests
	ActiveRequests prometheus.Gauge

	// DatabaseQueryDuration tracks database query durations
	DatabaseQueryDuration *prometheus.HistogramVec

	// CacheHitTotal tracks cache hits
	CacheHitTotal prometheus.Counter

	// CacheMissTotal tracks cache misses
	CacheMissTotal prometheus.Counter

	// ServiceOperationDuration tracks service operation durations
	ServiceOperationDuration *prometheus.HistogramVec

	// QueueSize tracks message queue size
	QueueSize *prometheus.GaugeVec

	// QueueProcessingTime tracks message processing time
	QueueProcessingTime *prometheus.HistogramVec

	// Visitor metrics
	CurrentVisitors prometheus.Gauge

	// VisitorsByPeriod tracks visitors by time period
	VisitorsByPeriod *prometheus.GaugeVec

	// PageViewsByPeriod tracks page views by time period
	PageViewsByPeriod *prometheus.GaugeVec

	// VisitorTrend tracks visitor trend percentages
	VisitorTrend *prometheus.GaugeVec

	// VisitorSessionDuration tracks visitor session durations
	VisitorSessionDuration *prometheus.HistogramVec

	// VisitorBounceRate tracks bounce rate
	VisitorBounceRate prometheus.Gauge

	// VisitorsByCountry tracks visitors by country
	VisitorsByCountry *prometheus.GaugeVec
)

// InitMetrics initializes all Prometheus metrics
func InitMetrics(cfg *config.MetricsConfig) {
	if !cfg.Enabled {
		return
	}

	// HTTP metrics
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "handler", "method", "status"},
	)

	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "handler", "method", "status"},
	)

	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"service", "type"},
	)

	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current number of active HTTP requests",
		},
	)

	// Database metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation", "table"},
	)

	// Cache metrics
	CacheHitTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total number of cache misses",
		},
	)

	// Service metrics
	ServiceOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_operation_duration_seconds",
			Help:    "Duration of service operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation"},
	)

	// Queue metrics
	QueueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_size",
			Help: "Current size of the message queue",
		},
		[]string{"queue"},
	)

	QueueProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "queue_processing_time_seconds",
			Help:    "Time to process a message from the queue",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"queue", "message_type"},
	)

	// Visitor metrics
	CurrentVisitors = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "visitor_current_active",
			Help: "Current number of active visitors",
		},
	)

	VisitorsByPeriod = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "visitor_count_by_period",
			Help: "Number of visitors by time period",
		},
		[]string{"period"},
	)

	PageViewsByPeriod = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "page_views_by_period",
			Help: "Number of page views by time period",
		},
		[]string{"period"},
	)

	VisitorTrend = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "visitor_trend_percentage",
			Help: "Visitor trend percentage change",
		},
		[]string{"period", "direction"},
	)

	VisitorSessionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "visitor_session_duration_seconds",
			Help:    "Duration of visitor sessions in seconds",
			Buckets: []float64{30, 60, 120, 300, 600, 1200, 1800, 3600},
		},
		[]string{"service"},
	)

	VisitorBounceRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "visitor_bounce_rate",
			Help: "Percentage of visitors who leave after viewing only one page",
		},
	)

	VisitorsByCountry = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "visitors_by_country",
			Help: "Number of visitors by country",
		},
		[]string{"country_code"},
	)
}

// PrometheusMiddleware creates a middleware for HTTP request metrics
func PrometheusMiddleware(service string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Track active requests
			ActiveRequests.Inc()
			defer ActiveRequests.Dec()

			// Create a response writer that captures the status code
			rww := &responseWriterWrapper{ResponseWriter: w, statusCode: 200}

			// Process the request
			next.ServeHTTP(rww, r)

			// Duration in seconds
			duration := time.Since(start).Seconds()

			// Record metrics
			labels := prometheus.Labels{
				"service": service,
				"handler": r.URL.Path,
				"method":  r.Method,
				"status":  fmt.Sprintf("%d", rww.statusCode),
			}

			RequestDuration.With(labels).Observe(duration)
			RequestsTotal.With(labels).Inc()

			// Record errors
			if rww.statusCode >= 500 {
				ErrorsTotal.With(prometheus.Labels{
					"service": service,
					"type":    "server_error",
				}).Inc()
			} else if rww.statusCode >= 400 {
				ErrorsTotal.With(prometheus.Labels{
					"service": service,
					"type":    "client_error",
				}).Inc()
			}
		})
	}
}

// TrackDatabaseQuery tracks a database query duration
func TrackDatabaseQuery(service, operation, table string, f func() error) error {
	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()

	DatabaseQueryDuration.With(prometheus.Labels{
		"service":   service,
		"operation": operation,
		"table":     table,
	}).Observe(duration)

	if err != nil {
		ErrorsTotal.With(prometheus.Labels{
			"service": service,
			"type":    "database_error",
		}).Inc()
	}

	return err
}

// TrackCacheOperation tracks cache hit/miss
func TrackCacheOperation(hit bool) {
	if hit {
		CacheHitTotal.Inc()
	} else {
		CacheMissTotal.Inc()
	}
}

// TrackServiceOperation tracks a service operation duration
func TrackServiceOperation(service, operation string, f func() error) error {
	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()

	ServiceOperationDuration.With(prometheus.Labels{
		"service":   service,
		"operation": operation,
	}).Observe(duration)

	if err != nil {
		ErrorsTotal.With(prometheus.Labels{
			"service": service,
			"type":    "service_error",
		}).Inc()
	}

	return err
}

// UpdateQueueSize updates the current queue size
func UpdateQueueSize(queue string, size float64) {
	QueueSize.With(prometheus.Labels{
		"queue": queue,
	}).Set(size)
}

// TrackQueueProcessingTime tracks message processing time
func TrackQueueProcessingTime(queue, messageType string, duration time.Duration) {
	QueueProcessingTime.With(prometheus.Labels{
		"queue":        queue,
		"message_type": messageType,
	}).Observe(duration.Seconds())
}

// MetricsHandler returns the Prometheus metrics HTTP handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// responseWriterWrapper captures the status code from a response
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rww *responseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}
