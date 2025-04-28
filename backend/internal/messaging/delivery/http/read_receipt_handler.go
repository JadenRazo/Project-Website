package http

import (
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ReadReceiptHandler handles read receipt operations
type ReadReceiptHandler struct {
	readReceiptUseCase domain.ReadReceiptUseCase
	validator          *validator.Validate
}

// NewReadReceiptHandler creates a new read receipt handler
func NewReadReceiptHandler(readReceiptUseCase domain.ReadReceiptUseCase, validator *validator.Validate) *ReadReceiptHandler {
	if readReceiptUseCase == nil {
		panic("readReceiptUseCase cannot be nil")
	}
	if validator == nil {
		panic("validator cannot be nil")
	}

	return &ReadReceiptHandler{
		readReceiptUseCase: readReceiptUseCase,
		validator:          validator,
	}
}

// RegisterRoutes registers the read receipt routes with the given router group
func (h *ReadReceiptHandler) RegisterRoutes(r *gin.RouterGroup) {
	if r == nil {
		panic("router cannot be nil")
	}

	r.POST("/channels/:channelId/read", h.MarkChannelAsRead)
	r.POST("/messages/:messageId/read", h.MarkMessageAsRead)
	r.GET("/unread", h.GetUserUnreadCount)
	r.POST("/mark-all-read", h.MarkAllAsRead)
}

// MarkChannelAsRead marks all messages in a channel as read for the current user
func (h *ReadReceiptHandler) MarkChannelAsRead(c *gin.Context) {
	channelID := c.Param("channelId")
	if channelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel ID is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.readReceiptUseCase.MarkChannelAsRead(c.Request.Context(), channelID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel marked as read"})
}

// MarkMessageAsRead marks a specific message as read for the current user
func (h *ReadReceiptHandler) MarkMessageAsRead(c *gin.Context) {
	messageID := c.Param("messageId")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.readReceiptUseCase.MarkMessageAsRead(c.Request.Context(), messageID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}

// GetUserUnreadCount gets the count of unread messages for the current user
func (h *ReadReceiptHandler) GetUserUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	counts, err := h.readReceiptUseCase.GetUnreadCounts(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GetTotalUnreadCount gets the total count of unread messages for the current user
func (h *ReadReceiptHandler) GetTotalUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	count, err := h.readReceiptUseCase.GetTotalUnreadCount(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_unread": count})
}

// MarkAllAsRead marks all messages as read for the current user
func (h *ReadReceiptHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.readReceiptUseCase.MarkAllAsRead(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All messages marked as read"})
}
