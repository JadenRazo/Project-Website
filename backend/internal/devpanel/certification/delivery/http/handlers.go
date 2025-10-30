package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/certification"
)

type Handler struct {
	service certification.Service
}

func NewHandler(service certification.Service) *Handler {
	return &Handler{service: service}
}

// Certification handlers

func (h *Handler) CreateCertification(c *gin.Context) {
	var req struct {
		Name             string `json:"name"`
		Issuer           string `json:"issuer"`
		CredentialID     string `json:"credential_id"`
		IssueDate        string `json:"issue_date"`
		ExpiryDate       string `json:"expiry_date"`
		VerificationURL  string `json:"verification_url"`
		VerificationText string `json:"verification_text"`
		BadgeURL         string `json:"badge_url"`
		Description      string `json:"description"`
		CategoryID       string `json:"category_id"`
		IsFeatured       bool   `json:"is_featured"`
		IsVisible        bool   `json:"is_visible"`
		SortOrder        int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cert := &certification.Certification{
		Name:       req.Name,
		Issuer:     req.Issuer,
		IsFeatured: req.IsFeatured,
		IsVisible:  req.IsVisible,
		SortOrder:  req.SortOrder,
	}

	// Parse dates
	if req.IssueDate != "" {
		issueDate, err := time.Parse("2006-01-02", req.IssueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue date format"})
			return
		}
		cert.IssueDate = issueDate
	}

	if req.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
			return
		}
		cert.ExpiryDate = &expiryDate
	}

	// Set optional fields
	if req.CredentialID != "" {
		cert.CredentialID = &req.CredentialID
	}

	if req.VerificationURL != "" {
		cert.VerificationURL = &req.VerificationURL
	}

	if req.VerificationText != "" {
		cert.VerificationText = &req.VerificationText
	}

	if req.BadgeURL != "" {
		cert.BadgeURL = &req.BadgeURL
	}

	if req.Description != "" {
		cert.Description = &req.Description
	}

	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		cert.CategoryID = &catID
	}

	if err := h.service.CreateCertification(cert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cert)
}

func (h *Handler) GetAllCertifications(c *gin.Context) {
	includeHidden := c.Query("include_hidden") == "true"

	certs, err := h.service.GetAllCertifications(includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certs)
}

func (h *Handler) GetVisibleCertifications(c *gin.Context) {
	certs, err := h.service.GetVisibleCertifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certs)
}

func (h *Handler) GetCertificationByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certification ID"})
		return
	}

	cert, err := h.service.GetCertificationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certification not found"})
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *Handler) UpdateCertification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certification ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse dates if present
	if issueDateStr, ok := updates["issue_date"].(string); ok && issueDateStr != "" {
		issueDate, err := time.Parse("2006-01-02", issueDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue date format"})
			return
		}
		updates["issue_date"] = issueDate
	}

	if expiryDateStr, ok := updates["expiry_date"].(string); ok && expiryDateStr != "" {
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
			return
		}
		updates["expiry_date"] = expiryDate
	}

	if err := h.service.UpdateCertification(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteCertification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certification ID"})
		return
	}

	if err := h.service.DeleteCertification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Category handlers

func (h *Handler) CreateCategory(c *gin.Context) {
	var cat certification.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.CreateCategory(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	includeHidden := c.Query("include_hidden") == "true"

	cats, err := h.service.GetAllCategories(includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetVisibleCategories(c *gin.Context) {
	cats, err := h.service.GetVisibleCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateCategory(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Certification routes
	router.POST("/certifications", h.CreateCertification)
	router.GET("/certifications", h.GetAllCertifications)
	router.GET("/certifications/:id", h.GetCertificationByID)
	router.PUT("/certifications/:id", h.UpdateCertification)
	router.DELETE("/certifications/:id", h.DeleteCertification)

	// Category routes
	router.POST("/certification-categories", h.CreateCategory)
	router.GET("/certification-categories", h.GetAllCategories)
	router.GET("/certification-categories/:id", h.GetCategoryByID)
	router.PUT("/certification-categories/:id", h.UpdateCategory)
	router.DELETE("/certification-categories/:id", h.DeleteCategory)
}
