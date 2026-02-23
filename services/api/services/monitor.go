package services

import (
	"context"
	"fmt"
	"time"

	"rdp-platform/rdp-api/models"

	"gorm.io/gorm"
)

// MonitorService handles system and API monitoring
type MonitorService struct {
	db *gorm.DB
}

// NewMonitorService creates a new MonitorService
func NewMonitorService(db *gorm.DB) *MonitorService {
	return &MonitorService{db: db}
}

// CreateSystemMetric creates a new system metric record
func (s *MonitorService) CreateSystemMetric(ctx context.Context, metric *models.SystemMetric) error {
	return s.db.WithContext(ctx).Create(metric).Error
}

// GetSystemMetrics retrieves system metrics with time range
func (s *MonitorService) GetSystemMetrics(ctx context.Context, startTime, endTime time.Time, limit int) ([]models.SystemMetric, error) {
	var metrics []models.SystemMetric
	query := s.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&metrics).Error
	return metrics, err
}

// GetLatestSystemMetric gets the most recent system metric
func (s *MonitorService) GetLatestSystemMetric(ctx context.Context) (*models.SystemMetric, error) {
	var metric models.SystemMetric
	err := s.db.WithContext(ctx).
		Order("timestamp DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// GetSystemMetricStats returns aggregated system metrics
func (s *MonitorService) GetSystemMetricStats(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	var result struct {
		AvgCPU          float64
		MaxCPU          float64
		AvgMemory       float64
		MaxMemory       float64
		AvgDisk         float64
		MaxDisk         float64
		TotalNetworkIn  int64
		TotalNetworkOut int64
	}

	err := s.db.WithContext(ctx).Model(&models.SystemMetric{}).
		Select(
			"AVG(cpu_usage) as avg_cpu",
			"MAX(cpu_usage) as max_cpu",
			"AVG(memory_usage) as avg_memory",
			"MAX(memory_usage) as max_memory",
			"AVG(disk_usage) as avg_disk",
			"MAX(disk_usage) as max_disk",
			"SUM(network_in) as total_network_in",
			"SUM(network_out) as total_network_out",
		).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"cpu": map[string]float64{
			"avg": result.AvgCPU,
			"max": result.MaxCPU,
		},
		"memory": map[string]float64{
			"avg": result.AvgMemory,
			"max": result.MaxMemory,
		},
		"disk": map[string]float64{
			"avg": result.AvgDisk,
			"max": result.MaxDisk,
		},
		"network": map[string]int64{
			"in":  result.TotalNetworkIn,
			"out": result.TotalNetworkOut,
		},
	}, nil
}

// CreateAPIMetric creates a new API metric record
func (s *MonitorService) CreateAPIMetric(ctx context.Context, metric *models.APIMetric) error {
	return s.db.WithContext(ctx).Create(metric).Error
}

// GetAPIMetrics retrieves API metrics with filters
func (s *MonitorService) GetAPIMetrics(ctx context.Context, endpoint string, startTime, endTime time.Time, limit int) ([]models.APIMetric, error) {
	var metrics []models.APIMetric
	query := s.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp DESC")

	if endpoint != "" {
		query = query.Where("endpoint = ?", endpoint)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&metrics).Error
	return metrics, err
}

// GetAPIMetricStats returns aggregated API metrics
func (s *MonitorService) GetAPIMetricStats(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	type EndpointStats struct {
		Endpoint    string
		Count       int64
		AvgDuration float64
		MaxDuration int64
		MinDuration int64
		ErrorCount  int64
	}

	var stats []EndpointStats
	err := s.db.WithContext(ctx).Model(&models.APIMetric{}).
		Select(
			"endpoint",
			"COUNT(*) as count",
			"AVG(duration) as avg_duration",
			"MAX(duration) as max_duration",
			"MIN(duration) as min_duration",
			"SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) as error_count",
		).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("endpoint").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	var totalRequests int64
	var totalErrors int64
	var totalDuration float64

	for _, s := range stats {
		totalRequests += s.Count
		totalErrors += s.ErrorCount
		totalDuration += s.AvgDuration * float64(s.Count)
	}

	avgDuration := float64(0)
	if totalRequests > 0 {
		avgDuration = totalDuration / float64(totalRequests)
	}

	errorRate := float64(0)
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"total_requests": totalRequests,
		"total_errors":   totalErrors,
		"error_rate":     errorRate,
		"avg_duration":   avgDuration,
		"endpoints":      stats,
	}, nil
}

// CreateLogEntry creates a new log entry
func (s *MonitorService) CreateLogEntry(ctx context.Context, entry *models.LogEntry) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

// GetLogEntries retrieves log entries with filters
func (s *MonitorService) GetLogEntries(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]models.LogEntry, int64, error) {
	var entries []models.LogEntry
	var total int64

	query := s.db.WithContext(ctx).Model(&models.LogEntry{})

	// Apply filters
	if level, ok := filters["level"].(string); ok && level != "" {
		query = query.Where("level = ?", level)
	}
	if source, ok := filters["source"].(string); ok && source != "" {
		query = query.Where("source = ?", source)
	}
	if module, ok := filters["module"].(string); ok && module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
		query = query.Where("message ILIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if startTime, ok := filters["start_time"].(time.Time); ok {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime, ok := filters["end_time"].(time.Time); ok {
		query = query.Where("timestamp <= ?", endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Order("timestamp DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&entries).Error

	return entries, total, err
}

// GetLogSources returns unique log sources
func (s *MonitorService) GetLogSources(ctx context.Context) ([]string, error) {
	var sources []string
	err := s.db.WithContext(ctx).Model(&models.LogEntry{}).
		Select("DISTINCT source").
		Pluck("source", &sources).Error
	return sources, err
}

// CleanupOldMetrics removes metrics older than retention period
func (s *MonitorService) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	// Cleanup system metrics
	if err := s.db.WithContext(ctx).Where("timestamp < ?", cutoff).Delete(&models.SystemMetric{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup system metrics: %w", err)
	}

	// Cleanup API metrics
	if err := s.db.WithContext(ctx).Where("timestamp < ?", cutoff).Delete(&models.APIMetric{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup api metrics: %w", err)
	}

	// Cleanup log entries
	if err := s.db.WithContext(ctx).Where("timestamp < ?", cutoff).Delete(&models.LogEntry{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup log entries: %w", err)
	}

	return nil
}
