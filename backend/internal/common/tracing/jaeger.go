package tracing

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// Config holds configuration for the tracer
type Config struct {
	Enabled     bool    `json:"enabled" mapstructure:"enabled"`
	ServiceName string  `json:"service_name" mapstructure:"service_name"`
	AgentHost   string  `json:"agent_host" mapstructure:"agent_host"`
	AgentPort   string  `json:"agent_port" mapstructure:"agent_port"`
	LogSpans    bool    `json:"log_spans" mapstructure:"log_spans"`
	SampleRate  float64 `json:"sample_rate" mapstructure:"sample_rate"`
}

// DefaultConfig returns the default configuration for tracing
func DefaultConfig() *Config {
	return &Config{
		Enabled:     false,
		ServiceName: "backend-api",
		AgentHost:   "localhost",
		AgentPort:   "6831",
		LogSpans:    false,
		SampleRate:  0.1,
	}
}

// Tracer wraps an opentracing.Tracer with additional functionality
type Tracer struct {
	opentracing.Tracer
	config     *Config
	closer     io.Closer
	isDisabled bool
}

// NewTracer creates a new Jaeger tracer
func NewTracer(config *Config) (*Tracer, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// If tracing is disabled, return a no-op tracer
	if !config.Enabled {
		return &Tracer{
			Tracer:     opentracing.NoopTracer{},
			config:     config,
			isDisabled: true,
		}, nil
	}

	// Configure Jaeger tracer
	jcfg := jaegercfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: config.SampleRate,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            config.LogSpans,
			LocalAgentHostPort:  fmt.Sprintf("%s:%s", config.AgentHost, config.AgentPort),
			BufferFlushInterval: 1 * time.Second,
		},
	}

	// Initialize tracer with metrics and logger
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := jcfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jaeger tracer: %w", err)
	}

	// Set as global tracer
	opentracing.SetGlobalTracer(tracer)

	return &Tracer{
		Tracer:     tracer,
		config:     config,
		closer:     closer,
		isDisabled: false,
	}, nil
}

// StartSpan starts a new span with the given operation name
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.Tracer.StartSpan(operationName, opts...)
}

// StartSpanFromContext starts a span with the given operation name using the parent span from context if available
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if t.isDisabled {
		return opentracing.NoopTracer{}.StartSpan(operationName), ctx
	}

	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = t.StartSpan(operationName, opts...)
	} else {
		span = t.StartSpan(operationName, opts...)
	}

	return span, opentracing.ContextWithSpan(ctx, span)
}

// Inject takes the SpanContext and injects it into the carrier using the format
func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	if t.isDisabled {
		return nil
	}
	return t.Tracer.Inject(sm, format, carrier)
}

// Extract extracts a SpanContext from the carrier using the format
func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	if t.isDisabled {
		return nil, opentracing.ErrSpanContextNotFound
	}
	return t.Tracer.Extract(format, carrier)
}

// IsEnabled returns whether tracing is enabled
func (t *Tracer) IsEnabled() bool {
	return !t.isDisabled
}

// Shutdown closes the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.isDisabled || t.closer == nil {
		return nil
	}

	// Create a channel to signal completion
	done := make(chan struct{})
	var err error

	// Close the tracer in a goroutine
	go func() {
		err = t.closer.Close()
		close(done)
	}()

	// Wait for close to complete or timeout
	select {
	case <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
