package devpanel

import (
    "sync"
    "time"
)

// MetricsCollector handles collection and storage of service metrics
type MetricsCollector struct {
    metrics map[string][]MetricPoint
    mu      sync.RWMutex
    config  Config
}

// MetricPoint represents a single metric data point
type MetricPoint struct {
    Timestamp time.Time
    Value     float64
    Labels    map[string]string
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config Config) *MetricsCollector {
    return &MetricsCollector{
        metrics: make(map[string][]MetricPoint),
        config:  config,
    }
}

// Collect collects metrics for a service
func (mc *MetricsCollector) Collect(serviceName string, metrics map[string]float64) {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    now := time.Now()
    for name, value := range metrics {
        key := serviceName + ":" + name
        point := MetricPoint{
            Timestamp: now,
            Value:     value,
            Labels:    map[string]string{"service": serviceName},
        }

        mc.metrics[key] = append(mc.metrics[key], point)

        // Trim old data points
        mc.trimOldData(key)
    }
}

// GetMetrics returns metrics for a service within a time range
func (mc *MetricsCollector) GetMetrics(serviceName string, start, end time.Time) map[string][]MetricPoint {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    result := make(map[string][]MetricPoint)
    for key, points := range mc.metrics {
        if !mc.isServiceMetric(key, serviceName) {
            continue
        }

        filtered := mc.filterPointsByTimeRange(points, start, end)
        if len(filtered) > 0 {
            result[key] = filtered
        }
    }

    return result
}

// GetLatestMetrics returns the most recent metrics for a service
func (mc *MetricsCollector) GetLatestMetrics(serviceName string) map[string]float64 {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    result := make(map[string]float64)
    for key, points := range mc.metrics {
        if !mc.isServiceMetric(key, serviceName) {
            continue
        }

        if len(points) > 0 {
            result[mc.getMetricName(key)] = points[len(points)-1].Value
        }
    }

    return result
}

// trimOldData removes data points older than the retention period
func (mc *MetricsCollector) trimOldData(key string) {
    points := mc.metrics[key]
    cutoff := time.Now().Add(-mc.config.MetricsInterval)
    
    var valid []MetricPoint
    for _, point := range points {
        if point.Timestamp.After(cutoff) {
            valid = append(valid, point)
        }
    }
    
    mc.metrics[key] = valid
}

// filterPointsByTimeRange returns points within the specified time range
func (mc *MetricsCollector) filterPointsByTimeRange(points []MetricPoint, start, end time.Time) []MetricPoint {
    var filtered []MetricPoint
    for _, point := range points {
        if point.Timestamp.After(start) && point.Timestamp.Before(end) {
            filtered = append(filtered, point)
        }
    }
    return filtered
}

// isServiceMetric checks if a metric belongs to a specific service
func (mc *MetricsCollector) isServiceMetric(key, serviceName string) bool {
    return len(key) > len(serviceName)+1 && key[:len(serviceName)+1] == serviceName+":"
}

// getMetricName extracts the metric name from a key
func (mc *MetricsCollector) getMetricName(key string) string {
    for i := len(key) - 1; i >= 0; i-- {
        if key[i] == ':' {
            return key[i+1:]
        }
    }
    return key
}

// StartCollecting begins collecting metrics for all services
func (mc *MetricsCollector) StartCollecting(serviceManager *core.ServiceManager) {
    go func() {
        ticker := time.NewTicker(mc.config.MetricsInterval)
        defer ticker.Stop()

        for range ticker.C {
            for name, service := range serviceManager.GetAllServices() {
                if status, err := serviceManager.GetServiceStatus(name); err == nil && status.Running {
                    metrics := mc.collectServiceMetrics(service)
                    mc.Collect(name, metrics)
                }
            }
        }
    }()
}

// collectServiceMetrics gathers metrics for a specific service
func (mc *MetricsCollector) collectServiceMetrics(service core.Service) map[string]float64 {
    metrics := make(map[string]float64)
    
    // Get process stats
    if process, err := process.NewProcess(int32(os.Getpid())); err == nil {
        if memPercent, err := process.MemoryPercent(); err == nil {
            metrics["memory_usage"] = memPercent
        }
        if cpuPercent, err := process.CPUPercent(); err == nil {
            metrics["cpu_usage"] = cpuPercent
        }
    }
    
    // Get service-specific metrics
    if status, err := service.Status(); err == nil {
        metrics["uptime_seconds"] = status.Uptime.Seconds()
        metrics["error_count"] = float64(len(status.Errors))
    }
    
    return metrics
} 