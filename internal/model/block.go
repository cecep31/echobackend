package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm" // Import GORM
)

// Block represents a content block within a page
type Block struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PageID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"page_id"` // Foreign key to Page
	Type      string     `gorm:"not null" json:"type"`                    // paragraph, heading, list, image, etc.
	Props     string     `gorm:"type:jsonb" json:"props"`                 // Block properties
	Content   string     `gorm:"type:jsonb" json:"content"`               // InlineContent array or TableContent
	ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`        // For nested blocks
	Position  float64    `gorm:"not null" json:"position"`                // For ordering blocks
	CreatedBy string     `gorm:"type:uuid;not null" json:"created_by"`    // Assuming this is a User UUID

	// Optional: Define relationships explicitly if needed
	// Page        Page      `gorm:"foreignKey:PageID"`
	// ParentBlock *Block    `gorm:"foreignKey:ParentID"`
	// ChildBlocks []Block   `gorm:"foreignKey:ParentID"`
	// Creator     User      `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Use gorm.DeletedAt for soft delete
}

// TableName specifies the table name for GORM
func (Block) TableName() string {
	return "blocks"
}
