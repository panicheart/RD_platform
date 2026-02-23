package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"rdp-platform/rdp-api/models"
)

type ForumService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

func NewForumService(db *gorm.DB, notificationService *NotificationService) *ForumService {
	return &ForumService{
		db:                  db,
		notificationService: notificationService,
	}
}

type ListBoardsQuery struct {
	Category string
	Page     int
	PageSize int
}

func (s *ForumService) ListBoards(query ListBoardsQuery) ([]models.ForumBoard, int64, error) {
	var boards []models.ForumBoard
	var total int64

	db := s.db.Model(&models.ForumBoard{}).Where("is_active = ?", true)

	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (s *ForumService) GetBoardByID(id string) (*models.ForumBoard, error) {
	var board models.ForumBoard
	if err := s.db.First(&board, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}
	return &board, nil
}

type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

func (s *ForumService) CreateBoard(req CreateBoardRequest, createdBy string) (*models.ForumBoard, error) {
	board := &models.ForumBoard{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if err := s.db.Create(board).Error; err != nil {
		return nil, err
	}

	return board, nil
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	IsActive    *bool  `json:"is_active"`
}

func (s *ForumService) UpdateBoard(id string, req UpdateBoardRequest) (*models.ForumBoard, error) {
	var board models.ForumBoard
	if err := s.db.First(&board, "id = ?", id).Error; err != nil {
		return nil, errors.New("board not found")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	updates["sort_order"] = req.SortOrder
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.db.Model(&board).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &board, nil
}

func (s *ForumService) DeleteBoard(id string) error {
	var board models.ForumBoard
	if err := s.db.First(&board, "id = ?", id).Error; err != nil {
		return errors.New("board not found")
	}

	var postCount int64
	if err := s.db.Model(&models.ForumPost{}).Where("board_id = ?", id).Count(&postCount).Error; err != nil {
		return err
	}

	if postCount > 0 {
		return s.db.Model(&board).Update("is_active", false).Error
	}

	return s.db.Delete(&board).Error
}

type ListPostsQuery struct {
	BoardID  string
	AuthorID string
	Tag      string
	Search   string
	IsPinned *bool
	Page     int
	PageSize int
}

func (s *ForumService) ListPosts(query ListPostsQuery) ([]models.ForumPost, int64, error) {
	var posts []models.ForumPost
	var total int64

	db := s.db.Model(&models.ForumPost{})

	if query.BoardID != "" {
		db = db.Where("board_id = ?", query.BoardID)
	}
	if query.AuthorID != "" {
		db = db.Where("author_id = ?", query.AuthorID)
	}
	if query.IsPinned != nil {
		db = db.Where("is_pinned = ?", *query.IsPinned)
	}
	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("is_pinned DESC, last_reply_at DESC NULLS LAST, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *ForumService) GetPostByID(id string) (*models.ForumPost, error) {
	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	go s.incrementPostViewCount(id)

	return &post, nil
}

func (s *ForumService) incrementPostViewCount(id string) {
	s.db.Model(&models.ForumPost{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
}

type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	BoardID     string   `json:"board_id" binding:"required"`
	Tags        []string `json:"tags"`
	KnowledgeID *string  `json:"knowledge_id"`
}

func (s *ForumService) CreatePost(req CreatePostRequest, authorID, authorName string) (*models.ForumPost, error) {
	var board models.ForumBoard
	if err := s.db.First(&board, "id = ? AND is_active = ?", req.BoardID, true).Error; err != nil {
		return nil, errors.New("board not found or inactive")
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	post := &models.ForumPost{
		Title:       req.Title,
		Content:     req.Content,
		BoardID:     req.BoardID,
		AuthorID:    authorID,
		AuthorName:  authorName,
		Tags:        string(tagsJSON),
		KnowledgeID: req.KnowledgeID,
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}

	s.db.Model(&board).Updates(map[string]interface{}{
		"topic_count":  gorm.Expr("topic_count + 1"),
		"last_post_at": time.Now(),
	})

	mentions := s.extractMentions(req.Content)
	if len(mentions) > 0 {
		go s.notifyMentions(mentions, authorName, post.ID, post.Title)
	}

	return post, nil
}

type UpdatePostRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	KnowledgeID *string  `json:"knowledge_id"`
}

func (s *ForumService) UpdatePost(id string, authorID string, isAdmin bool, req UpdatePostRequest) (*models.ForumPost, error) {
	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, errors.New("post not found")
	}

	if post.AuthorID != authorID && !isAdmin {
		return nil, errors.New("permission denied")
	}

	if post.IsLocked && !isAdmin {
		return nil, errors.New("post is locked")
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
		mentions := s.extractMentions(req.Content)
		if len(mentions) > 0 {
			tagsJSON, _ := json.Marshal(mentions)
			updates["tags"] = string(tagsJSON)
		}
	}
	if req.KnowledgeID != nil {
		updates["knowledge_id"] = *req.KnowledgeID
	}

	updates["updated_at"] = time.Now()

	if err := s.db.Model(&post).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.db.First(&post, "id = ?", id)

	return &post, nil
}

func (s *ForumService) DeletePost(id string, authorID string, isAdmin bool) error {
	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", id).Error; err != nil {
		return errors.New("post not found")
	}

	if post.AuthorID != authorID && !isAdmin {
		return errors.New("permission denied")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", id).Delete(&models.ForumReply{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		return tx.Model(&models.ForumBoard{}).Where("id = ?", post.BoardID).
			Updates(map[string]interface{}{
				"topic_count": gorm.Expr("topic_count - 1"),
			}).Error
	})
}

func (s *ForumService) PinPost(id string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("permission denied")
	}

	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", id).Error; err != nil {
		return errors.New("post not found")
	}

	return s.db.Model(&post).Update("is_pinned", !post.IsPinned).Error
}

