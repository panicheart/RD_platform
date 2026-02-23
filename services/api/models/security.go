package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID             string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	UserID         *string   `json:"user_id" gorm:"type:varchar(26);index"`
	Username       *string   `json:"username" gorm:"type:varchar(50)"`
	IPAddress      *string   `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent      *string   `json:"user_agent" gorm:"type:varchar(500)"`
	Action         string    `json:"action" gorm:"type:varchar(100);not null;index"`
	Resource       string    `json:"resource" gorm:"type:varchar(100);not null;index"`
	ResourceID     *string   `json:"resource_id" gorm:"type:varchar(50);index"`
	Method         *string   `json:"method" gorm:"type:varchar(10)"`
	Path           *string   `json:"path" gorm:"type:varchar(500)"`
	RequestBody    *string   `json:"request_body" gorm:"type:text"`
	ResponseCode   *int      `json:"response_code" gorm:"type:integer"`
	ErrorMessage   *string   `json:"error_message" gorm:"type:text"`
	Classification string    `json:"classification" gorm:"type:classification_level;default:'internal'"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// BeforeCreate generates ULID before insert
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// DataClassification represents the classification level of data
type DataClassification struct {
	Level       string `json:"level" gorm:"type:classification_level;primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(50);not null"`
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"type:varchar(20)"`
	Icon        string `json:"icon" gorm:"type:varchar(50)"`
}

// TableName specifies the table name
func (DataClassification) TableName() string {
	return "data_classifications"
}

// LoginLog tracks user login attempts
type LoginLog struct {
	ID            string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	UserID        *string   `json:"user_id" gorm:"type:varchar(26);index"`
	Username      string    `json:"username" gorm:"type:varchar(50);not null"`
	IPAddress     *string   `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent     *string   `json:"user_agent" gorm:"type:varchar(500)"`
	Success       bool      `json:"success" gorm:"default:false"`
	FailureReason *string   `json:"failure_reason" gorm:"type:varchar(100)"`
	Provider      string    `json:"provider" gorm:"type:varchar(50);default:'local'"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// BeforeCreate generates ULID before insert
func (l *LoginLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (LoginLog) TableName() string {
	return "login_logs"
}

// Session represents an active user session
type Session struct {
	ID           string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	UserID       string    `json:"user_id" gorm:"type:varchar(26);not null;index"`
	Token        string    `json:"token" gorm:"type:varchar(500);uniqueIndex;not null"`
	IPAddress    *string   `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent    *string   `json:"user_agent" gorm:"type:varchar(500)"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	LastActiveAt time.Time `json:"last_active_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsRevoked    bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// BeforeCreate generates ULID before insert
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (Session) TableName() string {
	return "sessions"
}
