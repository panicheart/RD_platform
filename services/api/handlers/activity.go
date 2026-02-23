package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"
)

// ActivityHandler handles activity-related HTTP requests
type ActivityHandler struct {
	activityService *services.ActivityService
}

// NewActivityHandler creates a new activity handler
func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

// CreateActivityRequest represents the request body for creating an activity
type CreateActivityRequest struct {
	WorkflowID  string              `json:"workflow_id" binding:"required"`
	Name        string              `json:"name" binding:"required,max=200"`
	Description string              `json:"description"`
	Type        models.ActivityType `json:"type" binding:"required"`
	Sequence    int                 `json:"sequence"`
	Priority    int                 `json:"priority"`
	AssigneeID  string              `json:"assignee_id"`
}

// CreateActivity creates a new activity
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	projectID := c.Param("projectId")
	userID, _ := c.Get("userID")

	activity, err := h.activityService.CreateActivity(
		req.WorkflowID,
		projectID,
		req.Name,
		req.Description,
		req.Type,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Activity created successfully", "data": activity})
}

// GetActivity retrieves an activity by ID
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	activityID := c.Param("id")

	activity, err := h.activityService.GetActivity(activityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Activity not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Success", "data": activity})
}

// ListActivities retrieves activities for a workflow
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	workflowID := c.Query("workflow_id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "workflow_id is required", "data": nil})
		return
	}

	status := models.ActivityStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	activities, total, err := h.activityService.ListActivities(workflowID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"items":     activities,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// StartActivityRequest represents the request body for starting an activity
type StartActivityRequest struct {
	AssigneeID string `json:"assignee_id"`
}

// StartActivity starts an activity
func (h *ActivityHandler) StartActivity(c *gin.Context) {
	activityID := c.Param("id")
	userID, _ := c.Get("userID")

	var req StartActivityRequest
	c.ShouldBindJSON(&req)

	if req.AssigneeID != "" {
		h.activityService.AssignActivity(activityID, req.AssigneeID)
	}

	if err := h.activityService.StartActivity(activityID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Activity started successfully", "data": nil})
}

// CompleteActivity completes an activity
func (h *ActivityHandler) CompleteActivity(c *gin.Context) {
	activityID := c.Param("id")
	userID, _ := c.Get("userID")

	if err := h.activityService.CompleteActivity(activityID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Activity completed successfully", "data": nil})
}

// UpdateActivityProgressRequest represents the request body for updating progress
type UpdateActivityProgressRequest struct {
	Progress int `json:"progress" binding:"required,min=0,max=100"`
}

// UpdateActivityProgress updates activity progress
func (h *ActivityHandler) UpdateActivityProgress(c *gin.Context) {
	activityID := c.Param("id")

	var req UpdateActivityProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	if err := h.activityService.UpdateActivityProgress(activityID, req.Progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Progress updated successfully", "data": nil})
}

// AssignActivity assigns an activity to a user
func (h *ActivityHandler) AssignActivity(c *gin.Context) {
	activityID := c.Param("id")

	var req struct {
		AssigneeID string `json:"assignee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	if err := h.activityService.AssignActivity(activityID, req.AssigneeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Activity assigned successfully", "data": nil})
}

// AddDependencyRequest represents the request body for adding a dependency
type AddDependencyRequest struct {
	DependsOnID    string `json:"depends_on_id" binding:"required"`
	DependencyType string `json:"dependency_type"`
}

// AddDependency adds a dependency to an activity
func (h *ActivityHandler) AddDependency(c *gin.Context) {
	activityID := c.Param("id")

	var req AddDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	depType := req.DependencyType
	if depType == "" {
		depType = "finish_to_start"
	}

	if err := h.activityService.AddDependency(activityID, req.DependsOnID, depType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Dependency added successfully", "data": nil})
}
