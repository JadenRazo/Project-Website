package devpanel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/visitor"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
)

// EnhancedService handles devpanel operations with real data
type EnhancedService struct {
	*core.BaseService
	serviceManager *core.ServiceManager
	visitorService *visitor.Service
	db             *gorm.DB
	logManager     *LogManager
	metricsManager *metrics.Manager
	config         Config
	processMap     map[string]int32 // Maps service names to process IDs
}

// NewEnhancedService creates a new enhanced devpanel service
func NewEnhancedService(
	serviceManager *core.ServiceManager,
	visitorService *visitor.Service,
	db *gorm.DB,
	logManager *LogManager,
	metricsManager *metrics.Manager,
	config Config,
) *EnhancedService {
	return &EnhancedService{
		BaseService:    core.NewBaseService("devpanel"),
		serviceManager: serviceManager,
		visitorService: visitorService,
		db:             db,
		logManager:     logManager,
		metricsManager: metricsManager,
		config:         config,
		processMap:     make(map[string]int32),
	}
}

// collectSystemStats gathers real system-wide statistics
func (s *EnhancedService) collectRealSystemStats() (SystemStats, error) {
	stats := SystemStats{
		GoVersion:     runtime.Version(),
		NumGoroutines: runtime.NumGoroutine(),
		LastUpdate:    time.Now(),
	}

	// Get real CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPUUsage = cpuPercent[0]
	}

	// Get real memory usage
	if memInfo, err := mem.VirtualMemory(); err == nil {
		stats.MemoryUsage = memInfo.UsedPercent
	}

	// Get real disk usage
	if diskInfo, err := disk.Usage("/"); err == nil {
		stats.DiskUsage = diskInfo.UsedPercent
	}

	// Get system uptime
	if uptimeSeconds, err := getSystemUptime(); err == nil {
		stats.Uptime = formatDuration(time.Duration(uptimeSeconds) * time.Second)
	}

	return stats, nil
}

// collectRealServiceStats gathers real statistics for a specific service
func (s *EnhancedService) collectRealServiceStats(name string, service core.Service) (ServiceStats, error) {
	stats := ServiceStats{
		Name:       name,
		Status:     "unknown",
		LastUpdate: time.Now(),
	}

	// Get real service status
	status, err := s.serviceManager.GetServiceStatus(name)
	if err == nil {
		stats.Status = "running"
		if !status.Running {
			stats.Status = "stopped"
		}
		stats.Uptime = formatDuration(status.Uptime)
		if len(status.Errors) > 0 {
			stats.LastError = status.Errors[len(status.Errors)-1]
			stats.ErrorCount = int64(len(status.Errors))
		}
	}

	// Get process stats if we have a process ID
	if pid, exists := s.processMap[name]; exists && pid > 0 {
		if proc, err := process.NewProcess(pid); err == nil {
			// Get memory usage
			if memInfo, err := proc.MemoryInfo(); err == nil {
				stats.MemoryUsage = float64(memInfo.RSS) / (1024 * 1024) // Convert to MB
			}

			// Get CPU usage
			if cpuPercent, err := proc.CPUPercent(); err == nil {
				stats.CPUUsage = cpuPercent
			}
		}
	}

	// Get metrics from the metrics manager if available
	// Note: This would need adjustment based on the actual metrics manager API
	// For now, we'll skip this part as it depends on the metrics manager implementation

	// Get service-specific metrics from database
	if s.db != nil {
		var dbStats struct {
			RequestCount int64
			ErrorCount   int64
			AvgLatency   float64
		}

		s.db.Raw(`
			SELECT
				COUNT(*) as request_count,
				SUM(CASE WHEN error IS NOT NULL THEN 1 ELSE 0 END) as error_count,
				AVG(latency_ms) as avg_latency
			FROM metric_data
			WHERE service = ? AND created_at > ?
		`, name, time.Now().Add(-24*time.Hour)).Scan(&dbStats)

		if dbStats.RequestCount > 0 {
			stats.RequestCount = dbStats.RequestCount
			stats.ErrorCount = dbStats.ErrorCount
			stats.AverageResponseTime = dbStats.AvgLatency
		}
	}

	return stats, nil
}

// getLogsForServiceReal retrieves real logs for a specific service
func (s *EnhancedService) getLogsForServiceReal(serviceName string, limit int) ([]string, error) {
	if s.logManager != nil {
		return s.logManager.GetLogs(serviceName, limit)
	}

	// Fallback to reading from standard log location
	logPath := filepath.Join("logs", "services", serviceName+".log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// Try alternative log paths
		logPath = filepath.Join("logs", serviceName+".log")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			return []string{}, nil
		}
	}

	// Read the log file
	logs, err := readLastNLines(logPath, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	return logs, nil
}

// streamLogsForServiceReal streams real logs for a specific service
func (s *EnhancedService) streamLogsForServiceReal(serviceName string, logChan chan<- string) {
	if s.logManager != nil {
		stream := s.logManager.StreamLogs(serviceName)
		go func() {
			for {
				select {
				case line := <-stream.Lines:
					logChan <- line
				case <-stream.Done:
					return
				}
			}
		}()
		return
	}

	// Fallback implementation for streaming logs
	logPath := filepath.Join("logs", "services", serviceName+".log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = filepath.Join("logs", serviceName+".log")
	}

	go tailFile(logPath, logChan)
}

