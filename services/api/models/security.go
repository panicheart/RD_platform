package models

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	Username      *string    `json:"username" gorm:"type:varchar(50)"`
	IPAddress     *string    `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent     *string    `json:"user_agent" gorm:"type:varchar(500)"`
	Action        string    `json:"action" gorm:"type:varchar(100);not null;index"`
	Resource      string    `json:"resource" gorm:"type:varchar(100);not null;index"`
	ResourceID    *string   `json:"resource_id" gorm:"type:varchar(50);index"`
	Method        *string   `json:"method" gorm:"type:varchar(10)"`
	Path          *string   `json:"path" gorm:"type:varchar(500)"`
	RequestBody   *string   `json:"request_body" gorm:"type:text"`
	ResponseCode  *int      `json:"response_code" gorm:"type:integer"`
	ErrorMessage  *string   `json:"error_message" gorm:"type:text"`
	Classification string   `json:"classification" gorm:"type:classification_level;default:'internal'"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
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
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	Username   string    `json:"username" gorm:"type:varchar(50);not null"`
	IPAddress  *string   `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent  *string   `json:"user_agent" gorm:"type:varchar(500)"`
	Success    bool      `json:"success" gorm:"default:false"`
	FailureReason *string `json:"failure_reason" gorm:"type:varchar(100)"`
	Provider   string    `json:"provider" gorm:"type:varchar(50);default:'local'"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (LoginLog) TableName() string {
	return "login_logs"
}

// Session represents an active user session
type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Token        string     `json:"token" gorm:"type:varchar(500);uniqueIndex;not null"`
	IPAddress    *string    `json:"ip_address" gorm:"type:varchar(50)"`
	UserAgent    *string    `json:"user_agent" gorm:"type:varchar(500)"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"not null"`
	LastActiveAt time.Time  `json:"last_active_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsRevoked    bool       `json:"is_revoked" gorm:"default:false"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (Session) TableName() string {
	return "sessions"
}
