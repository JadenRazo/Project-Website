package devpanel

import (
    "fmt"
    "runtime"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/process"
    
    "github.com/JadenRazo/Project-Website/backend/internal/core"
)

// Service handles devpanel operations
type Service struct {
    *core.BaseService
    serviceManager *core.ServiceManager
    config        Config
}

// Config holds devpanel configuration
type Config struct {
    AdminToken    string
    MetricsInterval time.Duration
    MaxLogLines    int
    LogRetention   time.Duration
}

// SystemStats represents system-wide statistics
type SystemStats struct {
    CPUUsage    float64   `json:"cpu_usage"`
    MemoryUsage float64   `json:"memory_usage"`
    DiskUsage   float64   `json:"disk_usage"`
    Uptime      string    `json:"uptime"`
    GoVersion   string    `json:"go_version"`
    NumGoroutines int     `json:"num_goroutines"`
    LastUpdate   time.Time `json:"last_update"`
}

// ServiceStats represents statistics for a specific service
type ServiceStats struct {
    Name         string    `json:"name"`
    Status       string    `json:"status"`
    Uptime       string    `json:"uptime"`
    MemoryUsage  float64   `json:"memory_usage"`
    CPUUsage     float64   `json:"cpu_usage"`
    RequestCount int64     `json:"request_count"`
    ErrorCount   int64     `json:"error_count"`
    LastError    string    `json:"last_error"`
    LastUpdate   time.Time `json:"last_update"`
}

// NewService creates a new devpanel service
func NewService(serviceManager *core.ServiceManager, config Config) *Service {
    return &Service{
        BaseService:    core.NewBaseService("devpanel"),
        serviceManager: serviceManager,
        config:        config,
    }
}

// RegisterRoutes registers the devpanel routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
    // Authentication middleware
    router.Use(s.authMiddleware())
    
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
}

// authMiddleware ensures only authenticated admin requests are allowed
func (s *Service) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != s.config.AdminToken {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
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
    service, err := s.serviceManager.GetServiceStatus(name)
    if err != nil {
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

// collectSystemStats gathers system-wide statistics
func (s *Service) collectSystemStats() (SystemStats, error) {
    stats := SystemStats{
        GoVersion:     runtime.Version(),
        NumGoroutines: runtime.NumGoroutine(),
        LastUpdate:    time.Now(),
    }
    
    // CPU usage
    cpuPercent, err := cpu.Percent(time.Second, false)
    if err == nil && len(cpuPercent) > 0 {
        stats.CPUUsage = cpuPercent[0]
    }
    
    // Memory usage
    if memInfo, err := mem.VirtualMemory(); err == nil {
        stats.MemoryUsage = memInfo.UsedPercent
    }
    
    // Uptime
    if uptime, err := time.ParseDuration("0s"); err == nil {
        stats.Uptime = uptime.String()
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
                stats.MemoryUsage = memPercent
            }
            if cpuPercent, err := process.CPUPercent(); err == nil {
                stats.CPUUsage = cpuPercent
            }
        }
    }
    
    return stats, nil
}

// getServiceLogs returns logs for a specific service
func (s *Service) getServiceLogs(c *gin.Context) {
    service := c.Param("service")
    // Implementation for retrieving service logs
}

// streamServiceLogs streams logs for a specific service
func (s *Service) streamServiceLogs(c *gin.Context) {
    service := c.Param("service")
    // Implementation for streaming service logs
}

// getServiceConfig returns configuration for a specific service
func (s *Service) getServiceConfig(c *gin.Context) {
    service := c.Param("service")
    // Implementation for retrieving service configuration
}

// updateServiceConfig updates configuration for a specific service
func (s *Service) updateServiceConfig(c *gin.Context) {
    service := c.Param("service")
    // Implementation for updating service configuration
}

// getServiceMetrics returns performance metrics for a specific service
func (s *Service) getServiceMetrics(c *gin.Context) {
    service := c.Param("service")
    // Implementation for retrieving service metrics
}

// getServiceMetricsHistory returns historical performance metrics
func (s *Service) getServiceMetricsHistory(c *gin.Context) {
    service := c.Param("service")
    // Implementation for retrieving historical metrics
}

// getServiceHealth returns health status for a specific service
func (s *Service) getServiceHealth(c *gin.Context) {
    service := c.Param("service")
    // Implementation for retrieving service health status
}

// runHealthCheck triggers a health check for a specific service
func (s *Service) runHealthCheck(c *gin.Context) {
    service := c.Param("service")
    // Implementation for running service health check
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