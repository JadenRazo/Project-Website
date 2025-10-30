package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service        *codestats.Service
	lastUpdate     time.Time
	updateInterval time.Duration
}

func NewHandler(service *codestats.Service) *Handler {
	return &Handler{
		service:        service,
		updateInterval: 5 * time.Minute, // Cache stats for 5 minutes
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/stats", h.GetStats)
	router.POST("/stats/update", h.UpdateStats)
}

func (h *Handler) GetStats(c *gin.Context) {
	requestID := c.GetString("RequestID")
	forceUpdate := c.Query("force") == "true"

	logger.Info("Getting code statistics",
		"request_id", requestID,
		"force_update", forceUpdate)

	// Check if we need to update stats (cache expired or forced)
	now := time.Now()
	shouldUpdate := forceUpdate || now.Sub(h.lastUpdate) > h.updateInterval

	if shouldUpdate {
		logger.Info("Updating code statistics",
			"request_id", requestID,
			"reason", map[bool]string{true: "forced", false: "cache_expired"}[forceUpdate])

		if err := h.service.UpdateStats(); err != nil {
			logger.Error("Failed to update code statistics",
				"request_id", requestID,
				"error", err)

			// Check if this is a configuration error (missing project paths)
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "no active project paths configured") || 
			   strings.Contains(errorMsg, "Please add project paths in the DevPanel") {
				// Configuration error - return error to user immediately
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      errorMsg,
					"request_id": requestID,
				})
				return
			}

			// Try to return cached stats if update fails (for temporary errors)
			stats, getErr := h.service.GetLatestStats()
			if getErr != nil || stats == nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":      "Failed to update code statistics. Please try again later.",
					"request_id": requestID,
				})
				return
			}

			// Return cached stats with warning (for temporary errors only)
			logger.Warn("Returning cached statistics due to update failure",
				"request_id", requestID,
				"stats_age", time.Since(stats.UpdatedAt))

			c.JSON(http.StatusOK, gin.H{
				"totalLines":  stats.TotalLines,
				"cached":      true,
				"lastUpdated": stats.UpdatedAt,
			})
			return
		}

		h.lastUpdate = now
	}

	stats, err := h.service.GetLatestStats()
	if err != nil {
		logger.Error("Failed to get code statistics",
			"request_id", requestID,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to retrieve code statistics. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	if stats == nil {
		logger.Warn("No code statistics available",
			"request_id", requestID)

		c.JSON(http.StatusNotFound, gin.H{
			"error":      "No code statistics available. Statistics will be generated shortly.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Code statistics retrieved successfully",
		"request_id", requestID,
		"total_lines", stats.TotalLines,
		"stats_age", time.Since(stats.UpdatedAt))

	// Return in the format expected by frontend - only total lines
	c.JSON(http.StatusOK, gin.H{
		"totalLines":  stats.TotalLines,
		"lastUpdated": stats.UpdatedAt,
	})
}

func (h *Handler) UpdateStats(c *gin.Context) {
	requestID := c.GetString("RequestID")

	logger.Info("Manually updating code statistics",
		"request_id", requestID)

	startTime := time.Now()
	if err := h.service.UpdateStats(); err != nil {
		logger.Error("Failed to update code statistics",
			"request_id", requestID,
			"error", err,
			"duration", time.Since(startTime))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to update code statistics. Please check logs for details.",
			"request_id": requestID,
		})
		return
	}

	// Update cache timestamp
	h.lastUpdate = time.Now()

	// Get the updated stats to return
	stats, err := h.service.GetLatestStats()
	if err != nil {
		logger.Warn("Updated stats but failed to retrieve them",
			"request_id", requestID,
			"error", err)
	}

	var totalLines int64
	if stats != nil {
		totalLines = stats.TotalLines
	}

	logger.Info("Code statistics updated successfully",
		"request_id", requestID,
		"duration", time.Since(startTime),
		"total_lines", totalLines)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Code statistics updated successfully",
		"totalLines": totalLines,
		"request_id": requestID,
	})
}
