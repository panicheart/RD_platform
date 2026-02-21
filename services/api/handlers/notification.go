package handlers

import (
	"net/http"
	"strconv"

	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification HTTP requests
type NotificationHandler struct {
	notificationService *services.NotificationService
	announcementService *services.AnnouncementService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService *services.NotificationService, announcementService *services.AnnouncementService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		announcementService: announcementService,
	}
}

// ListNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	notifications, total, err := h.notificationService.ListUserNotifications(c.Request.Context(), userID.(string), unreadOnly, page, pageSize)
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
		"data": gin.H{
			"items":     notifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNotification handles GET /api/v1/notifications/:id
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.notificationService.GetNotificationByID(c.Request.Context(), id)
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
		"data":    notification,
	})
}

// MarkAsRead handles PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	err := h.notificationService.MarkAsRead(c.Request.Context(), id)
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
		"message": "notification marked as read",
		"data":    nil,
	})
}

// MarkAllAsRead handles PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID.(string))
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
		"message": "all notifications marked as read",
		"data":    nil,
	})
}

// DeleteNotification handles DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	err := h.notificationService.DeleteNotification(c.Request.Context(), id)
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
		"message": "notification deleted",
		"data":    nil,
	})
}

// GetUnreadCount handles GET /api/v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID.(string))
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
		"data": gin.H{
			"count": count,
		},
	})
}

// ListAnnouncements handles GET /api/v1/announcements
func (h *NotificationHandler) ListAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	announcements, total, err := h.announcementService.ListActiveAnnouncements(c.Request.Context(), page, pageSize)
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
		"data": gin.H{
			"items":     announcements,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAnnouncement handles GET /api/v1/announcements/:id
func (h *NotificationHandler) GetAnnouncement(c *gin.Context) {
	id := c.Param("id")

	announcement, err := h.announcementService.GetAnnouncementByID(c.Request.Context(), id)
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
		"data":    announcement,
	})
}

// CreateAnnouncement handles POST /api/v1/announcements
func (h *NotificationHandler) CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	err := h.announcementService.CreateAnnouncement(c.Request.Context(), &announcement)
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
		"message": "announcement created successfully",
		"data":    announcement,
	})
}

// UpdateAnnouncement handles PUT /api/v1/announcements/:id
func (h *NotificationHandler) UpdateAnnouncement(c *gin.Context) {
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

	announcement, err := h.announcementService.UpdateAnnouncement(c.Request.Context(), id, updates)
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
		"message": "announcement updated successfully",
		"data":    announcement,
	})
}

// DeleteAnnouncement handles DELETE /api/v1/announcements/:id
func (h *NotificationHandler) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	err := h.announcementService.DeleteAnnouncement(c.Request.Context(), id)
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
		"message": "announcement deleted successfully",
		"data":    nil,
	})
}
