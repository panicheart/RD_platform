package indexers

import (
	"fmt"
	"strings"

	"rdp-platform/rdp-api/clients"
	"rdp-platform/rdp-api/models"
	"gorm.io/gorm"
)

// KnowledgeIndexer handles indexing knowledge items to MeiliSearch
type KnowledgeIndexer struct {
	db          *gorm.DB
	meiliClient *clients.MeiliSearchClient
}

// NewKnowledgeIndexer creates a new knowledge indexer
func NewKnowledgeIndexer(db *gorm.DB, meiliClient *clients.MeiliSearchClient) *KnowledgeIndexer {
	return &KnowledgeIndexer{
		db:          db,
		meiliClient: meiliClient,
	}
}

// IndexableKnowledge represents a knowledge document for indexing
type IndexableKnowledge struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	CategoryID  string   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	AuthorID    string   `json:"author_id"`
	AuthorName  string   `json:"author_name"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	Source      string   `json:"source"`
	ViewCount   int      `json:"view_count"`
	Version     int      `json:"version"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// IndexAll indexes all knowledge items
func (i *KnowledgeIndexer) IndexAll() error {
	var items []models.Knowledge
	if err := i.db.Preload("Tags").Find(&items).Error; err != nil {
		return fmt.Errorf("failed to fetch knowledge items: %w", err)
	}

	docs := make([]IndexableKnowledge, 0, len(items))
	for _, item := range items {
		docs = append(docs, i.toIndexableDocument(item))
	}

	if len(docs) > 0 {
		return i.meiliClient.Index(clients.IndexNames.Knowledge, docs)
	}
	
	return nil
}

// IndexSingle indexes a single knowledge item
func (i *KnowledgeIndexer) IndexSingle(item models.Knowledge) error {
	// Reload with tags
	i.db.Preload("Tags").First(&item, "id = ?", item.ID)
	
	doc := i.toIndexableDocument(item)
	return i.meiliClient.Index(clients.IndexNames.Knowledge, []IndexableKnowledge{doc})
}

// Delete removes a knowledge item from the index
func (i *KnowledgeIndexer) Delete(itemID string) error {
	return i.meiliClient.Delete(clients.IndexNames.Knowledge, []string{itemID})
}

func (i *KnowledgeIndexer) toIndexableDocument(k models.Knowledge) IndexableKnowledge {
	// Fetch category name
	var categoryName string
	var category models.Category
	if err := i.db.First(&category, "id = ?", k.CategoryID).Error; err == nil {
		categoryName = category.Name
	}

	// Fetch author name
	var authorName string
	var author models.User
	if err := i.db.First(&author, "id = ?", k.AuthorID).Error; err == nil {
		authorName = author.DisplayName
	}

	// Extract tags
	tags := make([]string, 0, len(k.Tags))
	for _, tag := range k.Tags {
		tags = append(tags, tag.Name)
	}

	// Clean content (remove markdown syntax for better search)
	content := cleanMarkdown(k.Content)

	return IndexableKnowledge{
		ID:           k.ID,
		Title:        k.Title,
		Content:      content,
		CategoryID:   k.CategoryID,
		CategoryName: categoryName,
		AuthorID:     k.AuthorID,
		AuthorName:   authorName,
		Tags:         tags,
		Status:       k.Status,
		Source:       k.Source,
		ViewCount:    k.ViewCount,
		Version:      k.Version,
		CreatedAt:    k.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    k.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func cleanMarkdown(content string) string {
	// Remove wiki links [[link]]
	content = strings.ReplaceAll(content, "[[", " ")
	content = strings.ReplaceAll(content, "]", " ")
	
	// Remove markdown headers
	content = strings.ReplaceAll(content, "#", " ")
	
	// Remove emphasis markers
	content = strings.ReplaceAll(content, "**", " ")
	content = strings.ReplaceAll(content, "*", " ")
	content = strings.ReplaceAll(content, "__", " ")
	content = strings.ReplaceAll(content, "_", " ")
	
	// Clean up extra whitespace
	content = strings.Join(strings.Fields(content), " ")
	
	return content
}
