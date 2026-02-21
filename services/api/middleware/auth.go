package middleware

import (
	"net/http"
	"strings"

	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	userService *services.UserService
	jwtSecret  string
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(userService *services.UserService, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		jwtSecret:  jwtSecret,
	}
}

// Authenticate validates JWT token and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4010,
				"message": "authorization header required",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4011,
				"message": "invalid authorization header format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token (simplified - in production use proper JWT library)
		claims, err := m.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4012,
				"message": "invalid or expired token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// validateToken validates JWT token and returns claims
func (m *AuthMiddleware) validateToken(token string) (map[string]interface{}, error) {
	// TODO: Implement proper JWT validation with jwt-go library
	// This is a placeholder that needs to be replaced with actual JWT validation
	
	// For now, return a placeholder claims map
	// In production: use jwt-go to parse and validate the token
	claims := map[string]interface{}{
		"user_id":  "placeholder",
		"username": "placeholder",
		"role":     "designer",
	}

	return claims, nil
}

// OptionalAuth authenticates if token is provided, otherwise continues without auth
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := m.validateToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
