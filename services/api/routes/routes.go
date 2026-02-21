package routes

import (
	"github.com/gin-gonic/gin"

	"rdp/services/api/handlers"
	"rdp/services/api/middleware"
	"rdp/services/api/models"
	"rdp/services/api/services"
)

// Router manages all application routes
type Router struct {
	engine          *gin.Engine
	userService     *services.UserService
	projectService  *services.ProjectService
	authMiddleware  *middleware.AuthMiddleware
}

// NewRouter creates a new Router
func NewRouter(
	engine *gin.Engine,
	userService *services.UserService,
	projectService *services.ProjectService,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		engine:          engine,
		userService:     userService,
		projectService:  projectService,
		authMiddleware:  authMiddleware,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	// Global middleware
	r.setupGlobalMiddleware()

	// Health check (public)
	r.setupHealthRoutes()

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes (public)
		r.setupAuthRoutes(v1)

		// User routes (authenticated)
		r.setupUserRoutes(v1)

		// Project routes (authenticated)
		r.setupProjectRoutes(v1)
	}
}

// setupGlobalMiddleware configures global middleware
func (r *Router) setupGlobalMiddleware() {
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.SecurityHeaders())
	r.engine.Use(gin.Recovery())
}

// setupHealthRoutes configures health check routes
func (r *Router) setupHealthRoutes() {
	healthHandler := handlers.NewHealthHandler()
	r.engine.GET("/api/v1/health", healthHandler.HealthCheck)
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler(r.userService, nil)

	auth := group.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
		auth.POST("/logout", r.authMiddleware.Authenticate(), userHandler.Logout)
	}
}

// setupUserRoutes configures user routes
func (r *Router) setupUserRoutes(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler(r.userService, nil)

	users := group.Group("/users")
	users.Use(r.authMiddleware.Authenticate())
	{
		// List users
		users.GET("", userHandler.ListUsers)

		// Current user
		users.GET("/me", userHandler.CurrentUser)
		users.PUT("/me", userHandler.UpdateCurrentUser)

		// User projects
		users.GET("/me/projects", r.projectHandler().GetUserProjects)

		// Create user (admin only)
		users.POST("", r.requireRole("admin"), userHandler.CreateUser)

		// Single user routes
		user := users.Group("/:id")
		{
			user.GET("", userHandler.GetUser)
			user.PUT("", r.requireRoleOrSelf("admin"), userHandler.UpdateUser)
			user.DELETE("", r.requireRole("admin"), userHandler.DeleteUser)
		}
	}
}

// setupProjectRoutes configures project routes
func (r *Router) setupProjectRoutes(group *gin.RouterGroup) {
	projectHandler := r.projectHandler()

	projects := group.Group("/projects")
	projects.Use(r.authMiddleware.Authenticate())
	{
		// List and create projects
		projects.GET("", projectHandler.GetProjects)
		projects.POST("", projectHandler.CreateProject)

		// Project stats
		projects.GET("/stats", projectHandler.GetProjectStats)

		// Single project routes
		project := projects.Group("/:id")
		{
			project.GET("", projectHandler.GetProject)
			project.PUT("", projectHandler.UpdateProject)
			project.DELETE("", projectHandler.DeleteProject)

			// Progress
			project.PUT("/progress", projectHandler.UpdateProgress)

			// Gantt chart data
			project.GET("/gantt", projectHandler.GetProjectGantt)

			// Members
			project.GET("/members", projectHandler.GetMembers)
			project.POST("/members", projectHandler.AddMember)
			project.DELETE("/members/:userId", projectHandler.RemoveMember)
			project.PUT("/members/:userId/role", projectHandler.UpdateMemberRole)

			// Activities
			project.GET("/activities", projectHandler.GetProjectActivities)
			project.POST("/activities", projectHandler.CreateActivity)
		}
	}
}

// projectHandler creates a new ProjectHandler instance
func (r *Router) projectHandler() *handlers.ProjectHandler {
	return handlers.NewProjectHandler(r.projectService)
}

// requireRole middleware requires specific role
func (r *Router) requireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"code": 401, "message": "unauthorized"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if string(r) == userRole {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"code": 403, "message": "forbidden"})
		c.Abort()
	}
}

// requireRoleOrSelf middleware requires specific role or self
func (r *Router) requireRoleOrSelf(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"code": 401, "message": "unauthorized"})
			c.Abort()
			return
		}

		// Check if accessing self
		paramID := c.Param("id")
		if userID == paramID {
			c.Next()
			return
		}

		// Check role
		role, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"code": 401, "message": "unauthorized"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if string(r) == userRole {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"code": 403, "message": "forbidden"})
		c.Abort()
	}
}

// RoleHierarchy defines role hierarchy for permission checking
var RoleHierarchy = map[string]int{
	"other":       0,
	"designer":    1,
	"team_leader": 2,
	"dept_leader": 3,
	"admin":       4,
}

// HasMinimumRole checks if user has minimum required role
func HasMinimumRole(userRole string, minRole string) bool {
	return RoleHierarchy[userRole] >= RoleHierarchy[minRole]
}
