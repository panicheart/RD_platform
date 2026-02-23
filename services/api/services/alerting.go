package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rdp-platform/rdp-api/models"

	"gorm.io/gorm"
)

// AlertingService handles alert rules and notifications
type AlertingService struct {
	db                  *gorm.DB
	monitorService      *MonitorService
	notificationService *NotificationService
}

// NewAlertingService creates a new AlertingService
func NewAlertingService(db *gorm.DB, monitorService *MonitorService, notificationService *NotificationService) *AlertingService {
	return &AlertingService{
		db:                  db,
		monitorService:      monitorService,
		notificationService: notificationService,
	}
}

// CreateAlertRule creates a new alert rule
func (s *AlertingService) CreateAlertRule(ctx context.Context, rule *models.AlertRule) error {
	return s.db.WithContext(ctx).Create(rule).Error
}

// GetAlertRule retrieves an alert rule by ID
func (s *AlertingService) GetAlertRule(ctx context.Context, id string) (*models.AlertRule, error) {
	var rule models.AlertRule
	err := s.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListAlertRules retrieves all alert rules with pagination
func (s *AlertingService) ListAlertRules(ctx context.Context, page, pageSize int) ([]models.AlertRule, int64, error) {
	var rules []models.AlertRule
	var total int64

	if err := s.db.WithContext(ctx).Model(&models.AlertRule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := s.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rules).Error

	return rules, total, err
}

// ListActiveAlertRules retrieves all active alert rules
func (s *AlertingService) ListActiveAlertRules(ctx context.Context) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := s.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&rules).Error
	return rules, err
}

// UpdateAlertRule updates an existing alert rule
func (s *AlertingService) UpdateAlertRule(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&models.AlertRule{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteAlertRule deletes an alert rule
func (s *AlertingService) DeleteAlertRule(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&models.AlertRule{}, "id = ?", id).Error
}

// EvaluateRules checks all active rules against current metrics
func (s *AlertingService) EvaluateRules(ctx context.Context) error {
	rules, err := s.ListActiveAlertRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list active rules: %w", err)
	}

	for _, rule := range rules {
		if err := s.evaluateRule(ctx, &rule); err != nil {
			// Log error but continue evaluating other rules
			fmt.Printf("Failed to evaluate rule %s: %v\n", rule.ID, err)
		}
	}

	return nil
}

// evaluateRule evaluates a single rule against current metrics
func (s *AlertingService) evaluateRule(ctx context.Context, rule *models.AlertRule) error {
	// Get the current metric value based on rule.Metric
	var currentValue float64

	switch rule.Metric {
	case "cpu_usage":
		metric, err := s.monitorService.GetLatestSystemMetric(ctx)
		if err != nil {
			return err
		}
		currentValue = metric.CPUUsage

	case "memory_usage":
		metric, err := s.monitorService.GetLatestSystemMetric(ctx)
		if err != nil {
			return err
		}
		currentValue = metric.MemoryUsage

	case "disk_usage":
		metric, err := s.monitorService.GetLatestSystemMetric(ctx)
		if err != nil {
			return err
		}
		currentValue = metric.DiskUsage

	case "db_connections":
		metric, err := s.monitorService.GetLatestSystemMetric(ctx)
		if err != nil {
			return err
		}
		currentValue = float64(metric.DBConnections)

	default:
		return fmt.Errorf("unknown metric: %s", rule.Metric)
	}

	// Evaluate condition
	triggered := false
	switch rule.Condition {
	case ">":
		triggered = currentValue > rule.Threshold
	case ">=":
		triggered = currentValue >= rule.Threshold
	case "<":
		triggered = currentValue < rule.Threshold
	case "<=":
		triggered = currentValue <= rule.Threshold
	case "==":
		triggered = currentValue == rule.Threshold
	case "!=":
		triggered = currentValue != rule.Threshold
	}

	if triggered {
		// Check if there's already a firing alert for this rule
		var existingAlert models.AlertHistory
		err := s.db.WithContext(ctx).
			Where("rule_id = ? AND status = ?", rule.ID, "firing").
			First(&existingAlert).Error

		if err == gorm.ErrRecordNotFound {
			// Create new alert
			alert := models.AlertHistory{
				RuleID:    rule.ID,
				RuleName:  rule.Name,
				Severity:  rule.Severity,
				Message:   fmt.Sprintf("%s is %s %.2f (current: %.2f)", rule.Metric, rule.Condition, rule.Threshold, currentValue),
				Value:     currentValue,
				Threshold: rule.Threshold,
				Status:    "firing",
			}
			if err := s.db.WithContext(ctx).Create(&alert).Error; err != nil {
				return fmt.Errorf("failed to create alert: %w", err)
			}

			// Send notification
			s.sendAlertNotification(ctx, &alert, rule)
		}
	} else {
		// Check if there's a firing alert that should be resolved
		var existingAlert models.AlertHistory
		err := s.db.WithContext(ctx).
			Where("rule_id = ? AND status = ?", rule.ID, "firing").
			First(&existingAlert).Error

		if err == nil {
			// Resolve the alert
			now := time.Now()
			existingAlert.Status = "resolved"
			existingAlert.ResolvedAt = &now
			if err := s.db.WithContext(ctx).Save(&existingAlert).Error; err != nil {
				return fmt.Errorf("failed to resolve alert: %w", err)
			}
		}
	}

	return nil
}

