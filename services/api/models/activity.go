package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ActivityStatus represents the status of an activity
type ActivityStatus string

const (
	ActivityStatusPending    ActivityStatus = "pending"
	ActivityStatusReady      ActivityStatus = "ready"
	ActivityStatusRunning    ActivityStatus = "running"
	ActivityStatusCompleted  ActivityStatus = "completed"
	ActivityStatusReviewing  ActivityStatus = "reviewing"
	ActivityStatusApproved   ActivityStatus = "approved"
	ActivityStatusRejected   ActivityStatus = "rejected"
	ActivityStatusSkipped    ActivityStatus = "skipped"
	ActivityStatusBlocked    ActivityStatus = "blocked"
)

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityTypeTask      ActivityType = "task"
	ActivityTypeMilestone ActivityType = "milestone"
	ActivityTypeDCP       ActivityType = "dcp"
	ActivityTypeReview    ActivityType = "review"
	ActivityTypeApproval  ActivityType = "approval"
)

// Activity represents a workflow activity
type Activity struct {
	ID           string         `json:"id" gorm:"primaryKey;type:char(26)"`
	WorkflowID   string         `json:"workflow_id" gorm:"index;not null;type:char(26)"`
	ProjectID    string         `json:"project_id" gorm:"index;not null;type:char(26)"`
	ParentID     *string        `json:"parent_id" gorm:"index;type:char(26)"`
	Name         string         `json:"name" gorm:"not null;size:200"`
	Description  string         `json:"description" gorm:"type:text"`
	Type         ActivityType   `json:"type" gorm:"not null;default:'task'"`
	Status       ActivityStatus `json:"status" gorm:"not null;default:'pending'"`
	Sequence     int            `json:"sequence" gorm:"default:0"`
	Priority     int            `json:"priority" gorm:"default:0"`
	AssigneeID   *string        `json:"assignee_id" gorm:"type:char(26)"`
	PlannedStart *time.Time     `json:"planned_start"`
	PlannedEnd   *time.Time     `json:"planned_end"`
	ActualStart  *time.Time     `json:"actual_start"`
	ActualEnd    *time.Time     `json:"actual_end"`
	Progress     int            `json:"progress" gorm:"default:0"`
	Deliverables []Deliverable  `json:"deliverables,omitempty" gorm:"foreignKey:ActivityID"`
	Dependencies []Dependency   `json:"dependencies,omitempty" gorm:"foreignKey:ActivityID"`
	Reviews      []Review       `json:"reviews,omitempty" gorm:"foreignKey:ActivityID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	CreatedBy    string         `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Workflow *Workflow `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	Assignee *User     `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
}

// TableName returns the table name for the model
func (Activity) TableName() string {
	return "activities"
}

// BeforeCreate generates ULID before insert
func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = ulid.Make().String()
	}
	return nil
}

// IsReady checks if the activity is ready to start
func (a *Activity) IsReady() bool {
	return a.Status == ActivityStatusReady || a.Status == ActivityStatusPending
}

// IsCompleted checks if the activity is completed
func (a *Activity) IsCompleted() bool {
	return a.Status == ActivityStatusCompleted || a.Status == ActivityStatusApproved
}

// CanStart checks if the activity can be started
func (a *Activity) CanStart() bool {
	return a.Status == ActivityStatusReady || a.Status == ActivityStatusPending
}

// Start marks the activity as started
func (a *Activity) Start() {
	now := time.Now()
	a.Status = ActivityStatusRunning
	a.ActualStart = &now
}

// Complete marks the activity as completed
func (a *Activity) Complete() {
	now := time.Now()
	a.Status = ActivityStatusCompleted
	a.ActualEnd = &now
	a.Progress = 100
}

// Deliverable represents a deliverable item for an activity
type Deliverable struct {
	ID          string    `json:"id" gorm:"primaryKey;type:char(26)"`
	ActivityID  string    `json:"activity_id" gorm:"index;not null;type:char(26)"`
	Name        string    `json:"name" gorm:"not null;size:200"`
	Description string    `json:"description" gorm:"type:text"`
	Type        string    `json:"type" gorm:"size:50"`
	FilePath    string    `json:"file_path" gorm:"size:500"`
	Status      string    `json:"status" gorm:"default:'pending';size:50"`
	SubmittedAt *time.Time `json:"submitted_at"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for the model
func (Deliverable) TableName() string {
	return "deliverables"
}

// BeforeCreate generates ULID before insert
func (d *Deliverable) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

// Dependency represents an activity dependency
type Dependency struct {
	ID             string `json:"id" gorm:"primaryKey;type:char(26)"`
	ActivityID     string `json:"activity_id" gorm:"index;not null;type:char(26)"`
	DependsOnID    string `json:"depends_on_id" gorm:"index;not null;type:char(26)"`
	DependencyType string `json:"dependency_type" gorm:"default:'finish_to_start';size:50"`
}

// TableName returns the table name for the model
func (Dependency) TableName() string {
	return "activity_dependencies"
}

// BeforeCreate generates ULID before insert
func (d *Dependency) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}
