package services

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	if memberID, ok := filters["member_id"].(string); ok && memberID != "" {
		query = query.Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ?", memberID)
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

// GenerateProjectCode generates a unique project code in format: RDP-{CATEGORY}-{YYYYMMDD}-{SEQ}
// Example: RDP-PD-20240221-001
func (s *ProjectService) GenerateProjectCode(ctx context.Context, category string) (string, error) {
	// Map category to short code
	categoryCode := getCategoryCode(category)
	
	// Get today's date
	today := time.Now().Format("20060102")
	
	// Pattern for today's codes
	prefix := fmt.Sprintf("RDP-%s-%s-", categoryCode, today)
	
	// Get the max sequence number for today
	var maxCode string
	err := s.db.Model(&models.Project{}).
		Where("code LIKE ?", prefix+"%").
		Order("code DESC").
		Pluck("code", &maxCode).Error
	
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	
	// Calculate next sequence
	seq := 1
	if maxCode != "" {
		// Parse sequence from maxCode (format: RDP-XX-YYYYMMDD-NNN)
		var existingSeq int
		if _, err := fmt.Sscanf(maxCode, prefix+"%03d", &existingSeq); err == nil {
			seq = existingSeq + 1
		}
	}
	
	// Generate code
	code := fmt.Sprintf("%s%03d", prefix, seq)
	
	// Verify uniqueness
	var count int64
	if err := s.db.Model(&models.Project{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("generated project code already exists")
	}
	
	return code, nil
}

// getCategoryCode returns the short code for a category
func getCategoryCode(category string) string {
	categoryCodes := map[string]string{
		"pd_project":     "PD",  // 产品开发项目
		"pre_research":   "PR",  // 预研项目
		"tech_research":  "TR",  // 技术攻关
		"platform":       "PL",  // 平台项目
		"customization":  "CU",  // 定制项目
		"improvement":    "IM",  // 改进项目
		"others":         "OT",  // 其他
	}
	
	if code, ok := categoryCodes[category]; ok {
		return code
	}
	return "XX"
}

// CreateProject creates a new project with auto-generated code
func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project, userID string) error {
	// Auto-generate project code if not provided
	if project.Code == "" {
		code, err := s.GenerateProjectCode(ctx, project.Category)
		if err != nil {
			return fmt.Errorf("failed to generate project code: %w", err)
		}
		project.Code = code
	} else {
		// Check if provided code exists
		var count int64
		if err := s.db.Model(&models.Project{}).Where("code = ?", project.Code).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("project code already exists")
		}
	}

	// Set defaults
	if project.Status == "" {
		project.Status = "draft"
	}
	if project.Progress < 0 {
		project.Progress = 0
	}
	if project.Progress > 100 {
		project.Progress = 100
	}
	if project.ClassificationLevel == "" {
		project.ClassificationLevel = "internal"
	}

	// Set created_by
	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err == nil {
			project.CreatedBy = &uid
		}
	}

	project.ID = uuid.New()
	
	// Use transaction to create project and add creator as member
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			return err
		}
		
		// Add creator as project manager if userID provided
		if userID != "" {
			uid, err := uuid.Parse(userID)
			if err == nil {
				member := models.ProjectMember{
					ID:        uuid.New(),
					ProjectID: project.ID,
					UserID:    uid,
					Role:      "manager",
				}
				if err := tx.Create(&member).Error; err != nil {
					return err
				}
			}
		}
		
		return nil
	})
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, id string, updates map[string]interface{}, userID string) (*models.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	// Check if user has permission to update (project manager, leader, or admin)
	if userID != "" {
		hasPermission, err := s.checkProjectPermission(ctx, id, userID, []string{"manager", "leader", "admin"})
		if err != nil {
			return nil, err
		}
		if !hasPermission {
			return nil, errors.New("insufficient permissions to update project")
		}
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "code")
	delete(updates, "created_by")

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
func (s *ProjectService) DeleteProject(ctx context.Context, id string, userID string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid project ID")
	}

	// Check if user has permission to delete (project manager or admin)
	if userID != "" {
		hasPermission, err := s.checkProjectPermission(ctx, id, userID, []string{"manager", "admin"})
		if err != nil {
			return err
		}
		if !hasPermission {
			return errors.New("insufficient permissions to delete project")
		}
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

// checkProjectPermission checks if user has required role in project or is admin
func (s *ProjectService) checkProjectPermission(ctx context.Context, projectID, userID string, allowedRoles []string) (bool, error) {
	// Check if user is admin
	var user models.User
	uid, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}
	
	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
		return false, err
	}
	
	if user.Role == "admin" {
		return true, nil
	}
	
	// Check if user is project member with required role
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return false, err
	}
	
	var member models.ProjectMember
	if err := s.db.First(&member, "project_id = ? AND user_id = ?", projectUID, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	
	for _, role := range allowedRoles {
		if member.Role == role {
			return true, nil
		}
	}
	
	// Check if user is project leader
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectUID).Error; err != nil {
		return false, err
	}
	
	if project.LeaderID != nil && *project.LeaderID == uid {
		return true, nil
	}
	
	return false, nil
}

