package handlers

import (
	"encoding/json"
	"net/http"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles product shelf HTTP requests
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// ==================== Product Operations ====================

// ListProductsRequest represents the query parameters for listing products
type ListProductsRequest struct {
	Page         int    `form:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" binding:"min=1,max=100"`
	Category     string `form:"category"`
	TRLLevel     int    `form:"trl_level"`
	TRLMin       int    `form:"trl_min"`
	TRLMax       int    `form:"trl_max"`
	IsPublished  *bool  `form:"is_published"`
	Search       string `form:"search"`
}

// ListProducts handles GET /api/v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BadRequestResponse(c, "invalid query parameters: "+err.Error())
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.Category != "" {
		filters["category"] = req.Category
	}
	if req.TRLLevel > 0 {
		filters["trl_level"] = req.TRLLevel
	}
	if req.TRLMin > 0 {
		filters["trl_min"] = req.TRLMin
	}
	if req.TRLMax > 0 {
		filters["trl_max"] = req.TRLMax
	}
	if req.IsPublished != nil {
		filters["is_published"] = *req.IsPublished
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	products, total, err := h.productService.ListProducts(c.Request.Context(), req.Page, req.PageSize, filters)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items":     products,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetProduct handles GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	
	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 404, err.Error())
		return
	}

	SuccessResponse(c, product)
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name            string                 `json:"name" binding:"required,max=200"`
	Description     string                 `json:"description"`
	TRLLevel        int                    `json:"trl_level" binding:"min=1,max=9"`
	Category        string                 `json:"category"`
	Version         string                 `json:"version"`
	SourceProjectID string                 `json:"source_project_id"`
	OwnerID         string                 `json:"owner_id"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	// Build product
	product := &models.Product{
		Name:            req.Name,
		Description:     req.Description,
		TRLLevel:        req.TRLLevel,
		Category:        req.Category,
		Version:         req.Version,
		SourceProjectID: req.SourceProjectID,
		OwnerID:         req.OwnerID,
	}

	if req.Metadata != nil {
		metadataJSON, _ := json.Marshal(req.Metadata)
		product.Metadata = string(metadataJSON)
	}

	if err := h.productService.CreateProduct(c.Request.Context(), product, userIDStr); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, product)
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string                 `json:"name" binding:"max=200"`
	Description string                 `json:"description"`
	TRLLevel    int                    `json:"trl_level" binding:"min=1,max=9"`
	Category    string                 `json:"category"`
	Version     string                 `json:"version"`
	OwnerID     string                 `json:"owner_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateProduct handles PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Build updates
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.TRLLevel > 0 {
		updates["trl_level"] = req.TRLLevel
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Version != "" {
		updates["version"] = req.Version
	}
	if req.OwnerID != "" {
		updates["owner_id"] = req.OwnerID
	}
	if req.Metadata != nil {
		metadataJSON, _ := json.Marshal(req.Metadata)
		updates["metadata"] = string(metadataJSON)
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, updates)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, product)
}

// DeleteProduct handles DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// PublishProduct handles POST /api/v1/products/:id/publish
func (h *ProductHandler) PublishProduct(c *gin.Context) {
	id := c.Param("id")
	
	product, err := h.productService.PublishProduct(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, product)
}

// UnpublishProduct handles POST /api/v1/products/:id/unpublish
func (h *ProductHandler) UnpublishProduct(c *gin.Context) {
	id := c.Param("id")
	
	product, err := h.productService.UnpublishProduct(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, product)
}

// GetProductCategories handles GET /api/v1/products/categories
func (h *ProductHandler) GetProductCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories(c.Request.Context())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, categories)
}

// ==================== Product Version Operations ====================

// CreateProductVersionRequest represents the request body for creating a product version
type CreateProductVersionRequest struct {
	Version         string                 `json:"version" binding:"required"`
	ParentVersionID string                 `json:"parent_version_id"`
	Changelog       string                 `json:"changelog"`
	Files           map[string]interface{} `json:"files"`
}

// CreateProductVersion handles POST /api/v1/products/:id/versions
func (h *ProductHandler) CreateProductVersion(c *gin.Context) {
	productID := c.Param("id")
	
	var req CreateProductVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	version := &models.ProductVersion{
		ProductID:       productID,
		Version:         req.Version,
		ParentVersionID: nil,
		Changelog:       req.Changelog,
	}

	if req.ParentVersionID != "" {
		version.ParentVersionID = &req.ParentVersionID
	}

	if req.Files != nil {
		filesJSON, _ := json.Marshal(req.Files)
		version.Files = string(filesJSON)
	}

	if err := h.productService.CreateProductVersion(c.Request.Context(), version, userIDStr); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, version)
}

// GetProductVersions handles GET /api/v1/products/:id/versions
func (h *ProductHandler) GetProductVersions(c *gin.Context) {
	productID := c.Param("id")
	
	versions, err := h.productService.GetProductVersions(c.Request.Context(), productID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, versions)
}

