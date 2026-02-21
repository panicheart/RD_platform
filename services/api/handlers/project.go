package handlers

import (
	"net/http"
	"strconv"

	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectHandler handles project HTTP requests
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// ListProjects handles GET /api/v1/projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if team := c.Query("team"); team != "" {
		filters["team"] = team
	}
	if productLine := c.Query("product_line"); productLine != "" {
		filters["product_line"] = productLine
	}
	if leaderID := c.Query("leader_id"); leaderID != "" {
		filters["leader_id"] = leaderID
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	projects, total, err := h.projectService.ListProjects(c.Request.Context(), page, pageSize, filters)
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
		"data": gin.H{
			"items":     projects,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProject handles GET /api/v1/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetProjectByID(c.Request.Context(), id)
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
		"data":    project,
	})
}

// CreateProject handles POST /api/v1/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Set defaults
	if project.Status == "" {
		project.Status = "draft"
	}
	if project.Progress < 0 {
		project.Progress = 0
	}
	if project.Progress > 100 {
		project.Progress = 100
	}

	// Get current user as creator
	if userID, exists := c.Get("user_id"); exists {
		userIDStr := userID.(string)
		project.CreatedBy = &userIDStr
	}

	err := h.projectService.CreateProject(c.Request.Context(), &project)
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
		"message": "project created successfully",
		"data":    project,
	})
}

// UpdateProject handles PUT /api/v1/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
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

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "code")
	delete(updates, "created_by")

	project, err := h.projectService.UpdateProject(c.Request.Context(), id, updates)
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
		"message": "project updated successfully",
		"data":    project,
	})
}

// DeleteProject handles DELETE /api/v1/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	err := h.projectService.DeleteProject(c.Request.Context(), id)
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
		"message": "project deleted successfully",
		"data":    nil,
	})
}

// AddMember handles POST /api/v1/projects/:id/members
func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	if req.Role == "" {
		req.Role = "member"
	}

	member, err := h.projectService.AddMember(c.Request.Context(), projectID, req.UserID, req.Role)
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
		"message": "member added successfully",
		"data":    member,
	})
}

// RemoveMember handles DELETE /api/v1/projects/:id/members/:userId
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")

	err := h.projectService.RemoveMember(c.Request.Context(), projectID, userID)
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
		"message": "member removed successfully",
		"data":    nil,
	})
}

// GetUserProjects handles GET /api/v1/users/me/projects
func (h *ProjectHandler) GetUserProjects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	projects, err := h.projectService.GetUserProjects(c.Request.Context(), userID.(string))
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
		"data":    projects,
	})
}

// GetProjectActivities handles GET /api/v1/projects/:id/activities
func (h *ProjectHandler) GetProjectActivities(c *gin.Context) {
	projectID := c.Param("id")

	activities, err := h.projectService.GetProjectActivities(c.Request.Context(), projectID)
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
		"data":    activities,
	})
}

// CreateActivity handles POST /api/v1/projects/:id/activities
func (h *ProjectHandler) CreateActivity(c *gin.Context) {
	projectID := c.Param("id")

	// Validate project ID
	if _, err := uuid.Parse(projectID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid project ID",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	projectUID, _ := uuid.Parse(projectID)
	activity.ProjectID = projectUID

	if activity.Status == "" {
		activity.Status = "pending"
	}

	err := h.projectService.CreateActivity(c.Request.Context(), &activity)
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
		"message": "activity created successfully",
		"data":    activity,
	})
}
