package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"rdp-platform/rdp-api/clients"
	"rdp-platform/rdp-api/models"
)

// ZoteroService provides business logic for Zotero integration
type ZoteroService struct {
	db *gorm.DB
}

// ZoteroConnection stores user Zotero API credentials
type ZoteroConnection struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	UserID       string     `json:"user_id" gorm:"uniqueIndex;not null"`
	APIKey       string     `json:"api_key" gorm:"not null"` // Encrypted
	ZoteroUserID string     `json:"zotero_user_id" gorm:"not null"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastSyncAt   *time.Time `json:"last_sync_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName specifies the table name for ZoteroConnection
func (ZoteroConnection) TableName() string {
	return "zotero_connections"
}

// ZoteroSyncResult represents the result of a sync operation
type ZoteroSyncResult struct {
	Created int      `json:"created"`
	Updated int      `json:"updated"`
	Deleted int      `json:"deleted"`
	Errors  []string `json:"errors,omitempty"`
}

// NewZoteroService creates a new Zotero service
func NewZoteroService(db *gorm.DB) *ZoteroService {
	return &ZoteroService{db: db}
}

// ==================== Connection Management ====================

// SaveConnection saves or updates Zotero connection for a user
func (s *ZoteroService) SaveConnection(userID, apiKey, zoteroUserID string) (*ZoteroConnection, error) {
	// Validate credentials by testing connection
	client, err := clients.NewZoteroClient(clients.ZoteroConfig{
		APIKey: apiKey,
		UserID: zoteroUserID,
	})
	if err != nil {
		return nil, err
	}

	if err := client.TestConnection(); err != nil {
		return nil, fmt.Errorf("failed to validate Zotero credentials: %w", err)
	}

	// Find existing connection
	var connection ZoteroConnection
	result := s.db.Where("user_id = ?", userID).First(&connection)

	if result.Error == nil {
		// Update existing
		connection.APIKey = apiKey
		connection.ZoteroUserID = zoteroUserID
		connection.IsActive = true
		if err := s.db.Save(&connection).Error; err != nil {
			return nil, fmt.Errorf("failed to update connection: %w", err)
		}
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new
		connection = ZoteroConnection{
			UserID:       userID,
			APIKey:       apiKey,
			ZoteroUserID: zoteroUserID,
			IsActive:     true,
		}
		if err := s.db.Create(&connection).Error; err != nil {
			return nil, fmt.Errorf("failed to create connection: %w", err)
		}
	} else {
		return nil, result.Error
	}

	return &connection, nil
}

// GetConnection retrieves Zotero connection for a user
func (s *ZoteroService) GetConnection(userID string) (*ZoteroConnection, error) {
	var connection ZoteroConnection
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).First(&connection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Zotero connection not found")
		}
		return nil, err
	}
	return &connection, nil
}

// DeleteConnection removes Zotero connection for a user
func (s *ZoteroService) DeleteConnection(userID string) error {
	result := s.db.Where("user_id = ?", userID).Delete(&ZoteroConnection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Zotero connection not found")
	}
	return nil
}

// TestConnection tests Zotero connection without saving
func (s *ZoteroService) TestConnection(apiKey, zoteroUserID string) error {
	client, err := clients.NewZoteroClient(clients.ZoteroConfig{
		APIKey: apiKey,
		UserID: zoteroUserID,
	})
	if err != nil {
		return err
	}
	return client.TestConnection()
}

// ==================== Item Sync Operations ====================

// SyncItems performs a full or incremental sync of Zotero items
func (s *ZoteroService) SyncItems(userID string, incremental bool) (*ZoteroSyncResult, error) {
	connection, err := s.GetConnection(userID)
	if err != nil {
		return nil, err
	}

	client, err := clients.NewZoteroClient(clients.ZoteroConfig{
		APIKey: connection.APIKey,
		UserID: connection.ZoteroUserID,
	})
	if err != nil {
		return nil, err
	}

	options := &clients.ItemQueryOptions{
		Limit:     100,
		Sort:      "dateModified",
		Direction: "desc",
	}

	if incremental && connection.LastSyncAt != nil {
		options.Since = connection.LastSyncAt.Unix()
	}

	result := &ZoteroSyncResult{
		Created: 0,
		Updated: 0,
		Errors:  []string{},
	}

	// Fetch all items with pagination
	offset := 0
	for {
		options.Start = offset
		items, _, err := client.GetItems(options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch items at offset %d: %v", offset, err))
			break
		}

		for _, item := range items {
			if err := s.syncItem(item, client); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to sync item %s: %v", item.Key, err))
			} else {
				// Check if item existed
				var existing models.ZoteroItem
				if err := s.db.Where("zotero_key = ?", item.Key).First(&existing).Error; err == nil {
					result.Updated++
				} else {
					result.Created++
				}
			}
		}

		if len(items) < options.Limit {
			break
		}
		offset += len(items)

		// Safety limit to prevent infinite loops
		if offset >= 10000 {
			result.Errors = append(result.Errors, "Sync stopped at 10000 items to prevent timeout")
			break
		}
	}

	// Update last sync time
	now := time.Now()
	connection.LastSyncAt = &now
	s.db.Save(connection)

	return result, nil
}