// DeleteProductVersion handles DELETE /api/v1/products/versions/:versionId
func (h *ProductHandler) DeleteProductVersion(c *gin.Context) {
	versionID := c.Param("versionId")
	
	if err := h.productService.DeleteProductVersion(c.Request.Context(), versionID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// ==================== Cart Operations ====================

// GetCartItems handles GET /api/v1/cart
func (h *ProductHandler) GetCartItems(c *gin.Context) {
	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	items, err := h.productService.GetCartItems(c.Request.Context(), userIDStr)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, items)
}

// AddToCartRequest represents the request body for adding to cart
type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"min=1"`
	Notes     string `json:"notes"`
}

// AddToCart handles POST /api/v1/cart
func (h *ProductHandler) AddToCart(c *gin.Context) {
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	if req.Quantity == 0 {
		req.Quantity = 1
	}

	item, err := h.productService.AddToCart(c.Request.Context(), userIDStr, req.ProductID, req.Quantity, req.Notes)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// UpdateCartItemRequest represents the request body for updating cart item
type UpdateCartItemRequest struct {
	Quantity int    `json:"quantity" binding:"min=1"`
	Notes    string `json:"notes"`
}

// UpdateCartItem handles PUT /api/v1/cart/:itemId
func (h *ProductHandler) UpdateCartItem(c *gin.Context) {
	itemID := c.Param("itemId")
	
	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	item, err := h.productService.UpdateCartItem(c.Request.Context(), itemID, req.Quantity, req.Notes)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// RemoveFromCart handles DELETE /api/v1/cart/:itemId
func (h *ProductHandler) RemoveFromCart(c *gin.Context) {
	itemID := c.Param("itemId")
	
	if err := h.productService.RemoveFromCart(c.Request.Context(), itemID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// ClearCart handles DELETE /api/v1/cart
func (h *ProductHandler) ClearCart(c *gin.Context) {
	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	if err := h.productService.ClearCart(c.Request.Context(), userIDStr); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// ==================== Technology Operations ====================

// ListTechnologiesRequest represents the query parameters for listing technologies
type ListTechnologiesRequest struct {
	Page        int    `form:"page" binding:"min=1"`
	PageSize    int    `form:"page_size" binding:"min=1,max=100"`
	Category    string `form:"category"`
	TRLLevel    int    `form:"trl_level"`
	IsPublished *bool  `form:"is_published"`
	Search      string `form:"search"`
}

// ListTechnologies handles GET /api/v1/technologies
func (h *ProductHandler) ListTechnologies(c *gin.Context) {
	var req ListTechnologiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BadRequestResponse(c, "invalid query parameters: "+err.Error())
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.Category != "" {
		filters["category"] = req.Category
	}
	if req.TRLLevel > 0 {
		filters["trl_level"] = req.TRLLevel
	}
	if req.IsPublished != nil {
		filters["is_published"] = *req.IsPublished
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	technologies, total, err := h.productService.ListTechnologies(c.Request.Context(), req.Page, req.PageSize, filters)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items":     technologies,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetTechnology handles GET /api/v1/technologies/:id
func (h *ProductHandler) GetTechnology(c *gin.Context) {
	id := c.Param("id")
	
	technology, err := h.productService.GetTechnologyByID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 404, err.Error())
		return
	}

	SuccessResponse(c, technology)
}

// CreateTechnologyRequest represents the request body for creating a technology
type CreateTechnologyRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	TRLLevel    int    `json:"trl_level" binding:"min=1,max=9"`
	Category    string `json:"category"`
	ParentID    string `json:"parent_id"`
	OwnerID     string `json:"owner_id"`
}

// CreateTechnology handles POST /api/v1/technologies
func (h *ProductHandler) CreateTechnology(c *gin.Context) {
	var req CreateTechnologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Get current user ID
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	technology := &models.Technology{
		Name:        req.Name,
		Description: req.Description,
		TRLLevel:    req.TRLLevel,
		Category:    req.Category,
		OwnerID:     req.OwnerID,
	}

	if req.ParentID != "" {
		technology.ParentID = &req.ParentID
	}

	if err := h.productService.CreateTechnology(c.Request.Context(), technology, userIDStr); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, technology)
}

// UpdateTechnologyRequest represents the request body for updating a technology
type UpdateTechnologyRequest struct {
	Name        string `json:"name" binding:"max=200"`
	Description string `json:"description"`
	TRLLevel    int    `json:"trl_level" binding:"min=1,max=9"`
	Category    string `json:"category"`
	OwnerID     string `json:"owner_id"`
	IsPublished *bool  `json:"is_published"`
}

// UpdateTechnology handles PUT /api/v1/technologies/:id
func (h *ProductHandler) UpdateTechnology(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateTechnologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "invalid request body: "+err.Error())
		return
	}

	// Build updates
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.TRLLevel > 0 {
		updates["trl_level"] = req.TRLLevel
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.OwnerID != "" {
		updates["owner_id"] = req.OwnerID
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}

	technology, err := h.productService.UpdateTechnology(c.Request.Context(), id, updates)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, technology)
}

// DeleteTechnology handles DELETE /api/v1/technologies/:id
func (h *ProductHandler) DeleteTechnology(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.productService.DeleteTechnology(c.Request.Context(), id); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, nil)
}

// ==================== Helper Functions ====================

// GetTRLLevels handles GET /api/v1/trl-levels
func (h *ProductHandler) GetTRLLevels(c *gin.Context) {
	levels := []gin.H{
		{ "level": 1, "name": "基本原理发现", "color": "red" },
		{ "level": 2, "name": "技术概念形成", "color": "red" },
		{ "level": 3, "name": "概念验证", "color": "red" },
		{ "level": 4, "name": "实验室验证", "color": "orange" },
		{ "level": 5, "name": "相关环境验证", "color": "orange" },
		{ "level": 6, "name": "系统/子系统验证", "color": "orange" },
		{ "level": 7, "name": "系统原型验证", "color": "green" },
		{ "level": 8, "name": "系统完成验证", "color": "green" },
		{ "level": 9, "name": "实际应用验证", "color": "green" },
	}

	SuccessResponse(c, levels)
}
