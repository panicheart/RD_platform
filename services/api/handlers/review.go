package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	reviewService *services.ReviewService
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReviewRequest represents the request body for creating a review
type CreateReviewRequest struct {
	ActivityID string           `json:"activity_id" binding:"required"`
	Type       models.ReviewType `json:"type" binding:"required"`
}

// CreateReview creates a new review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	projectID := c.Param("projectId")
	userID, _ := c.Get("userID")

	review, err := h.reviewService.CreateReview(
		req.ActivityID,
		projectID,
		req.Type,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Review created successfully", "data": review})
}

// GetReview retrieves a review by ID
func (h *ReviewHandler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")

	review, err := h.reviewService.GetReview(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Review not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Success", "data": review})
}

// ListReviews retrieves reviews for an activity
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	activityID := c.Query("activity_id")
	if activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "activity_id is required", "data": nil})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.reviewService.ListReviews(activityID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"items":     reviews,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// SubmitReviewRequest represents the request body for submitting a review
type SubmitReviewRequest struct {
	Comments string `json:"comments"`
	Score    int    `json:"score" binding:"required,min=0,max=100"`
}

// SubmitReview submits a review for approval
func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	reviewID := c.Param("id")

	var req SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	if err := h.reviewService.SubmitReview(reviewID, req.Comments, req.Score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Review submitted successfully", "data": nil})
}

// ApproveReview approves a review
func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")

	if err := h.reviewService.ApproveReview(reviewID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Review approved successfully", "data": nil})
}

// RejectReview rejects a review
func (h *ReviewHandler) RejectReview(c *gin.Context) {
	reviewID := c.Param("id")

	if err := h.reviewService.RejectReview(reviewID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Review rejected successfully", "data": nil})
}

// RequestRevision requests revision for a review
func (h *ReviewHandler) RequestRevision(c *gin.Context) {
	reviewID := c.Param("id")

	if err := h.reviewService.RequestRevision(reviewID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Revision requested successfully", "data": nil})
}
