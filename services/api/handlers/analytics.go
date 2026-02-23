package handlers

import (
	"net/http"
	"time"

	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics HTTP requests
type AnalyticsHandler struct {
	analyticsService   *services.AnalyticsService
	aggregationService *services.AggregationService
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService, aggregationService *services.AggregationService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService:   analyticsService,
		aggregationService: aggregationService,
	}
}

// parseDateRange parses start_date and end_date query parameters
func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		// Default to 30 days ago
		startDate = time.Now().AddDate(0, 0, -30)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if endDateStr == "" {
		// Default to today
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	// Set end date to end of day
	endDate = endDate.Add(24 * time.Hour).Add(-time.Second)

	return startDate, endDate, nil
}

// GetDashboard returns the dashboard overview data
// GET /api/v1/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	overview, err := h.analyticsService.GetDashboardOverview(c.Request.Context())
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, overview)
}

// GetDashboardWidgets returns all dashboard widgets data
// GET /api/v1/analytics/dashboard/widgets
func (h *AnalyticsHandler) GetDashboardWidgets(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	widgets, err := h.aggregationService.GetDashboardWidgets(c.Request.Context(), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, widgets)
}

// GetProjectStats returns project statistics
// GET /api/v1/analytics/projects
func (h *AnalyticsHandler) GetProjectStats(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	stats, err := h.analyticsService.GetProjectStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetUserStats returns user statistics
// GET /api/v1/analytics/users
func (h *AnalyticsHandler) GetUserStats(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	stats, err := h.analyticsService.GetUserStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetShelfStats returns shelf statistics
// GET /api/v1/analytics/shelf
func (h *AnalyticsHandler) GetShelfStats(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	stats, err := h.analyticsService.GetShelfStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetKnowledgeStats returns knowledge base statistics
// GET /api/v1/analytics/knowledge
func (h *AnalyticsHandler) GetKnowledgeStats(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	stats, err := h.analyticsService.GetKnowledgeStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetDashboardConfigs returns all dashboard configurations
// GET /api/v1/analytics/dashboards
func (h *AnalyticsHandler) GetDashboardConfigs(c *gin.Context) {
	configs, err := h.analyticsService.GetDashboardConfigs(c.Request.Context())
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, configs)
}

// GetDashboardConfig returns a specific dashboard configuration
// GET /api/v1/analytics/dashboards/:id
func (h *AnalyticsHandler) GetDashboardConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequestResponse(c, "dashboard ID is required")
		return
	}

	config, err := h.analyticsService.GetDashboardConfig(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "dashboard config not found" {
			NotFoundResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, config)
}

// CreateDashboardConfigRequest represents the request body for creating a dashboard config
type CreateDashboardConfigRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Layout      string `json:"layout" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

// CreateDashboardConfig creates a new dashboard configuration
// POST /api/v1/analytics/dashboards
func (h *AnalyticsHandler) CreateDashboardConfig(c *gin.Context) {
	var req CreateDashboardConfigRequest
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

	config := &services.DashboardConfig{
		Name:        req.Name,
		Description: req.Description,
		Layout:      req.Layout,
		IsDefault:   req.IsDefault,
		CreatedBy:   userIDStr,
	}

	if err := h.analyticsService.CreateDashboardConfig(c.Request.Context(), config); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "dashboard config created successfully",
		"data":    config,
	})
}

// UpdateDashboardConfigRequest represents the request body for updating a dashboard config
type UpdateDashboardConfigRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Layout      string `json:"layout" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

// UpdateDashboardConfig updates a dashboard configuration
// PUT /api/v1/analytics/dashboards/:id
func (h *AnalyticsHandler) UpdateDashboardConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequestResponse(c, "dashboard ID is required")
		return
	}

	var req UpdateDashboardConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	config := &services.DashboardConfig{
		Name:        req.Name,
		Description: req.Description,
		Layout:      req.Layout,
		IsDefault:   req.IsDefault,
	}

	if err := h.analyticsService.UpdateDashboardConfig(c.Request.Context(), id, config); err != nil {
		if err.Error() == "dashboard config not found" {
			NotFoundResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "dashboard config updated successfully"})
}

// DeleteDashboardConfig deletes a dashboard configuration
// DELETE /api/v1/analytics/dashboards/:id
func (h *AnalyticsHandler) DeleteDashboardConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequestResponse(c, "dashboard ID is required")
		return
	}

	if err := h.analyticsService.DeleteDashboardConfig(c.Request.Context(), id); err != nil {
		if err.Error() == "dashboard config not found" {
			NotFoundResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "dashboard config deleted successfully"})
}

