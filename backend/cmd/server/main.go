package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Define missing types
// ReadReceipt represents a record of when a user reads a message
type ReadReceipt struct {
	MessageID uint      `json:"message_id"`
	UserID    uint      `json:"user_id"`
	ReadAt    time.Time `json:"read_at"`
}

// ReadReceiptRepository defines the interface for read receipt operations
type ReadReceiptRepository interface {
	CreateReadReceipt(ctx context.Context, receipt *ReadReceipt) error
	CreateBulkReadReceipts(ctx context.Context, receipts []ReadReceipt) error
	ReadReceiptExists(ctx context.Context, messageID, userID uint) (bool, error)
	GetMessageReadReceipts(ctx context.Context, messageID uint) ([]ReadReceipt, error)
	GetUnreadCount(ctx context.Context, channelID, userID uint) (int64, error)
	DeleteReadReceipts(ctx context.Context, messageID uint) error
	Create(ctx context.Context, receipt *ReadReceipt) error
}

// Mock dependencies to enable compilation without all dependencies
type mockReadReceiptRepository struct{}

func (m *mockReadReceiptRepository) CreateReadReceipt(ctx context.Context, receipt *ReadReceipt) error {
	return nil
}

func (m *mockReadReceiptRepository) CreateBulkReadReceipts(ctx context.Context, receipts []ReadReceipt) error {
	return nil
}

func (m *mockReadReceiptRepository) ReadReceiptExists(ctx context.Context, messageID, userID uint) (bool, error) {
	return false, nil
}

func (m *mockReadReceiptRepository) GetMessageReadReceipts(ctx context.Context, messageID uint) ([]ReadReceipt, error) {
	return []ReadReceipt{}, nil
}

func (m *mockReadReceiptRepository) GetUnreadCount(ctx context.Context, channelID, userID uint) (int64, error) {
	return 0, nil
}

func (m *mockReadReceiptRepository) DeleteReadReceipts(ctx context.Context, messageID uint) error {
	return nil
}

func (m *mockReadReceiptRepository) Create(ctx context.Context, receipt *ReadReceipt) error {
	return nil
}

type mockMessageRepository struct{}

func (m *mockMessageRepository) GetMessage(ctx context.Context, messageID uint) (*Message, error) {
	return &Message{
		ID:        messageID,
		ChannelID: 1,
		SenderID:  2,
		Content:   "Test message",
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockMessageRepository) GetUnreadMessages(ctx context.Context, channelID, userID uint, upToTime time.Time) ([]Message, error) {
	return []Message{}, nil
}

// Minimal Message struct to avoid dependencies
type Message struct {
	ID        uint      `json:"id"`
	ChannelID uint      `json:"channel_id"`
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Define WebSocket Hub
type Hub struct{}

// Auth package mock
var auth = struct {
	GetUserID func(ctx context.Context) (uint, bool)
}{
	GetUserID: func(ctx context.Context) (uint, bool) {
		// Mock implementation that always returns user ID 1
		return 1, true
	},
}

// Auth middleware mock
var authMiddleware = struct {
	AuthMiddleware func() gin.HandlerFunc
}{
	AuthMiddleware: func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Mock implementation that always authenticates
			c.Next()
		}
	},
}

// ReadReceiptService handles read receipt operations
type ReadReceiptService struct {
	receiptRepo ReadReceiptRepository
	messageRepo *mockMessageRepository
	hub         *Hub
}

// NewReadReceiptService creates a new read receipt service
func NewReadReceiptService(
	receiptRepo ReadReceiptRepository,
	messageRepo *mockMessageRepository,
	hub *Hub,
) *ReadReceiptService {
	return &ReadReceiptService{
		receiptRepo: receiptRepo,
		messageRepo: messageRepo,
		hub:         hub,
	}
}

// MarkAsRead marks a message as read by a user
func (s *ReadReceiptService) MarkAsRead(ctx context.Context, messageID, userID uint) error {
	// Get the message to verify it exists and get sender ID
	message, err := s.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Don't create read receipts for your own messages
	if message.SenderID == userID {
		return nil
	}

	// Check if receipt already exists to avoid duplicates
	exists, err := s.receiptRepo.ReadReceiptExists(ctx, messageID, userID)
	if err != nil {
		return fmt.Errorf("failed to check read receipt: %w", err)
	}

	if !exists {
		// Create the read receipt
		receipt := &ReadReceipt{
			MessageID: messageID,
			UserID:    userID,
			ReadAt:    time.Now(),
		}

		if err := s.receiptRepo.CreateReadReceipt(ctx, receipt); err != nil {
			return fmt.Errorf("failed to create read receipt: %w", err)
		}

		// Notify through WebSocket if hub exists
		if s.hub != nil {
			// In a real implementation, we would notify the sender that their message was read
			log.Printf("User %d read message %d in channel %d", userID, messageID, message.ChannelID)
		}
	}

	return nil
}