// collectMetricsForServiceReal collects real metrics for a specific service
func (s *EnhancedService) collectMetricsForServiceReal(serviceName string) (map[string]interface{}, error) {
	// Get real service stats
	services := s.serviceManager.GetAllServices()
	service, exists := services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}

	serviceStats, err := s.collectRealServiceStats(serviceName, service)
	if err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"service":        serviceName,
		"timestamp":      time.Now(),
		"cpu_usage":      serviceStats.CPUUsage,
		"memory_usage":   serviceStats.MemoryUsage,
		"request_count":  serviceStats.RequestCount,
		"error_count":    serviceStats.ErrorCount,
		"status":         serviceStats.Status,
		"uptime":         serviceStats.Uptime,
		"response_time":  serviceStats.AverageResponseTime,
	}

	// Add network stats if available
	if netStats, err := getServiceNetworkStats(serviceName); err == nil {
		metrics["network"] = netStats
	}

	// Add database connection pool stats if available
	if s.db != nil {
		sqlDB, _ := s.db.DB()
		if sqlDB != nil {
			dbStats := sqlDB.Stats()
			metrics["database"] = map[string]interface{}{
				"open_connections": dbStats.OpenConnections,
				"in_use":          dbStats.InUse,
				"idle":            dbStats.Idle,
				"wait_count":      dbStats.WaitCount,
				"wait_duration":   dbStats.WaitDuration.Milliseconds(),
			}
		}
	}

	return metrics, nil
}

// getMetricsHistoryForServiceReal retrieves real historical metrics for a service
func (s *EnhancedService) getMetricsHistoryForServiceReal(serviceName string, duration string) ([]map[string]interface{}, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Parse duration
	var timeRange time.Duration
	switch duration {
		case "1h":
			timeRange = time.Hour
		case "6h":
			timeRange = 6 * time.Hour
		case "24h":
			timeRange = 24 * time.Hour
		case "7d":
			timeRange = 7 * 24 * time.Hour
		case "30d":
			timeRange = 30 * 24 * time.Hour
		default:
			timeRange = time.Hour
	}

	cutoff := time.Now().Add(-timeRange)

	// Query historical metrics from database
	var metrics []struct {
		Timestamp    time.Time
		CPUUsage     float64
		MemoryUsage  float64
		RequestCount int64
		ErrorCount   int64
		AvgLatency   float64
	}

	err := s.db.Raw(`
		SELECT
			created_at as timestamp,
			AVG(cpu_usage) as cpu_usage,
			AVG(memory_usage) as memory_usage,
			COUNT(*) as request_count,
			SUM(CASE WHEN error IS NOT NULL THEN 1 ELSE 0 END) as error_count,
			AVG(latency_ms) as avg_latency
		FROM metric_data
		WHERE service = ? AND created_at > ?
		GROUP BY DATE_TRUNC('minute', created_at)
		ORDER BY timestamp DESC
		LIMIT 100
	`, serviceName, cutoff).Scan(&metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query metrics history: %w", err)
	}

	// Convert to response format
	history := make([]map[string]interface{}, 0, len(metrics))
	for _, m := range metrics {
		history = append(history, map[string]interface{}{
			"timestamp":     m.Timestamp,
			"cpu_usage":     m.CPUUsage,
			"memory_usage":  m.MemoryUsage,
			"request_count": m.RequestCount,
			"error_count":   m.ErrorCount,
			"avg_latency":   m.AvgLatency,
		})
	}

	return history, nil
}

// Helper functions

func getSystemUptime() (int64, error) {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return 0, err
	}
	return info.Uptime, nil
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func readLastNLines(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	return lines, scanner.Err()
}

func tailFile(filename string, output chan<- string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	// Seek to end of file
	file.Seek(0, io.SeekEnd)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			output <- strings.TrimSpace(line)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func getServiceNetworkStats(serviceName string) (map[string]interface{}, error) {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return nil, err
	}

	return map[string]interface{}{
		"bytes_sent":   stats[0].BytesSent,
		"bytes_recv":   stats[0].BytesRecv,
		"packets_sent": stats[0].PacketsSent,
		"packets_recv": stats[0].PacketsRecv,
		"errors_in":    stats[0].Errin,
		"errors_out":   stats[0].Errout,
	}, nil
}

// PersistMetrics saves metrics to the database
func (s *EnhancedService) PersistMetrics(serviceName string, metrics map[string]interface{}) error {
	if s.db == nil {
		return fmt.Errorf("database not available")
	}

	// Convert metrics to database format
	metricData := map[string]interface{}{
		"service":      serviceName,
		"cpu_usage":    metrics["cpu_usage"],
		"memory_usage": metrics["memory_usage"],
		"latency_ms":   metrics["response_time"],
		"created_at":   time.Now(),
	}

	// Check if there was an error
	if lastError, exists := metrics["last_error"]; exists && lastError != "" {
		metricData["error"] = lastError
	}

	// Insert into metric_data table
	return s.db.Table("metric_data").Create(metricData).Error
}

// StartMetricsCollection starts periodic metrics collection
func (s *EnhancedService) StartMetricsCollection() {
	go func() {
		ticker := time.NewTicker(s.config.MetricsInterval)
		defer ticker.Stop()

		for range ticker.C {
			services := s.serviceManager.GetAllServices()
			for name := range services {
				if metrics, err := s.collectMetricsForServiceReal(name); err == nil {
					s.PersistMetrics(name, metrics)
				}
			}
		}
	}()
}