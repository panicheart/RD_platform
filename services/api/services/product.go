package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"rdp-platform/rdp-api/models"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ProductService handles product shelf business logic
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a new ProductService
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// ListProducts returns paginated products with filters
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).Preload("Owner").Preload("SourceProject")

	// Apply filters
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if trlLevel, ok := filters["trl_level"].(int); ok && trlLevel > 0 {
		query = query.Where("trl_level = ?", trlLevel)
	}
	if trlMin, ok := filters["trl_min"].(int); ok && trlMin > 0 {
		query = query.Where("trl_level >= ?", trlMin)
	}
	if trlMax, ok := filters["trl_max"].(int); ok && trlMax > 0 {
		query = query.Where("trl_level <= ?", trlMax)
	}
	if isPublished, ok := filters["is_published"].(bool); ok {
		query = query.Where("is_published = ?", isPublished)
	}
	if ownerID, ok := filters["owner_id"].(string); ok && ownerID != "" {
		query = query.Where("owner_id = ?", ownerID)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductByID returns a product by ID with all relations
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	
	if err := s.db.Preload("Owner").Preload("SourceProject").Preload("Versions").First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Increment download count
	s.db.Model(&product).UpdateColumn("download_count", gorm.Expr("download_count + 1"))

	return &product, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product, createdBy string) error {
	// Generate UUID if not provided
	if product.ID == "" {
		product.ID = ulid.Make().String()
	}
	
	product.CreatedBy = createdBy
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.DownloadCount = 0

	// Validate TRL level
	if product.TRLLevel < 1 || product.TRLLevel > 9 {
		return errors.New("TRL level must be between 1 and 9")
	}

	return s.db.Create(product).Error
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) (*models.Product, error) {
	var product models.Product
	
	if err := s.db.First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Validate TRL level if provided
	if trlLevel, ok := updates["trl_level"].(int); ok {
		if trlLevel < 1 || trlLevel > 9 {
			return nil, errors.New("TRL level must be between 1 and 9")
		}
	}

	// Update timestamp
	updates["updated_at"] = time.Now()

	if err := s.db.Model(&product).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct deletes a product and its versions
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	var product models.Product
	
	if err := s.db.First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Delete related cart items first
	if err := s.db.Where("product_id = ?", id).Delete(&models.CartItem{}).Error; err != nil {
		return err
	}

	// Delete product versions
	if err := s.db.Where("product_id = ?", id).Delete(&models.ProductVersion{}).Error; err != nil {
		return err
	}

	// Delete product
	return s.db.Delete(&product).Error
}

// PublishProduct publishes a product
func (s *ProductService) PublishProduct(ctx context.Context, id string) (*models.Product, error) {
	now := time.Now()
	updates := map[string]interface{}{
		"is_published": true,
		"published_at": &now,
		"updated_at":   now,
	}

	return s.UpdateProduct(ctx, id, updates)
}

// UnpublishProduct unpublishes a product
func (s *ProductService) UnpublishProduct(ctx context.Context, id string) (*models.Product, error) {
	updates := map[string]interface{}{
		"is_published": false,
		"published_at": nil,
		"updated_at":   time.Now(),
	}

	return s.UpdateProduct(ctx, id, updates)
}

// GetCategories returns all distinct product categories
func (s *ProductService) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	
	if err := s.db.Model(&models.Product{}).
		Where("category IS NOT NULL AND category != ''").
		Distinct().
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// ==================== Product Version Operations ====================

// CreateProductVersion creates a new product version
func (s *ProductService) CreateProductVersion(ctx context.Context, version *models.ProductVersion, createdBy string) error {
	// Check if product exists
	var product models.Product
	if err := s.db.First(&product, "id = ?", version.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Generate UUID
	version.ID = ulid.Make().String()
	version.CreatedBy = createdBy
	version.CreatedAt = time.Now()

	return s.db.Create(version).Error
}

// GetProductVersions returns all versions of a product
func (s *ProductService) GetProductVersions(ctx context.Context, productID string) ([]models.ProductVersion, error) {
	var versions []models.ProductVersion
	
	if err := s.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&versions).Error; err != nil {
		return nil, err
	}

	return versions, nil
}

// DeleteProductVersion deletes a product version
func (s *ProductService) DeleteProductVersion(ctx context.Context, versionID string) error {
	return s.db.Delete(&models.ProductVersion{}, "id = ?", versionID).Error
}

// ==================== Cart Operations ====================

// GetCartItems returns cart items for a user
func (s *ProductService) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	var items []models.CartItem
	
	if err := s.db.Where("user_id = ?", userID).
		Preload("Product").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// AddToCart adds a product to user's cart
