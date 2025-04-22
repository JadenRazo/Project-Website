package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(cfg *config.TracingConfig) (func(), error) {
	if !cfg.Enabled {
		// Return a no-op shutdown function
		return func() {}, nil
	}

	// Create the exporter
	exporter, err := createExporter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create a resource describing this application
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create a tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(createSampler(cfg)),
	)

	// Set the global trace provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Get a tracer
	tracer = tp.Tracer(cfg.ServiceName)

	// Return a function to shutdown the exporter when the application exits
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

// createExporter creates the appropriate exporter
func createExporter(cfg *config.TracingConfig) (sdktrace.SpanExporter, error) {
	ctx := context.Background()

	// Use OTLP exporter with gRPC
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(), // For development, consider secure connections for production
	)

	return otlptrace.New(ctx, client)
}

// createSampler creates a sampler based on configuration
func createSampler(cfg *config.TracingConfig) sdktrace.Sampler {
	switch cfg.SamplingStrategy {
	case "always":
		return sdktrace.AlwaysSample()
	case "never":
		return sdktrace.NeverSample()
	case "ratio":
		return sdktrace.TraceIDRatioBased(cfg.SamplingRatio)
	default:
		// Default to parent-based sampling with ratio
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplingRatio))
	}
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	if tracer == nil {
		// Return a no-op span if tracer is not initialized
		return ctx, trace.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, name)
}

// AddAttribute adds an attribute to the current span
func AddAttribute(ctx context.Context, key string, value interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case int64:
		span.SetAttributes(attribute.Int64(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	default:
		span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attributes ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.AddEvent(name, trace.WithAttributes(attributes...))
}

// RecordError records an error to the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() || err == nil {
		return
	}
	span.RecordError(err)
}

// TraceHTTPMiddleware creates a middleware that adds tracing to HTTP requests
func TraceHTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract context from the request headers
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start a new span
			ctx, span := tracer.Start(
				ctx,
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				trace.WithAttributes(
					semconv.HTTPMethodKey.String(r.Method),
					semconv.HTTPTargetKey.String(r.URL.Path),
					semconv.HTTPURLKey.String(r.URL.String()),
					semconv.HTTPHostKey.String(r.Host),
					semconv.HTTPUserAgentKey.String(r.UserAgent()),
					attribute.String("service.name", serviceName),
				),
			)
			defer span.End()

			// Create a wrapped response writer to capture status code
			wrapper := newResponseWriterWrapper(w)

			// Call the next handler with the augmented context
			next.ServeHTTP(wrapper, r.WithContext(ctx))

			// Add response attributes
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(wrapper.statusCode),
			)

			// Mark the span as error if status code is 5xx
			if wrapper.statusCode >= 500 {
				span.SetStatus(trace.StatusCodeError, fmt.Sprintf("HTTP %d", wrapper.statusCode))
			}
		})
	}
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

// newResponseWriterWrapper creates a new wrapper for response writer
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}