// SetDefaultDashboard sets a dashboard as the default
// PUT /api/v1/analytics/dashboards/:id/default
func (h *AnalyticsHandler) SetDefaultDashboard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequestResponse(c, "dashboard ID is required")
		return
	}

	if err := h.analyticsService.SetDefaultDashboard(c.Request.Context(), id); err != nil {
		if err.Error() == "invalid dashboard ID" {
			BadRequestResponse(c, err.Error())
			return
		}
		if err.Error() == "dashboard config not found" {
			NotFoundResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "default dashboard set successfully"})
}

// GetProjectProgressTrends returns project progress trends
// GET /api/v1/analytics/projects/trends
func (h *AnalyticsHandler) GetProjectProgressTrends(c *gin.Context) {
	// TODO: Implement project progress trends
	// This would require a project_progress_history table

	SuccessResponse(c, gin.H{
		"message": "not implemented - requires project progress history",
		"data":    []interface{}{},
	})
}

// GetComparisonData compares data between two periods
// GET /api/v1/analytics/compare
func (h *AnalyticsHandler) GetComparisonData(c *gin.Context) {
	// Parse current period
	currentStartStr := c.Query("current_start")
	currentEndStr := c.Query("current_end")
	prevStartStr := c.Query("prev_start")
	prevEndStr := c.Query("prev_end")

	if currentStartStr == "" || currentEndStr == "" || prevStartStr == "" || prevEndStr == "" {
		BadRequestResponse(c, "all date parameters are required: current_start, current_end, prev_start, prev_end")
		return
	}

	currentStart, err := time.Parse("2006-01-02", currentStartStr)
	if err != nil {
		BadRequestResponse(c, "invalid current_start date format")
		return
	}

	currentEnd, err := time.Parse("2006-01-02", currentEndStr)
	if err != nil {
		BadRequestResponse(c, "invalid current_end date format")
		return
	}

	prevStart, err := time.Parse("2006-01-02", prevStartStr)
	if err != nil {
		BadRequestResponse(c, "invalid prev_start date format")
		return
	}

	prevEnd, err := time.Parse("2006-01-02", prevEndStr)
	if err != nil {
		BadRequestResponse(c, "invalid prev_end date format")
		return
	}

	// Set end dates to end of day
	currentEnd = currentEnd.Add(24 * time.Hour).Add(-time.Second)
	prevEnd = prevEnd.Add(24 * time.Hour).Add(-time.Second)

	comparison, err := h.aggregationService.ComparePeriods(c.Request.Context(), currentStart, currentEnd, prevStart, prevEnd)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, comparison)
}

// ExportStatsRequest represents the request body for exporting statistics
type ExportStatsRequest struct {
	Type   string `json:"type" binding:"required,oneof=projects users shelf knowledge"`
	Format string `json:"format" binding:"required,oneof=json csv excel"`
}

// ExportStats exports statistics in various formats
// POST /api/v1/analytics/export
func (h *AnalyticsHandler) ExportStats(c *gin.Context) {
	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid date format: "+err.Error())
		return
	}

	exportType := c.DefaultQuery("type", "projects")
	format := c.DefaultQuery("format", "json")

	var exportData interface{}

	switch exportType {
	case "projects":
		exportData, err = h.aggregationService.ExportProjectStats(c.Request.Context(), startDate, endDate)
	case "users":
		exportData, err = h.aggregationService.ExportUserStats(c.Request.Context(), startDate, endDate)
	case "shelf":
		exportData, err = h.aggregationService.ExportShelfStats(c.Request.Context(), startDate, endDate)
	default:
		BadRequestResponse(c, "invalid export type")
		return
	}

	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	// Handle different export formats
	switch format {
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=analytics_export.csv")
		// TODO: Implement CSV conversion
		c.String(http.StatusOK, "CSV export not yet implemented")
	case "excel":
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=analytics_export.xlsx")
		// TODO: Implement Excel conversion
		c.String(http.StatusOK, "Excel export not yet implemented")
	default:
		SuccessResponse(c, exportData)
	}
}

// GenerateSnapshot generates a statistics snapshot for the current date
// POST /api/v1/analytics/snapshot
func (h *AnalyticsHandler) GenerateSnapshot(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)

	if err := h.analyticsService.GenerateProjectStatsSnapshot(c.Request.Context(), today); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"message": "snapshot generated successfully",
		"date":    today.Format("2006-01-02"),
	})
}