func (s *ForumService) LockPost(id string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("permission denied")
	}

	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", id).Error; err != nil {
		return errors.New("post not found")
	}

	return s.db.Model(&post).Update("is_locked", !post.IsLocked).Error
}

func (s *ForumService) MarkBestAnswer(postID string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("permission denied")
	}

	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", postID).Error; err != nil {
		return errors.New("post not found")
	}

	return s.db.Model(&post).Update("is_best_answer", true).Error
}

type ListRepliesQuery struct {
	PostID   string
	ParentID *string
	Page     int
	PageSize int
}

func (s *ForumService) ListReplies(query ListRepliesQuery) ([]models.ForumReply, int64, error) {
	var replies []models.ForumReply
	var total int64

	db := s.db.Model(&models.ForumReply{}).Where("post_id = ?", query.PostID)

	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	} else {
		db = db.Where("parent_id IS NULL")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at ASC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	return replies, total, nil
}

func (s *ForumService) GetReplyByID(id string) (*models.ForumReply, error) {
	var reply models.ForumReply
	if err := s.db.First(&reply, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reply not found")
		}
		return nil, err
	}
	return &reply, nil
}

type CreateReplyRequest struct {
	Content  string  `json:"content" binding:"required"`
	ParentID *string `json:"parent_id"`
}

func (s *ForumService) CreateReply(postID string, req CreateReplyRequest, authorID, authorName string) (*models.ForumReply, error) {
	var post models.ForumPost
	if err := s.db.First(&post, "id = ?", postID).Error; err != nil {
		return nil, errors.New("post not found")
	}

	if post.IsLocked {
		return nil, errors.New("post is locked")
	}

	if req.ParentID != nil {
		var parentReply models.ForumReply
		if err := s.db.First(&parentReply, "id = ?", *req.ParentID).Error; err != nil {
			return nil, errors.New("parent reply not found")
		}
		if parentReply.PostID != postID {
			return nil, errors.New("parent reply does not belong to this post")
		}
	}

	mentions := s.extractMentions(req.Content)
	mentionsJSON, _ := json.Marshal(mentions)

	reply := &models.ForumReply{
		PostID:     postID,
		ParentID:   req.ParentID,
		Content:    req.Content,
		AuthorID:   authorID,
		AuthorName: authorName,
		Mentions:   string(mentionsJSON),
	}

	if err := s.db.Create(reply).Error; err != nil {
		return nil, err
	}

	s.db.Model(&post).Updates(map[string]interface{}{
		"reply_count":    gorm.Expr("reply_count + 1"),
		"last_reply_at":  time.Now(),
	})

	if len(mentions) > 0 {
		go s.notifyMentions(mentions, authorName, postID, post.Title)
	}

	return reply, nil
}

type UpdateReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

