package middleware

import (
	"context"
	"time"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records API metrics for monitoring
type MetricsMiddleware struct {
	monitorService *services.MonitorService
	buffer         chan *models.APIMetric
	bufferSize     int
}

// NewMetricsMiddleware creates a new MetricsMiddleware with async processing
func NewMetricsMiddleware(monitorService *services.MonitorService) *MetricsMiddleware {
	bufferSize := 1000
	m := &MetricsMiddleware{
		monitorService: monitorService,
		buffer:         make(chan *models.APIMetric, bufferSize),
		bufferSize:     bufferSize,
	}

	// Start background worker to process metrics
	go m.processMetrics()

	return m
}

// RecordMetrics returns a gin middleware function that records API calls
func (m *MetricsMiddleware) RecordMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Get user ID from context if available
		userID, _ := c.Get("userID")
		userIDStr := ""
		if id, ok := userID.(string); ok {
			userIDStr = id
		}

		// Create metric record
		metric := &models.APIMetric{
			Timestamp:  time.Now(),
			Endpoint:   c.FullPath(),
			Method:     c.Request.Method,
			Duration:   duration,
			StatusCode: c.Writer.Status(),
			UserID:     userIDStr,
			IPAddress:  c.ClientIP(),
		}

		// Send to buffer (non-blocking)
		select {
		case m.buffer <- metric:
			// Successfully queued
		default:
			// Buffer full, drop metric to avoid blocking
			// In production, this should be logged
		}
	}
}

// processMetrics is a background worker that saves metrics to database
func (m *MetricsMiddleware) processMetrics() {
	for metric := range m.buffer {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := m.monitorService.CreateAPIMetric(ctx, metric)
		cancel()

		if err != nil {
			// Log error but continue processing
			// In production, use proper logging
			println("Failed to save API metric:", err.Error())
		}
	}
}

// Close gracefully shuts down the metrics middleware
func (m *MetricsMiddleware) Close() {
	close(m.buffer)
}

// MetricsConfig holds configuration for metrics middleware
type MetricsConfig struct {
	Enabled       bool
	BufferSize    int
	BatchSize     int
	FlushInterval time.Duration
}

// DefaultMetricsConfig returns default configuration
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:       true,
		BufferSize:    1000,
		BatchSize:     100,
		FlushInterval: 30 * time.Second,
	}
}
