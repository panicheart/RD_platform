package indexers

import (
	"fmt"

	"rdp-platform/rdp-api/clients"
	"rdp-platform/rdp-api/models"
	"gorm.io/gorm"
)

// ProductIndexer handles indexing products to MeiliSearch
type ProductIndexer struct {
	db          *gorm.DB
	meiliClient *clients.MeiliSearchClient
}

// NewProductIndexer creates a new product indexer
func NewProductIndexer(db *gorm.DB, meiliClient *clients.MeiliSearchClient) *ProductIndexer {
	return &ProductIndexer{
		db:          db,
		meiliClient: meiliClient,
	}
}

// IndexableProduct represents a product document for indexing
type IndexableProduct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Maturity    string `json:"maturity"`
	OwnerID     string `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
	TRL         string `json:"trl"`
	Version     string `json:"version"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// IndexAll indexes all products
func (i *ProductIndexer) IndexAll() error {
	var products []models.Product
	if err := i.db.Find(&products).Error; err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	docs := make([]IndexableProduct, 0, len(products))
	for _, p := range products {
		docs = append(docs, i.toIndexableDocument(p))
	}

	if len(docs) > 0 {
		return i.meiliClient.Index(clients.IndexNames.Products, docs)
	}
	
	return nil
}

// IndexSingle indexes a single product
func (i *ProductIndexer) IndexSingle(product models.Product) error {
	doc := i.toIndexableDocument(product)
	return i.meiliClient.Index(clients.IndexNames.Products, []IndexableProduct{doc})
}

// Delete removes a product from the index
func (i *ProductIndexer) Delete(productID string) error {
	return i.meiliClient.Delete(clients.IndexNames.Products, []string{productID})
}

func (i *ProductIndexer) toIndexableDocument(p models.Product) IndexableProduct {
	// Fetch owner name
	var ownerName string
	var owner models.User
	if err := i.db.First(&owner, "id = ?", p.OwnerID).Error; err == nil {
		ownerName = owner.DisplayName
	}

	return IndexableProduct{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Type:        p.Type,
		Maturity:    p.Maturity,
		OwnerID:     p.OwnerID,
		OwnerName:   ownerName,
		TRL:         p.TRL,
		Version:     p.Version,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
