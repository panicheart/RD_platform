package models

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// WorkflowState represents the state of a workflow
type WorkflowState string

const (
	WorkflowStateDraft      WorkflowState = "draft"
	WorkflowStatePlanning   WorkflowState = "planning"
	WorkflowStateExecuting  WorkflowState = "executing"
	WorkflowStateReviewing  WorkflowState = "reviewing"
	WorkflowStateCompleted  WorkflowState = "completed"
	WorkflowStatePaused     WorkflowState = "paused"
	WorkflowStateCancelled  WorkflowState = "cancelled"
)

// ValidWorkflowStates contains all valid workflow states
var ValidWorkflowStates = []WorkflowState{
	WorkflowStateDraft,
	WorkflowStatePlanning,
	WorkflowStateExecuting,
	WorkflowStateReviewing,
	WorkflowStateCompleted,
	WorkflowStatePaused,
	WorkflowStateCancelled,
}

// Workflow represents a project workflow instance
type Workflow struct {
	ID          string        `json:"id" gorm:"primaryKey;type:char(26)"`
	ProjectID   string        `json:"project_id" gorm:"index;not null;type:char(26)"`
	TemplateID  string        `json:"template_id" gorm:"index;type:char(26)"`
	Name        string        `json:"name" gorm:"not null;size:200"`
	Description string        `json:"description" gorm:"type:text"`
	State       WorkflowState `json:"state" gorm:"not null;default:'draft'"`
	StartedAt   *time.Time    `json:"started_at"`
	CompletedAt *time.Time    `json:"completed_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   string        `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Project    *Project     `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Activities []Activity   `json:"activities,omitempty" gorm:"foreignKey:WorkflowID"`
}

// TableName returns the table name for the model
func (Workflow) TableName() string {
	return "workflows"
}

// BeforeCreate generates ULID before insert
func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = ulid.Make().String()
	}
	return nil
}

// IsValidState checks if the current state is valid
func (w *Workflow) IsValidState() bool {
	for _, s := range ValidWorkflowStates {
		if w.State == s {
			return true
		}
	}
	return false
}

// CanTransitionTo checks if the workflow can transition to the target state
func (w *Workflow) CanTransitionTo(target WorkflowState) error {
	if w.State == target {
		return errors.New("workflow is already in the target state")
	}

	switch w.State {
	case WorkflowStateDraft:
		if target != WorkflowStatePlanning && target != WorkflowStateCancelled {
			return errors.New("draft workflow can only transition to planning or cancelled")
		}
	case WorkflowStatePlanning:
		if target != WorkflowStateExecuting && target != WorkflowStatePaused && target != WorkflowStateCancelled {
			return errors.New("planning workflow can only transition to executing, paused, or cancelled")
		}
	case WorkflowStateExecuting:
		if target != WorkflowStateReviewing && target != WorkflowStatePaused && target != WorkflowStateCancelled {
			return errors.New("executing workflow can only transition to reviewing, paused, or cancelled")
		}
	case WorkflowStateReviewing:
		if target != WorkflowStateCompleted && target != WorkflowStateExecuting && target != WorkflowStateCancelled {
			return errors.New("reviewing workflow can only transition to completed, executing, or cancelled")
		}
	case WorkflowStatePaused:
		if target != WorkflowStateExecuting && target != WorkflowStateCancelled {
			return errors.New("paused workflow can only transition to executing or cancelled")
		}
	case WorkflowStateCompleted, WorkflowStateCancelled:
		return errors.New("completed or cancelled workflows cannot transition")
	default:
		return errors.New("invalid current state")
	}

	return nil
}

// TransitionTo transitions the workflow to a new state
func (w *Workflow) TransitionTo(target WorkflowState) error {
	if err := w.CanTransitionTo(target); err != nil {
		return err
	}

	w.State = target
	now := time.Now()

	switch target {
	case WorkflowStateExecuting:
		if w.StartedAt == nil {
			w.StartedAt = &now
		}
	case WorkflowStateCompleted, WorkflowStateCancelled:
		w.CompletedAt = &now
	}

	return nil
}
