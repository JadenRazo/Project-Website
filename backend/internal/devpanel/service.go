package devpanel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	projectservice "github.com/JadenRazo/Project-Website/backend/internal/projects/service"
	"github.com/JadenRazo/Project-Website/backend/internal/visitor"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// Service handles devpanel operations
type Service struct {
	*core.BaseService
	serviceManager   *core.ServiceManager
	visitorService   *visitor.Service
	projectService   *project.Service
	metricsCollector *MetricsCollector
	config           Config
}

// Config holds devpanel configuration
type Config struct {
	AdminToken      string
	MetricsInterval time.Duration
	MaxLogLines     int
	LogRetention    time.Duration
}

// SystemStats represents system-wide statistics
type SystemStats struct {
	CPUUsage      float64   `json:"cpuUsage"`
	MemoryUsage   float64   `json:"memoryUsage"`
	DiskUsage     float64   `json:"diskUsage"`
	Uptime        string    `json:"uptime"`
	GoVersion     string    `json:"goVersion"`
	NumGoroutines int       `json:"numGoroutines"`
	LastUpdate    time.Time `json:"lastUpdate"`
}

// ServiceStats represents statistics for a specific service
type ServiceStats struct {
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	Uptime              string    `json:"uptime"`
	MemoryUsage         float64   `json:"memoryUsage"`
	CPUUsage            float64   `json:"cpuUsage"`
	RequestCount        int64     `json:"requestCount"`
	ErrorCount          int64     `json:"errorCount"`
	AverageResponseTime float64   `json:"averageResponseTime"`
	LastError           string    `json:"lastError"`
	LastUpdate          time.Time `json:"lastUpdate"`
}

// NewService creates a new devpanel service
func NewService(serviceManager *core.ServiceManager, visitorService *visitor.Service, metricsCollector *MetricsCollector, config Config) *Service {
	// Initialize in-memory project service
	memProjectService := projectservice.NewMemoryProjectService()
	projectService := project.NewService(memProjectService)

	return &Service{
		BaseService:      core.NewBaseService("devpanel"),
		serviceManager:   serviceManager,
		visitorService:   visitorService,
		projectService:   projectService,
		metricsCollector: metricsCollector,
		config:           config,
	}
}

// RegisterRoutes registers the devpanel routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	// System overview
	router.GET("/system", s.getSystemStats)

	// Service management
	router.GET("/services", s.getServices)
	router.GET("/services/:name", s.getServiceStats)
	router.POST("/services/:name/start", s.startService)
	router.POST("/services/:name/stop", s.stopService)
	router.POST("/services/:name/restart", s.restartService)

	// Logs
	router.GET("/logs/:service", s.getServiceLogs)
	router.GET("/logs/:service/stream", s.streamServiceLogs)

	// Configuration
	router.GET("/config/:service", s.getServiceConfig)
	router.PUT("/config/:service", s.updateServiceConfig)

	// Performance metrics
	router.GET("/metrics/:service", s.getServiceMetrics)
	router.GET("/metrics/:service/history", s.getServiceMetricsHistory)

	// Health checks
	router.GET("/health/:service", s.getServiceHealth)
	router.POST("/health/:service/check", s.runHealthCheck)

	// Visitor analytics
	router.GET("/visitors/stats", s.getVisitorStats)
	router.GET("/visitors/timeline", s.getVisitorTimeline)
	router.GET("/visitors/realtime", s.getVisitorRealtime)
	router.GET("/visitors/locations", s.getVisitorLocations)
	router.GET("/visitors/breakdown", s.getVisitorBreakdown)

	// Project management
	router.GET("/projects", s.listProjects)
	router.POST("/projects", s.createProject)
	router.GET("/projects/:id", s.getProject)
	router.PUT("/projects/:id", s.updateProject)
	router.DELETE("/projects/:id", s.deleteProject)
	router.POST("/projects/:id/archive", s.archiveProject)
	router.POST("/projects/:id/activate", s.activateProject)
}

// getSystemStats returns system-wide statistics
func (s *Service) getSystemStats(c *gin.Context) {
	stats, err := s.collectSystemStats()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}

