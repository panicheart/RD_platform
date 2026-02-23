package handlers

import (
	"net/http"
	"strconv"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
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

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Name                string  `json:"name" binding:"required,max=200"`
	Description         *string `json:"description"`
	Category            string  `json:"category" binding:"required"`
	ProductLine         *string `json:"product_line"`
	Team                *string `json:"team"`
	ProcessTemplateID   *string `json:"process_template_id"`
	StartDate           *string `json:"start_date"`
	EndDate             *string `json:"end_date"`
	LeaderID            *string `json:"leader_id"`
	TechLeaderID        *string `json:"tech_leader_id"`
	ProductLeaderID     *string `json:"product_leader_id"`
	ClassificationLevel string  `json:"classification_level"`
}

// CreateProject handles POST /api/v1/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	// Build project model
	project := models.Project{
		Name:                req.Name,
		Description:         req.Description,
		Category:            req.Category,
		Status:              "draft",
		ClassificationLevel: req.ClassificationLevel,
	}

	if req.ProductLine != nil {
		project.ProductLine = req.ProductLine
	}
	if req.Team != nil {
		project.Team = req.Team
	}
	if req.ProcessTemplateID != nil {
		project.ProcessTemplateID = req.ProcessTemplateID
	}
	if req.LeaderID != nil {
		project.LeaderID = req.LeaderID
	}
	if req.TechLeaderID != nil {
		project.TechLeaderID = req.TechLeaderID
	}
	if req.ProductLeaderID != nil {
		project.ProductLeaderID = req.ProductLeaderID
	}

	if project.ClassificationLevel == "" {
		project.ClassificationLevel = "internal"
	}

	// Create project
	if err := h.projectService.CreateProject(c.Request.Context(), &project, userIDStr); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 6101, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "project created successfully",
		"data":    project,
	})
}

// GetProjects handles GET /api/v1/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
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
	if memberID := c.Query("member_id"); memberID != "" {
		filters["member_id"] = memberID
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	projects, total, err := h.projectService.ListProjects(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
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
		NotFoundResponse(c, err.Error())
		return
	}

	SuccessResponse(c, project)
}

// UpdateProject handles PUT /api/v1/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "code")
	delete(updates, "created_by")

	// Get current user ID for permission check
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), id, updates, userIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to update project" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6102, err.Error())
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

	// Get current user ID for permission check
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	err := h.projectService.DeleteProject(c.Request.Context(), id, userIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to delete project" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6103, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// AddMemberRequest represents the request body for adding a member
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}

// AddMember handles POST /api/v1/projects/:id/members
func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID := c.Param("id")

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	if req.Role == "" {
		req.Role = "member"
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	member, err := h.projectService.AddProjectMember(c.Request.Context(), projectID, req.UserID, req.Role, userIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to add members" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6201, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "member added successfully",
		"data":    member,
	})
}

// GetMembers handles GET /api/v1/projects/:id/members
func (h *ProjectHandler) GetMembers(c *gin.Context) {
	projectID := c.Param("id")

	members, err := h.projectService.GetProjectMembers(c.Request.Context(), projectID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 6202, err.Error())
		return
	}

	SuccessResponse(c, members)
}

// RemoveMember handles DELETE /api/v1/projects/:id/members/:userId
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")

	// Get current user ID
	currentUserID, _ := c.Get("user_id")
	currentUserIDStr := ""
	if currentUserID != nil {
		currentUserIDStr = currentUserID.(string)
	}

	err := h.projectService.RemoveMember(c.Request.Context(), projectID, userID, currentUserIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to remove members" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6203, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// UpdateMemberRole handles PUT /api/v1/projects/:id/members/:userId/role
func (h *ProjectHandler) UpdateMemberRole(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	currentUserID, _ := c.Get("user_id")
	currentUserIDStr := ""
	if currentUserID != nil {
		currentUserIDStr = currentUserID.(string)
	}

	err := h.projectService.UpdateMemberRole(c.Request.Context(), projectID, userID, req.Role, currentUserIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to update member roles" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6204, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// UpdateProgressRequest represents the request body for updating progress
type UpdateProgressRequest struct {
	Progress int `json:"progress" binding:"min=0,max=100"`
}

// UpdateProgress handles PUT /api/v1/projects/:id/progress
func (h *ProjectHandler) UpdateProgress(c *gin.Context) {
	projectID := c.Param("id")

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	project, err := h.projectService.UpdateProjectProgress(c.Request.Context(), projectID, req.Progress, userIDStr)
	if err != nil {
		if err.Error() == "insufficient permissions to update progress" {
			ForbiddenResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusBadRequest, 6301, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "progress updated successfully",
		"data":    project,
	})
}

// GetProjectActivities handles GET /api/v1/projects/:id/activities
func (h *ProjectHandler) GetProjectActivities(c *gin.Context) {
	projectID := c.Param("id")

	activities, err := h.projectService.GetProjectActivities(c.Request.Context(), projectID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 6401, err.Error())
		return
	}

	SuccessResponse(c, activities)
}

// GetUserProjects handles GET /api/v1/users/me/projects
func (h *ProjectHandler) GetUserProjects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		UnauthorizedResponse(c, "unauthorized")
		return
	}

	projects, err := h.projectService.GetUserProjects(c.Request.Context(), userID.(string))
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, projects)
}

// GetProjectStats handles GET /api/v1/projects/stats
func (h *ProjectHandler) GetProjectStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	stats, err := h.projectService.GetProjectStats(c.Request.Context(), userIDStr)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetProjectGantt handles GET /api/v1/projects/:id/gantt
// Returns project activities in Gantt chart format
func (h *ProjectHandler) GetProjectGantt(c *gin.Context) {
	projectID := c.Param("id")

	// Get project details
	project, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		NotFoundResponse(c, err.Error())
		return
	}

	// Get activities
	activities, err := h.projectService.GetProjectActivities(c.Request.Context(), projectID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 6501, err.Error())
		return
	}

	// Format for Gantt chart
	type GanttTask struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		StartDate  *string `json:"start_date,omitempty"`
		EndDate    *string `json:"end_date,omitempty"`
		Progress   int     `json:"progress"`
		Status     string  `json:"status"`
		AssigneeID *string `json:"assignee_id,omitempty"`
		Sequence   int     `json:"sequence"`
	}

	tasks := make([]GanttTask, 0, len(activities))
	for _, activity := range activities {
		task := GanttTask{
			ID:       activity.ID,
			Name:     activity.Name,
			Progress: activity.Progress,
			Status:   string(activity.Status),
			Sequence: activity.Sequence,
		}

		if activity.PlannedStart != nil {
			startStr := activity.PlannedStart.Format("2006-01-02")
			task.StartDate = &startStr
		}
		if activity.PlannedEnd != nil {
			endStr := activity.PlannedEnd.Format("2006-01-02")
			task.EndDate = &endStr
		}
		if activity.AssigneeID != nil {
			task.AssigneeID = activity.AssigneeID
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"project": gin.H{
				"id":       project.ID,
				"name":     project.Name,
				"code":     project.Code,
				"status":   project.Status,
				"progress": project.Progress,
			},
			"tasks": tasks,
		},
	})
}
