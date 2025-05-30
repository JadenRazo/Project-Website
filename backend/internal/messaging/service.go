package messaging

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    "github.com/JadenRazo/Project-Website/backend/internal/core"
    "github.com/JadenRazo/Project-Website/backend/internal/messaging/delivery/websocket"
)

// Service handles messaging operations
type Service struct {
    *core.BaseService
    db     *gorm.DB
    config Config
    hub    *Hub
}

// Config holds messaging service configuration
type Config struct {
    WebSocketPort int
    MaxMessageSize int
    MaxAttachments int
    AllowedFileTypes []string
}

// Hub is an alias for the websocket Hub
type Hub = websocket.Hub

// Client is an alias for the websocket Client  
type Client = websocket.Client

// NewService creates a new messaging service
func NewService(db *gorm.DB, config Config) *Service {
    // Create connection manager with default config
    connManager := websocket.NewConnectionManager(&websocket.WebSocketConfig{
        MaxConnections:     1000,
        RateLimitPerMinute: 100,
        ConnectionTimeout:  5 * time.Minute,
    })
    
    hub := websocket.NewHub(connManager)

    service := &Service{
        BaseService: core.NewBaseService("messaging"),
        db:     db,
        config: config,
        hub:    hub,
    }

    go hub.Run()
    return service
}

// RegisterRoutes registers the service's HTTP routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/channels", s.CreateChannel)
    router.GET("/channels", s.GetChannels)
    router.GET("/channels/:id", s.GetChannel)
    router.POST("/channels/:id/messages", s.SendMessage)
    router.GET("/channels/:id/messages", s.GetMessages)
    router.POST("/channels/:id/messages/:messageId/reactions", s.AddReaction)
    router.DELETE("/channels/:id/messages/:messageId/reactions", s.RemoveReaction)
    router.GET("/ws", s.HandleWebSocket)
}

// CreateChannel creates a new chat channel
func (s *Service) CreateChannel(c *gin.Context) {
    // Implementation
}

// GetChannels retrieves all channels for the current user
func (s *Service) GetChannels(c *gin.Context) {
    // Implementation
}

// GetChannel retrieves a specific channel
func (s *Service) GetChannel(c *gin.Context) {
    // Implementation
}

// SendMessage sends a message to a channel
func (s *Service) SendMessage(c *gin.Context) {
    // Implementation
}

// GetMessages retrieves messages from a channel
func (s *Service) GetMessages(c *gin.Context) {
    // Implementation
}

// AddReaction adds a reaction to a message
func (s *Service) AddReaction(c *gin.Context) {
    // Implementation
}

// RemoveReaction removes a reaction from a message
func (s *Service) RemoveReaction(c *gin.Context) {
    // Implementation
}

// HandleWebSocket handles WebSocket connections
func (s *Service) HandleWebSocket(c *gin.Context) {
    // Implementation
}

// HealthCheck performs service-specific health checks
func (s *Service) HealthCheck() error {
    if err := s.BaseService.HealthCheck(); err != nil {
        return err
    }

    // Check database connection
    sqlDB, err := s.db.DB()
    if err != nil {
        s.AddError(err)
        return err
    }

    if err := sqlDB.Ping(); err != nil {
        s.AddError(err)
        return err
    }

    // Check WebSocket server
    if s.hub == nil {
        err := fmt.Errorf("WebSocket hub not initialized")
        s.AddError(err)
        return err
    }

    return nil
}

// Start initializes the WebSocket server
func (s *Service) Start() error {
    if err := s.BaseService.Start(); err != nil {
        return err
    }

    // Start WebSocket server
    go func() {
        addr := fmt.Sprintf(":%d", s.config.WebSocketPort)
        if err := http.ListenAndServe(addr, nil); err != nil {
            s.AddError(err)
        }
    }()

    return nil
}

// Stop shuts down the WebSocket server
func (s *Service) Stop() error {
    if err := s.BaseService.Stop(); err != nil {
        return err
    }

    // The hub will handle closing connections when it stops
    return nil
} 