// getServices returns information about all services
func (s *Service) getServices(c *gin.Context) {
	services := s.serviceManager.GetAllServices()
	stats := make([]ServiceStats, 0, len(services))

	for name, service := range services {
		serviceStats, err := s.collectServiceStats(name, service)
		if err != nil {
			continue
		}
		stats = append(stats, serviceStats)
	}

	c.JSON(200, stats)
}

// getServiceStats returns detailed statistics for a specific service
func (s *Service) getServiceStats(c *gin.Context) {
	name := c.Param("name")
	services := s.serviceManager.GetAllServices()
	service, exists := services[name]
	if !exists {
		c.JSON(404, gin.H{"error": "service not found"})
		return
	}

	stats, err := s.collectServiceStats(name, service)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, stats)
}

// collectSystemStats gathers real system-wide statistics
func (s *Service) collectSystemStats() (SystemStats, error) {
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
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err == nil {
		uptime := time.Duration(info.Uptime) * time.Second
		days := int(uptime.Hours() / 24)
		hours := int(uptime.Hours()) % 24
		minutes := int(uptime.Minutes()) % 60

		if days > 0 {
			stats.Uptime = fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
		} else if hours > 0 {
			stats.Uptime = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			stats.Uptime = fmt.Sprintf("%dm", minutes)
		}
	}

	return stats, nil
}

// collectServiceStats gathers statistics for a specific service
func (s *Service) collectServiceStats(name string, service core.Service) (ServiceStats, error) {
	stats := ServiceStats{
		Name:       name,
		Status:     "unknown",
		LastUpdate: time.Now(),
	}

	// Get service status
	status, err := s.serviceManager.GetServiceStatus(name)
	if err == nil {
		stats.Status = "running"
		if !status.Running {
			stats.Status = "stopped"
		}
		stats.Uptime = status.Uptime.String()
		if len(status.Errors) > 0 {
			stats.LastError = status.Errors[len(status.Errors)-1]
		}
	}

	// Get process stats if running
	if status.Running {
		if process, err := process.NewProcess(int32(os.Getpid())); err == nil {
			if memPercent, err := process.MemoryPercent(); err == nil {
				stats.MemoryUsage = float64(memPercent)
			}
			if cpuPercent, err := process.CPUPercent(); err == nil {
				stats.CPUUsage = cpuPercent
			}
		}
	}

	if baseService, ok := service.(*core.BaseService); ok {
		stats.RequestCount = baseService.GetRequestCount()
		stats.ErrorCount = baseService.GetErrorCount()
	}

	return stats, nil
}

// getServiceLogs returns logs for a specific service
func (s *Service) getServiceLogs(c *gin.Context) {
	serviceName := c.Param("service")

	// Get limit from query params
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	// Get logs from the service
	logs, err := s.getLogsForService(serviceName, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"service": serviceName,
		"logs":    logs,
		"limit":   limit,
	})
}

// streamServiceLogs streams logs for a specific service
func (s *Service) streamServiceLogs(c *gin.Context) {
	serviceName := c.Param("service")

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Create a channel for log streaming
	logChan := make(chan string, 100)
	defer close(logChan)

	// Start streaming logs
	go s.streamLogsForService(serviceName, logChan)

	// Send logs to client
	c.Stream(func(w io.Writer) bool {
		select {
		case log := <-logChan:
			c.SSEvent("log", log)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// getServiceConfig returns configuration for a specific service
func (s *Service) getServiceConfig(c *gin.Context) {
	serviceName := c.Param("service")

	// Get service configuration
	config, err := s.getConfigForService(serviceName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"service": serviceName,
		"config":  config,
	})
}

// updateServiceConfig updates configuration for a specific service
func (s *Service) updateServiceConfig(c *gin.Context) {
	serviceName := c.Param("service")

	var configUpdate map[string]interface{}
	if err := c.ShouldBindJSON(&configUpdate); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Update service configuration
	if err := s.updateConfigForService(serviceName, configUpdate); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("Configuration updated for service %s", serviceName),
	})
}

// getServiceMetrics returns performance metrics for a specific service
func (s *Service) getServiceMetrics(c *gin.Context) {
	serviceName := c.Param("service")

	metrics, err := s.collectMetricsForService(serviceName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, metrics)
}

// getServiceMetricsHistory returns historical performance metrics
func (s *Service) getServiceMetricsHistory(c *gin.Context) {
	serviceName := c.Param("service")

	// Get time range from query params
	duration := c.DefaultQuery("duration", "1h")

	history, err := s.getMetricsHistoryForService(serviceName, duration)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"service":  serviceName,
		"duration": duration,
		"history":  history,
	})
}

