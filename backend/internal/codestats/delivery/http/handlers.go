package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
)

type Handler struct {
	service *codestats.Service
}

func NewHandler(service *codestats.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/stats", h.GetStats)
	router.POST("/stats/update", h.UpdateStats)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetLatestStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get code statistics",
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No statistics available",
		})
		return
	}

	// Return in the format expected by frontend - only total lines
	c.JSON(http.StatusOK, gin.H{
		"totalLines": stats.TotalLines,
	})
}

func (h *Handler) UpdateStats(c *gin.Context) {
	if err := h.service.UpdateStats(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update code statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Code statistics updated successfully",
	})
}