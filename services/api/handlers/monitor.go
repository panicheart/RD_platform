package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rdp-platform/rdp-api/collectors"
	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// MonitorHandler handles monitoring HTTP requests
type MonitorHandler struct {
	monitorService  *services.MonitorService
	alertingService *services.AlertingService
}

// NewMonitorHandler creates a new MonitorHandler
func NewMonitorHandler(monitorService *services.MonitorService, alertingService *services.AlertingService) *MonitorHandler {
	return &MonitorHandler{
		monitorService:  monitorService,
		alertingService: alertingService,
	}
}

// GetSystemMetrics returns system metrics with time range
// GET /api/v1/monitor/metrics/system
func (h *MonitorHandler) GetSystemMetrics(c *gin.Context) {
	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid time range: "+err.Error())
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	metrics, err := h.monitorService.GetSystemMetrics(c.Request.Context(), startTime, endTime, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, metrics)
}

// GetSystemMetricStats returns aggregated system metrics statistics
// GET /api/v1/monitor/metrics/system/stats
func (h *MonitorHandler) GetSystemMetricStats(c *gin.Context) {
	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid time range: "+err.Error())
		return
	}

	stats, err := h.monitorService.GetSystemMetricStats(c.Request.Context(), startTime, endTime)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetAPIMetrics returns API metrics with filters
// GET /api/v1/monitor/metrics/api
func (h *MonitorHandler) GetAPIMetrics(c *gin.Context) {
	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid time range: "+err.Error())
		return
	}

	endpoint := c.Query("endpoint")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	metrics, err := h.monitorService.GetAPIMetrics(c.Request.Context(), endpoint, startTime, endTime, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, metrics)
}

// GetAPIMetricStats returns aggregated API metrics statistics
// GET /api/v1/monitor/metrics/api/stats
func (h *MonitorHandler) GetAPIMetricStats(c *gin.Context) {
	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		BadRequestResponse(c, "invalid time range: "+err.Error())
		return
	}

	stats, err := h.monitorService.GetAPIMetricStats(c.Request.Context(), startTime, endTime)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetPrometheusMetrics returns metrics in Prometheus format
// GET /api/v1/monitor/metrics/prometheus
func (h *MonitorHandler) GetPrometheusMetrics(c *gin.Context) {
	// Get latest system metric
	latestMetric, err := h.monitorService.GetLatestSystemMetric(c.Request.Context())
	if err != nil {
		// Return basic info even if no metrics yet
		c.String(http.StatusOK, "# No metrics available yet\n")
		return
	}

	// Format as Prometheus exposition format
	output := "# HELP rdp_cpu_usage CPU usage percentage\n"
	output += "# TYPE rdp_cpu_usage gauge\n"
	output += sprintf("rdp_cpu_usage %.2f\n", latestMetric.CPUUsage)

	output += "# HELP rdp_memory_usage Memory usage percentage\n"
	output += "# TYPE rdp_memory_usage gauge\n"
	output += sprintf("rdp_memory_usage %.2f\n", latestMetric.MemoryUsage)

	output += "# HELP rdp_memory_used Memory used in bytes\n"
	output += "# TYPE rdp_memory_used gauge\n"
	output += sprintf("rdp_memory_used %d\n", latestMetric.MemoryUsed)

	output += "# HELP rdp_disk_usage Disk usage percentage\n"
	output += "# TYPE rdp_disk_usage gauge\n"
	output += sprintf("rdp_disk_usage %.2f\n", latestMetric.DiskUsage)

	c.Header("Content-Type", "text/plain; version=0.0.4")
	c.String(http.StatusOK, output)
}

// GetLogEntries returns log entries with pagination and filters
// GET /api/v1/monitor/logs
func (h *MonitorHandler) GetLogEntries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := map[string]interface{}{
		"level":   c.Query("level"),
		"source":  c.Query("source"),
		"module":  c.Query("module"),
		"keyword": c.Query("keyword"),
	}

	// Parse time range
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			filters["start_time"] = t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			filters["end_time"] = t
		}
	}

	entries, total, err := h.monitorService.GetLogEntries(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items": entries,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateLogEntry creates a new log entry (for internal use)
// POST /api/v1/monitor/logs
func (h *MonitorHandler) CreateLogEntry(c *gin.Context) {
	var entry models.LogEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	entry.Timestamp = time.Now()
	if err := h.monitorService.CreateLogEntry(c.Request.Context(), &entry); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    201,
		Message: "created",
		Data:    entry,
	})
}

// GetLogSources returns unique log sources
// GET /api/v1/monitor/logs/sources
func (h *MonitorHandler) GetLogSources(c *gin.Context) {
	sources, err := h.monitorService.GetLogSources(c.Request.Context())
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, sources)
}

// ListAlertRules returns all alert rules
// GET /api/v1/monitor/alerts/rules
func (h *MonitorHandler) ListAlertRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	rules, total, err := h.alertingService.ListAlertRules(c.Request.Context(), page, pageSize)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items": rules,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateAlertRule creates a new alert rule
// POST /api/v1/monitor/alerts/rules
func (h *MonitorHandler) CreateAlertRule(c *gin.Context) {
	var rule models.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	if err := h.alertingService.CreateAlertRule(c.Request.Context(), &rule); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    201,
		Message: "created",
		Data:    rule,
	})
}