func (s *ProductService) AddToCart(ctx context.Context, userID, productID string, quantity int, notes string) (*models.CartItem, error) {
	// Check if product exists and is published
	var product models.Product
	if err := s.db.First(&product, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	if !product.IsPublished {
		return nil, errors.New("product is not published")
	}

	// Check if already in cart
	var existingItem models.CartItem
	if err := s.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingItem).Error; err == nil {
		// Update quantity
		existingItem.Quantity += quantity
		if err := s.db.Save(&existingItem).Error; err != nil {
			return nil, err
		}
		return &existingItem, nil
	}

	// Create new cart item
	item := &models.CartItem{
		ID:        ulid.Make().String(),
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Notes:     notes,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}

	// Load product relation
	s.db.Preload("Product").First(item, "id = ?", item.ID)

	return item, nil
}

// UpdateCartItem updates a cart item
func (s *ProductService) UpdateCartItem(ctx context.Context, itemID string, quantity int, notes string) (*models.CartItem, error) {
	var item models.CartItem
	
	if err := s.db.First(&item, "id = ?", itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart item not found")
		}
		return nil, err
	}

	item.Quantity = quantity
	item.Notes = notes

	if err := s.db.Save(&item).Error; err != nil {
		return nil, err
	}

	s.db.Preload("Product").First(&item, "id = ?", item.ID)

	return &item, nil
}

// RemoveFromCart removes an item from cart
func (s *ProductService) RemoveFromCart(ctx context.Context, itemID string) error {
	return s.db.Delete(&models.CartItem{}, "id = ?", itemID).Error
}

// ClearCart clears all items from user's cart
func (s *ProductService) ClearCart(ctx context.Context, userID string) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}

// ==================== Technology Operations ====================

// ListTechnologies returns paginated technologies with filters
func (s *ProductService) ListTechnologies(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	query := s.db.Model(&models.Technology{}).Preload("Owner")

	// Apply filters
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if trlLevel, ok := filters["trl_level"].(int); ok && trlLevel > 0 {
		query = query.Where("trl_level = ?", trlLevel)
	}
	if isPublished, ok := filters["is_published"].(bool); ok {
		query = query.Where("is_published = ?", isPublished)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&technologies).Error; err != nil {
		return nil, 0, err
	}

	return technologies, total, nil
}

// GetTechnologyByID returns a technology by ID
func (s *ProductService) GetTechnologyByID(ctx context.Context, id string) (*models.Technology, error) {
	var technology models.Technology
	
	if err := s.db.Preload("Owner").Preload("Children").First(&technology, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("technology not found")
		}
		return nil, err
	}

	return &technology, nil
}

// CreateTechnology creates a new technology
func (s *ProductService) CreateTechnology(ctx context.Context, technology *models.Technology, createdBy string) error {
	if technology.ID == "" {
		technology.ID = ulid.Make().String()
	}
	
	technology.CreatedBy = createdBy
	technology.CreatedAt = time.Now()
	technology.UpdatedAt = time.Now()

	if technology.TRLLevel < 1 || technology.TRLLevel > 9 {
		return errors.New("TRL level must be between 1 and 9")
	}

	return s.db.Create(technology).Error
}

// UpdateTechnology updates a technology
func (s *ProductService) UpdateTechnology(ctx context.Context, id string, updates map[string]interface{}) (*models.Technology, error) {
	var technology models.Technology
	
	if err := s.db.First(&technology, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("technology not found")
		}
		return nil, err
	}

	if trlLevel, ok := updates["trl_level"].(int); ok {
		if trlLevel < 1 || trlLevel > 9 {
			return nil, errors.New("TRL level must be between 1 and 9")
		}
	}

	updates["updated_at"] = time.Now()

	if err := s.db.Model(&technology).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &technology, nil
}

// DeleteTechnology deletes a technology
func (s *ProductService) DeleteTechnology(ctx context.Context, id string) error {
	var technology models.Technology
	
	if err := s.db.First(&technology, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("technology not found")
		}
		return err
	}

	return s.db.Delete(&technology).Error
}

// GetTRLLevelName returns the name for a TRL level
func GetTRLLevelName(level int) string {
	names := map[int]string{
		1: "基本原理发现",
		2: "技术概念形成",
		3: "概念验证",
		4: "实验室验证",
		5: "相关环境验证",
		6: "系统/子系统验证",
		7: "系统原型验证",
		8: "系统完成验证",
		9: "实际应用验证",
	}
	
	if name, ok := names[level]; ok {
		return name
	}
	return fmt.Sprintf("TRL %d", level)
}

// GetTRLColor returns the color code for a TRL level
func GetTRLColor(level int) string {
	if level <= 3 {
		return "red"
	} else if level <= 6 {
		return "orange"
	}
	return "green"
}

// ParseMetadata parses product metadata JSON
func (s *ProductService) ParseMetadata(metadata string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if metadata == "" {
		return result, nil
	}
	
	if err := json.Unmarshal([]byte(metadata), &result); err != nil {
		return nil, err
	}
	
	return result, nil
}
