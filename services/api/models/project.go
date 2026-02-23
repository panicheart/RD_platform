package models

import (
	"time"
)

// Project represents a research/development project
type Project struct {
	ID          string  `json:"id" gorm:"type:varchar(26);primaryKey"`
	Code        string  `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name        string  `json:"name" gorm:"type:varchar(200);not null"`
	Description *string `json:"description" gorm:"type:text"`
	Category    string  `json:"category" gorm:"type:project_category;not null"`
	Status      string  `json:"status" gorm:"type:project_status;default:'draft'"`
	ProductLine *string `json:"product_line" gorm:"type:pd_product_line"`
	Team        *string `json:"team" gorm:"type:team_type"`

	// Process binding
	ProcessTemplateID *string `json:"process_template_id" gorm:"type:varchar(26)"`

	// Dates
	StartDate       *time.Time `json:"start_date" gorm:"type:date"`
	EndDate         *time.Time `json:"end_date" gorm:"type:date"`
	ActualStartDate *time.Time `json:"actual_start_date" gorm:"type:date"`
	ActualEndDate   *time.Time `json:"actual_end_date" gorm:"type:date"`

	// Progress
	Progress int `json:"progress" gorm:"default:0"`

	// Git repository
	GitRepoID  *string `json:"git_repo_id" gorm:"type:varchar(100)"`
	GitRepoURL *string `json:"git_repo_url" gorm:"type:varchar(500)"`

	// Classification
	ClassificationLevel string `json:"classification_level" gorm:"type:classification_level;default:'internal'"`

	// Team
	LeaderID        *string `json:"leader_id" gorm:"type:varchar(26)"`
	TechLeaderID    *string `json:"tech_leader_id" gorm:"type:varchar(26)"`
	ProductLeaderID *string `json:"product_leader_id" gorm:"type:varchar(26)"`

	// Metadata
	CreatedBy *string   `json:"created_by" gorm:"type:varchar(26)"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Associations
	Members []ProjectMember `json:"members,omitempty" gorm:"foreignKey:ProjectID"`
}

// TableName specifies the table name
func (Project) TableName() string {
	return "projects"
}

// ProjectMember represents a user participating in a project
type ProjectMember struct {
	ID        string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	ProjectID string    `json:"project_id" gorm:"type:varchar(26);not null;index"`
	UserID    string    `json:"user_id" gorm:"type:varchar(26);not null;index"`
	Role      string    `json:"role" gorm:"type:varchar(50);default:'member'"`
	JoinedAt  time.Time `json:"joined_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Association
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name
func (ProjectMember) TableName() string {
	return "project_members"
}

// ProcessTemplate represents a workflow template
type ProcessTemplate struct {
	ID          string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null"`
	Code        string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Category    string    `json:"category" gorm:"type:project_category;not null"`
	Description *string   `json:"description" gorm:"type:text"`
	Activities  string    `json:"activities" gorm:"type:jsonb;not null"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedBy   *string   `json:"created_by" gorm:"type:varchar(26)"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (ProcessTemplate) TableName() string {
	return "process_templates"
}

// ProjectFile represents a file in a project
type ProjectFile struct {
	ID          string    `json:"id" gorm:"type:varchar(26);primaryKey"`
	ProjectID   string    `json:"project_id" gorm:"type:varchar(26);not null;index"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Path        string    `json:"path" gorm:"type:varchar(500);not null"`
	Size        int64     `json:"size" gorm:"default:0"`
	ContentType string    `json:"content_type" gorm:"type:varchar(100)"`
	IsDirectory bool      `json:"is_directory" gorm:"default:false"`
	StoragePath string    `json:"storage_path" gorm:"type:varchar(1000)"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (ProjectFile) TableName() string {
	return "project_files"
}

// File is a generic file representation for API responses
type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	IsDirectory bool      `json:"is_directory"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
