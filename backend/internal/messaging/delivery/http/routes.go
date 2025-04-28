package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler contains the HTTP handlers for the messaging module
type Handler struct {
	channelHandler     *ChannelHandler
	messageHandler     *MessageHandler
	readReceiptHandler *ReadReceiptHandler
	authMiddleware     func(next http.Handler) http.Handler
	channelRepo        domain.ChannelRepository
	messageRepo        domain.MessageRepository
	websocketHandler   http.HandlerFunc
}

// NewHandler creates a new HTTP handler for the messaging module
func NewHandler(
	channelHandler *ChannelHandler,
	messageHandler *MessageHandler,
	authMiddleware *auth.Middleware,
	channelRepo domain.ChannelRepository,
	messageRepo domain.MessageRepository,
	readReceiptHandler *ReadReceiptHandler,
	websocketHandler http.HandlerFunc,
) *Handler {
	if channelHandler == nil {
		panic("channelHandler cannot be nil")
	}
	if messageHandler == nil {
		panic("messageHandler cannot be nil")
	}
	if authMiddleware == nil {
		panic("authMiddleware cannot be nil")
	}
	if channelRepo == nil {
		panic("channelRepo cannot be nil")
	}
	if messageRepo == nil {
		panic("messageRepo cannot be nil")
	}

	return &Handler{
		channelHandler:     channelHandler,
		messageHandler:     messageHandler,
		authMiddleware:     authMiddleware.RequireAuth,
		channelRepo:        channelRepo,
		messageRepo:        messageRepo,
		readReceiptHandler: readReceiptHandler,
		websocketHandler:   websocketHandler,
	}
}

// RegisterGinRoutes registers the HTTP routes for the messaging module with Gin
func (h *Handler) RegisterGinRoutes(r *gin.RouterGroup) {
	if h.readReceiptHandler != nil {
		h.readReceiptHandler.RegisterRoutes(r)
	}
}

