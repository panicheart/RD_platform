package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Knowledge represents a knowledge base entry
type Knowledge struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Content     string     `json:"content" gorm:"type:text"`
	CategoryID  string     `json:"category_id" gorm:"index"`
	AuthorID    string     `json:"author_id" gorm:"index"`
	Tags        []Tag      `json:"tags" gorm:"many2many:knowledge_tags;"`
	Status      string     `json:"status" gorm:"default:'draft'"` // draft, published, archived
	Version     int        `json:"version" gorm:"default:1"`
	ParentID    *string    `json:"parent_id" gorm:"index"` // For version history
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	Source      string     `json:"source"`    // obsidian, zotero, forum, manual
	SourceID    string     `json:"source_id"` // Original ID from source system
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// BeforeCreate generates ULID before inserting
func (k *Knowledge) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (Knowledge) TableName() string {
	return "knowledge"
}

// Category represents a knowledge category
type Category struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	ParentID    *string   `json:"parent_id" gorm:"index"`
	Level       int       `json:"level" gorm:"default:1"` // 1, 2, 3 for tree depth
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate generates ULID before inserting
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (Category) TableName() string {
	return "categories"
}

// Tag represents a knowledge tag
type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Color     string    `json:"color" gorm:"default:'#1890ff'"`
	Count     int       `json:"count" gorm:"default:0"` // Usage count
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate generates ULID before inserting
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (Tag) TableName() string {
	return "tags"
}

// ZoteroItem represents a Zotero library item
type ZoteroItem struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ZoteroKey   string    `json:"zotero_key" gorm:"uniqueIndex"`
	Title       string    `json:"title"`
	ItemType    string    `json:"item_type"`                // book, journalArticle, etc.
	Authors     string    `json:"authors" gorm:"type:text"` // JSON array
	Abstract    string    `json:"abstract" gorm:"type:text"`
	Publication string    `json:"publication"`
	Volume      string    `json:"volume"`
	Issue       string    `json:"issue"`
	Pages       string    `json:"pages"`
	Date        string    `json:"date"`
	DOI         string    `json:"doi"`
	URL         string    `json:"url"`
	PDFPath     string    `json:"pdf_path"`
	Tags        string    `json:"tags" gorm:"type:text"` // JSON array
	SyncedAt    time.Time `json:"synced_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (ZoteroItem) TableName() string {
	return "zotero_items"
}

// KnowledgeReview represents a knowledge review/approval record
type KnowledgeReview struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	KnowledgeID string    `json:"knowledge_id" gorm:"index"`
	ReviewerID  string    `json:"reviewer_id"`
	Status      string    `json:"status"` // pending, approved, rejected
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

// BeforeCreate generates ULID before inserting
func (kr *KnowledgeReview) BeforeCreate(tx *gorm.DB) error {
	if kr.ID == "" {
		kr.ID = ulid.Make().String()
	}
	return nil
}

// TableName specifies the table name
func (KnowledgeReview) TableName() string {
	return "knowledge_reviews"
}

// ObsidianMapping stores Obsidian vault path mappings
type ObsidianMapping struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	VaultPath  string     `json:"vault_path" gorm:"not null"`
	LocalPath  string     `json:"local_path" gorm:"not null"`
	CategoryID string     `json:"category_id"`
	AutoSync   bool       `json:"auto_sync" gorm:"default:false"`
	LastSyncAt *time.Time `json:"last_sync_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName specifies the table name
func (ObsidianMapping) TableName() string {
	return "obsidian_mappings"
}