// MarkChannelAsRead marks all messages in a channel as read up to a certain time
func (s *ReadReceiptService) MarkChannelAsRead(ctx context.Context, channelID, userID uint, upToTime time.Time) error {
	// Get unread messages in channel before the specified time
	messages, err := s.messageRepo.GetUnreadMessages(ctx, channelID, userID, upToTime)
	if err != nil {
		return fmt.Errorf("failed to get unread messages: %w", err)
	}

	// Nothing to mark as read
	if len(messages) == 0 {
		return nil
	}

	// Create batch read receipts
	var receipts []ReadReceipt
	for _, message := range messages {
		// Skip messages sent by the user
		if message.SenderID == userID {
			continue
		}

		receipts = append(receipts, ReadReceipt{
			MessageID: message.ID,
			UserID:    userID,
			ReadAt:    time.Now(),
		})
	}

	// Save read receipts in batch
	if len(receipts) > 0 {
		if err := s.receiptRepo.CreateBulkReadReceipts(ctx, receipts); err != nil {
			return fmt.Errorf("failed to create read receipts: %w", err)
		}

		// Notify through WebSocket for the latest message only
		if s.hub != nil && len(messages) > 0 {
			// In a real implementation, we would send a notification
			// Just log it for now to avoid linter errors
			log.Printf("User %d marked channel %d as read", userID, channelID)
		}
	}

	return nil
}

// GetUnreadMessageCount gets the count of unread messages for a user in a channel
func (s *ReadReceiptService) GetUnreadMessageCount(ctx context.Context, channelID, userID uint) (int, error) {
	count, err := s.receiptRepo.GetUnreadCount(ctx, channelID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return int(count), nil
}

func main() {
	// Create router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))

	// Initialize database - use the db package
	_, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize mock dependencies
	readReceiptRepo := &mockReadReceiptRepository{}
	messageRepo := &mockMessageRepository{}

	// Initialize WebSocket hub
	hub := &Hub{}

	// Initialize read receipt service
	readReceiptService := NewReadReceiptService(readReceiptRepo, messageRepo, hub)

	// Setup API routes with auth middleware
	api := router.Group("/api")
	api.Use(authMiddleware.AuthMiddleware())

	// Read Receipt routes
	readReceipts := api.Group("/read-receipts")
	{
		// Mark a message as read
		readReceipts.POST("/messages/:messageID", func(c *gin.Context) {
			// Get message ID from URL parameter
			messageIDStr := c.Param("messageID")
			messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
				return
			}

			// Get user ID from context (set by auth middleware)
			userID, exists := auth.GetUserID(c.Request.Context())
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
				return
			}

			// Mark the message as read
			err = readReceiptService.MarkAsRead(c.Request.Context(), uint(messageID), userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Status(http.StatusNoContent)
		})

		// Mark all messages in a channel as read
		readReceipts.POST("/channels/:channelID", func(c *gin.Context) {
			// Get channel ID from URL parameter
			channelIDStr := c.Param("channelID")
			channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
				return
			}

			// Get user ID from context (set by auth middleware)
			userID, exists := auth.GetUserID(c.Request.Context())
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
				return
			}

			// Mark all messages in the channel as read up to the current time
			err = readReceiptService.MarkChannelAsRead(c.Request.Context(), uint(channelID), userID, time.Now())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Status(http.StatusNoContent)
		})

		// Get unread count for a channel
		readReceipts.GET("/channels/:channelID/unread", func(c *gin.Context) {
			// Get channel ID from URL parameter
			channelIDStr := c.Param("channelID")
			channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
				return
			}

			// Get user ID from context (set by auth middleware)
			userID, exists := auth.GetUserID(c.Request.Context())
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
				return
			}

			// Get unread count
			count, err := readReceiptService.GetUnreadMessageCount(c.Request.Context(), uint(channelID), userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"unread_count": count})
		})
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		c.String(http.StatusOK, "WebSocket connection point")
	})

	// Start the server
	port := 8080 // Replace with configured port
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
