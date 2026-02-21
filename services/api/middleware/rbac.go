package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RBACMiddleware handles role-based access control
type RBACMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewRBACMiddleware creates a new RBACMiddleware
func NewRBACMiddleware(enforcer *casbin.Enforcer) *RBACMiddleware {
	return &RBACMiddleware{
		enforcer: enforcer,
	}
}

// RequirePermission checks if user has the required permission
func (m *RBACMiddleware) RequirePermission(object, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    4030,
				"message": "access denied: no role found",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// Check permission using Casbin
		allowed, err := m.enforcer.Enforce(role, object, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "error checking permissions",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    4031,
				"message": "access denied: insufficient permissions",
				"data":    gin.H{
					"required": gin.H{
						"object": object,
						"action": action,
					},
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole checks if user has one of the required roles
func (m *RBACMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    4030,
				"message": "access denied: no role found",
				"data":    nil,
			})
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    4032,
			"message": "access denied: insufficient role",
			"data": gin.H{
				"required_roles": roles,
				"current_role":  userRole,
			},
		})
		c.Abort()
	}
}

// RequireAdmin requires admin role
func (m *RBACMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// RequireManager requires manager or admin role
func (m *RBACMiddleware) RequireManager() gin.HandlerFunc {
	return m.RequireRole("admin", "manager")
}
