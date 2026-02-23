package clients

import (
	"fmt"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// MeiliSearchClient wraps the MeiliSearch SDK
type MeiliSearchClient struct {
	client meilisearch.ServiceManager
	config MeiliSearchConfig
}

// MeiliSearchConfig holds configuration for MeiliSearch
type MeiliSearchConfig struct {
	Host   string
	APIKey string
}

// NewMeiliSearchClient creates a new MeiliSearch client
func NewMeiliSearchClient(config MeiliSearchConfig) (*MeiliSearchClient, error) {
	client := meilisearch.New(config.Host, meilisearch.WithAPIKey(config.APIKey))
	
	// Test connection
	if _, err := client.Health(); err != nil {
		return nil, fmt.Errorf("failed to connect to MeiliSearch: %w", err)
	}
	
	return &MeiliSearchClient{
		client: client,
		config: config,
	}, nil
}

// IndexNames contains the names of all search indexes
var IndexNames = struct {
	Projects   string
	Knowledge  string
	Products   string
	ForumPosts string
}{
	Projects:   "projects",
	Knowledge:  "knowledge",
	Products:   "products",
	ForumPosts: "forum_posts",
}

// InitializeIndexes creates all required indexes with settings
func (c *MeiliSearchClient) InitializeIndexes() error {
	indexes := []struct {
		name    string
		options *meilisearch.IndexConfig
	}{
		{
			name: IndexNames.Projects,
			options: &meilisearch.IndexConfig{
				Uid:        IndexNames.Projects,
				PrimaryKey: "id",
			},
		},
		{
			name: IndexNames.Knowledge,
			options: &meilisearch.IndexConfig{
				Uid:        IndexNames.Knowledge,
				PrimaryKey: "id",
			},
		},
		{
			name: IndexNames.Products,
			options: &meilisearch.IndexConfig{
				Uid:        IndexNames.Products,
				PrimaryKey: "id",
			},
		},
		{
			name: IndexNames.ForumPosts,
			options: &meilisearch.IndexConfig{
				Uid:        IndexNames.ForumPosts,
				PrimaryKey: "id",
			},
		},
	}

	for _, idx := range indexes {
		if err := c.createIndexWithSettings(idx.name, idx.options); err != nil {
			return err
		}
	}

	return nil
}

func (c *MeiliSearchClient) createIndexWithSettings(name string, config *meilisearch.IndexConfig) error {
	// Check if index exists
	_, err := c.client.GetIndex(name)
	if err == nil {
		// Index already exists
		return nil
	}

	// Create index
	task, err := c.client.CreateIndex(config)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", name, err)
	}

	// Wait for task completion
	if err := c.waitForTask(task.TaskUID); err != nil {
		return err
	}

	// Configure index settings
	if err := c.configureIndex(name); err != nil {
		return err
	}

	return nil
}

func (c *MeiliSearchClient) configureIndex(name string) error {
	index := c.client.Index(name)

	var settings *meilisearch.Settings

	switch name {
	case IndexNames.Projects:
		settings = &meilisearch.Settings{
			SearchableAttributes: []string{"name", "description", "code"},
			FilterableAttributes: []string{"status", "category", "leader_id"},
			SortableAttributes:   []string{"created_at", "updated_at"},
		}
	case IndexNames.Knowledge:
		settings = &meilisearch.Settings{
			SearchableAttributes: []string{"title", "content", "tags"},
			FilterableAttributes: []string{"category_id", "status", "source", "tags"},
			SortableAttributes:   []string{"created_at", "updated_at", "view_count"},
		}
	case IndexNames.Products:
		settings = &meilisearch.Settings{
			SearchableAttributes: []string{"name", "description", "type"},
			FilterableAttributes: []string{"type", "maturity", "owner_id"},
			SortableAttributes:   []string{"created_at"},
		}
	case IndexNames.ForumPosts:
		settings = &meilisearch.Settings{
			SearchableAttributes: []string{"title", "content", "tags"},
			FilterableAttributes: []string{"board_id", "author_id", "is_best_answer"},
			SortableAttributes:   []string{"created_at", "reply_count", "view_count"},
		}
	}

	if settings != nil {
		task, err := index.UpdateSettings(settings)
		if err != nil {
			return fmt.Errorf("failed to update settings for %s: %w", name, err)
		}
		return c.waitForTask(task.TaskUID)
	}

	return nil
}

func (c *MeiliSearchClient) waitForTask(taskUID int64) error {
	for {
		task, err := c.client.GetTask(taskUID)
		if err != nil {
			return err
		}

		if task.Status == meilisearch.TaskStatusSucceeded {
			return nil
		}

		if task.Status == meilisearch.TaskStatusFailed {
			return fmt.Errorf("task failed: %s", task.Error)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// Index adds or updates documents in an index
func (c *MeiliSearchClient) Index(indexName string, documents interface{}) error {
	index := c.client.Index(indexName)
	task, err := index.AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to index documents: %w", err)
	}
	return c.waitForTask(task.TaskUID)
}

// Delete removes documents from an index
func (c *MeiliSearchClient) Delete(indexName string, documentIDs []string) error {
	index := c.client.Index(indexName)
	task, err := index.DeleteDocuments(documentIDs, nil)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	return c.waitForTask(task.TaskUID)
}

// Search performs a search query
func (c *MeiliSearchClient) Search(indexName string, query string, options *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	index := c.client.Index(indexName)
	
	if options == nil {
		options = &meilisearch.SearchRequest{}
	}
	
	result, err := index.Search(query, options)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	
	return result, nil
}

// MultiSearch performs searches across multiple indexes
func (c *MeiliSearchClient) MultiSearch(queries *meilisearch.MultiSearchRequest) (*meilisearch.MultiSearchResponse, error) {
	result, err := c.client.MultiSearch(queries)
	if err != nil {
		return nil, fmt.Errorf("multi-search failed: %w", err)
	}
	return result, nil
}

// GetStats returns index statistics
func (c *MeiliSearchClient) GetStats(indexName string) (*meilisearch.StatsIndex, error) {
	return c.client.Index(indexName).GetStats()
}

// Health checks if MeiliSearch is healthy
func (c *MeiliSearchClient) Health() (*meilisearch.Health, error) {
	return c.client.Health()
}

// GetClient returns the underlying MeiliSearch client
func (c *MeiliSearchClient) GetClient() meilisearch.ServiceManager {
	return c.client
}
