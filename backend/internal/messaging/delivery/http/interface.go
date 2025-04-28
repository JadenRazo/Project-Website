package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
)

// HTTPHandler is the interface for HTTP handlers
type HTTPHandler interface {
	Register(router chi.Router)
}

// GinHTTPHandler is the interface for Gin HTTP handlers
type GinHTTPHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// ReadReceiptHTTPHandler is the interface for read receipt HTTP handlers
type ReadReceiptHTTPHandler interface {
	GinHTTPHandler
	MarkMessageAsRead(c *gin.Context)
	MarkChannelAsRead(c *gin.Context)
	GetUnreadCount(c *gin.Context)
	GetTotalUnreadCount(c *gin.Context)
}

// ChannelHTTPHandler is the interface for channel HTTP handlers
type ChannelHTTPHandler interface {
	HTTPHandler
	ListChannels(w http.ResponseWriter, r *http.Request)
	CreateChannel(w http.ResponseWriter, r *http.Request)
	GetChannel(w http.ResponseWriter, r *http.Request)
	UpdateChannel(w http.ResponseWriter, r *http.Request)
	DeleteChannel(w http.ResponseWriter, r *http.Request)
}

// MessageHTTPHandler is the interface for message HTTP handlers
type MessageHTTPHandler interface {
	HTTPHandler
	ListChannelMessages(w http.ResponseWriter, r *http.Request)
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetMessage(w http.ResponseWriter, r *http.Request)
	UpdateMessage(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
}