// AddProjectMember adds a member to a project
func (s *ProjectService) AddProjectMember(ctx context.Context, projectID, userID string, role string, addedBy string) (*models.ProjectMember, error) {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if adder has permission
	if addedBy != "" {
		hasPermission, err := s.checkProjectPermission(ctx, projectID, addedBy, []string{"manager", "admin"})
		if err != nil {
			return nil, err
		}
		if !hasPermission {
			return nil, errors.New("insufficient permissions to add members")
		}
	}

	// Check if already a member
	var existing models.ProjectMember
	if err := s.db.First(&existing, "project_id = ? AND user_id = ?", projectUID, userUID).Error; err == nil {
		return nil, errors.New("user is already a member of this project")
	}

	// Validate role
	validRoles := []string{"member", "developer", "tester", "manager", "leader", "observer"}
	isValidRole := false
	for _, r := range validRoles {
		if role == r {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		role = "member"
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

	// Load user info
	var user models.User
	if err := s.db.First(&user, "id = ?", userUID).Error; err == nil {
		member.User = &user
	}

	return &member, nil
}

// GetProjectMembers returns all members of a project
func (s *ProjectService) GetProjectMembers(ctx context.Context, projectID string) ([]models.ProjectMember, error) {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	var members []models.ProjectMember
	if err := s.db.Where("project_id = ?", projectUID).
		Preload("User").
		Order("role ASC, joined_at ASC").
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// RemoveMember removes a member from a project
func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID string, removedBy string) error {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if remover has permission (can't remove self unless admin)
	if removedBy != "" && removedBy != userID {
		hasPermission, err := s.checkProjectPermission(ctx, projectID, removedBy, []string{"manager", "admin"})
		if err != nil {
			return err
		}
		if !hasPermission {
			return errors.New("insufficient permissions to remove members")
		}
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
func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectID, userID, newRole string, updatedBy string) error {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return errors.New("invalid project ID")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if updater has permission
	if updatedBy != "" {
		hasPermission, err := s.checkProjectPermission(ctx, projectID, updatedBy, []string{"manager", "admin"})
		if err != nil {
			return err
		}
		if !hasPermission {
			return errors.New("insufficient permissions to update member roles")
		}
	}

	// Validate role
	validRoles := []string{"member", "developer", "tester", "manager", "leader", "observer"}
	isValidRole := false
	for _, r := range validRoles {
		if newRole == r {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return errors.New("invalid role")
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

// UpdateProjectProgress updates project progress based on activities completion
func (s *ProjectService) UpdateProjectProgress(ctx context.Context, projectID string, progress int, userID string) (*models.Project, error) {
	// Validate progress
	if progress < 0 || progress > 100 {
		return nil, errors.New("progress must be between 0 and 100")
	}

	// Check permission
	if userID != "" {
		hasPermission, err := s.checkProjectPermission(ctx, projectID, userID, []string{"manager", "leader", "admin"})
		if err != nil {
			return nil, err
		}
		if !hasPermission {
			return nil, errors.New("insufficient permissions to update progress")
		}
	}

	// Get activities to calculate progress if auto-calculate
	activities, err := s.GetProjectActivities(ctx, projectID)
	if err == nil && len(activities) > 0 {
		// Calculate progress from activities
		calculatedProgress := calculateProgressFromActivities(activities)
		// Use the higher of manual input or calculated progress
		if calculatedProgress > progress {
			progress = calculatedProgress
		}
	}

	// Update progress
	updates := map[string]interface{}{
		"progress": progress,
	}

	// Auto-update status based on progress
	if progress == 0 {
		updates["status"] = "draft"
	} else if progress == 100 {
		updates["status"] = "completed"
		now := time.Now()
		updates["actual_end_date"] = &now
	} else if progress > 0 {
		updates["status"] = "in_progress"
		// Set actual start date if not set
		var project models.Project
		uid, _ := uuid.Parse(projectID)
		if err := s.db.First(&project, "id = ?", uid).Error; err == nil {
			if project.ActualStartDate == nil {
				now := time.Now()
				updates["actual_start_date"] = &now
			}
		}
	}

	return s.UpdateProject(ctx, projectID, updates, "")
}

// calculateProgressFromActivities calculates project progress from activities
func calculateProgressFromActivities(activities []models.Activity) int {
	if len(activities) == 0 {
		return 0
	}
	
	totalWeight := 0
	completedWeight := 0
	
	for _, activity := range activities {
		// Each activity has equal weight by default
		weight := 1
		totalWeight += weight
		
		switch activity.Status {
		case "completed":
			completedWeight += weight
		case "in_progress":
			// Partial credit based on activity progress
			completedWeight += weight * activity.Progress / 100
		}
	}
	
	if totalWeight == 0 {
		return 0
	}
	
	return completedWeight * 100 / totalWeight
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

// GetProjectStats returns project statistics
func (s *ProjectService) GetProjectStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Count by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	s.db.Model(&models.Project{}).Select("status, count(*) as count").Group("status").Scan(&statusCounts)
	
	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}
	stats["by_status"] = statusMap
	
	// Count by category
	var categoryCounts []struct {
		Category string
		Count    int64
	}
	s.db.Model(&models.Project{}).Select("category, count(*) as count").Group("category").Scan(&categoryCounts)
	
	categoryMap := make(map[string]int64)
	for _, cc := range categoryCounts {
		categoryMap[cc.Category] = cc.Count
	}
	stats["by_category"] = categoryMap
	
	// User's projects count
	if userID != "" {
		userUID, err := uuid.Parse(userID)
		if err == nil {
			var myProjects int64
			s.db.Model(&models.Project{}).
				Joins("JOIN project_members ON project_members.project_id = projects.id").
				Where("project_members.user_id = ?", userUID).
				Count(&myProjects)
			stats["my_projects"] = myProjects
		}
	}
	
	// Total count
	var total int64
	s.db.Model(&models.Project{}).Count(&total)
	stats["total"] = total
	
	return stats, nil
}
