package messaging

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/delivery/websocket"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/entity"
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
	WebSocketPort    int
	MaxMessageSize   int
	MaxAttachments   int
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
		db:          db,
		config:      config,
		hub:         hub,
	}

	go hub.Run()
	return service
}

// RegisterRoutes registers the service's HTTP routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	messaging := router.Group("/messaging")
	{
		// Channel management
		messaging.POST("/channels", s.CreateChannel)
		messaging.GET("/channels", s.GetChannels)
		messaging.GET("/channels/:id", s.GetChannel)

		// Message management
		messaging.POST("/channels/:id/messages", s.SendMessage)
		messaging.GET("/channels/:id/messages", s.GetMessages)

		// Reactions
		messaging.POST("/channels/:id/messages/:messageId/reactions", s.AddReaction)
		messaging.DELETE("/channels/:id/messages/:messageId/reactions", s.RemoveReaction)

		// WebSocket endpoint
		messaging.GET("/ws", s.HandleWebSocket)
	}
}

// CreateChannel creates a new chat channel
func (s *Service) CreateChannel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Parse user ID
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Set default channel type if not specified
	if req.Type == "" {
		req.Type = "public"
	}

	// Create channel
	channel := &entity.Channel{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatedBy:   uid,
	}

	err = s.db.Create(channel).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	// Add creator as owner
	member := &entity.ChannelMember{
		ChannelID: channel.ID,
		UserID:    uid,
		Role:      "owner",
	}

	err = s.db.Create(member).Error
	if err != nil {
		// Rollback channel creation
		s.db.Delete(channel)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel membership"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Channel created successfully",
		"channel": channel,
	})
}

// GetChannels retrieves all channels for the current user
func (s *Service) GetChannels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination
	var pagination entity.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get channels where user is a member
	var channels []entity.Channel
	offset := (pagination.Page - 1) * pagination.PageSize

	query := s.db.Table("channels").
		Select("channels.*").
		Joins("JOIN channel_members ON channels.id = channel_members.channel_id").
		Where("channel_members.user_id = ? AND channels.is_archived = false", uid).
		Order("channels.updated_at DESC")

	var total int64
	query.Count(&total)

	err = query.Offset(offset).Limit(pagination.PageSize).Find(&channels).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels"})
		return
	}

	// Build response with additional data
	var channelResponses []entity.ChannelResponse
	for _, channel := range channels {
		var memberCount int64
		var messageCount int64

		s.db.Model(&entity.ChannelMember{}).Where("channel_id = ?", channel.ID).Count(&memberCount)
		s.db.Model(&entity.Message{}).Where("channel_id = ? AND is_deleted = false", channel.ID).Count(&messageCount)

		channelResponses = append(channelResponses, entity.ChannelResponse{
			Channel:      &channel,
			MemberCount:  int(memberCount),
			MessageCount: int(messageCount),
		})
	}

	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	response := entity.ChannelListResponse{
		Channels: channelResponses,
		Pagination: entity.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    pagination.Page < totalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetChannel retrieves a specific channel
func (s *Service) GetChannel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelIDParam := c.Param("id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is a member of the channel
	var membership entity.ChannelMember
	err = s.db.Where("channel_id = ? AND user_id = ?", channelID, uid).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this channel"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify channel access"})
		}
		return
	}

	// Get channel details
	var channel entity.Channel
	err = s.db.Where("id = ? AND is_archived = false", channelID).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel"})
		}
		return
	}

	// Get additional channel data
	var memberCount int64
	var messageCount int64

	s.db.Model(&entity.ChannelMember{}).Where("channel_id = ?", channelID).Count(&memberCount)
	s.db.Model(&entity.Message{}).Where("channel_id = ? AND is_deleted = false", channelID).Count(&messageCount)

	response := entity.ChannelResponse{
		Channel:      &channel,
		MemberCount:  int(memberCount),
		MessageCount: int(messageCount),
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage sends a message to a channel
func (s *Service) SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelIDParam := c.Param("id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var req entity.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message data", "details": err.Error()})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is a member of the channel
	var membership entity.ChannelMember
	err = s.db.Where("channel_id = ? AND user_id = ?", channelID, uid).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this channel"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify channel access"})
		}
		return
	}

	// Set default message type if not specified
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	// Create message
	message := &entity.Message{
		ChannelID:   channelID,
		UserID:      uid,
		Content:     strings.TrimSpace(req.Content),
		MessageType: req.MessageType,
	}

	// Handle parent message for replies
	if req.ParentID != "" {
		parentID, err := uuid.Parse(req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent message ID"})
			return
		}

		// Verify parent message exists and is in the same channel
		var parentMessage entity.Message
		err = s.db.Where("id = ? AND channel_id = ? AND is_deleted = false", parentID, channelID).First(&parentMessage).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent message not found or invalid"})
			return
		}

		message.ParentID = &parentID
	}

	err = s.db.Create(message).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Update channel's updated_at timestamp
	s.db.Model(&entity.Channel{}).Where("id = ?", channelID).Update("updated_at", time.Now())


	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
		"data":    message,
	})
}

