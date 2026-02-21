package handlers

import (
	"net/http"
	"strconv"
	"time"

	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
)

// SecurityHandler handles security HTTP requests
type SecurityHandler struct {
	securityService *services.SecurityService
	sessionService  *services.SessionService
}

// NewSecurityHandler creates a new SecurityHandler
func NewSecurityHandler(securityService *services.SecurityService, sessionService *services.SessionService) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
		sessionService:  sessionService,
	}
}

// ListAuditLogs handles GET /api/v1/audit-logs
func (h *SecurityHandler) ListAuditLogs(c *gin.Context) {
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
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}
	if classification := c.Query("classification"); classification != "" {
		filters["classification"] = classification
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t
		}
	}

	logs, total, err := h.securityService.ListAuditLogs(c.Request.Context(), page, pageSize, filters)
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
			"items":     logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAuditLog handles GET /api/v1/audit-logs/:id
func (h *SecurityHandler) GetAuditLog(c *gin.Context) {
	id := c.Param("id")

	log, err := h.securityService.GetAuditLogByID(c.Request.Context(), id)
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
		"data":    log,
	})
}

// ListLoginLogs handles GET /api/v1/login-logs
func (h *SecurityHandler) ListLoginLogs(c *gin.Context) {
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
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if username := c.Query("username"); username != "" {
		filters["username"] = username
	}
	if success := c.Query("success"); success != "" {
		filters["success"] = success == "true"
	}

	logs, total, err := h.securityService.ListLoginLogs(c.Request.Context(), page, pageSize, filters)
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
			"items":     logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDataClassifications handles GET /api/v1/data-classifications
func (h *SecurityHandler) GetDataClassifications(c *gin.Context) {
	classifications, err := h.securityService.GetDataClassifications(c.Request.Context())
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
		"data":    classifications,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *SecurityHandler) Logout(c *gin.Context) {
	// Get token from header
	token := c.GetHeader("Authorization")
	if token != "" {
		// Remove "Bearer " prefix
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		_ = h.sessionService.RevokeSession(c.Request.Context(), token)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "logged out successfully",
		"data":    nil,
	})
}

// RevokeAllSessions handles POST /api/v1/auth/revoke-all-sessions
func (h *SecurityHandler) RevokeAllSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	err := h.sessionService.RevokeAllUserSessions(c.Request.Context(), userID.(string))
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
		"message": "all sessions revoked",
		"data":    nil,
	})
}