// getServiceHealth returns health status for a specific service
func (s *Service) getServiceHealth(c *gin.Context) {
	serviceName := c.Param("service")

	health, err := s.checkHealthForService(serviceName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, health)
}

// runHealthCheck triggers a health check for a specific service
func (s *Service) runHealthCheck(c *gin.Context) {
	serviceName := c.Param("service")

	result, err := s.performHealthCheckForService(serviceName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"service":   serviceName,
		"timestamp": time.Now(),
		"result":    result,
	})
}

// startService starts a specific service
func (s *Service) startService(c *gin.Context) {
	name := c.Param("name")
	if err := s.serviceManager.StartService(name); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to start service: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("Service %s started successfully", name),
	})
}

// stopService stops a specific service
func (s *Service) stopService(c *gin.Context) {
	name := c.Param("name")
	if err := s.serviceManager.StopService(name); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to stop service: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("Service %s stopped successfully", name),
	})
}

// restartService restarts a specific service
func (s *Service) restartService(c *gin.Context) {
	name := c.Param("name")

	// First stop the service
	if err := s.serviceManager.StopService(name); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to stop service: %v", err),
		})
		return
	}

	// Then start it again
	if err := s.serviceManager.StartService(name); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to restart service: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("Service %s restarted successfully", name),
	})
}

// Helper methods for service operations

// getLogsForService retrieves real logs for a specific service
func (s *Service) getLogsForService(serviceName string, limit int) ([]string, error) {
	// Try using LogManager if available
	logPath := filepath.Join("logs", "services", serviceName+".log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// Try alternative log paths
		logPath = filepath.Join("logs", serviceName+".log")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			// Try the main API log for the service
			logPath = filepath.Join("logs", "api.log")
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				return []string{}, nil
			}
		}
	}

	// Read the log file
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter logs for the specific service if reading from a shared log
		if serviceName != "" && !strings.Contains(line, serviceName) {
			continue
		}
		lines = append(lines, line)
		if len(lines) > limit*2 {
			// Keep a sliding window
			lines = lines[1:]
		}
	}

	// Return the last N lines
	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}

	return lines, scanner.Err()
}

// streamLogsForService streams real logs for a specific service
func (s *Service) streamLogsForService(serviceName string, logChan chan<- string) {
	logPath := filepath.Join("logs", "services", serviceName+".log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = filepath.Join("logs", serviceName+".log")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			logPath = filepath.Join("logs", "api.log")
		}
	}

	// Open file for tailing
	file, err := os.Open(logPath)
	if err != nil {
		logChan <- fmt.Sprintf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	// Seek to end of file to start streaming new entries
	file.Seek(0, io.SeekEnd)
	reader := bufio.NewReader(file)

	// Stream new log entries
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(30 * time.Second) // Stop after 30 seconds
	for {
		select {
		case <-ticker.C:
			line, err := reader.ReadString('\n')
			if err == nil {
				// Filter for service if reading shared log
				if serviceName == "" || strings.Contains(line, serviceName) {
					logChan <- strings.TrimSpace(line)
				}
			}
		case <-timeout:
			return
		}
	}
}

// getConfigForService retrieves configuration for a specific service
func (s *Service) getConfigForService(serviceName string) (map[string]interface{}, error) {
	// Return mock configuration for demonstration
	config := map[string]interface{}{
		"port":           8080,
		"debug":          false,
		"maxConnections": 100,
		"timeout":        "30s",
	}
	return config, nil
}

// updateConfigForService updates configuration for a specific service
func (s *Service) updateConfigForService(serviceName string, config map[string]interface{}) error {
	// Mock configuration update for demonstration
	return nil
}

