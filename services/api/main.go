package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rdp-platform/rdp-api/config"
	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/routes"
	"rdp-platform/rdp-api/services"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库连接
	db, err := initDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库模型
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化服务
	userService := services.NewUserService(db, cfg.Auth)

	// 创建默认管理员（如果不存在）
	if err := createDefaultAdmin(userService); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	// 创建Gin引擎
	router := gin.New()

	// 配置路由
	routerManager := routes.NewRouter(router, userService)
	routerManager.SetupRoutes()

	// 开发环境启用测试路由
	if cfg.IsDevelopment() {
		routerManager.SetupTestRoutes()
		log.Println("Test routes enabled at /test/*")
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 优雅启动和关闭
	go func() {
		log.Printf("RDP API Server starting on %s:%s in %s mode", cfg.Server.Host, cfg.Server.Port, cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("Server exited gracefully")
}

// initDatabase 初始化数据库连接
func initDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
	logLevel := logger.Silent
	if os.Getenv("RDP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// autoMigrate 自动迁移数据库模型
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.TokenBlacklist{},
	)
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin(userService *services.UserService) error {
	ctx := context.Background()

	// 检查是否已存在管理员
	_, err := userService.GetUserByUsername(ctx, "admin")
	if err == nil {
		// 管理员已存在
		return nil
	}
	if err != services.ErrUserNotFound {
		return err
	}

	// 创建默认管理员
	req := models.CreateUserRequest{
		Username:    "admin",
		Password:    "Admin@123",
		DisplayName: "系统管理员",
		Email:       "admin@rdp.local",
		Role:        models.RoleAdmin,
		Team:        models.TeamGeneralMgmt,
	}

	_, err = userService.CreateUser(ctx, req)
	if err != nil {
		return err
	}

	log.Println("Default admin user created (username: admin, password: Admin@123)")
	return nil
}
