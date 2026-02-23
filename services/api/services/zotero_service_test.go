package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"rdp-platform/rdp-api/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate tables
	err = db.AutoMigrate(
		&models.ZoteroItem{},
		&ZoteroConnection{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestZoteroService_SaveConnection(t *testing.T) {
	db := setupTestDB(t)
	service := NewZoteroService(db)

	// Note: This test will fail without valid Zotero credentials
	// In real testing, you should mock the Zotero client or use test credentials
	t.Run("save connection with empty credentials should fail", func(t *testing.T) {
		_, err := service.SaveConnection("user-123", "", "")
		assert.Error(t, err)
	})
}

func TestZoteroService_GetConnection(t *testing.T) {
	db := setupTestDB(t)
	service := NewZoteroService(db)

	t.Run("get non-existent connection should fail", func(t *testing.T) {
		_, err := service.GetConnection("user-123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("get existing connection should succeed", func(t *testing.T) {
		// Create connection directly in DB
		connection := ZoteroConnection{
			UserID:       "user-456",
			APIKey:       "test-api-key",
			ZoteroUserID: "12345678",
			IsActive:     true,
		}
		err := db.Create(&connection).Error
		assert.NoError(t, err)

		// Retrieve
		result, err := service.GetConnection("user-456")
		assert.NoError(t, err)
		assert.Equal(t, "user-456", result.UserID)
		assert.Equal(t, "12345678", result.ZoteroUserID)
	})
}

func TestZoteroService_ListItems(t *testing.T) {
	db := setupTestDB(t)
	service := NewZoteroService(db)

	// Create test items
	items := []models.ZoteroItem{
		{
			ZoteroKey:   "ITEM001",
			Title:       "Test Article 1",
			ItemType:    "journalArticle",
			Authors:     `[{"creatorType":"author","firstName":"John","lastName":"Doe"}]`,
			Publication: "Test Journal",
			SyncedAt:    time.Now(),
		},
		{
			ZoteroKey: "ITEM002",
			Title:     "Test Book 1",
			ItemType:  "book",
			Authors:   `[{"creatorType":"author","firstName":"Jane","lastName":"Smith"}]`,
			SyncedAt:  time.Now(),
		},
	}

	for _, item := range items {
		err := db.Create(&item).Error
		assert.NoError(t, err)
	}

	t.Run("list all items", func(t *testing.T) {
		result, total, err := service.ListItems(ListItemsQuery{
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("list with item type filter", func(t *testing.T) {
		result, total, err := service.ListItems(ListItemsQuery{
			ItemType: "book",
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "Test Book 1", result[0].Title)
	})

	t.Run("list with search", func(t *testing.T) {
		result, total, err := service.ListItems(ListItemsQuery{
			Search:   "Article",
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "Test Article 1", result[0].Title)
	})
}

func TestZoteroService_GenerateCitation(t *testing.T) {
	db := setupTestDB(t)
	service := NewZoteroService(db)

	// Create test item
	item := models.ZoteroItem{
		ZoteroKey:   "ITEM003",
		Title:       "Research on Microwave Systems",
		ItemType:    "journalArticle",
		Authors:     `[{"creatorType":"author","firstName":"John","lastName":"Smith"},{"creatorType":"author","firstName":"Jane","lastName":"Doe"}]`,
		Publication: "IEEE Transactions",
		Volume:      "45",
		Issue:       "3",
		Pages:       "123-135",
		Date:        "2025-06",
		DOI:         "10.1109/example.2025.1234567",
		SyncedAt:    time.Now(),
	}
	err := db.Create(&item).Error
	assert.NoError(t, err)

	t.Run("generate GB/T 7714 citation", func(t *testing.T) {
		citation, err := service.GenerateCitation(item.ID, "gb7714")
		assert.NoError(t, err)
		assert.Contains(t, citation, "Smith John")
		assert.Contains(t, citation, "Doe Jane")
		assert.Contains(t, citation, "Research on Microwave Systems")
		assert.Contains(t, citation, "IEEE Transactions")
		assert.Contains(t, citation, "DOI:10.1109/example.2025.1234567")
	})

	t.Run("generate APA citation", func(t *testing.T) {
		citation, err := service.GenerateCitation(item.ID, "apa")
		assert.NoError(t, err)
		assert.Contains(t, citation, "Smith, J.")
		assert.Contains(t, citation, "Doe, J.")
		assert.Contains(t, citation, "(2025)")
		assert.Contains(t, citation, "Research on Microwave Systems")
	})

	t.Run("generate MLA citation", func(t *testing.T) {
		citation, err := service.GenerateCitation(item.ID, "mla")
		assert.NoError(t, err)
		assert.Contains(t, citation, "Smith, John")
		assert.Contains(t, citation, "et al")
		assert.Contains(t, citation, "\"Research on Microwave Systems.\"")
	})
}

func TestZoteroService_DeleteConnection(t *testing.T) {
	db := setupTestDB(t)
	service := NewZoteroService(db)

	t.Run("delete non-existent connection should fail", func(t *testing.T) {
		err := service.DeleteConnection("user-999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("delete existing connection should succeed", func(t *testing.T) {
		// Create connection
		connection := ZoteroConnection{
			UserID:       "user-to-delete",
			APIKey:       "test-api-key",
			ZoteroUserID: "12345678",
			IsActive:     true,
		}
		err := db.Create(&connection).Error
		assert.NoError(t, err)

		// Delete
		err = service.DeleteConnection("user-to-delete")
		assert.NoError(t, err)

		// Verify deletion
		_, err = service.GetConnection("user-to-delete")
		assert.Error(t, err)
	})
}
