package services

import (
	"context"
	"errors"

	"rdp/services/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationService handles notification business logic
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// ListUserNotifications returns notifications for a user
func (s *NotificationService) ListUserNotifications(ctx context.Context, userID string, unreadOnly bool, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", uid)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetNotificationByID returns a notification by ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, id string) (*models.Notification, error) {
	var notification models.Notification
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid notification ID")
	}

	if err := s.db.First(&notification, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}

	return &notification, nil
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
	notification.ID = uuid.New()
	return s.db.Create(notification).Error
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid notification ID")
	}

	result := s.db.Model(&models.Notification{}).Where("id = ?", uid).Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", uid, false).Update("is_read", true).Error
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid notification ID")
	}

	result := s.db.Where("id = ?", uid).Delete(&models.Notification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("notification not found")
	}

	return nil
}

// GetUnreadCount returns the count of unread notifications for a user
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}

	if err := s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", uid, false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// AnnouncementService handles announcement business logic
type AnnouncementService struct {
	db *gorm.DB
}

// NewAnnouncementService creates a new AnnouncementService
func NewAnnouncementService(db *gorm.DB) *AnnouncementService {
	return &AnnouncementService{db: db}
}

// ListActiveAnnouncements returns active announcements
func (s *AnnouncementService) ListActiveAnnouncements(ctx context.Context, page, pageSize int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	query := s.db.Model(&models.Announcement{}).
		Where("published_at IS NOT NULL AND published_at <= NOW()").
		Where("(expires_at IS NULL OR expires_at > NOW())")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("is_pinned DESC, published_at DESC").Find(&announcements).Error; err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

// GetAnnouncementByID returns an announcement by ID
func (s *AnnouncementService) GetAnnouncementByID(ctx context.Context, id string) (*models.Announcement, error) {
	var announcement models.Announcement
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid announcement ID")
	}

	if err := s.db.First(&announcement, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("announcement not found")
		}
		return nil, err
	}

	return &announcement, nil
}

// CreateAnnouncement creates a new announcement
func (s *AnnouncementService) CreateAnnouncement(ctx context.Context, announcement *models.Announcement) error {
	announcement.ID = uuid.New()
	return s.db.Create(announcement).Error
}

// UpdateAnnouncement updates an announcement
func (s *AnnouncementService) UpdateAnnouncement(ctx context.Context, id string, updates map[string]interface{}) (*models.Announcement, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid announcement ID")
	}

	result := s.db.Model(&models.Announcement{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("announcement not found")
	}

	return s.GetAnnouncementByID(ctx, id)
}

// DeleteAnnouncement deletes an announcement
func (s *AnnouncementService) DeleteAnnouncement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid announcement ID")
	}

	result := s.db.Where("id = ?", uid).Delete(&models.Announcement{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("announcement not found")
	}

	return nil
}
