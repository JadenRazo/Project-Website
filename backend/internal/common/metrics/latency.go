package metrics

import (
	"sort"
	"sync"
	"time"
)

// LatencyMetric represents a latency measurement at a specific time
type LatencyMetric struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   float64   `json:"latency"` // in milliseconds
	Endpoint  string    `json:"endpoint,omitempty"`
}

// LatencyTracker tracks latency metrics over time
type LatencyTracker struct {
	mu      sync.RWMutex
	metrics []LatencyMetric
	maxSize int
}

// NewLatencyTracker creates a new latency tracker
func NewLatencyTracker(maxSize int) *LatencyTracker {
	if maxSize <= 0 {
		maxSize = 10080 // Default: 1 week of 1-minute intervals
	}
	
	tracker := &LatencyTracker{
		metrics: make([]LatencyMetric, 0, maxSize),
		maxSize: maxSize,
	}
	
	return tracker
}

// AddMetric adds a new latency metric
func (lt *LatencyTracker) AddMetric(latency float64, endpoint string) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	
	metric := LatencyMetric{
		Timestamp: time.Now(),
		Latency:   latency,
		Endpoint:  endpoint,
	}
	
	lt.metrics = append(lt.metrics, metric)
	
	// Remove old metrics if we exceed max size
	if len(lt.metrics) > lt.maxSize {
		// Remove oldest 10% to avoid frequent array copies
		removeCount := lt.maxSize / 10
		if removeCount < 1 {
			removeCount = 1
		}
		lt.metrics = lt.metrics[removeCount:]
	}
}

// GetMetrics returns metrics for the specified time period
func (lt *LatencyTracker) GetMetrics(period TimePeriod) []LatencyMetric {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	
	now := time.Now()
	var cutoff time.Time
	
	switch period {
	case TimePeriodDay:
		cutoff = now.Add(-24 * time.Hour)
	case TimePeriodWeek:
		cutoff = now.Add(-7 * 24 * time.Hour)
	case TimePeriodMonth:
		cutoff = now.Add(-30 * 24 * time.Hour)
	default:
		cutoff = now.Add(-24 * time.Hour)
	}
	
	var result []LatencyMetric
	for _, metric := range lt.metrics {
		if metric.Timestamp.After(cutoff) {
			result = append(result, metric)
		}
	}
	
	// Sort by timestamp chronologically
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	
	return result
}

// GetAggregatedMetrics returns aggregated metrics for the specified period
func (lt *LatencyTracker) GetAggregatedMetrics(period TimePeriod) []LatencyMetric {
	metrics := lt.GetMetrics(period)
	if len(metrics) == 0 {
		return []LatencyMetric{}
	}
	
	var interval time.Duration
	switch period {
	case TimePeriodDay:
		interval = 10 * time.Minute // 10-minute intervals for day view (12:00, 12:10, 12:20, etc.)
	case TimePeriodWeek:
		interval = 30 * time.Minute // 30-minute intervals for week view
	case TimePeriodMonth:
		interval = 2 * time.Hour // 2-hour intervals for month view
	default:
		interval = 10 * time.Minute
	}
	
	return lt.aggregateByInterval(metrics, interval)
}

// aggregateByInterval groups metrics by time intervals and calculates averages
func (lt *LatencyTracker) aggregateByInterval(metrics []LatencyMetric, interval time.Duration) []LatencyMetric {
	if len(metrics) == 0 {
		return []LatencyMetric{}
	}
	
	buckets := make(map[int64][]float64)
	bucketTimeSet := make(map[int64]bool)
	var bucketTimes []int64
	
	// Group metrics into time buckets aligned to clean intervals
	for _, metric := range metrics {
		// Align timestamp to the interval boundary
		intervalSeconds := int64(interval.Seconds())
		alignedTime := (metric.Timestamp.Unix() / intervalSeconds) * intervalSeconds
		
		if !bucketTimeSet[alignedTime] {
			bucketTimes = append(bucketTimes, alignedTime)
			bucketTimeSet[alignedTime] = true
		}
		buckets[alignedTime] = append(buckets[alignedTime], metric.Latency)
	}
	
	// Sort bucket times chronologically
	sort.Slice(bucketTimes, func(i, j int) bool {
		return bucketTimes[i] < bucketTimes[j]
	})
	
	// Calculate averages for each bucket in chronological order
	var result []LatencyMetric
	for _, bucketTime := range bucketTimes {
		latencies := buckets[bucketTime]
		if len(latencies) == 0 {
			continue
		}
		
		var sum float64
		for _, latency := range latencies {
			sum += latency
		}
		avg := sum / float64(len(latencies))
		
		// Use the aligned timestamp directly
		timestamp := time.Unix(bucketTime, 0)
		result = append(result, LatencyMetric{
			Timestamp: timestamp,
			Latency:   avg,
		})
	}
	
	return result
}

// GetLatestMetric returns the most recent latency metric
func (lt *LatencyTracker) GetLatestMetric() *LatencyMetric {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	
	if len(lt.metrics) == 0 {
		return nil
	}
	
	latest := lt.metrics[len(lt.metrics)-1]
	return &latest
}

// GetAverageLatency returns the average latency for the specified period
func (lt *LatencyTracker) GetAverageLatency(period TimePeriod) float64 {
	metrics := lt.GetMetrics(period)
	if len(metrics) == 0 {
		return 0
	}
	
	var sum float64
	for _, metric := range metrics {
		sum += metric.Latency
	}
	
	return sum / float64(len(metrics))
}

// HasSufficientData checks if there's enough data for the specified period
func (lt *LatencyTracker) HasSufficientData(period TimePeriod) bool {
	metrics := lt.GetMetrics(period)
	
	switch period {
	case TimePeriodDay:
		return len(metrics) >= 6 // At least 30 minutes of 30-second intervals
	case TimePeriodWeek:
		return len(metrics) >= 48 // At least 2 hours worth of data
	case TimePeriodMonth:
		return len(metrics) >= 168 // At least 12 hours worth of data
	default:
		return len(metrics) >= 6
	}
}

// TimePeriod represents different time periods for metrics
type TimePeriod string

const (
	TimePeriodDay   TimePeriod = "day"
	TimePeriodWeek  TimePeriod = "week"
	TimePeriodMonth TimePeriod = "month"
)

// LatencyStats represents statistical information about latency
type LatencyStats struct {
	Average float64 `json:"average"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Count   int     `json:"count"`
	Period  string  `json:"period"`
}

// GetLatencyStats calculates statistics for the specified period
func (lt *LatencyTracker) GetLatencyStats(period TimePeriod) LatencyStats {
	metrics := lt.GetMetrics(period)
	
	stats := LatencyStats{
		Period: string(period),
		Count:  len(metrics),
	}
	
	if len(metrics) == 0 {
		return stats
	}
	
	var sum, min, max float64
	min = metrics[0].Latency
	max = metrics[0].Latency
	
	for _, metric := range metrics {
		sum += metric.Latency
		if metric.Latency < min {
			min = metric.Latency
		}
		if metric.Latency > max {
			max = metric.Latency
		}
	}
	
	stats.Average = sum / float64(len(metrics))
	stats.Min = min
	stats.Max = max
	
	return stats
}

