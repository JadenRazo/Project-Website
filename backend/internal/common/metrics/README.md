# Metrics System

This package provides a metrics collection and reporting system for the backend application.

## Overview

The metrics system is designed to be extensible, supporting multiple metric providers. Currently, it includes:

- **Grafana/Prometheus Integration**: Collects and exposes metrics in Prometheus format that can be visualized in Grafana dashboards.

## Configuration

The metrics system can be configured via the `Config` struct:

```go
config := metrics.DefaultConfig()

// Enable or disable Grafana metrics
config.Grafana.Enabled = true

// Configure the HTTP endpoint for metrics
config.Grafana.Endpoint = "/metrics"

// Set a prefix for all metrics
config.Grafana.MetricsPrefix = "myapp"

// Add default labels to all metrics
config.Grafana.DefaultLabels = map[string]string{
    "environment": "production",
}
```

## Initialization

To initialize the metrics system:

```go
metricsManager, err := metrics.NewManager(config)
if err != nil {
    log.Fatalf("Failed to initialize metrics: %v", err)
}
```

## HTTP Middleware

The metrics system includes middleware for HTTP request instrumentation:

```go
// Example with standard library
mux := http.NewServeMux()
mux.Handle("/api/", metricsManager.Middleware()(apiHandler))

// Or for a global router
router.Use(metricsManager.Middleware())
```

## Exposing Metrics Endpoint

To expose the metrics endpoint:

```go
// Register metrics handlers on a mux
mux := http.NewServeMux()
metricsManager.RegisterHandlers(mux)

// Start a server for metrics
go http.ListenAndServe(":9090", mux)
```

## Business Metrics

The metrics system includes methods to track business metrics:

```go
// Set application information
metricsManager.SetAppInfo("1.0.0", "go1.16", "2023-01-01")

// Track user count
metricsManager.SetUserCount(1000)

// Track URL count
metricsManager.SetURLCount(5000)

// Track message count
metricsManager.SetMessageCount(10000)
```

## Extending

To add a new metrics provider:

1. Create a new file for the provider (e.g., `newrelic.go`)
2. Add the provider's configuration to the `Config` struct
3. Initialize the provider in the `NewManager` function
4. Add provider-specific methods to the `Manager` struct

## Grafana Dashboard

A sample Grafana dashboard configuration is available in the `dashboards` directory. To use it:

1. Import the JSON file into your Grafana instance
2. Configure a Prometheus data source that scrapes your application's metrics endpoint
3. Select the data source in the dashboard settings

## Best Practices

1. Use concise metric names that follow a consistent naming pattern
2. Use labels to differentiate metric dimensions, but avoid high cardinality
3. Document all custom metrics in your application's documentation
4. Monitor the performance impact of metrics collection 