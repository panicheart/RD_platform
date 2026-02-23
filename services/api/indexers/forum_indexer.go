package indexers

import (
	"fmt"

	"rdp-platform/rdp-api/clients"
	"rdp-platform/rdp-api/models"
	"gorm.io/gorm"
)

// ForumIndexer handles indexing forum posts to MeiliSearch
type ForumIndexer struct {
	db          *gorm.DB
	meiliClient *clients.MeiliSearchClient
}

// NewForumIndexer creates a new forum indexer
func NewForumIndexer(db *gorm.DB, meiliClient *clients.MeiliSearchClient) *ForumIndexer {
	return &ForumIndexer{
		db:          db,
		meiliClient: meiliClient,
	}
}

// IndexableForumPost represents a forum post document for indexing
type IndexableForumPost struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	BoardID       string   `json:"board_id"`
	BoardName     string   `json:"board_name"`
	AuthorID      string   `json:"author_id"`
	AuthorName    string   `json:"author_name"`
	Tags          []string `json:"tags"`
	ReplyCount    int      `json:"reply_count"`
	ViewCount     int      `json:"view_count"`
	IsPinned      bool     `json:"is_pinned"`
	IsBestAnswer  bool     `json:"is_best_answer"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// IndexAll indexes all forum posts
func (i *ForumIndexer) IndexAll() error {
	var posts []models.ForumPost
	if err := i.db.Preload("Tags").Find(&posts).Error; err != nil {
		return fmt.Errorf("failed to fetch forum posts: %w", err)
	}

	docs := make([]IndexableForumPost, 0, len(posts))
	for _, p := range posts {
		docs = append(docs, i.toIndexableDocument(p))
	}

	if len(docs) > 0 {
		return i.meiliClient.Index(clients.IndexNames.ForumPosts, docs)
	}
	
	return nil
}

// IndexSingle indexes a single forum post
func (i *ForumIndexer) IndexSingle(post models.ForumPost) error {
	i.db.Preload("Tags").First(&post, "id = ?", post.ID)
	
	doc := i.toIndexableDocument(post)
	return i.meiliClient.Index(clients.IndexNames.ForumPosts, []IndexableForumPost{doc})
}

// Delete removes a forum post from the index
func (i *ForumIndexer) Delete(postID string) error {
	return i.meiliClient.Delete(clients.IndexNames.ForumPosts, []string{postID})
}

func (i *ForumIndexer) toIndexableDocument(p models.ForumPost) IndexableForumPost {
	// Fetch board name
	var boardName string
	var board models.ForumBoard
	if err := i.db.First(&board, "id = ?", p.BoardID).Error; err == nil {
		boardName = board.Name
	}

	// Fetch author name
	var authorName string
	var author models.User
	if err := i.db.First(&author, "id = ?", p.AuthorID).Error; err == nil {
		authorName = author.DisplayName
	}

	// Extract tags
	tags := make([]string, 0, len(p.Tags))
	for _, tag := range p.Tags {
		tags = append(tags, tag.Name)
	}

	return IndexableForumPost{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		BoardID:      p.BoardID,
		BoardName:    boardName,
		AuthorID:     p.AuthorID,
		AuthorName:   authorName,
		Tags:         tags,
		ReplyCount:   p.ReplyCount,
		ViewCount:    p.ViewCount,
		IsPinned:     p.IsPinned,
		IsBestAnswer: p.IsBestAnswer,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
