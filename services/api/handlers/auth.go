package handlers

import (
	"net/http"

	"rdp-platform/rdp-api/middleware"
	"rdp-platform/rdp-api/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	userService   *services.UserService
	jwtMiddleware *middleware.JWTMiddleware
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService *services.UserService, jwtMiddleware *middleware.JWTMiddleware) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		jwtMiddleware: jwtMiddleware,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	User         interface{} `json:"user"`
}

// RefreshRequest represents refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, "invalid request: "+err.Error())
		return
	}

	// Get user by username
	user, err := h.userService.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, 4013, "invalid credentials")
		return
	}

	// Check password
	if user.PasswordHash == nil || !h.checkPassword(req.Password, *user.PasswordHash) {
		ErrorResponse(c, http.StatusUnauthorized, 4013, "invalid credentials")
		return
	}

	// Update last login
	h.userService.UpdateLastLogin(c.Request.Context(), user.ID)

	// Generate tokens
	accessToken, refreshToken, err := h.jwtMiddleware.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, "failed to generate tokens")
		return
	}

	SuccessResponse(c, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2 hours
		TokenType:    "Bearer",
		User: gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"role":         user.Role,
			"email":        user.Email,
			"avatar_url":   user.AvatarURL,
		},
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, "invalid request")
		return
	}

	// Get user function for refresh
	getUserFunc := func(userID string) (string, string, error) {
		user, err := h.userService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			return "", "", err
		}
		return user.Username, user.Role, nil
	}

	newAccessToken, newRefreshToken, err := h.jwtMiddleware.RefreshTokens(req.RefreshToken, getUserFunc)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, 4014, "invalid refresh token")
		return
	}

	SuccessResponse(c, LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a complete implementation, you would blacklist the token here
	// For now, we just return success and let the client delete the token
	SuccessResponse(c, gin.H{
		"message": "logged out successfully",
	})
}

// GetCurrentUser returns current user info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 4015, "not authenticated")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 4041, "user not found")
		return
	}

	SuccessResponse(c, user)
}

// checkPassword compares password with hash
func (h *AuthHandler) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, 4015, "not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, 4001, "invalid request")
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, 4041, "user not found")
		return
	}

	// Verify old password
	if user.PasswordHash == nil || !h.checkPassword(req.OldPassword, *user.PasswordHash) {
		ErrorResponse(c, http.StatusUnauthorized, 4013, "invalid old password")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, "failed to hash password")
		return
	}

	hashStr := string(hashedPassword)
	updates := map[string]interface{}{
		"password_hash": &hashStr,
	}

	_, err = h.userService.UpdateUser(c.Request.Context(), userID.(string), updates)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, 5001, "failed to update password")
		return
	}

	SuccessResponse(c, gin.H{
		"message": "password changed successfully",
	})
}