// GetAlertRule returns a specific alert rule
// GET /api/v1/monitor/alerts/rules/:id
func (h *MonitorHandler) GetAlertRule(c *gin.Context) {
	id := c.Param("id")

	rule, err := h.alertingService.GetAlertRule(c.Request.Context(), id)
	if err != nil {
		NotFoundResponse(c, "alert rule not found")
		return
	}

	SuccessResponse(c, rule)
}

// UpdateAlertRule updates an existing alert rule
// PUT /api/v1/monitor/alerts/rules/:id
func (h *MonitorHandler) UpdateAlertRule(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	if err := h.alertingService.UpdateAlertRule(c.Request.Context(), id, updates); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "alert rule updated"})
}

// DeleteAlertRule deletes an alert rule
// DELETE /api/v1/monitor/alerts/rules/:id
func (h *MonitorHandler) DeleteAlertRule(c *gin.Context) {
	id := c.Param("id")

	if err := h.alertingService.DeleteAlertRule(c.Request.Context(), id); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "alert rule deleted"})
}

// GetAlertHistory returns alert history with filters
// GET /api/v1/monitor/alerts/history
func (h *MonitorHandler) GetAlertHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	alerts, total, err := h.alertingService.GetAlertHistory(c.Request.Context(), status, page, pageSize)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items": alerts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ResolveAlert marks an alert as resolved
// PUT /api/v1/monitor/alerts/history/:id/resolve
func (h *MonitorHandler) ResolveAlert(c *gin.Context) {
	id := c.Param("id")

	if err := h.alertingService.ResolveAlert(c.Request.Context(), id); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "alert resolved"})
}

// GetAlertStats returns alert statistics
// GET /api/v1/monitor/alerts/stats
func (h *MonitorHandler) GetAlertStats(c *gin.Context) {
	stats, err := h.alertingService.GetAlertStats(c.Request.Context())
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}

// GetHealth returns detailed health status
// GET /api/v1/monitor/health
func (h *MonitorHandler) GetHealth(c *gin.Context) {
	// Get latest system metric
	latestMetric, err := h.monitorService.GetLatestSystemMetric(c.Request.Context())

	// Get system info
	sysInfo := collectors.GetSystemInfo()

	// Determine overall health
	status := "healthy"
	checks := gin.H{}

	if err == nil {
		// Check CPU
		if latestMetric.CPUUsage > 90 {
			checks["cpu"] = gin.H{"status": "critical", "value": latestMetric.CPUUsage}
			status = "degraded"
		} else if latestMetric.CPUUsage > 70 {
			checks["cpu"] = gin.H{"status": "warning", "value": latestMetric.CPUUsage}
		} else {
			checks["cpu"] = gin.H{"status": "healthy", "value": latestMetric.CPUUsage}
		}

		// Check Memory
		if latestMetric.MemoryUsage > 90 {
			checks["memory"] = gin.H{"status": "critical", "value": latestMetric.MemoryUsage}
			status = "degraded"
		} else if latestMetric.MemoryUsage > 80 {
			checks["memory"] = gin.H{"status": "warning", "value": latestMetric.MemoryUsage}
		} else {
			checks["memory"] = gin.H{"status": "healthy", "value": latestMetric.MemoryUsage}
		}

		// Check Disk
		if latestMetric.DiskUsage > 90 {
			checks["disk"] = gin.H{"status": "critical", "value": latestMetric.DiskUsage}
			status = "degraded"
		} else if latestMetric.DiskUsage > 80 {
			checks["disk"] = gin.H{"status": "warning", "value": latestMetric.DiskUsage}
		} else {
			checks["disk"] = gin.H{"status": "healthy", "value": latestMetric.DiskUsage}
		}
	} else {
		checks["metrics"] = gin.H{"status": "unknown", "message": "no metrics available"}
	}

	health := gin.H{
		"status":      status,
		"timestamp":   time.Now().UTC(),
		"system_info": sysInfo,
		"checks":      checks,
	}

	SuccessResponse(c, health)
}

// GetSystemInfo returns static system information
// GET /api/v1/monitor/system/info
func (h *MonitorHandler) GetSystemInfo(c *gin.Context) {
	info := collectors.GetSystemInfo()

	// Get disk partitions
	partitions, _ := collectors.GetDiskPartitions()
	info["partitions"] = partitions

	SuccessResponse(c, info)
}

// parseTimeRange parses start_time and end_time query parameters
func parseTimeRange(c *gin.Context) (time.Time, time.Time, error) {
	startStr := c.DefaultQuery("start_time", "")
	endStr := c.DefaultQuery("end_time", "")

	var startTime, endTime time.Time
	var err error

	if startStr == "" {
		// Default to 1 hour ago
		startTime = time.Now().Add(-1 * time.Hour)
	} else {
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			// Try parsing as date only
			startTime, err = time.Parse("2006-01-02", startStr)
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
		}
	}

	if endStr == "" {
		// Default to now
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			// Try parsing as date only
			endTime, err = time.Parse("2006-01-02", endStr)
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
		}
	}

	return startTime, endTime, nil
}

// sprintf is a helper function for formatting strings
func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