// syncItem synchronizes a single Zotero item to the database
func (s *ZoteroService) syncItem(item clients.ZoteroItem, client *clients.ZoteroClient) error {
	// Convert creators to JSON
	authorsJSON, err := json.Marshal(item.Creators)
	if err != nil {
		return fmt.Errorf("failed to marshal authors: %w", err)
	}

	// Convert tags to JSON
	tagsJSON, err := json.Marshal(item.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Find PDF attachment
	pdfPath := ""
	children, _ := client.GetItemChildren(item.Key)
	for _, child := range children {
		if child.ItemType == "attachment" && child.ContentType == "application/pdf" {
			pdfPath = child.Key
			break
		}
	}

	// Check if item exists
	var existing models.ZoteroItem
	err = s.db.Where("zotero_key = ?", item.Key).First(&existing).Error

	zoteroItem := models.ZoteroItem{
		ZoteroKey:   item.Key,
		Title:       item.Title,
		ItemType:    item.ItemType,
		Authors:     string(authorsJSON),
		Abstract:    item.AbstractNote,
		Publication: item.Publication,
		Volume:      item.Volume,
		Issue:       item.Issue,
		Pages:       item.Pages,
		Date:        item.Date,
		DOI:         item.DOI,
		URL:         item.URL,
		PDFPath:     pdfPath,
		Tags:        string(tagsJSON),
		SyncedAt:    time.Now(),
	}

	if err == nil {
		// Update existing
		zoteroItem.ID = existing.ID
		return s.db.Save(&zoteroItem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new
		return s.db.Create(&zoteroItem).Error
	}

	return err
}

// ==================== Item Query Operations ====================

// ListItemsQuery represents query parameters for listing Zotero items
type ListItemsQuery struct {
	ItemType string
	Tag      string
	Search   string
	Page     int
	PageSize int
}

// ListItems returns paginated Zotero items
func (s *ZoteroService) ListItems(query ListItemsQuery) ([]models.ZoteroItem, int64, error) {
	var items []models.ZoteroItem
	var total int64

	db := s.db.Model(&models.ZoteroItem{})

	// Apply filters
	if query.ItemType != "" {
		db = db.Where("item_type = ?", query.ItemType)
	}
	if query.Tag != "" {
		db = db.Where("tags LIKE ?", fmt.Sprintf("%%%s%%", query.Tag))
	}
	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("title LIKE ? OR authors LIKE ? OR abstract LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	offset := (query.Page - 1) * query.PageSize

	if err := db.Order("synced_at DESC").Offset(offset).Limit(query.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetItemByID retrieves a Zotero item by its internal ID
func (s *ZoteroService) GetItemByID(id string) (*models.ZoteroItem, error) {
	var item models.ZoteroItem
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// GetItemByZoteroKey retrieves a Zotero item by its Zotero key
func (s *ZoteroService) GetItemByZoteroKey(key string) (*models.ZoteroItem, error) {
	var item models.ZoteroItem
	if err := s.db.First(&item, "zotero_key = ?", key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// DeleteItem deletes a Zotero item from the database
func (s *ZoteroService) DeleteItem(id string) error {
	var item models.ZoteroItem
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item not found")
		}
		return err
	}

	return s.db.Delete(&item).Error
}

// ==================== PDF Operations ====================

// GetPDFURL returns the URL for viewing a PDF
func (s *ZoteroService) GetPDFURL(userID, itemID string) (string, error) {
	connection, err := s.GetConnection(userID)
	if err != nil {
		return "", err
	}

	item, err := s.GetItemByID(itemID)
	if err != nil {
		return "", err
	}

	if item.PDFPath == "" {
		return "", errors.New("no PDF attachment found for this item")
	}

	client, err := clients.NewZoteroClient(clients.ZoteroConfig{
		APIKey: connection.APIKey,
		UserID: connection.ZoteroUserID,
	})
	if err != nil {
		return "", err
	}

	return client.GetPDFURL(item.PDFPath)
}

// ==================== Citation Operations ====================

// GenerateCitation generates a citation in the specified format
func (s *ZoteroService) GenerateCitation(itemID string, format string) (string, error) {
	item, err := s.GetItemByID(itemID)
	if err != nil {
		return "", err
	}

	// Parse authors from JSON
	var creators []clients.ZoteroCreator
	if err := json.Unmarshal([]byte(item.Authors), &creators); err != nil {
		// If unmarshal fails, use empty authors
		creators = []clients.ZoteroCreator{}
	}

	switch format {
	case "gb7714", "gb/t 7714", "GB/T 7714":
		return s.generateGB7714Citation(item, creators)
	case "apa":
		return s.generateAPACitation(item, creators)
	case "mla":
		return s.generateMLACitation(item, creators)
	default:
		return s.generateGB7714Citation(item, creators) // Default to GB/T 7714
	}
}

// generateGB7714Citation generates GB/T 7714-2015 format citation
func (s *ZoteroService) generateGB7714Citation(item *models.ZoteroItem, creators []clients.ZoteroCreator) (string, error) {
	var citation string

	// Format authors
	if len(creators) > 0 {
		for i, creator := range creators {
			if i > 0 {
				citation += ", "
			}
			if creator.Name != "" {
				// Institutional author
				citation += creator.Name
			} else {
				citation += creator.LastName + " " + creator.FirstName
			}
		}
		citation += ". "
	}

	// Title
	citation += item.Title

	// Publication info
	if item.Publication != "" {
		citation += "[J]. " + item.Publication
		if item.Volume != "" {
			citation += ", " + item.Volume
			if item.Issue != "" {
				citation += "(" + item.Issue + ")"
			}
		}
		if item.Pages != "" {
			citation += ": " + item.Pages
		}
		citation += "."
	}

	// Date
	if item.Date != "" {
		citation += " " + item.Date + "."
	}

	// DOI
	if item.DOI != "" {
		citation += " DOI:" + item.DOI + "."
	}

	return citation, nil
}

// generateAPACitation generates APA format citation
func (s *ZoteroService) generateAPACitation(item *models.ZoteroItem, creators []clients.ZoteroCreator) (string, error) {
	var citation string

	// Format authors
	if len(creators) > 0 {
		if len(creators) == 1 {
			citation += creators[0].LastName + ", " + string(creators[0].FirstName[0]) + "."
		} else if len(creators) == 2 {
			citation += creators[0].LastName + ", " + string(creators[0].FirstName[0]) + "., & " +
				creators[1].LastName + ", " + string(creators[1].FirstName[0]) + "."
		} else if len(creators) > 2 {
			citation += creators[0].LastName + ", " + string(creators[0].FirstName[0]) + ". et al."
		}
		citation += " "
	}

	// Year
	if item.Date != "" {
		// Extract year from date
		year := item.Date
		if len(year) > 4 {
			year = year[:4]
		}
		citation += "(" + year + "). "
	}

	// Title
	citation += item.Title + ". "

	// Publication
	if item.Publication != "" {
		citation += "*" + item.Publication + "*"
		if item.Volume != "" {
			citation += ", " + item.Volume
			if item.Issue != "" {
				citation += "(" + item.Issue + ")"
			}
		}
		if item.Pages != "" {
			citation += ", " + item.Pages
		}
		citation += "."
	}

	// DOI
	if item.DOI != "" {
		citation += " https://doi.org/" + item.DOI
	}

	return citation, nil
}

// generateMLACitation generates MLA format citation
func (s *ZoteroService) generateMLACitation(item *models.ZoteroItem, creators []clients.ZoteroCreator) (string, error) {
	var citation string

	// Format authors
	if len(creators) > 0 {
		citation += creators[0].LastName + ", " + creators[0].FirstName
		if len(creators) > 1 {
			citation += ", et al"
		}
		citation += ". "
	}

	// Title
	citation += "\"" + item.Title + ".\" "

	// Publication
	if item.Publication != "" {
		citation += "*" + item.Publication + "*"
		if item.Volume != "" {
			citation += " " + item.Volume
			if item.Issue != "" {
				citation += "." + item.Issue
			}
		}
		if item.Date != "" {
			citation += " (" + item.Date + ")"
		}
		if item.Pages != "" {
			citation += ": " + item.Pages
		}
		citation += "."
	}

	return citation, nil
}

// ==================== Collection Operations ====================

// GetCollections retrieves collections from Zotero
func (s *ZoteroService) GetCollections(userID string) ([]clients.ZoteroCollection, error) {
	connection, err := s.GetConnection(userID)
	if err != nil {
		return nil, err
	}

	client, err := clients.NewZoteroClient(clients.ZoteroConfig{
		APIKey: connection.APIKey,
		UserID: connection.ZoteroUserID,
	})
	if err != nil {
		return nil, err
	}

	return client.GetCollections()
}
