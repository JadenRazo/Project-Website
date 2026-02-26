package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JadenRazo/Project-Website/backend/internal/blog"
)

type Handler struct {
	service blog.Service
}

func NewHandler(service blog.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListPublished)
	rg.GET("/featured", h.GetFeatured)
	rg.GET("/:slug", h.GetBySlug)
	rg.POST("/:slug/view", h.IncrementViewCount)
}

func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListAll)
	rg.POST("", h.Create)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

func (h *Handler) ListPublished(c *gin.Context) {
	var params blog.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	result, err := h.service.ListPublished(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetFeatured(c *gin.Context) {
	posts, err := h.service.GetFeatured(3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch featured posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	post, err := h.service.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) IncrementViewCount(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	if err := h.service.IncrementViewCount(slug, c.Request); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListAll(c *gin.Context) {
	var params blog.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	result, err := h.service.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Create(c *gin.Context) {
	var post blog.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.CreatePost(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdatePost(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, post)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	if err := h.service.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
