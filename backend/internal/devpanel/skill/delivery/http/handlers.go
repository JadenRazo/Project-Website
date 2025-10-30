package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/skill"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SkillService defines the interface for skill operations
type SkillService interface {
	Create(ctx context.Context, skill *skill.Skill) error
	Get(ctx context.Context, id uuid.UUID) (*skill.Skill, error)
	GetByCategory(ctx context.Context, category skill.Category) ([]*skill.Skill, error)
	GetFeatured(ctx context.Context) ([]*skill.Skill, error)
	Update(ctx context.Context, skill *skill.Skill) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *skill.Filter, pagination *skill.Pagination) ([]*skill.Skill, int, error)
	UpdateSortOrder(ctx context.Context, skillID uuid.UUID, sortOrder int) error
	ToggleFeatured(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service SkillService
}

func NewHandler(service SkillService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetSkills)
	router.GET("/featured", h.GetFeaturedSkills)
	router.GET("/categories", h.GetCategories)
	router.GET("/proficiency-levels", h.GetProficiencyLevels)
	router.GET("/category/:category", h.GetSkillsByCategory)
	router.GET("/:id", h.GetSkill)
	router.POST("", h.CreateSkill)
	router.PUT("/:id", h.UpdateSkill)
	router.DELETE("/:id", h.DeleteSkill)
	router.PUT("/:id/sort-order", h.UpdateSortOrder)
	router.PUT("/:id/toggle-featured", h.ToggleFeatured)
}

func (h *Handler) GetSkills(c *gin.Context) {
	// Parse query parameters for filtering and pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	category := c.Query("category")
	proficiencyLevel := c.Query("proficiency_level")
	isFeaturedStr := c.Query("is_featured")
	tag := c.Query("tag")
	search := c.Query("search")

	filter := &skill.Filter{
		Tag:    tag,
		Search: search,
	}

	if category != "" {
		cat := skill.Category(category)
		filter.Category = &cat
	}

	if proficiencyLevel != "" {
		pl := skill.ProficiencyLevel(proficiencyLevel)
		filter.ProficiencyLevel = &pl
	}

	if isFeaturedStr != "" {
		isFeatured, _ := strconv.ParseBool(isFeaturedStr)
		filter.IsFeatured = &isFeatured
	}

	pagination := &skill.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	skills, total, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get skills",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills":   skills,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetFeaturedSkills(c *gin.Context) {
	skills, err := h.service.GetFeatured(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get featured skills",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": skills,
	})
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories := skill.GetAllCategories()
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func (h *Handler) GetProficiencyLevels(c *gin.Context) {
	levels := skill.GetAllProficiencyLevels()
	c.JSON(http.StatusOK, gin.H{
		"proficiency_levels": levels,
	})
}

func (h *Handler) GetSkillsByCategory(c *gin.Context) {
	categoryStr := c.Param("category")
	category := skill.Category(categoryStr)

	skills, err := h.service.GetByCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get skills by category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": skills,
	})
}

func (h *Handler) GetSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid skill ID",
		})
		return
	}

	skillResult, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if err == skill.ErrSkillNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Skill not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get skill",
		})
		return
	}

	c.JSON(http.StatusOK, skillResult)
}

func (h *Handler) CreateSkill(c *gin.Context) {
	var req skill.Skill
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		if err == skill.ErrInvalidSkill {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create skill",
		})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *Handler) UpdateSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid skill ID",
		})
		return
	}

	var req skill.Skill
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	req.ID = id
	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		if err == skill.ErrSkillNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Skill not found",
			})
			return
		}
		if err == skill.ErrInvalidSkill {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update skill",
		})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid skill ID",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == skill.ErrSkillNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Skill not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete skill",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Skill deleted successfully",
	})
}

func (h *Handler) UpdateSortOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid skill ID",
		})
		return
	}

	var req struct {
		SortOrder int `json:"sort_order" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateSortOrder(c.Request.Context(), id, req.SortOrder); err != nil {
		if err == skill.ErrSkillNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Skill not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update sort order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sort order updated successfully",
	})
}

func (h *Handler) ToggleFeatured(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid skill ID",
		})
		return
	}

	if err := h.service.ToggleFeatured(c.Request.Context(), id); err != nil {
		if err == skill.ErrSkillNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Skill not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to toggle featured status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Featured status toggled successfully",
	})
}