// sendAlertNotification sends notification for an alert
func (s *AlertingService) sendAlertNotification(ctx context.Context, alert *models.AlertHistory, rule *models.AlertRule) {
	if s.notificationService == nil {
		return
	}

	// Parse notification channels
	var channels []string
	if err := json.Unmarshal([]byte(rule.NotifyChannels), &channels); err != nil {
		channels = []string{"in_app"}
	}

	// Build notification content
	title := fmt.Sprintf("[%s] %s", alert.Severity, alert.RuleName)
	content := alert.Message

	for _, channel := range channels {
		switch channel {
		case "in_app":
			// Create in-app notification
			// This would typically notify admin users
			// For now, we just log it
			fmt.Printf("Alert notification [%s]: %s\n", title, content)
		case "email":
			// TODO: Implement email notification
			fmt.Printf("Email notification would be sent: %s\n", title)
		}
	}
}

// GetAlertHistory retrieves alert history with filters
func (s *AlertingService) GetAlertHistory(ctx context.Context, status string, page, pageSize int) ([]models.AlertHistory, int64, error) {
	var alerts []models.AlertHistory
	var total int64

	query := s.db.WithContext(ctx).Model(&models.AlertHistory{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&alerts).Error

	return alerts, total, err
}

// ResolveAlert marks an alert as resolved
func (s *AlertingService) ResolveAlert(ctx context.Context, alertID string) error {
	var alert models.AlertHistory
	if err := s.db.WithContext(ctx).First(&alert, "id = ?", alertID).Error; err != nil {
		return err
	}

	if alert.Status != "firing" {
		return fmt.Errorf("alert is not in firing state")
	}

	now := time.Now()
	alert.Status = "resolved"
	alert.ResolvedAt = &now

	return s.db.WithContext(ctx).Save(&alert).Error
}

// GetAlertStats returns alert statistics
func (s *AlertingService) GetAlertStats(ctx context.Context) (map[string]interface{}, error) {
	var totalRules int64
	var activeRules int64
	var firingAlerts int64
	var resolvedAlerts int64

	if err := s.db.WithContext(ctx).Model(&models.AlertRule{}).Count(&totalRules).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&models.AlertRule{}).
		Where("is_active = ?", true).Count(&activeRules).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&models.AlertHistory{}).
		Where("status = ?", "firing").Count(&firingAlerts).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&models.AlertHistory{}).
		Where("status = ?", "resolved").Count(&resolvedAlerts).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_rules":     totalRules,
		"active_rules":    activeRules,
		"firing_alerts":   firingAlerts,
		"resolved_alerts": resolvedAlerts,
	}, nil
}

// CleanupOldAlerts removes resolved alerts older than retention period
func (s *AlertingService) CleanupOldAlerts(ctx context.Context, retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	return s.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", "resolved", cutoff).
		Delete(&models.AlertHistory{}).Error
}
