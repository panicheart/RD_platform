package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rdp-platform/rdp-api/services"
)

// ZoteroHandler handles Zotero integration HTTP requests
type ZoteroHandler struct {
	zoteroService *services.ZoteroService
}

// NewZoteroHandler creates a new Zotero handler
func NewZoteroHandler(zoteroService *services.ZoteroService) *ZoteroHandler {
	return &ZoteroHandler{zoteroService: zoteroService}
}

// ==================== Connection Management ====================

// SaveConnectionRequest represents the request body for saving Zotero connection
type SaveConnectionRequest struct {
	APIKey       string `json:"api_key" binding:"required"`
	ZoteroUserID string `json:"zotero_user_id" binding:"required"`
}

// SaveConnection saves or updates Zotero API credentials for the current user
func (h *ZoteroHandler) SaveConnection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	var req SaveConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, "invalid request: "+err.Error())
		return
	}

	connection, err := h.zoteroService.SaveConnection(userID.(string), req.APIKey, req.ZoteroUserID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	// Mask API key for security
	connection.APIKey = "***"
	SuccessResponse(c, connection)
}

// GetConnection retrieves the current user's Zotero connection status
func (h *ZoteroHandler) GetConnection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	connection, err := h.zoteroService.GetConnection(userID.(string))
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 404, err.Error())
		return
	}

	// Mask API key for security
	connection.APIKey = "***"
	SuccessResponse(c, connection)
}

// DeleteConnection removes the current user's Zotero connection
func (h *ZoteroHandler) DeleteConnection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	if err := h.zoteroService.DeleteConnection(userID.(string)); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "connection deleted"})
}

// TestConnectionRequest represents the request body for testing Zotero connection
type TestConnectionRequest struct {
	APIKey       string `json:"api_key" binding:"required"`
	ZoteroUserID string `json:"zotero_user_id" binding:"required"`
}

// TestConnection tests Zotero credentials without saving them
func (h *ZoteroHandler) TestConnection(c *gin.Context) {
	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.zoteroService.TestConnection(req.APIKey, req.ZoteroUserID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, "connection test failed: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"status": "connected", "message": "Zotero connection successful"})
}

// ==================== Sync Operations ====================

// SyncItemsRequest represents the request body for syncing items
type SyncItemsRequest struct {
	Incremental bool `json:"incremental"`
}

// SyncItems triggers a sync of Zotero items
func (h *ZoteroHandler) SyncItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	var req SyncItemsRequest
	c.ShouldBindJSON(&req)

	result, err := h.zoteroService.SyncItems(userID.(string), req.Incremental)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, result)
}

// ==================== Item Operations ====================

// ListItemsQuery represents query parameters for listing items
type ListItemsQuery struct {
	ItemType string `form:"item_type"`
	Tag      string `form:"tag"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// ListItems returns a paginated list of Zotero items
func (h *ZoteroHandler) ListItems(c *gin.Context) {
	var query services.ListItemsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, "invalid query parameters: "+err.Error())
		return
	}

	items, total, err := h.zoteroService.ListItems(query)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"items": items,
		"total": total,
		"page":  query.Page,
	})
}

// GetItem returns a single Zotero item by ID
func (h *ZoteroHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ErrorResponse(c, http.StatusBadRequest, 400, "item id is required")
		return
	}

	item, err := h.zoteroService.GetItemByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 404, err.Error())
		return
	}

	SuccessResponse(c, item)
}

// DeleteItem deletes a Zotero item from the database
func (h *ZoteroHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ErrorResponse(c, http.StatusBadRequest, 400, "item id is required")
		return
	}

	if err := h.zoteroService.DeleteItem(id); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "item deleted"})
}

// ==================== PDF Operations ====================

// GetPDFURL returns the URL for viewing a PDF attachment
func (h *ZoteroHandler) GetPDFURL(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		ErrorResponse(c, http.StatusBadRequest, 400, "item id is required")
		return
	}

	url, err := h.zoteroService.GetPDFURL(userID.(string), itemID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"url": url})
}

// ==================== Citation Operations ====================

// GenerateCitationRequest represents the request body for generating citations
type GenerateCitationRequest struct {
	Format string `json:"format" binding:"required"` // gb7714, apa, mla
}

// GenerateCitation generates a citation for a Zotero item
func (h *ZoteroHandler) GenerateCitation(c *gin.Context) {
	itemID := c.Param("id")
	if itemID == "" {
		ErrorResponse(c, http.StatusBadRequest, 400, "item id is required")
		return
	}

	var req GenerateCitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, "invalid request: "+err.Error())
		return
	}

	citation, err := h.zoteroService.GenerateCitation(itemID, req.Format)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"citation": citation,
		"format":   req.Format,
	})
}

// ==================== Collection Operations ====================

// GetCollections returns all Zotero collections for the current user
func (h *ZoteroHandler) GetCollections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	collections, err := h.zoteroService.GetCollections(userID.(string))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	SuccessResponse(c, collections)
}
