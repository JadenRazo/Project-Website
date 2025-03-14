package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/service"
	"github.com/gin-gonic/gin"
)

// ReadReceiptHandler handles API endpoints for read receipts
type ReadReceiptHandler struct {
	service *service.ReadReceiptService
}

// NewReadReceiptHandler creates a new read receipt handler
func NewReadReceiptHandler(service *service.ReadReceiptService) *ReadReceiptHandler {
	return &ReadReceiptHandler{
		service: service,
	}
}

// MarkMessageAsRead marks a message as read by the authenticated user
func (h *ReadReceiptHandler) MarkMessageAsRead(c *gin.Context) {
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
	err = h.service.MarkAsRead(c.Request.Context(), uint(messageID), userID)
	if err != nil {
		status, response := errors.HandleError(err)
		c.JSON(status, response)
		return
	}

	c.Status(http.StatusNoContent)
}

// MarkChannelAsRead marks all messages in a channel as read up to the current time
func (h *ReadReceiptHandler) MarkChannelAsRead(c *gin.Context) {
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
	err = h.service.MarkChannelAsRead(c.Request.Context(), uint(channelID), userID, time.Now())
	if err != nil {
		status, response := errors.HandleError(err)
		c.JSON(status, response)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUnreadCount gets the count of unread messages in a channel
func (h *ReadReceiptHandler) GetUnreadCount(c *gin.Context) {
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
	count, err := h.service.GetUnreadMessageCount(c.Request.Context(), uint(channelID), userID)
	if err != nil {
		status, response := errors.HandleError(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// RegisterRoutes registers the read receipt routes
func (h *ReadReceiptHandler) RegisterRoutes(router *gin.RouterGroup) {
	readReceipts := router.Group("/read-receipts")
	{
		// Mark a message as read
		readReceipts.POST("/messages/:messageID", h.MarkMessageAsRead)

		// Mark all messages in a channel as read
		readReceipts.POST("/channels/:channelID", h.MarkChannelAsRead)

		// Get unread count for a channel
		readReceipts.GET("/channels/:channelID/unread", h.GetUnreadCount)
	}
}
