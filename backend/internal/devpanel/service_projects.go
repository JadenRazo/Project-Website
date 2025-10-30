package devpanel

import (
	"net/http"
	"strconv"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Project Management Handlers

// listProjects returns a paginated list of projects
func (s *Service) listProjects(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Parse filter parameters
	filter := &project.Filter{
		Search: c.Query("search"),
		Tag:    c.Query("tag"),
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := project.Status(statusStr)
		filter.Status = &status
	}

	// Parse owner filter
	if ownerIDStr := c.Query("owner_id"); ownerIDStr != "" {
		if ownerID, err := uuid.Parse(ownerIDStr); err == nil {
			filter.OwnerID = ownerID
		}
	}

	pagination := &project.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	// Get projects from service
	projects, total, err := s.projectService.List(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects":   projects,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}

// createProject creates a new project
func (s *Service) createProject(c *gin.Context) {
	var proj project.Project
	if err := c.ShouldBindJSON(&proj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set owner ID (in real implementation, this would come from auth context)
	// For now, using a default owner
	if proj.OwnerID == uuid.Nil {
		proj.OwnerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if err := s.projectService.Create(c.Request.Context(), &proj); err != nil {
		if err == project.ErrInvalidProject {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, proj)
}

// getProject returns a specific project by ID
func (s *Service) getProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	proj, err := s.projectService.Get(c.Request.Context(), id)
	if err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proj)
}

// updateProject updates an existing project
func (s *Service) updateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	var proj project.Project
	if err := c.ShouldBindJSON(&proj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID matches
	proj.ID = id

	if err := s.projectService.Update(c.Request.Context(), &proj); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		if err == project.ErrInvalidProject {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proj)
}

// deleteProject soft-deletes a project
func (s *Service) deleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	if err := s.projectService.Delete(c.Request.Context(), id); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project deleted successfully"})
}

// archiveProject archives a project
func (s *Service) archiveProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	if err := s.projectService.Archive(c.Request.Context(), id); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project archived successfully"})
}

// activateProject activates a project
func (s *Service) activateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	if err := s.projectService.Activate(c.Request.Context(), id); err != nil {
		if err == project.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project activated successfully"})
}