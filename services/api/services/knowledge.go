package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"rdp-platform/rdp-api/models"
)

// KnowledgeService provides business logic for knowledge base
type KnowledgeService struct {
	db *gorm.DB
}

// NewKnowledgeService creates a new knowledge service
func NewKnowledgeService(db *gorm.DB) *KnowledgeService {
	return &KnowledgeService{db: db}
}

// ==================== Category Tree Operations ====================

// GetCategoryTree returns the full 3-level category tree
func (s *KnowledgeService) GetCategoryTree() ([]*CategoryNode, error) {
	var categories []models.Category
	if err := s.db.Order("sort_order, created_at").Find(&categories).Error; err != nil {
		return nil, err
	}

	// Build tree structure
	nodeMap := make(map[string]*CategoryNode)
	var roots []*CategoryNode

	// First pass: create all nodes
	for i := range categories {
		cat := &categories[i]
		node := &CategoryNode{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Level:       cat.Level,
			SortOrder:   cat.SortOrder,
			Children:    []*CategoryNode{},
		}
		nodeMap[cat.ID] = node
	}

	// Second pass: build parent-child relationships
	for i := range categories {
		cat := &categories[i]
		node := nodeMap[cat.ID]
		
		if cat.ParentID == nil {
			// Root level
			roots = append(roots, node)
		} else if parent, exists := nodeMap[*cat.ParentID]; exists {
			// Add to parent
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// CategoryNode represents a category in tree structure
type CategoryNode struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Level       int            `json:"level"`
	SortOrder   int            `json:"sort_order"`
	Children    []*CategoryNode `json:"children"`
}

// CreateCategory creates a new category with level validation
func (s *KnowledgeService) CreateCategory(name, description string, parentID *string) (*models.Category, error) {
	category := &models.Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}

	// Determine level based on parent
	if parentID != nil {
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", *parentID).Error; err != nil {
			return nil, errors.New("parent category not found")
		}
		if parent.Level >= 3 {
			return nil, errors.New("maximum depth (3 levels) reached")
		}
		category.Level = parent.Level + 1
	} else {
		category.Level = 1
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates a category
func (s *KnowledgeService) UpdateCategory(id, name, description string) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = name
	category.Description = description

	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// DeleteCategory deletes a category and all its children
func (s *KnowledgeService) DeleteCategory(id string) error {
	// Check if category has children
	var childCount int64
	if err := s.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("cannot delete category with children")
	}

	// Check if category has knowledge items
	var knowledgeCount int64
	if err := s.db.Model(&models.Knowledge{}).Where("category_id = ?", id).Count(&knowledgeCount).Error; err != nil {
		return err
	}
	if knowledgeCount > 0 {
		return errors.New("cannot delete category with associated knowledge items")
	}

	return s.db.Delete(&models.Category{}, "id = ?", id).Error
}

// MoveCategory moves a category to a new parent
func (s *KnowledgeService) MoveCategory(id string, newParentID *string) error {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return errors.New("category not found")
	}

	// Calculate new level
	newLevel := 1
	if newParentID != nil {
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", *newParentID).Error; err != nil {
			return errors.New("parent category not found")
		}
		newLevel = parent.Level + 1
		if newLevel > 3 {
			return errors.New("cannot move: would exceed maximum depth")
		}
	}

	// Update category and all children's levels
	return s.updateCategoryLevel(id, newParentID, newLevel)
}

func (s *KnowledgeService) updateCategoryLevel(id string, parentID *string, level int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update this category
		if err := tx.Model(&models.Category{}).Where("id = ?", id).Updates(map[string]interface{}{
			"parent_id": parentID,
			"level":     level,
		}).Error; err != nil {
			return err
		}

		// Update all children recursively
		var children []models.Category
		if err := tx.Where("parent_id = ?", id).Find(&children).Error; err != nil {
			return err
		}

		for _, child := range children {
			if err := s.updateCategoryLevelInTx(tx, child.ID, level+1); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *KnowledgeService) updateCategoryLevelInTx(tx *gorm.DB, id string, level int) error {
	if level > 3 {
		return errors.New("maximum depth exceeded")
	}

	if err := tx.Model(&models.Category{}).Where("id = ?", id).Update("level", level).Error; err != nil {
		return err
	}

	var children []models.Category
	if err := tx.Where("parent_id = ?", id).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := s.updateCategoryLevelInTx(tx, child.ID, level+1); err != nil {
			return err
		}
	}

	return nil
}

// ==================== Knowledge Item Operations ====================

// ListKnowledgeQuery represents query parameters for listing knowledge
type ListKnowledgeQuery struct {
	CategoryID string
	Status     string
	TagID      string
	Search     string
	Page       int
	PageSize   int
}

// ListKnowledge returns paginated knowledge items
func (s *KnowledgeService) ListKnowledge(query ListKnowledgeQuery) ([]models.Knowledge, int64, error) {
	var knowledge []models.Knowledge
	var total int64

	db := s.db.Model(&models.Knowledge{})

	// Apply filters
	if query.CategoryID != "" {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.TagID != "" {
		db = db.Joins("JOIN knowledge_tags ON knowledge_tags.knowledge_id = knowledge.id").
			Where("knowledge_tags.tag_id = ?", query.TagID)
	}
	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").
		Preload("Tags").
		Offset(offset).
		Limit(query.PageSize).
		Find(&knowledge).Error; err != nil {
		return nil, 0, err
	}

	return knowledge, total, nil
}

// GetKnowledgeByID returns a single knowledge item with all relations
func (s *KnowledgeService) GetKnowledgeByID(id string) (*models.Knowledge, error) {
	var knowledge models.Knowledge
	if err := s.db.Preload("Tags").First(&knowledge, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("knowledge not found")
		}
		return nil, err
	}

	// Increment view count asynchronously
	go s.incrementViewCount(id)

	return &knowledge, nil
}

