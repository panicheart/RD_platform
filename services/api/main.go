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
	"rdp-platform/rdp-api/middleware"
	"rdp-platform/rdp-api/routes"
	"rdp-platform/rdp-api/services"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := initDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userService := services.NewUserService(db)
	projectService := services.NewProjectService(db)
	notificationService := services.NewNotificationService(db)
	forumService := services.NewForumService(db, notificationService)
	knowledgeService := services.NewKnowledgeService(db)

	router := gin.New()

	jwtMiddleware := middleware.NewJWTMiddleware(cfg.Auth.JWTSecret, cfg.Auth.Issuer, cfg.Auth.Audience)

	routerManager := routes.NewRouter(router, userService, projectService, forumService, knowledgeService, jwtMiddleware)
	routerManager.SetupRoutes()

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("RDP API Server starting on %s:%s in %s mode", cfg.Server.Host, cfg.Server.Port, cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("Server exited gracefully")
}

func initDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
