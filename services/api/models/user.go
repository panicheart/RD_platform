package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID             string     `json:"id" gorm:"type:varchar(26);primaryKey"`
	Username       string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	DisplayName    string     `json:"display_name" gorm:"type:varchar(100);not null"`
	Email          *string    `json:"email" gorm:"type:varchar(100)"`
	Phone          *string    `json:"phone" gorm:"type:varchar(20)"`
	AvatarURL      *string    `json:"avatar_url" gorm:"type:varchar(500)"`
	Role           string     `json:"role" gorm:"type:user_role;default:'designer'"`
	Team           *string    `json:"team" gorm:"type:team_type"`
	Specialty      *string    `json:"specialty" gorm:"type:varchar(50)"`
	ProductLine    *string    `json:"product_line" gorm:"type:pd_product_line"`
	Title          *string    `json:"title" gorm:"type:title_level"`
	OrganizationID *string    `json:"organization_id" gorm:"type:varchar(26)"`
	Skills         []string   `json:"skills" gorm:"type:jsonb;serializer:json"`
	Honors         []string   `json:"honors" gorm:"type:jsonb;serializer:json"`
	Bio            *string    `json:"bio" gorm:"type:text"`
	PasswordHash   *string    `json:"-" gorm:"type:varchar(255)"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	CasdoorID      *string    `json:"casdoor_id" gorm:"type:varchar(100)"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// Organization represents an organization unit
type Organization struct {
	ID          string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null"`
	Code        string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	ParentID    *string   `json:"parent_id" gorm:"type:varchar(26)"`
	Level       int       `json:"level" gorm:"default:1"`
	Description *string   `json:"description" gorm:"type:text"`
	LeaderID    *string   `json:"leader_id" gorm:"type:varchar(26)"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (Organization) TableName() string {
	return "organizations"
}

// Notification represents a user notification
type Notification struct {
	ID          string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	UserID      string    `json:"user_id" gorm:"type:varchar(26);not null;index"`
	Type        string    `json:"type" gorm:"type:varchar(50);not null"`
	Title       string    `json:"title" gorm:"type:varchar(200);not null"`
	Content     *string   `json:"content" gorm:"type:text"`
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	RelatedID   *string   `json:"related_id" gorm:"type:varchar(50)"`
	RelatedType *string   `json:"related_type" gorm:"type:varchar(50)"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}

// Announcement represents a system announcement
type Announcement struct {
	ID          string     `json:"id" gorm:"type:varchar(26);primaryKey"`
	Title       string     `json:"title" gorm:"type:varchar(200);not null"`
	Content     string     `json:"content" gorm:"type:text;not null"`
	AuthorID    *string    `json:"author_id" gorm:"type:varchar(26)"`
	Priority    string     `json:"priority" gorm:"default:'normal'"`
	IsPinned    bool       `json:"is_pinned" gorm:"default:false"`
	PublishedAt *time.Time `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (Announcement) TableName() string {
	return "announcements"
}

// Honor represents a team honor/award
type Honor struct {
	ID          string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(200);not null"`
	Description *string   `json:"description" gorm:"type:text"`
	AwardYear   *int      `json:"award_year"`
	AwardMonth  *int      `json:"award_month"`
	RecipientID *string   `json:"recipient_id" gorm:"type:varchar(26)"`
	ImageURL    *string   `json:"image_url" gorm:"type:varchar(500)"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (Honor) TableName() string {
	return "honors"
}