func (s *ForumService) UpdateReply(id string, authorID string, isAdmin bool, req UpdateReplyRequest) (*models.ForumReply, error) {
	var reply models.ForumReply
	if err := s.db.First(&reply, "id = ?", id).Error; err != nil {
		return nil, errors.New("reply not found")
	}

	if reply.AuthorID != authorID && !isAdmin {
		return nil, errors.New("permission denied")
	}

	mentions := s.extractMentions(req.Content)
	mentionsJSON, _ := json.Marshal(mentions)

	updates := map[string]interface{}{
		"content":    req.Content,
		"mentions":   string(mentionsJSON),
		"updated_at": time.Now(),
	}

	if err := s.db.Model(&reply).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.db.First(&reply, "id = ?", id)

	return &reply, nil
}

func (s *ForumService) DeleteReply(id string, authorID string, isAdmin bool) error {
	var reply models.ForumReply
	if err := s.db.First(&reply, "id = ?", id).Error; err != nil {
		return errors.New("reply not found")
	}

	if reply.AuthorID != authorID && !isAdmin {
		return errors.New("permission denied")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&models.ForumReply{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&reply).Error; err != nil {
			return err
		}

		return tx.Model(&models.ForumPost{}).Where("id = ?", reply.PostID).
			Update("reply_count", gorm.Expr("reply_count - 1")).Error
	})
}

func (s *ForumService) MarkReplyBestAnswer(replyID string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("permission denied")
	}

	var reply models.ForumReply
	if err := s.db.First(&reply, "id = ?", replyID).Error; err != nil {
		return errors.New("reply not found")
	}

	return s.db.Model(&reply).Update("is_best_answer", true).Error
}

func (s *ForumService) ListTags() ([]models.ForumTag, error) {
	var tags []models.ForumTag
	if err := s.db.Order("name").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *ForumService) CreateTag(name, color string) (*models.ForumTag, error) {
	tag := &models.ForumTag{
		Name:  name,
		Color: color,
	}

	if err := s.db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *ForumService) DeleteTag(id string) error {
	return s.db.Delete(&models.ForumTag{}, "id = ?", id).Error
}

func (s *ForumService) SearchPosts(keyword string, boardID string, page, pageSize int) ([]models.ForumPost, int64, error) {
	var posts []models.ForumPost
	var total int64

	searchPattern := fmt.Sprintf("%%%s%%", keyword)
	db := s.db.Model(&models.ForumPost{}).
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern)

	if boardID != "" {
		db = db.Where("board_id = ?", boardID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *ForumService) extractMentions(content string) []string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_\u4e00-\u9fa5]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	mentions := make([]string, 0)
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			mentions = append(mentions, match[1])
			seen[match[1]] = true
		}
	}
	return mentions
}

func (s *ForumService) notifyMentions(mentions []string, authorName, postID, postTitle string) {
	if s.notificationService == nil {
		return
	}

	for _, mention := range mentions {
		var user models.User
		if err := s.db.Where("username = ?", mention).First(&user).Error; err != nil {
			continue
		}

		notification := &models.Notification{
			UserID:      user.ID,
			Type:        "mention",
			Title:       "你被提到了",
			Content:     &[]string{fmt.Sprintf("%s 在帖子《%s》中提到了你", authorName, postTitle)}[0],
			RelatedID:   &postID,
			RelatedType: &[]string{"forum_post"}[0],
		}

		ctx := context.Background()
		s.notificationService.CreateNotification(ctx, notification)
	}
}

func (s *ForumService) GetBoardStats(boardID string) (map[string]interface{}, error) {
	var topicCount, postCount int64

	if err := s.db.Model(&models.ForumPost{}).Where("board_id = ?", boardID).Count(&topicCount).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.ForumReply{}).
		Joins("JOIN forum_posts ON forum_replies.post_id = forum_posts.id").
		Where("forum_posts.board_id = ?", boardID).
		Count(&postCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"topic_count": topicCount,
		"post_count":  postCount,
		"total_count": topicCount + postCount,
	}, nil
}

func ParseTags(tagsJSON string) []string {
	if tagsJSON == "" {
		return []string{}
	}
	var tags []string
	json.Unmarshal([]byte(tagsJSON), &tags)
	return tags
}

func FormatTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(tags)
	return string(data)
}

type PaginatedResponse struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func StringToInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
