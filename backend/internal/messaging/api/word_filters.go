package api

import (
	"net/http"
	"strings"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WordFilterHandler struct {
	wordFilterUseCase domain.WordFilterUseCase
}

func NewWordFilterHandler(wordFilterUseCase domain.WordFilterUseCase) *WordFilterHandler {
	return &WordFilterHandler{
		wordFilterUseCase: wordFilterUseCase,
	}
}

// validateWordFilter validates the word filter input
func validateWordFilter(filter *domain.WordFilter) error {
	if strings.TrimSpace(filter.Word) == "" {
		return domain.NewValidationError("word is required")
	}
	if filter.Scope == "" {
		return domain.NewValidationError("scope is required")
	}
	if filter.Scope != "private" && filter.Scope != "public" && filter.Scope != "all" {
		return domain.NewValidationError("invalid scope")
	}
	return nil
}

// CreateWordFilter creates a new word filter
func (h *WordFilterHandler) CreateWordFilter(c *gin.Context) {
	var filter domain.WordFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateWordFilter(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.wordFilterUseCase.CreateFilter(c.Request.Context(), &filter); err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create word filter"})
		return
	}

	c.JSON(http.StatusCreated, filter)
}

// UpdateWordFilter updates an existing word filter
func (h *WordFilterHandler) UpdateWordFilter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var filter domain.WordFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	filter.ID = id
	if err := validateWordFilter(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.wordFilterUseCase.UpdateFilter(c.Request.Context(), &filter); err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "word filter not found"})
			return
		}
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update word filter"})
		return
	}

	c.JSON(http.StatusOK, filter)
}

// DeleteWordFilter deletes a word filter
func (h *WordFilterHandler) DeleteWordFilter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.wordFilterUseCase.DeleteFilter(c.Request.Context(), id); err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "word filter not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete word filter"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetWordFilter retrieves a word filter by ID
func (h *WordFilterHandler) GetWordFilter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	filter, err := h.wordFilterUseCase.GetFilter(c.Request.Context(), id)
	if err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "word filter not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get word filter"})
		return
	}

	c.JSON(http.StatusOK, filter)
}

// ListWordFilters retrieves all word filters for a server
func (h *WordFilterHandler) ListWordFilters(c *gin.Context) {
	serverID, err := uuid.Parse(c.Param("server_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid server id"})
		return
	}

	filters, err := h.wordFilterUseCase.ListFilters(c.Request.Context(), serverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list word filters"})
		return
	}

	c.JSON(http.StatusOK, filters)
}

// GetWordFiltersByScope retrieves word filters by scope
func (h *WordFilterHandler) GetWordFiltersByScope(c *gin.Context) {
	serverID, err := uuid.Parse(c.Param("server_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid server id"})
		return
	}

	scope := c.Param("scope")
	if scope != "private" && scope != "public" && scope != "all" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scope"})
		return
	}

	filters, err := h.wordFilterUseCase.GetFiltersByScope(
		c.Request.Context(),
		serverID,
		domain.WordFilterScope(scope),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get word filters by scope"})
		return
	}

	c.JSON(http.StatusOK, filters)
}
