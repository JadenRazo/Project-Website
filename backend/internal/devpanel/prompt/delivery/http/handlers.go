
package http

import (
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/prompt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service prompt.Service
}

func NewHandler(service prompt.Service) *Handler {
	return &Handler{service: service}
}

// Prompt handlers

func (h *Handler) CreatePrompt(c *gin.Context) {
	var p prompt.Prompt
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePrompt(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prompt"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *Handler) GetAllPrompts(c *gin.Context) {
	includeHidden := c.Query("include_hidden") == "true"
	prompts, err := h.service.GetAllPrompts(includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prompts"})
		return
	}
	c.JSON(http.StatusOK, prompts)
}

func (h *Handler) GetVisiblePrompts(c *gin.Context) {
	prompts, err := h.service.GetVisiblePrompts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prompts"})
		return
	}
	c.JSON(http.StatusOK, prompts)
}

func (h *Handler) GetPromptByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	p, err := h.service.GetPromptByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prompt not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *Handler) UpdatePrompt(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePrompt(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prompt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt updated successfully"})
}

func (h *Handler) DeletePrompt(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	if err := h.service.DeletePrompt(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prompt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt deleted successfully"})
}

// Category handlers

func (h *Handler) CreateCategory(c *gin.Context) {
	var cat prompt.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCategory(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	includeHidden := c.Query("include_hidden") == "true"
	cats, err := h.service.GetAllCategories(includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetVisibleCategories(c *gin.Context) {
	cats, err := h.service.GetVisibleCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	cat, err := h.service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateCategory(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Prompt routes
	router.POST("/prompts", h.CreatePrompt)
	router.GET("/prompts", h.GetAllPrompts)
	router.GET("/prompts/:id", h.GetPromptByID)
	router.PUT("/prompts/:id", h.UpdatePrompt)
	router.DELETE("/prompts/:id", h.DeletePrompt)

	// Category routes
	router.POST("/prompt-categories", h.CreateCategory)
	router.GET("/prompt-categories", h.GetAllCategories)
	router.GET("/prompt-categories/:id", h.GetCategoryByID)
	router.PUT("/prompt-categories/:id", h.UpdateCategory)
	router.DELETE("/prompt-categories/:id", h.DeleteCategory)
}