// GetMessages retrieves messages from a channel
func (s *Service) GetMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelIDParam := c.Param("id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	// Parse pagination
	var pagination entity.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is a member of the channel
	var membership entity.ChannelMember
	err = s.db.Where("channel_id = ? AND user_id = ?", channelID, uid).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this channel"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify channel access"})
		}
		return
	}

	// Get messages with pagination
	var messages []entity.Message
	offset := (pagination.Page - 1) * pagination.PageSize

	query := s.db.Where("channel_id = ? AND is_deleted = false", channelID).
		Order("created_at DESC")

	var total int64
	query.Model(&entity.Message{}).Count(&total)

	err = query.Offset(offset).Limit(pagination.PageSize).Find(&messages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	// Build response with additional data
	var messageResponses []entity.MessageResponse
	for _, message := range messages {
		var reactionCount int64
		var replyCount int64

		s.db.Model(&entity.MessageReaction{}).Where("message_id = ?", message.ID).Count(&reactionCount)
		s.db.Model(&entity.Message{}).Where("parent_id = ? AND is_deleted = false", message.ID).Count(&replyCount)

		messageResponses = append(messageResponses, entity.MessageResponse{
			Message:       &message,
			ReactionCount: int(reactionCount),
			ReplyCount:    int(replyCount),
		})
	}

	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	response := entity.MessageListResponse{
		Messages: messageResponses,
		Pagination: entity.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    pagination.Page < totalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	c.JSON(http.StatusOK, response)
}

// AddReaction adds a reaction to a message
func (s *Service) AddReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageIDParam := c.Param("messageId")
	messageID, err := uuid.Parse(messageIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var req entity.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reaction data", "details": err.Error()})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify message exists and user has access
	var message entity.Message
	err = s.db.Where("id = ? AND is_deleted = false", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve message"})
		}
		return
	}

	// Check if user is a member of the channel
	var membership entity.ChannelMember
	err = s.db.Where("channel_id = ? AND user_id = ?", message.ChannelID, uid).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this channel"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify channel access"})
		}
		return
	}

	// Check if reaction already exists
	var existingReaction entity.MessageReaction
	err = s.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, uid, req.Emoji).First(&existingReaction).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Reaction already exists"})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing reaction"})
		return
	}

	// Create reaction
	reaction := &entity.MessageReaction{
		MessageID: messageID,
		UserID:    uid,
		Emoji:     req.Emoji,
	}

	err = s.db.Create(reaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reaction added successfully",
		"reaction": reaction,
	})
}

// RemoveReaction removes a reaction from a message
func (s *Service) RemoveReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageIDParam := c.Param("messageId")
	messageID, err := uuid.Parse(messageIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	emoji := c.Query("emoji")
	if emoji == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Emoji parameter is required"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find and remove the reaction
	err = s.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, uid, emoji).Delete(&entity.MessageReaction{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction removed successfully",
	})
}

// HandleWebSocket handles WebSocket connections
func (s *Service) HandleWebSocket(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "WebSocket implementation pending",
		"message": "Real-time messaging via WebSocket will be implemented in the next phase",
	})
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
