package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm" // Import GORM
)

// Page represents a document/page in the workspace
type Page struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // Assuming uuid_generate_v4() is a DB func
	WorkspaceID uuid.UUID  `gorm:"type:uuid;not null;index" json:"workspace_id"`              // Added index for faster lookups
	ParentID    *uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`                          // For nested pages, added index
	Title       string     `gorm:"not null" json:"title"`
	Icon        string     `json:"icon"`
	Blocks      []Block    `gorm:"foreignKey:PageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"blocks"` // Assuming Block has PageID. Added constraints.
	CreatedBy   string     `gorm:"type:uuid;not null" json:"created_by"`                                          // Assuming this is a User UUID

	// Optional: Define relationships explicitly for easier access if needed
	// Workspace   Workspace `gorm:"foreignKey:WorkspaceID"`
	// Parent      *Page     `gorm:"foreignKey:ParentID"`
	// Children    []Page    `gorm:"foreignKey:ParentID"`
	// Creator     User      `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Use gorm.DeletedAt for soft delete
}

// TableName specifies the table name for GORM
func (Page) TableName() string {
	return "pages"
}
