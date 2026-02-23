package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// SystemMetric stores system monitoring metrics
type SystemMetric struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Timestamp     time.Time `json:"timestamp" gorm:"index"`
	CPUUsage      float64   `json:"cpu_usage"`      // Percentage
	MemoryUsage   float64   `json:"memory_usage"`   // Percentage
	MemoryTotal   int64     `json:"memory_total"`   // Bytes
	MemoryUsed    int64     `json:"memory_used"`    // Bytes
	DiskUsage     float64   `json:"disk_usage"`     // Percentage
	DiskTotal     int64     `json:"disk_total"`     // Bytes
	DiskUsed      int64     `json:"disk_used"`      // Bytes
	NetworkIn     int64     `json:"network_in"`     // Bytes
	NetworkOut    int64     `json:"network_out"`    // Bytes
	DBConnections int       `json:"db_connections"` // Active connections
	APIRequests   int64     `json:"api_requests"`   // Request count
	CreatedAt     time.Time `json:"created_at"`
}

func (sm *SystemMetric) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == "" {
		sm.ID = ulid.Make().String()
	}
	return nil
}

func (SystemMetric) TableName() string {
	return "system_metrics"
}

// APIMetric stores API performance metrics
type APIMetric struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Timestamp  time.Time `json:"timestamp" gorm:"index"`
	Endpoint   string    `json:"endpoint" gorm:"index"`
	Method     string    `json:"method"`
	Duration   int64     `json:"duration"` // Milliseconds
	StatusCode int       `json:"status_code"`
	UserID     string    `json:"user_id"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

func (am *APIMetric) BeforeCreate(tx *gorm.DB) error {
	if am.ID == "" {
		am.ID = ulid.Make().String()
	}
	return nil
}

func (APIMetric) TableName() string {
	return "api_metrics"
}

// LogEntry stores application log entries
type LogEntry struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	Level     string    `json:"level" gorm:"index"` // DEBUG, INFO, WARN, ERROR
	Message   string    `json:"message" gorm:"type:text"`
	Source    string    `json:"source"` // Service name
	Module    string    `json:"module"` // Module/component
	UserID    string    `json:"user_id"`
	RequestID string    `json:"request_id"`
	Metadata  string    `json:"metadata" gorm:"type:text"` // JSON additional data
	CreatedAt time.Time `json:"created_at"`
}

func (le *LogEntry) BeforeCreate(tx *gorm.DB) error {
	if le.ID == "" {
		le.ID = ulid.Make().String()
	}
	return nil
}

func (LogEntry) TableName() string {
	return "log_entries"
}

// AlertRule defines alert rules
type AlertRule struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	Metric         string    `json:"metric"`    // Metric to monitor
	Condition      string    `json:"condition"` // >, <, ==, !=
	Threshold      float64   `json:"threshold"` // Trigger value
	Duration       int       `json:"duration"`  // Duration in minutes
	Severity       string    `json:"severity"`  // warning, critical
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	NotifyChannels string    `json:"notify_channels" gorm:"type:text"` // JSON array
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ar *AlertRule) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == "" {
		ar.ID = ulid.Make().String()
	}
	return nil
}

func (AlertRule) TableName() string {
	return "alert_rules"
}

// AlertHistory stores triggered alerts
type AlertHistory struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	RuleID     string     `json:"rule_id" gorm:"index"`
	RuleName   string     `json:"rule_name"`
	Severity   string     `json:"severity"`
	Message    string     `json:"message"`
	Value      float64    `json:"value"`     // Actual value that triggered
	Threshold  float64    `json:"threshold"` // Threshold value
	Status     string     `json:"status"`    // firing, resolved
	ResolvedAt *time.Time `json:"resolved_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (ah *AlertHistory) BeforeCreate(tx *gorm.DB) error {
	if ah.ID == "" {
		ah.ID = ulid.Make().String()
	}
	return nil
}

func (AlertHistory) TableName() string {
	return "alert_history"
}
