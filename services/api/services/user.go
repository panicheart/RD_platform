package services

import (
	"context"
	"errors"

	"rdp/services/api/models"

	"github.com/google/uuid"
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
	if orgID, ok := filters["organization_id"].(string); ok && orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("username ILIKE ? OR display_name ILIKE ?", "%"+search+"%", "%"+search+"%")
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
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

// GetUserByCasdoorID returns a user by Casdoor ID
func (s *UserService) GetUserByCasdoorID(ctx context.Context, casdoorID string) (*models.User, error) {
	var user models.User

	if err := s.db.First(&user, "casdoor_id = ?", casdoorID).Error; err != nil {
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

	user.ID = uuid.New()
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return s.db.Create(user).Error
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*models.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	result := s.db.Model(&models.User{}).Where("id = ?", uid).Updates(updates)
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	result := s.db.Model(&models.User{}).Where("id = ?", uid).Update("is_active", false)
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.db.Model(&models.User{}).Where("id = ?", uid).Update("last_login_at", gorm.Expr("NOW()")).Error
}

// OrganizationService handles organization business logic
type OrganizationService struct {
	db *gorm.DB
}

// NewOrganizationService creates a new OrganizationService
func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// ListOrganizations returns all organizations
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid organization ID")
	}

	if err := s.db.First(&org, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	return &org, nil
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, org *models.Organization) error {
	org.ID = uuid.New()
	return s.db.Create(org).Error
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id string, updates map[string]interface{}) (*models.Organization, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid organization ID")
	}

	result := s.db.Model(&models.Organization{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("organization not found")
	}

	return s.GetOrganizationByID(ctx, id)
}

// DeleteOrganization soft deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid organization ID")
	}

	result := s.db.Model(&models.Organization{}).Where("id = ?", uid).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("organization not found")
	}

	return nil
}
