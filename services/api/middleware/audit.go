package middleware

import (
	"bytes"
	"io"
	"time"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs all API requests
type AuditMiddleware struct {
	securityService *services.SecurityService
}

// NewAuditMiddleware creates a new AuditMiddleware
func NewAuditMiddleware(securityService *services.SecurityService) *AuditMiddleware {
	return &AuditMiddleware{
		securityService: securityService,
	}
}

// Audit logs HTTP requests
func (m *AuditMiddleware) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health check and metrics endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Start time
		startTime := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Get user info if authenticated
		var userID *string
		var username *string
		if uid, exists := c.Get("user_id"); exists {
			uidStr := uid.(string)
			userID = &uidStr
		}
		if uname, exists := c.Get("username"); exists {
			unameStr := uname.(string)
			username = &unameStr
		}

		// Determine classification based on resource
		classification := m.determineClassification(c.Request.URL.Path)

		// Create audit log
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		log := &models.AuditLog{
			UserID:         nil, // Would need to parse UUID
			Username:       username,
			IPAddress:      m.getClientIP(c),
			UserAgent:      &userAgent,
			Action:         method,
			Resource:       m.getResource(path),
			ResourceID:     m.getResourceID(c),
			Method:         &method,
			Path:           &path,
			ResponseCode:   &statusCode,
			Classification: classification,
		}

		if userID != nil && *userID != "" {
			// userID would be parsed to UUID here
		}

		// Log asynchronously to not block request
		go func() {
			_ = m.securityService.CreateAuditLog(c.Request.Context(), log)
		}()

		_ = duration // Would be used for performance logging
	}
}

// determineClassification determines the data classification based on the resource
func (m *AuditMiddleware) determineClassification(path string) string {
	// Sensitive paths get higher classification
	sensitivePaths := []string{
		"/api/v1/users",
		"/api/v1/auth",
		"/api/v1/admin",
		"/api/v1/audit-logs",
	}

	for _, sp := range sensitivePaths {
		if len(path) >= len(sp) && path[:len(sp)] == sp {
			return "confidential"
		}
	}

	return "internal"
}

// getResource extracts the resource from the path
func (m *AuditMiddleware) getResource(path string) string {
	// Remove /api/v1/ prefix and get first segment
	if len(path) > 9 && path[:9] == "/api/v1/" {
		path = path[9:]
	}
	// Get first segment
	for i, c := range path {
		if c == '/' {
			return path[:i]
		}
	}
	return path
}

// getResourceID extracts the resource ID from the path
func (m *AuditMiddleware) getResourceID(c *gin.Context) *string {
	// Try to get ID from URL params
	if id := c.Param("id"); id != "" {
		return &id
	}
	if id := c.Param("userId"); id != "" {
		return &id
	}
	if id := c.Param("projectId"); id != "" {
		return &id
	}
	return nil
}

// getClientIP gets the client IP from the request
func (m *AuditMiddleware) getClientIP(c *gin.Context) *string {
	// Check X-Forwarded-For header first (for reverse proxy)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// Take the first IP (original client)
		for i, c := range xff {
			if c == ',' {
				ip := xff[:i]
				return &ip
			}
		}
		return &xff
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return &xri
	}

	// Fall back to RemoteAddr
	ip := c.Request.RemoteAddr
	// Remove port if present
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] == ':' {
			ip = ip[:i]
			break
		}
	}
	return &ip
}
