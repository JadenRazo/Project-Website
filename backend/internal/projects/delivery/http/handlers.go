package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectService defines the interface for project operations
type ProjectService interface {
	Create(ctx context.Context, project *project.Project) error
	Get(ctx context.Context, id uuid.UUID) (*project.Project, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*project.Project, error)
	Update(ctx context.Context, project *project.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *project.Filter, pagination *project.Pagination) ([]*project.Project, int, error)
}

type Handler struct {
	service ProjectService
}

func NewHandler(service ProjectService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetProjects)
	router.GET("/featured", h.GetFeaturedProjects)
	router.GET("/:id", h.GetProject)
	router.POST("", h.CreateProject)
	router.PUT("/:id", h.UpdateProject)
	router.DELETE("/:id", h.DeleteProject)
}

func (h *Handler) GetProjects(c *gin.Context) {
	requestID := c.GetString("RequestID")

	// Parse query parameters for filtering and pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	status := c.Query("status")
	tag := c.Query("tag")
	search := c.Query("search")

	filter := &project.Filter{
		Tag:    tag,
		Search: search,
	}

	if status != "" {
		s := project.Status(status)
		filter.Status = &s
	}

	pagination := &project.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	logger.Info("Getting projects list",
		"request_id", requestID,
		"page", page,
		"page_size", pageSize,
		"status", status,
		"tag", tag,
		"search", search)

	projects, total, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		logger.Error("Failed to get projects",
			"request_id", requestID,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to retrieve projects. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	// Ensure we return empty array instead of null
	if projects == nil {
		projects = []*project.Project{}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetFeaturedProjects(c *gin.Context) {
	requestID := c.GetString("RequestID")

	// Get active projects with a limit for featured display
	filter := &project.Filter{}
	activeStatus := project.StatusActive
	filter.Status = &activeStatus

	pagination := &project.Pagination{
		Page:     1,
		PageSize: 6, // Limit to 6 featured projects
	}

	logger.Info("Getting featured projects", "request_id", requestID)

	projects, _, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		logger.Error("Failed to get featured projects",
			"request_id", requestID,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to retrieve featured projects. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	// Ensure we return empty array instead of null
	if projects == nil {
		projects = []*project.Project{}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func (h *Handler) GetProject(c *gin.Context) {
	requestID := c.GetString("RequestID")
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("Invalid project ID format",
			"request_id", requestID,
			"id", idStr,
			"error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid project ID format. Please provide a valid UUID.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Getting project details",
		"request_id", requestID,
		"project_id", id)

	projectResult, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			logger.Warn("Project not found",
				"request_id", requestID,
				"project_id", id)

			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Project not found.",
				"request_id": requestID,
			})
			return
		}

		logger.Error("Failed to get project",
			"request_id", requestID,
			"project_id", id,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to retrieve project details. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	c.JSON(http.StatusOK, projectResult)
}

func (h *Handler) CreateProject(c *gin.Context) {
	requestID := c.GetString("RequestID")

	var req project.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body for project creation",
			"request_id", requestID,
			"error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid request body. Please check your input and try again.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Creating new project",
		"request_id", requestID,
		"project_name", req.Name)

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		logger.Error("Failed to create project",
			"request_id", requestID,
			"project_name", req.Name,
			"error", err)

		// Check for specific validation errors
		if errors.Is(err, project.ErrInvalidProject) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      err.Error(),
				"request_id": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to create project. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Project created successfully",
		"request_id", requestID,
		"project_id", req.ID,
		"project_name", req.Name)

	c.JSON(http.StatusCreated, req)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	requestID := c.GetString("RequestID")
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("Invalid project ID format for update",
			"request_id", requestID,
			"id", idStr,
			"error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid project ID format. Please provide a valid UUID.",
			"request_id": requestID,
		})
		return
	}

	var req project.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body for project update",
			"request_id", requestID,
			"project_id", id,
			"error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid request body. Please check your input and try again.",
			"request_id": requestID,
		})
		return
	}

	req.ID = id

	logger.Info("Updating project",
		"request_id", requestID,
		"project_id", id,
		"project_name", req.Name)

	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			logger.Warn("Project not found for update",
				"request_id", requestID,
				"project_id", id)

			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Project not found.",
				"request_id": requestID,
			})
			return
		}

		if errors.Is(err, project.ErrInvalidProject) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      err.Error(),
				"request_id": requestID,
			})
			return
		}

		logger.Error("Failed to update project",
			"request_id", requestID,
			"project_id", id,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to update project. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Project updated successfully",
		"request_id", requestID,
		"project_id", id)

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	requestID := c.GetString("RequestID")
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("Invalid project ID format for deletion",
			"request_id", requestID,
			"id", idStr,
			"error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid project ID format. Please provide a valid UUID.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Deleting project",
		"request_id", requestID,
		"project_id", id)

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			logger.Warn("Project not found for deletion",
				"request_id", requestID,
				"project_id", id)

			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Project not found.",
				"request_id": requestID,
			})
			return
		}

		logger.Error("Failed to delete project",
			"request_id", requestID,
			"project_id", id,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to delete project. Please try again later.",
			"request_id": requestID,
		})
		return
	}

	logger.Info("Project deleted successfully",
		"request_id", requestID,
		"project_id", id)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Project deleted successfully",
		"request_id": requestID,
	})
}
