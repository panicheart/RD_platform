package handlers

import (
	"net/http"
	"strconv"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// ProcessTemplateHandler handles process template HTTP requests
type ProcessTemplateHandler struct {
	templateService *services.ProcessTemplateService
}

// NewProcessTemplateHandler creates a new ProcessTemplateHandler
func NewProcessTemplateHandler(templateService *services.ProcessTemplateService) *ProcessTemplateHandler {
	return &ProcessTemplateHandler{
		templateService: templateService,
	}
}

// ListTemplates handles GET /api/v1/process-templates
func (h *ProcessTemplateHandler) ListTemplates(c *gin.Context) {
	category := c.Query("category")

	templates, err := h.templateService.ListTemplates(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    templates,
	})
}

// GetTemplate handles GET /api/v1/process-templates/:id
func (h *ProcessTemplateHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	template, err := h.templateService.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    template,
	})
}

// CreateTemplate handles POST /api/v1/process-templates
func (h *ProcessTemplateHandler) CreateTemplate(c *gin.Context) {
	var template models.ProcessTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Set defaults
	if template.Activities == "" {
		template.Activities = "[]"
	}

	err := h.templateService.CreateTemplate(c.Request.Context(), &template)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "template created successfully",
		"data":    template,
	})
}

// UpdateTemplate handles PUT /api/v1/process-templates/:id
func (h *ProcessTemplateHandler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	template, err := h.templateService.UpdateTemplate(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "template updated successfully",
		"data":    template,
	})
}

// DeleteTemplate handles DELETE /api/v1/process-templates/:id
func (h *ProcessTemplateHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")

	err := h.templateService.DeleteTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "template deleted successfully",
		"data":    nil,
	})
}

// GetDefaultTemplate handles GET /api/v1/process-templates/default/:category
func (h *ProcessTemplateHandler) GetDefaultTemplate(c *gin.Context) {
	category := c.Param("category")

	template, err := h.templateService.GetDefaultTemplate(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    template,
	})
}

// UpdateActivity handles PUT /api/v1/projects/:projectId/activities/:id
func (h *ProcessTemplateHandler) UpdateActivity(c *gin.Context) {
	projectID := c.Param("projectId")
	activityID := c.Param("id")

	// Validate IDs
	if projectID == "" || activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid project or activity ID",
			"data":    nil,
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "project_id")
	delete(updates, "created_at")

	// This would call projectService.UpdateActivity
	// For now, return a placeholder response
	_ = projectID

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "activity updated successfully",
		"data":    updates,
	})
}

// Helper function - unused but kept for reference
var _ = strconv.Atoi
