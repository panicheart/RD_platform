package services

import (
	"context"
	"errors"
	"time"

	"rdp-platform/rdp-api/models"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// ListUsers returns paginated users
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// Apply filters
	if role, ok := filters["role"].(string); ok && role != "" {
		query = query.Where("role = ?", role)
	}
	if team, ok := filters["team"].(string); ok && team != "" {
		query = query.Where("team = ?", team)
	}
	if productLine, ok := filters["product_line"].(string); ok && productLine != "" {
		query = query.Where("product_line = ?", productLine)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("username ILIKE ? OR display_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername returns a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail returns a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Check if username exists
	var count int64
	if err := s.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	// Check if email exists (if provided)
	if user.Email != nil && *user.Email != "" {
		if err := s.db.Model(&models.User{}).Where("email = ?", *user.Email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("email already exists")
		}
	}

	user.ID = ulid.Make().String()

	return s.db.Create(user).Error
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*models.User, error) {
	result := s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return s.GetUserByID(ctx, id)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	result := s.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", &now).Error
}

// GetUserStats returns user statistics
func (s *UserService) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total users
	var total int64
	s.db.Model(&models.User{}).Count(&total)
	stats["total"] = total

	// Active users
	var active int64
	s.db.Model(&models.User{}).Where("is_active = ?", true).Count(&active)
	stats["active"] = active

	// Users by role
	var roleCounts []struct {
		Role  string
		Count int64
	}
	s.db.Model(&models.User{}).Select("role, count(*) as count").Group("role").Scan(&roleCounts)
	stats["by_role"] = roleCounts

	// Users by team
	var teamCounts []struct {
		Team  string
		Count int64
	}
	s.db.Model(&models.User{}).Where("team IS NOT NULL").Select("team, count(*) as count").Group("team").Scan(&teamCounts)
	stats["by_team"] = teamCounts

	return stats, nil
}

// OrganizationService handles organization business logic
type OrganizationService struct {
	db *gorm.DB
}

// NewOrganizationService creates a new OrganizationService
func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// ListOrganizations returns all organizations in tree structure
func (s *OrganizationService) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	if err := s.db.Order("level ASC, sort_order ASC").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetOrganizationByID returns an organization by ID
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.First(&org, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &org, nil
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, org *models.Organization) error {
	// Check if code exists
	var count int64
	if err := s.db.Model(&models.Organization{}).Where("code = ?", org.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("organization code already exists")
	}

	org.ID = ulid.Make().String()

	return s.db.Create(org).Error
}

// UpdateOrganization updates an existing organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id string, updates map[string]interface{}) (*models.Organization, error) {
	result := s.db.Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("organization not found")
	}

	return s.GetOrganizationByID(ctx, id)
}

// DeleteOrganization deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	// Check if has children
	var count int64
	if err := s.db.Model(&models.Organization{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete organization with sub-organizations")
	}

	result := s.db.Delete(&models.Organization{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("organization not found")
	}

	return nil
}
