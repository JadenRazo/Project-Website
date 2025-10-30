package gateway

import (
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/gin-gonic/gin"
)

// Gateway handles API routing and middleware
type Gateway struct {
	router     *gin.Engine
	services   *core.ServiceManager
	middleware []gin.HandlerFunc
}

// NewGateway creates a new API gateway instance
func NewGateway(serviceManager *core.ServiceManager) *Gateway {
	return &Gateway{
		router:     gin.Default(),
		services:   serviceManager,
		middleware: make([]gin.HandlerFunc, 0),
	}
}

// RegisterService registers a new service with its routes
func (g *Gateway) RegisterService(name string, routes func(*gin.RouterGroup)) {
	group := g.router.Group("/api/v1/" + name)

	// Add common middleware for all service routes
	group.Use(g.middleware...)

	// Register service-specific routes
	routes(group)
}

// AddMiddleware adds middleware to be applied to all routes
func (g *Gateway) AddMiddleware(middleware gin.HandlerFunc) {
	g.middleware = append(g.middleware, middleware)
}

// Start starts the API gateway server
func (g *Gateway) Start(addr string) error {
	return g.router.Run(addr)
}

// Stop gracefully stops the API gateway
func (g *Gateway) Stop() error {
	// Implement graceful shutdown logic
	return nil
}

// GetRouter returns the underlying gin router
func (g *Gateway) GetRouter() *gin.Engine {
	return g.router
}

// RegisterHealthCheck adds a health check endpoint
func (g *Gateway) RegisterHealthCheck() {
	g.router.GET("/health", func(c *gin.Context) {
		status := make(map[string]interface{})

		// Check all services
		for name, _ := range g.services.GetAllServices() {
			serviceStatus, err := g.services.GetServiceStatus(name)
			if err != nil {
				status[name] = gin.H{
					"error": err.Error(),
				}
				continue
			}

			status[name] = gin.H{
				"running": serviceStatus.Running,
				"uptime":  serviceStatus.Uptime.String(),
				"errors":  serviceStatus.Errors,
			}
		}

		c.JSON(200, gin.H{
			"status":   "ok",
			"services": status,
		})
	})
}

// RegisterMetrics adds a metrics endpoint
func (g *Gateway) RegisterMetrics() {
	g.router.GET("/metrics", func(c *gin.Context) {
		// Implement metrics collection and reporting
		c.JSON(200, gin.H{
			"status":  "ok",
			"metrics": "implemented",
		})
	})
}
