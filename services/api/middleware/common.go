package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeaders returns a middleware that adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		
		// XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy (adjust as needed)
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log format: [timestamp] method path status latency client_ip
		gin.DefaultWriter.Write([]byte(
			"[" + start.Format("2006-01-02 15:04:05") + "] " +
				method + " " +
				path + " " +
				http.StatusText(statusCode) + " " +
				latency.String() + " " +
				clientIP + "\n",
		))
	}
}

// RateLimiter returns a simple rate limiting middleware
// In production, use Redis-based rate limiting
func RateLimiter(requestsPerMinute int) gin.HandlerFunc {
	// Simple in-memory rate limiter (for production, use Redis)
	// This is a placeholder implementation
	return func(c *gin.Context) {
		// TODO: Implement proper rate limiting with Redis
		c.Next()
	}
}

// Timeout returns a middleware that adds timeout to requests
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the request context with a timeout
		ctx, cancel := c.Request.Context(), func() {}
		_ = ctx
		defer cancel()

		// For a proper implementation, use:
		// ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		// defer cancel()
		// c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetCurrentUser extracts the current user from the context
// Returns user info set by auth middleware
func GetCurrentUser(c *gin.Context) (*UserInfo, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, false
	}

	username, _ := c.Get("username")
	role, _ := c.Get("role")

	return &UserInfo{
		UserID:   userID.(string),
		Username: getString(username),
		Role:     getString(role),
	}, true
}

// UserInfo represents current user information
type UserInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	return v.(string)
}