// Register registers the HTTP routes for the messaging module
func (h *Handler) Register(r chi.Router) {
	// Apply global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// API group with authentication
	r.Route("/api/messaging", func(r chi.Router) {
		// Apply auth middleware to all messaging routes
		r.Use(h.authMiddleware)

		// Health check
		r.Get("/health", h.HealthCheck)

		// Channels
		r.Route("/channels", func(r chi.Router) {
			// Channel routes
			r.Get("/", h.channelHandler.ListChannels)
			r.Post("/", h.channelHandler.CreateChannel)
			r.Get("/direct", h.channelHandler.ListDirectChannels)
			r.Post("/direct", h.channelHandler.CreateDirectChannel)

			// Specific channel routes
			r.Route("/{channelID}", func(r chi.Router) {
				// Apply channel membership middleware to channel routes
				r.Use(RequireChannelMembership(h.channelRepo))

				r.Get("/", h.channelHandler.GetChannel)
				r.Put("/", h.channelHandler.UpdateChannel)
				r.Delete("/", h.channelHandler.DeleteChannel)
				r.Get("/members", h.channelHandler.ListChannelMembers)
				r.Post("/members", h.channelHandler.AddChannelMember)
				r.Delete("/members/{userID}", h.channelHandler.RemoveChannelMember)
				r.Post("/leave", h.channelHandler.LeaveChannel)
				r.Post("/join", h.channelHandler.JoinChannel)
				r.Post("/mark-read", h.channelHandler.MarkChannelAsRead)

				// Admin-only routes
				r.Group(func(r chi.Router) {
					r.Use(RequireChannelAdmin(h.channelRepo))
					r.Post("/archive", h.channelHandler.ArchiveChannel)
					r.Post("/unarchive", h.channelHandler.UnarchiveChannel)
				})

				// Messages
				r.Route("/messages", func(r chi.Router) {
					r.Get("/", h.messageHandler.ListChannelMessages)
					r.Post("/", h.messageHandler.CreateMessage)

					// Specific message routes
					r.Route("/{messageID}", func(r chi.Router) {
						r.Get("/", h.messageHandler.GetMessage)
						r.Get("/thread", h.messageHandler.GetMessageThread)
						r.Post("/thread", h.messageHandler.CreateThreadReply)

						// Message owner or admin only routes
						r.Group(func(r chi.Router) {
							r.Use(ValidateMessageOwnership(h.messageRepo))
							r.Put("/", h.messageHandler.UpdateMessage)
							r.Delete("/", h.messageHandler.DeleteMessage)
						})

						// Reactions
						r.Get("/reactions", h.messageHandler.GetMessageReactions)
						r.Post("/reactions", h.messageHandler.AddReaction)
						r.Delete("/reactions/{emoji}", h.messageHandler.RemoveReaction)
					})
				})
			})
		})

		// User-specific routes
		r.Route("/me", func(r chi.Router) {
			r.Get("/channels", h.channelHandler.GetUserChannels)
			r.Get("/unread", func(w http.ResponseWriter, r *http.Request) {
				// This is a wrapper to convert the Gin handler to a standard http.HandlerFunc
				ctx := r.Context()
				user, err := GetUserFromContext(ctx)
				if err != nil {
					RespondWithError(w, http.StatusUnauthorized, "User not found in context")
					return
				}

				// Mock Gin context with required user info
				c := &gin.Context{
					Request: r,
				}
				c.Set("user_id", user.ID)

				// Call the Gin handler
				h.readReceiptHandler.GetUserUnreadCount(c)
			})
			r.Get("/mentions", h.messageHandler.GetUserMentions)
			r.Post("/mark-all-read", func(w http.ResponseWriter, r *http.Request) {
				// This is a wrapper to convert the Gin handler to a standard http.HandlerFunc
				ctx := r.Context()
				user, err := GetUserFromContext(ctx)
				if err != nil {
					RespondWithError(w, http.StatusUnauthorized, "User not found in context")
					return
				}

				// Mock Gin context with required user info
				c := &gin.Context{
					Request: r,
				}
				c.Set("user_id", user.ID)

				// Call the Gin handler
				h.readReceiptHandler.MarkAllAsRead(c)
			})
		})

		// Websocket endpoint for real-time messaging
		r.Get("/ws", h.HandleWebSocket)
	})
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "messaging",
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// HandleWebSocket handles websocket connections
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if h.websocketHandler != nil {
		h.websocketHandler(w, r)
	} else {
		RespondWithError(w, http.StatusNotImplemented, "WebSocket support not configured")
	}
}

// Pagination represents pagination parameters
type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
		TotalItems  int `json:"total_items"`
		TotalPages  int `json:"total_pages"`
	} `json:"pagination"`
}

// Helper functions

// getPaginationFromRequest parses pagination parameters from the request
func getPaginationFromRequest(r *http.Request) Pagination {
	page := 1
	perPage := 20

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse per_page parameter
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if parsedPerPage, err := strconv.Atoi(perPageStr); err == nil && parsedPerPage > 0 {
			if parsedPerPage > 100 {
				perPage = 100 // Cap at 100
			} else {
				perPage = parsedPerPage
			}
		}
	}

	// Calculate offset
	offset := (page - 1) * perPage

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

// createPaginatedResponse creates a paginated response
func createPaginatedResponse(data interface{}, pagination Pagination, totalCount int) PaginatedResponse {
	totalPages := calculateTotalPages(totalCount, pagination.PerPage)

	response := PaginatedResponse{
		Data: data,
	}

	response.Pagination.CurrentPage = pagination.Page
	response.Pagination.PerPage = pagination.PerPage
	response.Pagination.TotalItems = totalCount
	response.Pagination.TotalPages = totalPages

	return response
}

// calculateTotalPages calculates the total number of pages
func calculateTotalPages(totalItems, perPage int) int {
	if perPage <= 0 {
		return 0
	}

	totalPages := totalItems / perPage
	if totalItems%perPage > 0 {
		totalPages++
	}

	return totalPages
}

// Context helper functions
func getUserFromContext(ctx context.Context) (*User, error) {
	return GetUserFromContext(ctx)
}

func getChannelFromContext(ctx context.Context) (*domain.Channel, error) {
	return GetChannelFromContext(ctx)
}
