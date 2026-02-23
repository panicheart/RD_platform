package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rdp-platform/rdp-api/services"
)

// KnowledgeHandler handles knowledge base related requests
type KnowledgeHandler struct {
	knowledgeService *services.KnowledgeService
}

// NewKnowledgeHandler creates a new KnowledgeHandler
func NewKnowledgeHandler(knowledgeService *services.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{
		knowledgeService: knowledgeService,
	}
}

// ListKnowledge handles GET /api/v1/knowledge
func (h *KnowledgeHandler) ListKnowledge(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")
	status := c.DefaultQuery("status", "published")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := services.ListKnowledgeQuery{
		CategoryID: categoryID,
		Status:     status,
		Search:     search,
		Page:       page,
		PageSize:   pageSize,
	}

	items, total, err := h.knowledgeService.ListKnowledge(query)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetKnowledge handles GET /api/v1/knowledge/:id
func (h *KnowledgeHandler) GetKnowledge(c *gin.Context) {
	id := c.Param("id")

	item, err := h.knowledgeService.GetKnowledgeByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 4041, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// CreateKnowledgeRequest represents create knowledge request
type CreateKnowledgeRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content"`
	CategoryID string   `json:"category_id" binding:"required"`
	Tags       []string `json:"tags"`
}

// CreateKnowledge handles POST /api/v1/knowledge
func (h *KnowledgeHandler) CreateKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 4010, "unauthorized")
		return
	}

	var req CreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, err.Error())
		return
	}

	item, err := h.knowledgeService.CreateKnowledge(
		req.Title,
		req.Content,
		req.CategoryID,
		userID.(string),
		"manual",
		req.Tags,
	)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// UpdateKnowledgeRequest represents update knowledge request
type UpdateKnowledgeRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CategoryID string   `json:"category_id"`
	Tags       []string `json:"tags"`
}

// UpdateKnowledge handles PUT /api/v1/knowledge/:id
func (h *KnowledgeHandler) UpdateKnowledge(c *gin.Context) {
	id := c.Param("id")

	var req UpdateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, err.Error())
		return
	}

	item, err := h.knowledgeService.UpdateKnowledge(
		id,
		req.Title,
		req.Content,
		req.CategoryID,
		req.Tags,
	)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// DeleteKnowledge handles DELETE /api/v1/knowledge/:id
func (h *KnowledgeHandler) DeleteKnowledge(c *gin.Context) {
	id := c.Param("id")

	if err := h.knowledgeService.DeleteKnowledge(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "knowledge deleted successfully"})
}

// PublishKnowledge handles POST /api/v1/knowledge/:id/publish
func (h *KnowledgeHandler) PublishKnowledge(c *gin.Context) {
	id := c.Param("id")

	if err := h.knowledgeService.PublishKnowledge(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "knowledge published successfully"})
}

// GetCategoryTree handles GET /api/v1/categories/tree
func (h *KnowledgeHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.knowledgeService.GetCategoryTree()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, tree)
}

// CreateCategoryRequest represents create category request
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"`
}

// CreateCategory handles POST /api/v1/categories
func (h *KnowledgeHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, err.Error())
		return
	}

	category, err := h.knowledgeService.CreateCategory(req.Name, req.Description, req.ParentID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, category)
}

// UpdateCategoryRequest represents update category request
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateCategory handles PUT /api/v1/categories/:id
func (h *KnowledgeHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, err.Error())
		return
	}

	category, err := h.knowledgeService.UpdateCategory(id, req.Name, req.Description)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, category)
}

// DeleteCategory handles DELETE /api/v1/categories/:id
func (h *KnowledgeHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.knowledgeService.DeleteCategory(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "category deleted successfully"})
}

// ListTags handles GET /api/v1/tags
func (h *KnowledgeHandler) ListTags(c *gin.Context) {
	tags, err := h.knowledgeService.ListTags()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}
	SuccessResponse(c, tags)
}

// CreateTagRequest represents create tag request
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

// CreateTag handles POST /api/v1/tags
func (h *KnowledgeHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, err.Error())
		return
	}

	tag, err := h.knowledgeService.CreateTag(req.Name, req.Color)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, tag)
}

// SearchKnowledge handles GET /api/v1/knowledge/search
func (h *KnowledgeHandler) SearchKnowledge(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		ErrorResponse(c, http.StatusBadRequest, 4001, "search query required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Use ListKnowledge with search filter
	q := services.ListKnowledgeQuery{
		Status:   "published",
		Search:   query,
		Page:     page,
		PageSize: pageSize,
	}

	items, total, err := h.knowledgeService.ListKnowledge(q)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPopularKnowledge handles GET /api/v1/knowledge/popular
func (h *KnowledgeHandler) GetPopularKnowledge(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Use ListKnowledge with pagination
	q := services.ListKnowledgeQuery{
		Status:   "published",
		Page:     1,
		PageSize: limit,
	}

	items, _, err := h.knowledgeService.ListKnowledge(q)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, items)
}

// GetRecentKnowledge handles GET /api/v1/knowledge/recent
func (h *KnowledgeHandler) GetRecentKnowledge(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Use ListKnowledge with pagination
	q := services.ListKnowledgeQuery{
		Status:   "published",
		Page:     1,
		PageSize: limit,
	}

	items, _, err := h.knowledgeService.ListKnowledge(q)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, err.Error())
		return
	}

	SuccessResponse(c, items)
}
