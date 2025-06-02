package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
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
	// Parse query parameters for filtering and pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
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

	projects, total, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get projects",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetFeaturedProjects(c *gin.Context) {
	// Get active projects with a limit for featured display
	filter := &project.Filter{}
	activeStatus := project.StatusActive
	filter.Status = &activeStatus

	pagination := &project.Pagination{
		Page:     1,
		PageSize: 6, // Limit to 6 featured projects
	}

	projects, _, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get featured projects",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func (h *Handler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	projectResult, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get project",
		})
		return
	}

	c.JSON(http.StatusOK, projectResult)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req project.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create project",
		})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req project.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	req.ID = id
	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update project",
		})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete project",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
	})
}