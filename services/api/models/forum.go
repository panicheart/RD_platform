package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ForumBoard represents a forum board/section
type ForumBoard struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Category    string     `json:"category"` // tech, general, help, etc.
	Icon        string     `json:"icon"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	TopicCount  int        `json:"topic_count" gorm:"default:0"`
	PostCount   int        `json:"post_count" gorm:"default:0"`
	LastPostAt  *time.Time `json:"last_post_at"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (b *ForumBoard) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = ulid.Make().String()
	}
	return nil
}

func (ForumBoard) TableName() string {
	return "forum_boards"
}

// ForumPost represents a forum post/topic
type ForumPost struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	BoardID      string     `json:"board_id" gorm:"index"`
	Title        string     `json:"title" gorm:"not null"`
	Content      string     `json:"content" gorm:"type:text"`
	AuthorID     string     `json:"author_id" gorm:"index"`
	AuthorName   string     `json:"author_name"`
	ViewCount    int        `json:"view_count" gorm:"default:0"`
	ReplyCount   int        `json:"reply_count" gorm:"default:0"`
	IsPinned     bool       `json:"is_pinned" gorm:"default:false"`
	IsLocked     bool       `json:"is_locked" gorm:"default:false"`
	IsBestAnswer bool       `json:"is_best_answer" gorm:"default:false"`
	Tags         string     `json:"tags" gorm:"type:text"` // JSON array
	KnowledgeID  *string    `json:"knowledge_id"`          // Linked knowledge entry
	LastReplyAt  *time.Time `json:"last_reply_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (p *ForumPost) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = ulid.Make().String()
	}
	return nil
}

func (ForumPost) TableName() string {
	return "forum_posts"
}

// ForumReply represents a reply to a forum post
type ForumReply struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	PostID       string    `json:"post_id" gorm:"index"`
	ParentID     *string   `json:"parent_id" gorm:"index"` // For nested replies
	Content      string    `json:"content" gorm:"type:text"`
	AuthorID     string    `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	IsBestAnswer bool      `json:"is_best_answer" gorm:"default:false"`
	Mentions     string    `json:"mentions" gorm:"type:text"` // JSON array of user IDs
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (r *ForumReply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ulid.Make().String()
	}
	return nil
}

func (ForumReply) TableName() string {
	return "forum_replies"
}

// ForumTag represents a forum tag
type ForumTag struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Color     string    `json:"color" gorm:"default:'#1890ff'"`
	Count     int       `json:"count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *ForumTag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ulid.Make().String()
	}
	return nil
}

func (ForumTag) TableName() string {
	return "forum_tags"
}
