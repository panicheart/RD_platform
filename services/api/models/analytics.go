package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// AnalyticsDashboard stores dashboard configuration
type AnalyticsDashboard struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Layout      string    `json:"layout" gorm:"type:text"` // JSON configuration
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *AnalyticsDashboard) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

func (AnalyticsDashboard) TableName() string {
	return "analytics_dashboards"
}

// ProjectStats stores project statistics
type ProjectStats struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	Date              time.Time `json:"date" gorm:"index"`
	TotalProjects     int       `json:"total_projects"`
	ActiveProjects    int       `json:"active_projects"`
	CompletedProjects int       `json:"completed_projects"`
	DelayedProjects   int       `json:"delayed_projects"`
	AvgProgress       float64   `json:"avg_progress"`
	CreatedAt         time.Time `json:"created_at"`
}

func (ps *ProjectStats) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = ulid.Make().String()
	}
	return nil
}

func (ProjectStats) TableName() string {
	return "project_stats"
}

// UserStats stores user performance statistics
type UserStats struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	UserID         string    `json:"user_id" gorm:"index"`
	Date           time.Time `json:"date" gorm:"index"`
	TasksCompleted int       `json:"tasks_completed"`
	TasksCreated   int       `json:"tasks_created"`
	ProjectsJoined int       `json:"projects_joined"`
	ReviewsDone    int       `json:"reviews_done"`
	Contribution   float64   `json:"contribution"` // Calculated score
	WorkHours      float64   `json:"work_hours"`
	CreatedAt      time.Time `json:"created_at"`
}

func (us *UserStats) BeforeCreate(tx *gorm.DB) error {
	if us.ID == "" {
		us.ID = ulid.Make().String()
	}
	return nil
}

func (UserStats) TableName() string {
	return "user_stats"
}

// ReportTemplate stores report templates
type ReportTemplate struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type"`                      // project, user, system
	Format      string    `json:"format"`                    // pdf, excel
	Template    string    `json:"template" gorm:"type:text"` // Template configuration
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (rt *ReportTemplate) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == "" {
		rt.ID = ulid.Make().String()
	}
	return nil
}

func (ReportTemplate) TableName() string {
	return "report_templates"
}
