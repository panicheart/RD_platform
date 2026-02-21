package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// RequirementType represents the type of requirement
type RequirementType string

const (
	RequirementTypeFunctional    RequirementType = "functional"
	RequirementTypeNonFunctional RequirementType = "non_functional"
	RequirementTypeInterface     RequirementType = "interface"
	RequirementTypeSafety        RequirementType = "safety"
)

// RequirementStatus represents the status of requirement
type RequirementStatus string

const (
	RequirementStatusDraft     RequirementStatus = "draft"
	RequirementStatusReviewed  RequirementStatus = "reviewed"
	RequirementStatusApproved  RequirementStatus = "approved"
	RequirementStatusRejected  RequirementStatus = "rejected"
	RequirementStatusImplemented RequirementStatus = "implemented"
	RequirementStatusVerified  RequirementStatus = "verified"
)

// RequirementPriority represents the priority of requirement
type RequirementPriority int

const (
	RequirementPriorityLow    RequirementPriority = 1
	RequirementPriorityMedium RequirementPriority = 2
	RequirementPriorityHigh   RequirementPriority = 3
	RequirementPriorityCritical RequirementPriority = 4
)

// Requirement represents a project requirement
type Requirement struct {
	ID          string              `json:"id" gorm:"primaryKey;type:char(26)"`
	ProjectID   string              `json:"project_id" gorm:"index;not null;type:char(26)"`
	ParentID    *string             `json:"parent_id" gorm:"index;type:char(26)"`
	Title       string              `json:"title" gorm:"not null;size:200"`
	Description string              `json:"description" gorm:"type:text"`
	Type        RequirementType     `json:"type" gorm:"not null"`
	Priority    RequirementPriority `json:"priority" gorm:"default:2"`
	Status      RequirementStatus   `json:"status" gorm:"not null;default:'draft'"`
	Rationale   string              `json:"rationale" gorm:"type:text"`
	Source      string              `json:"source" gorm:"size:200"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CreatedBy   string              `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Project  *Project       `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Parent   *Requirement   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Requirement  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName returns the table name for the model
func (Requirement) TableName() string {
	return "requirements"
}

// BeforeCreate generates ULID before insert
func (r *Requirement) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ulid.Make().String()
	}
	return nil
}

// ChangeRequestType represents the type of change request
type ChangeRequestType string

const (
	ChangeRequestTypeECR ChangeRequestType = "ecr" // Engineering Change Request
	ChangeRequestTypeECO ChangeRequestType = "eco" // Engineering Change Order
)

// ChangeRequestStatus represents the status of change request
type ChangeRequestStatus string

const (
	ChangeRequestStatusDraft      ChangeRequestStatus = "draft"
	ChangeRequestStatusSubmitted  ChangeRequestStatus = "submitted"
	ChangeRequestStatusEvaluated  ChangeRequestStatus = "evaluated"
	ChangeRequestStatusApproved   ChangeRequestStatus = "approved"
	ChangeRequestStatusRejected   ChangeRequestStatus = "rejected"
	ChangeRequestStatusImplemented ChangeRequestStatus = "implemented"
	ChangeRequestStatusClosed     ChangeRequestStatus = "closed"
)

// ChangeRequest represents a change request (ECR/ECO)
type ChangeRequest struct {
	ID                string              `json:"id" gorm:"primaryKey;type:char(26)"`
	ProjectID         string              `json:"project_id" gorm:"index;not null;type:char(26)"`
	Type              ChangeRequestType   `json:"type" gorm:"not null"`
	Title             string              `json:"title" gorm:"not null;size:200"`
	Description       string              `json:"description" gorm:"type:text"`
	Reason            string              `json:"reason" gorm:"type:text"`
	ImpactAnalysis    string              `json:"impact_analysis" gorm:"type:text"`
	AffectedItems     string              `json:"affected_items" gorm:"type:jsonb"` // Array of requirement IDs, files, etc.
	Status            ChangeRequestStatus `json:"status" gorm:"not null;default:'draft'"`
	RequesterID       string              `json:"requester_id" gorm:"type:char(26)"`
	ApproverID        *string             `json:"approver_id" gorm:"type:char(26)"`
	ApprovedAt        *time.Time          `json:"approved_at"`
	ImplementedAt     *time.Time          `json:"implemented_at"`
	ClosedAt          *time.Time          `json:"closed_at"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`

	// Relations
	Project   *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Requester *User    `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Approver  *User    `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

// TableName returns the table name for the model
func (ChangeRequest) TableName() string {
	return "change_requests"
}

// BeforeCreate generates ULID before insert
func (cr *ChangeRequest) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == "" {
		cr.ID = ulid.Make().String()
	}
	return nil
}

// DefectSeverity represents the severity of a defect
type DefectSeverity string

const (
	DefectSeverityCritical DefectSeverity = "critical"
	DefectSeverityHigh     DefectSeverity = "high"
	DefectSeverityMedium   DefectSeverity = "medium"
	DefectSeverityLow      DefectSeverity = "low"
)

// DefectStatus represents the status of a defect
type DefectStatus string

const (
	DefectStatusNew        DefectStatus = "new"
	DefectStatusAssigned   DefectStatus = "assigned"
	DefectStatusInProgress DefectStatus = "in_progress"
	DefectStatusResolved   DefectStatus = "resolved"
	DefectStatusClosed     DefectStatus = "closed"
	DefectStatusReopened   DefectStatus = "reopened"
)

// Defect represents a project defect
type Defect struct {
	ID          string         `json:"id" gorm:"primaryKey;type:char(26)"`
	ProjectID   string         `json:"project_id" gorm:"index;not null;type:char(26)"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	Description string         `json:"description" gorm:"type:text"`
	Severity    DefectSeverity `json:"severity" gorm:"not null"`
	Status      DefectStatus   `json:"status" gorm:"not null;default:'new'"`
	ReporterID  string         `json:"reporter_id" gorm:"type:char(26)"`
	AssigneeID  *string        `json:"assignee_id" gorm:"type:char(26)"`
	Resolution  string         `json:"resolution" gorm:"type:text"`
	ReportedAt  time.Time      `json:"reported_at"`
	ResolvedAt  *time.Time     `json:"resolved_at"`
	ClosedAt    *time.Time     `json:"closed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relations
	Project  *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Reporter *User    `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Assignee *User    `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
}

// TableName returns the table name for the model
func (Defect) TableName() string {
	return "defects"
}

// BeforeCreate generates ULID before insert
func (d *Defect) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	d.ReportedAt = time.Now()
	return nil
}
