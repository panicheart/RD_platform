package services

import (
	"context"
	"errors"

	"rdp/services/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProcessTemplateService handles process template business logic
type ProcessTemplateService struct {
	db *gorm.DB
}

// NewProcessTemplateService creates a new ProcessTemplateService
func NewProcessTemplateService(db *gorm.DB) *ProcessTemplateService {
	return &ProcessTemplateService{db: db}
}

// ListTemplates returns all process templates
func (s *ProcessTemplateService) ListTemplates(ctx context.Context, category string) ([]models.ProcessTemplate, error) {
	var templates []models.ProcessTemplate

	query := s.db.Model(&models.ProcessTemplate{}).Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("is_default DESC, name ASC").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// GetTemplateByID returns a template by ID
func (s *ProcessTemplateService) GetTemplateByID(ctx context.Context, id string) (*models.ProcessTemplate, error) {
	var template models.ProcessTemplate
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid template ID")
	}

	if err := s.db.First(&template, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}

	return &template, nil
}

// GetTemplateByCode returns a template by code
func (s *ProcessTemplateService) GetTemplateByCode(ctx context.Context, code string) (*models.ProcessTemplate, error) {
	var template models.ProcessTemplate

	if err := s.db.First(&template, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}

	return &template, nil
}

// GetDefaultTemplate returns the default template for a category
func (s *ProcessTemplateService) GetDefaultTemplate(ctx context.Context, category string) (*models.ProcessTemplate, error) {
	var template models.ProcessTemplate

	if err := s.db.First(&template, "category = ? AND is_default = ?", category, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no default template found")
		}
		return nil, err
	}

	return &template, nil
}

// CreateTemplate creates a new template
func (s *ProcessTemplateService) CreateTemplate(ctx context.Context, template *models.ProcessTemplate) error {
	// Check if code exists
	var count int64
	if err := s.db.Model(&models.ProcessTemplate{}).Where("code = ?", template.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("template code already exists")
	}

	template.ID = uuid.New()

	// If this is set as default, unset other defaults
	if template.IsDefault {
		s.db.Model(&models.ProcessTemplate{}).
			Where("category = ? AND is_default = ?", template.Category, true).
			Update("is_default", false)
	}

	return s.db.Create(template).Error
}

// UpdateTemplate updates a template
func (s *ProcessTemplateService) UpdateTemplate(ctx context.Context, id string, updates map[string]interface{}) (*models.ProcessTemplate, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid template ID")
	}

	// If setting as default, unset other defaults
	if isDefault, ok := updates["is_default"].(bool); ok && isDefault {
		template, _ := s.GetTemplateByID(ctx, id)
		if template != nil {
			s.db.Model(&models.ProcessTemplate{}).
				Where("category = ? AND id != ? AND is_default = ?", template.Category, uid, true).
				Update("is_default", false)
		}
	}

	result := s.db.Model(&models.ProcessTemplate{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("template not found")
	}

	return s.GetTemplateByID(ctx, id)
}

// DeleteTemplate soft deletes a template
func (s *ProcessTemplateService) DeleteTemplate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid template ID")
	}

	result := s.db.Model(&models.ProcessTemplate{}).Where("id = ?", uid).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("template not found")
	}

	return nil
}
