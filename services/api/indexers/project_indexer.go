package indexers

import (
	"encoding/json"
	"fmt"

	"rdp-platform/rdp-api/clients"
	"rdp-platform/rdp-api/models"
	"gorm.io/gorm"
)

// ProjectIndexer handles indexing projects to MeiliSearch
type ProjectIndexer struct {
	db          *gorm.DB
	meiliClient *clients.MeiliSearchClient
}

// NewProjectIndexer creates a new project indexer
func NewProjectIndexer(db *gorm.DB, meiliClient *clients.MeiliSearchClient) *ProjectIndexer {
	return &ProjectIndexer{
		db:          db,
		meiliClient: meiliClient,
	}
}

// IndexableProject represents a project document for indexing
type IndexableProject struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	LeaderID    string `json:"leader_id"`
	LeaderName  string `json:"leader_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// IndexAll indexes all projects
func (i *ProjectIndexer) IndexAll() error {
	var projects []models.Project
	if err := i.db.Find(&projects).Error; err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	docs := make([]IndexableProject, 0, len(projects))
	for _, p := range projects {
		docs = append(docs, i.toIndexableDocument(p))
	}

	if len(docs) > 0 {
		return i.meiliClient.Index(clients.IndexNames.Projects, docs)
	}
	
	return nil
}

// IndexSingle indexes a single project
func (i *ProjectIndexer) IndexSingle(project models.Project) error {
	doc := i.toIndexableDocument(project)
	return i.meiliClient.Index(clients.IndexNames.Projects, []IndexableProject{doc})
}

// Delete removes a project from the index
func (i *ProjectIndexer) Delete(projectID string) error {
	return i.meiliClient.Delete(clients.IndexNames.Projects, []string{projectID})
}

func (i *ProjectIndexer) toIndexableDocument(p models.Project) IndexableProject {
	// Fetch leader name
	var leaderName string
	var leader models.User
	if err := i.db.First(&leader, "id = ?", p.LeaderID).Error; err == nil {
		leaderName = leader.DisplayName
	}

	return IndexableProject{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		Status:      p.Status,
		Category:    p.Category,
		LeaderID:    p.LeaderID,
		LeaderName:  leaderName,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
