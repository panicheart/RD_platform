package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ReviewStatus represents the status of a review
type ReviewStatus string

const (
	ReviewStatusPending   ReviewStatus = "pending"
	ReviewStatusSubmitted ReviewStatus = "submitted"
	ReviewStatusApproved  ReviewStatus = "approved"
	ReviewStatusRejected  ReviewStatus = "rejected"
	ReviewStatusRevision  ReviewStatus = "revision"
)

// ReviewType represents the type of review
type ReviewType string

const (
	ReviewTypeDCP   ReviewType = "dcp"
	ReviewTypeCode  ReviewType = "code"
	ReviewTypeDoc   ReviewType = "doc"
	ReviewTypeFinal ReviewType = "final"
)

// Review represents a review record for an activity
type Review struct {
	ID          string       `json:"id" gorm:"primaryKey;type:char(26)"`
	ActivityID  string       `json:"activity_id" gorm:"index;not null;type:char(26)"`
	ProjectID   string       `json:"project_id" gorm:"index;not null;type:char(26)"`
	Type        ReviewType   `json:"type" gorm:"not null;default:'dcp'"`
	Status      ReviewStatus `json:"status" gorm:"not null;default:'pending'"`
	ReviewerID  *string      `json:"reviewer_id" gorm:"type:char(26)"`
	Comments    string       `json:"comments" gorm:"type:text"`
	Score       *int         `json:"score"`
	SubmittedAt *time.Time   `json:"submitted_at"`
	ReviewedAt  *time.Time   `json:"reviewed_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CreatedBy   string       `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Activity *Activity `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	Reviewer *User     `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	Feedbacks []Feedback `json:"feedbacks,omitempty" gorm:"foreignKey:ReviewID"`
}

// TableName returns the table name for the model
func (Review) TableName() string {
	return "reviews"
}

// BeforeCreate generates ULID before insert
func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ulid.Make().String()
	}
	return nil
}

// IsPending checks if the review is pending
func (r *Review) IsPending() bool {
	return r.Status == ReviewStatusPending
}

// IsComplete checks if the review is complete
func (r *Review) IsComplete() bool {
	return r.Status == ReviewStatusApproved || r.Status == ReviewStatusRejected
}

// Submit marks the review as submitted
func (r *Review) Submit() {
	now := time.Now()
	r.Status = ReviewStatusSubmitted
	r.SubmittedAt = &now
}

// Approve marks the review as approved
func (r *Review) Approve() {
	now := time.Now()
	r.Status = ReviewStatusApproved
	r.ReviewedAt = &now
}

// Reject marks the review as rejected
func (r *Review) Reject() {
	now := time.Now()
	r.Status = ReviewStatusRejected
	r.ReviewedAt = &now
}

// RequestRevision marks the review as needing revision
func (r *Review) RequestRevision() {
	now := time.Now()
	r.Status = ReviewStatusRevision
	r.ReviewedAt = &now
}

// Feedback represents feedback on a review
type Feedback struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(26)"`
	ReviewID  string    `json:"review_id" gorm:"index;not null;type:char(26)"`
	ParentID  *string   `json:"parent_id" gorm:"index;type:char(26)"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	AuthorID  string    `json:"author_id" gorm:"not null;type:char(26)"`
	Mentions  string    `json:"mentions" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Review *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	Author *User   `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Parent *Feedback `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName returns the table name for the model
func (Feedback) TableName() string {
	return "feedbacks"
}

// BeforeCreate generates ULID before insert
func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = ulid.Make().String()
	}
	return nil
}
