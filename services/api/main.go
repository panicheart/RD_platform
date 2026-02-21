package main

import (
	"log"

	"rdp/services/api/handlers"
	"rdp/services/api/middleware"
	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	loadConfig()

	// Initialize database
	db := initDB()

	// Initialize services
	userService := services.NewUserService(db)
	orgService := services.NewOrganizationService(db)
	notificationService := services.NewNotificationService(db)
	announcementService := services.NewAnnouncementService(db)

	// Initialize Casbin enforcer
	enforcer := initCasbin()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, orgService)
	notificationHandler := handlers.NewNotificationHandler(notificationService, announcementService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userService, viper.GetString("jwt.secret"))
	rbacMiddleware := middleware.NewRBACMiddleware(enforcer)

	// Setup router
	router := gin.Default()
	setupRoutes(router, userHandler, notificationHandler, authMiddleware, rbacMiddleware)

	// Start server
	port := viper.GetString("server.port")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig() {
	// Try to load .env file
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("jwt.secret", "your-secret-key")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found, using defaults")
	}
}

func initDB() *gorm.DB {
	dsn := viper.GetString("database.dsn")
	if dsn == "" {
		dsn = "host=localhost user=rdp password=rdp dbname=rdp port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.Notification{},
		&models.Announcement{},
		&models.Honor{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func initCasbin() *casbin.Enforcer {
	// Use model file and policy file
	enforcer, err := casbin.NewEnforcer("./config/rbac_model.conf", "./config/rbac_policy.csv")
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}
	return enforcer
}

func setupRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	notificationHandler *handlers.NotificationHandler,
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
) {
	// Health check (no auth)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (optional auth)
		announcements := v1.Group("/announcements")
		{
			announcements.GET("", notificationHandler.ListAnnouncements)
			announcements.GET("/:id", notificationHandler.GetAnnouncement)
		}

		// Protected routes (require auth)
		protected := v1.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.CurrentUser)
				users.PUT("/me", userHandler.UpdateCurrentUser)
				users.GET("", rbacMiddleware.RequirePermission("user", "read"), userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("", rbacMiddleware.RequirePermission("user", "write"), userHandler.CreateUser)
				users.PUT("/:id", rbacMiddleware.RequirePermission("user", "write"), userHandler.UpdateUser)
				users.DELETE("/:id", rbacMiddleware.RequireAdmin(), userHandler.DeleteUser)
			}

			// Organization routes
			orgs := protected.Group("/organizations")
			{
				orgs.GET("", rbacMiddleware.RequirePermission("organization", "read"), userHandler.ListOrganizations)
				orgs.GET("/:id", userHandler.GetOrganization)
				orgs.POST("", rbacMiddleware.RequirePermission("organization", "write"), userHandler.CreateOrganization)
				orgs.PUT("/:id", rbacMiddleware.RequirePermission("organization", "write"), userHandler.UpdateOrganization)
				orgs.DELETE("/:id", rbacMiddleware.RequireAdmin(), userHandler.DeleteOrganization)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.ListNotifications)
				notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
				notifications.GET("/:id", notificationHandler.GetNotification)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notifications.DELETE("/:id", notificationHandler.DeleteNotification)
			}

			// Admin: Announcement management
			adminAnnouncements := protected.Group("/announcements")
			adminAnnouncements.Use(rbacMiddleware.RequireRole("admin", "manager"))
			{
				adminAnnouncements.POST("", notificationHandler.CreateAnnouncement)
				adminAnnouncements.PUT("/:id", notificationHandler.UpdateAnnouncement)
				adminAnnouncements.DELETE("/:id", notificationHandler.DeleteAnnouncement)
			}
		}
	}
}
