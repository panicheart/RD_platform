package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Product represents a product in the shelf
type Product struct {
	ID              string          `json:"id" gorm:"primaryKey;type:char(26)"`
	Name            string          `json:"name" gorm:"not null;size:200"`
	Description     string          `json:"description" gorm:"type:text"`
	TRLLevel        int             `json:"trl_level" gorm:"check:trl_level >= 1 AND trl_level <= 9"`
	Category        string          `json:"category" gorm:"size:100"`
	Version         string          `json:"version" gorm:"size:50"`
	SourceProjectID string          `json:"source_project_id" gorm:"type:char(26)"`
	OwnerID         string          `json:"owner_id" gorm:"type:char(26)"`
	IsPublished     bool            `json:"is_published" gorm:"default:false"`
	PublishedAt     *time.Time      `json:"published_at"`
	DownloadCount   int             `json:"download_count" gorm:"default:0"`
	Metadata        string          `json:"metadata" gorm:"type:jsonb"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedBy       string          `json:"created_by" gorm:"type:char(26)"`

	// Relations
	SourceProject *Project          `json:"source_project,omitempty" gorm:"foreignKey:SourceProjectID"`
	Owner         *User             `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Versions      []ProductVersion  `json:"versions,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for the model
func (Product) TableName() string {
	return "products"
}

// BeforeCreate generates ULID before insert
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = ulid.Make().String()
	}
	return nil
}

// GetTRLColor returns the color for TRL level
func (p *Product) GetTRLColor() string {
	if p.TRLLevel <= 3 {
		return "red"
	} else if p.TRLLevel <= 6 {
		return "yellow"
	}
	return "green"
}

// ProductVersion represents a product version
type ProductVersion struct {
	ID             string     `json:"id" gorm:"primaryKey;type:char(26)"`
	ProductID      string     `json:"product_id" gorm:"index;not null;type:char(26)"`
	Version        string     `json:"version" gorm:"not null;size:50"`
	ParentVersionID *string   `json:"parent_version_id" gorm:"type:char(26)"`
	Changelog      string     `json:"changelog" gorm:"type:text"`
	Files          string     `json:"files" gorm:"type:jsonb"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      string     `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Product       *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ParentVersion *ProductVersion `json:"parent_version,omitempty" gorm:"foreignKey:ParentVersionID"`
}

// TableName returns the table name for the model
func (ProductVersion) TableName() string {
	return "product_versions"
}

// BeforeCreate generates ULID before insert
func (pv *ProductVersion) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == "" {
		pv.ID = ulid.Make().String()
	}
	return nil
}

// CartItem represents a shopping cart item
type CartItem struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(26)"`
	UserID    string    `json:"user_id" gorm:"index;not null;type:char(26)"`
	ProductID string    `json:"product_id" gorm:"index;not null;type:char(26)"`
	ProjectID string    `json:"project_id" gorm:"type:char(26)"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for the model
func (CartItem) TableName() string {
	return "cart_items"
}

// BeforeCreate generates ULID before insert
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == "" {
		ci.ID = ulid.Make().String()
	}
	return nil
}

// Technology represents a technology in the shelf
type Technology struct {
	ID           string    `json:"id" gorm:"primaryKey;type:char(26)"`
	Name         string    `json:"name" gorm:"not null;size:200"`
	Description  string    `json:"description" gorm:"type:text"`
	TRLLevel     int       `json:"trl_level" gorm:"check:trl_level >= 1 AND trl_level <= 9"`
	Category     string    `json:"category" gorm:"size:100"`
	ParentID     *string   `json:"parent_id" gorm:"index;type:char(26)"`
	OwnerID      string    `json:"owner_id" gorm:"type:char(26)"`
	IsPublished  bool      `json:"is_published" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    string    `json:"created_by" gorm:"type:char(26)"`

	// Relations
	Parent       *Technology  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children     []Technology `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Owner        *User        `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
}

// TableName returns the table name for the model
func (Technology) TableName() string {
	return "technologies"
}

// BeforeCreate generates ULID before insert
func (t *Technology) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ulid.Make().String()
	}
	return nil
}

// GetTRLColor returns the color for TRL level
func (t *Technology) GetTRLColor() string {
	if t.TRLLevel <= 3 {
		return "red"
	} else if t.TRLLevel <= 6 {
		return "yellow"
	}
	return "green"
}
