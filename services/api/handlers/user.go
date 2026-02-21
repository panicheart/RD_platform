package handlers

import (
	"net/http"
	"strconv"

	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService         *services.UserService
	organizationService *services.OrganizationService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, orgService *services.OrganizationService) *UserHandler {
	return &UserHandler{
		userService:         userService,
		organizationService: orgService,
	}
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if team := c.Query("team"); team != "" {
		filters["team"] = team
	}
	if orgID := c.Query("organization_id"); orgID != "" {
		filters["organization_id"] = orgID
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Set defaults
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	if !user.IsActive {
		user.IsActive = true
	}

	err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "user created successfully",
		"data":    user,
	})
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "username")
	delete(updates, "password_hash")
	delete(updates, "casdoor_id")

	user, err := h.userService.UpdateUser(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "user updated successfully",
		"data":    user,
	})
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "user deleted successfully",
		"data":    nil,
	})
}

// ListOrganizations handles GET /api/v1/organizations
func (h *UserHandler) ListOrganizations(c *gin.Context) {
	orgs, err := h.organizationService.ListOrganizations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    orgs,
	})
}

// GetOrganization handles GET /api/v1/organizations/:id
func (h *UserHandler) GetOrganization(c *gin.Context) {
	id := c.Param("id")

	org, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    org,
	})
}

// CreateOrganization handles POST /api/v1/organizations
func (h *UserHandler) CreateOrganization(c *gin.Context) {
	var org models.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	err := h.organizationService.CreateOrganization(c.Request.Context(), &org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "organization created successfully",
		"data":    org,
	})
}

// UpdateOrganization handles PUT /api/v1/organizations/:id
func (h *UserHandler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")

	org, err := h.organizationService.UpdateOrganization(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "organization updated successfully",
		"data":    org,
	})
}

// DeleteOrganization handles DELETE /api/v1/organizations/:id
func (h *UserHandler) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")

	err := h.organizationService.DeleteOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "organization deleted successfully",
		"data":    nil,
	})
}

// CurrentUser handles GET /api/v1/users/me
func (h *UserHandler) CurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// UpdateCurrentUser handles PUT /api/v1/users/me
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Remove immutable fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "username")
	delete(updates, "password_hash")
	delete(updates, "casdoor_id")
	delete(updates, "role")
	delete(updates, "organization_id")

	user, err := h.userService.UpdateUser(c.Request.Context(), userID.(string), updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "profile updated successfully",
		"data":    user,
	})
}