// collectMetricsForService collects metrics for a specific service
func (s *Service) collectMetricsForService(serviceName string) (map[string]interface{}, error) {
	// Get service stats
	serviceStats, err := s.collectServiceStats(serviceName, nil)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"service":       serviceName,
		"timestamp":     time.Now(),
		"cpu_usage":     serviceStats.CPUUsage,
		"memory_usage":  serviceStats.MemoryUsage,
		"request_count": serviceStats.RequestCount,
		"error_count":   serviceStats.ErrorCount,
		"status":        serviceStats.Status,
		"uptime":        serviceStats.Uptime,
	}, nil
}

// getMetricsHistoryForService retrieves historical metrics for a service
func (s *Service) getMetricsHistoryForService(serviceName string, duration string) ([]map[string]interface{}, error) {
	if s.metricsCollector == nil {
		return []map[string]interface{}{}, nil
	}

	durationParsed := 1 * time.Hour
	switch duration {
	case "1h":
		durationParsed = 1 * time.Hour
	case "6h":
		durationParsed = 6 * time.Hour
	case "24h":
		durationParsed = 24 * time.Hour
	case "7d":
		durationParsed = 7 * 24 * time.Hour
	}

	endTime := time.Now()
	startTime := endTime.Add(-durationParsed)

	metricsData := s.metricsCollector.GetMetrics(serviceName, startTime, endTime)

	history := make([]map[string]interface{}, 0)

	pointsByTime := make(map[time.Time]map[string]interface{})

	for key, points := range metricsData {
		metricName := s.metricsCollector.getMetricName(key)
		for _, point := range points {
			timestamp := point.Timestamp.Truncate(time.Minute)
			if _, exists := pointsByTime[timestamp]; !exists {
				pointsByTime[timestamp] = map[string]interface{}{
					"timestamp": timestamp,
				}
			}
			pointsByTime[timestamp][metricName] = point.Value
		}
	}

	for _, point := range pointsByTime {
		history = append(history, point)
	}

	return history, nil
}

// checkHealthForService checks health status for a specific service
func (s *Service) checkHealthForService(serviceName string) (map[string]interface{}, error) {
	services := s.serviceManager.GetAllServices()
	_, exists := services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}

	status, err := s.serviceManager.GetServiceStatus(serviceName)
	if err != nil {
		return nil, err
	}

	health := map[string]interface{}{
		"service": serviceName,
		"healthy": status.Running && len(status.Errors) == 0,
		"status":  "unknown",
		"checks": map[string]interface{}{
			"running": status.Running,
			"errors":  len(status.Errors),
			"uptime":  status.Uptime.String(),
		},
	}

	if status.Running {
		if len(status.Errors) == 0 {
			health["status"] = "healthy"
		} else {
			health["status"] = "degraded"
		}
	} else {
		health["status"] = "unhealthy"
	}

	return health, nil
}

// performHealthCheckForService performs a health check for a specific service
func (s *Service) performHealthCheckForService(serviceName string) (map[string]interface{}, error) {
	// First get current health
	health, err := s.checkHealthForService(serviceName)
	if err != nil {
		return nil, err
	}

	// Return current health status

	return health, nil
}

// Visitor Analytics Handlers

// getVisitorStats returns comprehensive visitor statistics
func (s *Service) getVisitorStats(c *gin.Context) {
	stats, err := s.visitorService.GetVisitorStats(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}

// getVisitorTimeline returns visitor timeline data
func (s *Service) getVisitorTimeline(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	interval := c.DefaultQuery("interval", "day")

	data, err := s.visitorService.GetTimelineData(c.Request.Context(), period, interval)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"period":   period,
		"interval": interval,
		"data":     data,
	})
}

// getVisitorRealtime returns real-time visitor data
func (s *Service) getVisitorRealtime(c *gin.Context) {
	count := s.visitorService.GetRealTimeCount(c.Request.Context())
	c.JSON(200, gin.H{
		"activeVisitors": count,
	})
}

// getVisitorLocations returns visitor location distribution
func (s *Service) getVisitorLocations(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	locations, err := s.visitorService.GetLocationDistribution(c.Request.Context(), period)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"period":    period,
		"locations": locations,
	})
}

// getVisitorBreakdown returns device and browser breakdown
func (s *Service) getVisitorBreakdown(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	breakdown, err := s.visitorService.GetDeviceBrowserStats(c.Request.Context(), period)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, breakdown)
}
