package http

import (
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service projectpath.Service
}

func NewHandler(service projectpath.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetProjectPaths)
	router.GET("/:id", h.GetProjectPath)
	router.POST("", h.CreateProjectPath)
	router.PUT("/:id", h.UpdateProjectPath)
	router.DELETE("/:id", h.DeleteProjectPath)
}

func (h *Handler) GetProjectPaths(c *gin.Context) {
	var req projectpath.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := &projectpath.Filter{
		Name:     req.Name,
		IsActive: req.IsActive,
		Search:   req.Search,
	}

	pagination := &projectpath.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	paths, total, err := h.service.ListProjectPaths(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hasNext := pagination.Page*pagination.PageSize < total

	response := &projectpath.ListResponse{
		Data:        paths,
		TotalCount:  total,
		CurrentPage: pagination.Page,
		PageSize:    pagination.PageSize,
		HasNext:     hasNext,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetProjectPath(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	path, err := h.service.GetProjectPath(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, path)
}

func (h *Handler) CreateProjectPath(c *gin.Context) {
	var req projectpath.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := &projectpath.ProjectPath{
		Name:            req.Name,
		Path:            req.Path,
		Description:     req.Description,
		ExcludePatterns: req.ExcludePatterns,
		IsActive:        req.IsActive,
	}

	if err := h.service.CreateProjectPath(c.Request.Context(), path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, path)
}

func (h *Handler) UpdateProjectPath(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req projectpath.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := &projectpath.ProjectPath{
		ID:              id,
		Name:            req.Name,
		Path:            req.Path,
		Description:     req.Description,
		ExcludePatterns: req.ExcludePatterns,
		IsActive:        req.IsActive,
	}

	if err := h.service.UpdateProjectPath(c.Request.Context(), path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, path)
}

func (h *Handler) DeleteProjectPath(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteProjectPath(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}