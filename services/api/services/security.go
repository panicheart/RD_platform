package services

import (
	"context"
	"errors"
	"time"

	"rdp/services/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SecurityService handles security-related business logic
type SecurityService struct {
	db *gorm.DB
}

// NewSecurityService creates a new SecurityService
func NewSecurityService(db *gorm.DB) *SecurityService {
	return &SecurityService{db: db}
}

// CreateAuditLog creates a new audit log entry
func (s *SecurityService) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	log.ID = uuid.New()
	log.CreatedAt = time.Now().UTC()
	return s.db.Create(log).Error
}

// LogAction logs a user action for audit
func (s *SecurityService) LogAction(ctx context.Context, userID *uuid.UUID, username *string, action, resource, resourceID string, classification string) error {
	log := &models.AuditLog{
		UserID:        userID,
		Username:      username,
		Action:        action,
		Resource:      resource,
		ResourceID:    &resourceID,
		Classification: classification,
	}

	// Get IP from context if available
	if ip, exists := ctx.Value("ip_address").(string); exists {
		log.IPAddress = &ip
	}

	return s.CreateAuditLog(ctx, log)
}

// ListAuditLogs returns audit logs with filters
func (s *SecurityService) ListAuditLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{})

	// Apply filters
	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if action, ok := filters["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}
	if resource, ok := filters["resource"].(string); ok && resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if classification, ok := filters["classification"].(string); ok && classification != "" {
		query = query.Where("classification = ?", classification)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID returns an audit log by ID
func (s *SecurityService) GetAuditLogByID(ctx context.Context, id string) (*models.AuditLog, error) {
	var log models.AuditLog
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid audit log ID")
	}

	if err := s.db.First(&log, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}

	return &log, nil
}

// CreateLoginLog creates a login attempt log
func (s *SecurityService) CreateLoginLog(ctx context.Context, log *models.LoginLog) error {
	log.ID = uuid.New()
	log.CreatedAt = time.Now().UTC()
	return s.db.Create(log).Error
}

// ListLoginLogs returns login logs with filters
func (s *SecurityService) ListLoginLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.LoginLog, int64, error) {
	var logs []models.LoginLog
	var total int64

	query := s.db.Model(&models.LoginLog{})

	// Apply filters
	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if username, ok := filters["username"].(string); ok && username != "" {
		query = query.Where("username = ?", username)
	}
	if success, ok := filters["success"].(bool); ok {
		query = query.Where("success = ?", success)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetDataClassification returns all data classifications
func (s *SecurityService) GetDataClassifications(ctx context.Context) ([]models.DataClassification, error) {
	var classifications []models.DataClassification
	if err := s.db.Order("level ASC").Find(&classifications).Error; err != nil {
		return nil, err
	}
	return classifications, nil
}

// SeedDefaultClassifications seeds default data classifications
func (s *SecurityService) SeedDefaultClassifications(ctx context.Context) error {
	classifications := []models.DataClassification{
		{Level: "public", Name: "公开", Description: "可对外公开的信息", Color: "#52c41a", Icon: "global"},
		{Level: "internal", Name: "内部", Description: "仅限内部人员访问", Color: "#1890ff", Icon: "team"},
		{Level: "confidential", Name: "机密", Description: "仅限于特定人员访问", Color: "#faad14", Icon: "lock"},
		{Level: "secret", Name: "绝密", Description: "严格限制访问范围", Color: "#f5222d", Icon: "safety-certificate"},
	}

	// Use upsert to avoid duplicates
	for _, c := range classifications {
		err := s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "level"}},
			DoNothing: true,
		}).Create(&c).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// SessionService handles session management
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new SessionService
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// CreateSession creates a new session
func (s *SessionService) CreateSession(ctx context.Context, session *models.Session) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now().UTC()
	session.LastActiveAt = session.CreatedAt
	return s.db.Create(session).Error
}

// GetSessionByToken returns a session by token
func (s *SessionService) GetSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	if err := s.db.First(&session, "token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found or expired")
		}
		return nil, err
	}
	return &session, nil
}

// RevokeSession revokes a session
func (s *SessionService) RevokeSession(ctx context.Context, token string) error {
	result := s.db.Model(&models.Session{}).Where("token = ?", token).Update("is_revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *SessionService) RevokeAllUserSessions(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return s.db.Model(&models.Session{}).Where("user_id = ?", uid).Update("is_revoked", true).Error
}

// CleanExpiredSessions removes expired sessions
func (s *SessionService) CleanExpiredSessions(ctx context.Context) error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}
