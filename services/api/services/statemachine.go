package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"rdp/services/api/models"
)

// StateMachineService handles workflow state machine logic
type StateMachineService struct {
	db *gorm.DB
}

// NewStateMachineService creates a new state machine service
func NewStateMachineService(db *gorm.DB) *StateMachineService {
	return &StateMachineService{db: db}
}

// TransitionWorkflow transitions a workflow to a new state
func (s *StateMachineService) TransitionWorkflow(workflowID string, targetState models.WorkflowState, userID string) error {
	var workflow models.Workflow
	if err := s.db.First(&workflow, "id = ?", workflowID).Error; err != nil {
		return err
	}

	if err := workflow.TransitionTo(targetState); err != nil {
		return err
	}

	return s.db.Save(&workflow).Error
}

// GetWorkflow retrieves a workflow by ID
func (s *StateMachineService) GetWorkflow(workflowID string) (*models.Workflow, error) {
	var workflow models.Workflow
	if err := s.db.Preload("Activities").First(&workflow, "id = ?", workflowID).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

// ListWorkflows retrieves workflows with pagination
func (s *StateMachineService) ListWorkflows(projectID string, page, pageSize int) ([]models.Workflow, int64, error) {
	var workflows []models.Workflow
	var total int64

	query := s.db.Model(&models.Workflow{})
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Activities").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&workflows).Error; err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

// CreateWorkflow creates a new workflow from a template
func (s *StateMachineService) CreateWorkflow(projectID, templateID, name, description, userID string) (*models.Workflow, error) {
	workflow := &models.Workflow{
		ProjectID:   projectID,
		TemplateID:  templateID,
		Name:        name,
		Description: description,
		State:       models.WorkflowStateDraft,
		CreatedBy:   userID,
	}

	if err := s.db.Create(workflow).Error; err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetAvailableTransitions returns valid state transitions for a workflow
func (s *StateMachineService) GetAvailableTransitions(state models.WorkflowState) []models.WorkflowState {
	switch state {
	case models.WorkflowStateDraft:
		return []models.WorkflowState{models.WorkflowStatePlanning, models.WorkflowStateCancelled}
	case models.WorkflowStatePlanning:
		return []models.WorkflowState{models.WorkflowStateExecuting, models.WorkflowStatePaused, models.WorkflowStateCancelled}
	case models.WorkflowStateExecuting:
		return []models.WorkflowState{models.WorkflowStateReviewing, models.WorkflowStatePaused, models.WorkflowStateCancelled}
	case models.WorkflowStateReviewing:
		return []models.WorkflowState{models.WorkflowStateCompleted, models.WorkflowStateExecuting, models.WorkflowStateCancelled}
	case models.WorkflowStatePaused:
		return []models.WorkflowState{models.WorkflowStateExecuting, models.WorkflowStateCancelled}
	default:
		return []models.WorkflowState{}
	}
}

// ValidateTransition checks if a state transition is valid
func (s *StateMachineService) ValidateTransition(from, to models.WorkflowState) error {
	available := s.GetAvailableTransitions(from)
	for _, state := range available {
		if state == to {
			return nil
		}
	}
	return fmt.Errorf("invalid transition from %s to %s", from, to)
}

// ActivityService handles activity business logic
type ActivityService struct {
	db *gorm.DB
}

// NewActivityService creates a new activity service
func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

// CreateActivity creates a new activity
func (s *ActivityService) CreateActivity(workflowID, projectID, name, description string, activityType models.ActivityType, userID string) (*models.Activity, error) {
	activity := &models.Activity{
		WorkflowID:  workflowID,
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		Type:        activityType,
		Status:      models.ActivityStatusPending,
		CreatedBy:   userID,
	}

	if err := s.db.Create(activity).Error; err != nil {
		return nil, err
	}

	return activity, nil
}

// GetActivity retrieves an activity by ID
func (s *ActivityService) GetActivity(activityID string) (*models.Activity, error) {
	var activity models.Activity
	if err := s.db.Preload("Assignee").Preload("Deliverables").Preload("Reviews").First(&activity, "id = ?", activityID).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

// StartActivity starts an activity
func (s *ActivityService) StartActivity(activityID, userID string) error {
	activity, err := s.GetActivity(activityID)
	if err != nil {
		return err
	}

	if !activity.CanStart() {
		return errors.New("activity cannot be started")
	}

	activity.Start()
	return s.db.Save(activity).Error
}

// CompleteActivity completes an activity
func (s *ActivityService) CompleteActivity(activityID, userID string) error {
	activity, err := s.GetActivity(activityID)
	if err != nil {
		return err
	}

	if activity.Status != models.ActivityStatusRunning {
		return errors.New("activity must be running to complete")
	}

	activity.Complete()
	return s.db.Save(activity).Error
}

// AssignActivity assigns an activity to a user
func (s *ActivityService) AssignActivity(activityID, assigneeID string) error {
	return s.db.Model(&models.Activity{}).Where("id = ?", activityID).Update("assignee_id", assigneeID).Error
}

// UpdateActivityProgress updates activity progress
func (s *ActivityService) UpdateActivityProgress(activityID string, progress int) error {
	if progress < 0 || progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}
	return s.db.Model(&models.Activity{}).Where("id = ?", activityID).Update("progress", progress).Error
}

// ListActivities retrieves activities for a workflow
func (s *ActivityService) ListActivities(workflowID string, status models.ActivityStatus, page, pageSize int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := s.db.Model(&models.Activity{}).Where("workflow_id = ?", workflowID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Assignee").Order("sequence ASC, created_at ASC").Limit(pageSize).Offset(offset).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// AddDependency adds a dependency between activities
func (s *ActivityService) AddDependency(activityID, dependsOnID, depType string) error {
	dependency := &models.Dependency{
		ActivityID:     activityID,
		DependsOnID:    dependsOnID,
		DependencyType: depType,
	}
	return s.db.Create(dependency).Error
}

// CheckDependencies checks if all dependencies are completed
func (s *ActivityService) CheckDependencies(activityID string) (bool, error) {
	var dependencies []models.Dependency
	if err := s.db.Where("activity_id = ?", activityID).Find(&dependencies).Error; err != nil {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	for _, dep := range dependencies {
		var activity models.Activity
		if err := s.db.First(&activity, "id = ?", dep.DependsOnID).Error; err != nil {
			return false, err
		}
		if !activity.IsCompleted() {
			return false, nil
		}
	}

	return true, nil
}

// ReviewService handles review business logic
type ReviewService struct {
	db *gorm.DB
}

// NewReviewService creates a new review service
func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(activityID, projectID string, reviewType models.ReviewType, userID string) (*models.Review, error) {
	review := &models.Review{
		ActivityID: activityID,
		ProjectID:  projectID,
		Type:       reviewType,
		Status:     models.ReviewStatusPending,
		CreatedBy:  userID,
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

// SubmitReview submits a review for approval
func (s *ReviewService) SubmitReview(reviewID, comments string, score int) error {
	review, err := s.GetReview(reviewID)
	if err != nil {
		return err
	}

	review.Comments = comments
	review.Score = &score
	review.Submit()

	return s.db.Save(review).Error
}

// ApproveReview approves a review
func (s *ReviewService) ApproveReview(reviewID string) error {
	review, err := s.GetReview(reviewID)
	if err != nil {
		return err
	}

	if review.Status != models.ReviewStatusSubmitted {
		return errors.New("review must be submitted before approval")
	}

	review.Approve()
	return s.db.Save(review).Error
}

// RejectReview rejects a review
func (s *ReviewService) RejectReview(reviewID string) error {
	review, err := s.GetReview(reviewID)
	if err != nil {
		return err
	}

	if review.Status != models.ReviewStatusSubmitted {
		return errors.New("review must be submitted before rejection")
	}

	review.Reject()
	return s.db.Save(review).Error
}

// RequestRevision requests revision for a review
func (s *ReviewService) RequestRevision(reviewID string) error {
	review, err := s.GetReview(reviewID)
	if err != nil {
		return err
	}

	if review.Status != models.ReviewStatusSubmitted {
		return errors.New("review must be submitted before requesting revision")
	}

	review.RequestRevision()
	return s.db.Save(review).Error
}

// GetReview retrieves a review by ID
func (s *ReviewService) GetReview(reviewID string) (*models.Review, error) {
	var review models.Review
	if err := s.db.Preload("Activity").Preload("Reviewer").First(&review, "id = ?", reviewID).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// ListReviews retrieves reviews for an activity
func (s *ReviewService) ListReviews(activityID string, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{}).Where("activity_id = ?", activityID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Reviewer").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
