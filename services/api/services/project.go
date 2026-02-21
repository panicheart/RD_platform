package services

import (
	"context"
	"errors"

	"rdp/services/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectService handles project business logic
type ProjectService struct {
	db *gorm.DB
}

// NewProjectService creates a new ProjectService
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// ListProjects returns paginated projects
func (s *ProjectService) ListProjects(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := s.db.Model(&models.Project{})

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if team, ok := filters["team"].(string); ok && team != "" {
		query = query.Where("team = ?", team)
	}
	if productLine, ok := filters["product_line"].(string); ok && productLine != "" {
		query = query.Where("product_line = ?", productLine)
	}
	if leaderID, ok := filters["leader_id"].(string); ok && leaderID != "" {
		query = query.Where("leader_id = ?", leaderID)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// GetProjectByID returns a project by ID
func (s *ProjectService) GetProjectByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	if err := s.db.Preload("Members").Preload("Members.User").First(&project, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return &project, nil
}

// GetProjectByCode returns a project by code
func (s *ProjectService) GetProjectByCode(ctx context.Context, code string) (*models.Project, error) {
	var project models.Project

	if err := s.db.Preload("Members").First(&project, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return &project, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project) error {
	// Check if code exists
	var count int64
	if err := s.db.Model(&models.Project{}).Where("code = ?", project.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("project code already exists")
	}

	project.ID = uuid.New()
	return s.db.Create(project).Error
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, id string, updates map[string]interface{}) (*models.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	result := s.db.Model(&models.Project{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("project not found")
	}

	return s.GetProjectByID(ctx, id)
}

// DeleteProject soft deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid project ID")
	}

	// Set status to deleted instead of actually deleting
	result := s.db.Model(&models.Project{}).Where("id = ?", uid).Update("status", "deleted")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// AddMember adds a member to a project
func (s *ProjectService) AddMember(ctx context.Context, projectID, userID string, role string) (*models.ProjectMember, error) {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if already a member
	var existing models.ProjectMember
	if err := s.db.First(&existing, "project_id = ? AND user_id = ?", projectUID, userUID).Error; err == nil {
		return nil, errors.New("user is already a member of this project")
	}

	member := models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectUID,
		UserID:    userUID,
		Role:      role,
	}

	if err := s.db.Create(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// RemoveMember removes a member from a project
func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID string) error {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	result := s.db.Where("project_id = ? AND user_id = ?", projectUID, userUID).Delete(&models.ProjectMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}

	return nil
}

// UpdateMemberRole updates a member's role in a project
func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectID, userID, newRole string) error {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	result := s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectUID, userUID).
		Update("role", newRole)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}

	return nil
}

// GetUserProjects returns all projects a user is a member of
func (s *ProjectService) GetUserProjects(ctx context.Context, userID string) ([]models.Project, error) {
	var projects []models.Project
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if err := s.db.Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ?", userUID).
		Order("projects.created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProjectActivities returns all activities for a project
func (s *ProjectService) GetProjectActivities(ctx context.Context, projectID string) ([]models.Activity, error) {
	var activities []models.Activity
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	if err := s.db.Where("project_id = ?", projectUID).Order("sort_order ASC").Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

// CreateActivity creates a new activity for a project
func (s *ProjectService) CreateActivity(ctx context.Context, activity *models.Activity) error {
	activity.ID = uuid.New()
	return s.db.Create(activity).Error
}

// UpdateActivity updates an activity
func (s *ProjectService) UpdateActivity(ctx context.Context, id string, updates map[string]interface{}) (*models.Activity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid activity ID")
	}

	result := s.db.Model(&models.Activity{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("activity not found")
	}

	var activity models.Activity
	if err := s.db.First(&activity, "id = ?", uid).Error; err != nil {
		return nil, err
	}

	return &activity, nil
}