func (s *KnowledgeService) incrementViewCount(id string) {
	s.db.Model(&models.Knowledge{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1", ))
}

// CreateKnowledge creates a new knowledge item
func (s *KnowledgeService) CreateKnowledge(title, content, categoryID, authorID, source string, tagIDs []string) (*models.Knowledge, error) {
	// Validate category exists
	var category models.Category
	if err := s.db.First(&category, "id = ?", categoryID).Error; err != nil {
		return nil, errors.New("category not found")
	}

	knowledge := &models.Knowledge{
		Title:      title,
		Content:    content,
		CategoryID: categoryID,
		AuthorID:   authorID,
		Source:     source,
		Status:     "draft",
	}

	if err := s.db.Create(knowledge).Error; err != nil {
		return nil, err
	}

	// Associate tags
	if len(tagIDs) > 0 {
		var tags []models.Tag
		if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(knowledge).Association("Tags").Append(&tags); err != nil {
			return nil, err
		}
	}

	return knowledge, nil
}

// UpdateKnowledge updates a knowledge item
func (s *KnowledgeService) UpdateKnowledge(id, title, content, categoryID string, tagIDs []string) (*models.Knowledge, error) {
	var knowledge models.Knowledge
	if err := s.db.First(&knowledge, "id = ?", id).Error; err != nil {
		return nil, errors.New("knowledge not found")
	}

	// Validate new category if changed
	if categoryID != "" && categoryID != knowledge.CategoryID {
		var category models.Category
		if err := s.db.First(&category, "id = ?", categoryID).Error; err != nil {
			return nil, errors.New("category not found")
		}
		knowledge.CategoryID = categoryID
	}

	knowledge.Title = title
	knowledge.Content = content

	if err := s.db.Save(&knowledge).Error; err != nil {
		return nil, err
	}

	// Update tags
	if tagIDs != nil {
		var tags []models.Tag
		if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(&knowledge).Association("Tags").Replace(&tags); err != nil {
			return nil, err
		}
	}

	// Reload with tags
	if err := s.db.Preload("Tags").First(&knowledge, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &knowledge, nil
}

// DeleteKnowledge deletes a knowledge item
func (s *KnowledgeService) DeleteKnowledge(id string) error {
	var knowledge models.Knowledge
	if err := s.db.First(&knowledge, "id = ?", id).Error; err != nil {
		return errors.New("knowledge not found")
	}

	return s.db.Delete(&knowledge).Error
}

// PublishKnowledge publishes a knowledge item
func (s *KnowledgeService) PublishKnowledge(id string) error {
	return s.db.Model(&models.Knowledge{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": gorm.Expr("NOW()"),
		"version":      gorm.Expr("version + 1"),
	}).Error
}

// ArchiveKnowledge archives a knowledge item
func (s *KnowledgeService) ArchiveKnowledge(id string) error {
	return s.db.Model(&models.Knowledge{}).Where("id = ?", id).Update("status", "archived").Error
}

// ==================== Tag Operations ====================

// ListTags returns all tags
func (s *KnowledgeService) ListTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Order("name").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CreateTag creates a new tag
func (s *KnowledgeService) CreateTag(name, color string) (*models.Tag, error) {
	tag := &models.Tag{
		Name:  name,
		Color: color,
	}

	if err := s.db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag deletes a tag
func (s *KnowledgeService) DeleteTag(id string) error {
	return s.db.Delete(&models.Tag{}, "id = ?", id).Error
}

// ==================== Review Operations ====================

// SubmitReview submits a knowledge item for review
func (s *KnowledgeService) SubmitReview(knowledgeID, reviewerID string) (*models.KnowledgeReview, error) {
	review := &models.KnowledgeReview{
		KnowledgeID: knowledgeID,
		ReviewerID:  reviewerID,
		Status:      "pending",
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

// ApproveReview approves a knowledge review
func (s *KnowledgeService) ApproveReview(reviewID, comment string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update review
		if err := tx.Model(&models.KnowledgeReview{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
			"status":  "approved",
			"comment": comment,
		}).Error; err != nil {
			return err
		}

		// Get knowledge ID and publish
		var review models.KnowledgeReview
		if err := tx.First(&review, "id = ?", reviewID).Error; err != nil {
			return err
		}

		return tx.Model(&models.Knowledge{}).Where("id = ?", review.KnowledgeID).Updates(map[string]interface{}{
			"status":       "published",
			"published_at": gorm.Expr("NOW()"),
		}).Error
	})
}

// RejectReview rejects a knowledge review
func (s *KnowledgeService) RejectReview(reviewID, comment string) error {
	return s.db.Model(&models.KnowledgeReview{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
		"status":  "rejected",
		"comment": comment,
	}).Error
}

// GetPendingReviews returns pending reviews for a user
func (s *KnowledgeService) GetPendingReviews(reviewerID string) ([]models.KnowledgeReview, error) {
	var reviews []models.KnowledgeReview
	if err := s.db.Where("reviewer_id = ? AND status = ?", reviewerID, "pending").
		Preload("Knowledge").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
