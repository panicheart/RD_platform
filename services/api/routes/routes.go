package routes

import (
	"github.com/gin-gonic/gin"

	"rdp-platform/rdp-api/handlers"
	"rdp-platform/rdp-api/middleware"
	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"
)

// Router 路由管理器
type Router struct {
	engine      *gin.Engine
	userService *services.UserService
	authMW      *middleware.AuthMiddleware
}

// NewRouter 创建路由管理器
func NewRouter(engine *gin.Engine, userService *services.UserService) *Router {
	return &Router{
		engine:      engine,
		userService: userService,
		authMW:      middleware.NewAuthMiddleware(userService),
	}
}

// SetupRoutes 配置所有路由
func (r *Router) SetupRoutes() {
	// 全局中间件
	r.setupGlobalMiddleware()

	// 健康检查（公开）
	r.setupHealthRoutes()

	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	{
		// 认证路由（公开）
		r.setupAuthRoutes(v1)

		// 用户路由（需要认证）
		r.setupUserRoutes(v1)

		// 项目路由（需要认证）
		// TODO: 实现项目路由
		// r.setupProjectRoutes(v1)
	}
}

// setupGlobalMiddleware 配置全局中间件
func (r *Router) setupGlobalMiddleware() {
	// CORS
	r.engine.Use(middleware.CORS())
	// 安全头部
	r.engine.Use(middleware.SecurityHeaders())
	// 恢复
	r.engine.Use(gin.Recovery())
}

// setupHealthRoutes 配置健康检查路由
func (r *Router) setupHealthRoutes() {
	healthHandler := handlers.NewHealthHandler()
	r.engine.GET("/api/v1/health", healthHandler.HealthCheck)
}

// setupAuthRoutes 配置认证相关路由
func (r *Router) setupAuthRoutes(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler(r.userService)

	auth := group.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
		auth.POST("/logout", r.authMW.JWTAuth(), userHandler.Logout)
	}
}

// setupUserRoutes 配置用户相关路由
func (r *Router) setupUserRoutes(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler(r.userService)

	users := group.Group("/users")
	{
		// 公开/可选认证路由
		// 获取用户列表（需要认证）
		users.GET("", r.authMW.JWTAuth(), userHandler.GetUserList)

		// 获取当前用户信息（需要认证）
		users.GET("/me", r.authMW.JWTAuth(), userHandler.GetCurrentUser)

		// 创建用户（需要管理员权限）
		users.POST("", r.authMW.RequireAdmin(), userHandler.CreateUser)

		// 单个用户路由
		user := users.Group("/:id")
		{
			// 获取用户详情（需要认证）
			user.GET("", r.authMW.JWTAuth(), userHandler.GetUser)

			// 获取用户Profile（需要认证）
			user.GET("/profile", r.authMW.JWTAuth(), userHandler.GetUserProfile)

			// 更新用户（需要管理员或本人）
			user.PUT("", r.authMW.RequireAdminOrSelf(), userHandler.UpdateUser)

			// 删除用户（需要管理员权限）
			user.DELETE("", r.authMW.RequireAdmin(), userHandler.DeleteUser)

			// 修改密码（需要管理员或本人）
			user.PUT("/password", r.authMW.RequireAdminOrSelf(), userHandler.ChangePassword)
		}
	}
}

// SetupTestRoutes 配置测试路由（用于开发和测试）
func (r *Router) SetupTestRoutes() {
	// 测试端点，用于验证中间件
	r.engine.GET("/test/auth", r.authMW.JWTAuth(), func(c *gin.Context) {
		user, _ := middleware.GetCurrentUser(c)
		c.JSON(200, gin.H{
			"message":  "Authenticated",
			"user_id":  user.UserID,
			"username": user.Username,
			"role":     user.Role,
		})
	})

	r.engine.GET("/test/admin", r.authMW.RequireAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Admin only",
		})
	})

	r.engine.GET("/test/role", r.authMW.RequireRole(models.RoleAdmin, models.RoleDeptLeader), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Admin or DeptLeader",
		})
	})
}

// RoleHierarchy 角色层级定义（用于权限检查）
var RoleHierarchy = map[models.UserRole]int{
	models.RoleOther:      0,
	models.RoleDesigner:   1,
	models.RoleTeamLeader: 2,
	models.RoleDeptLeader: 3,
	models.RoleAdmin:      4,
}

// HasMinimumRole 检查角色是否满足最低要求
func HasMinimumRole(userRole models.UserRole, minRole models.UserRole) bool {
	return RoleHierarchy[userRole] >= RoleHierarchy[minRole]
}
