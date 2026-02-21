package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("RDP_ENV")
	if env == "" {
		env = "development"
	}

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	port := os.Getenv("RDP_API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("RDP API Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
