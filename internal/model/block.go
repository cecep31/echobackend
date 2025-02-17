package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Block represents a content block within a page
type Block struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PageID    uuid.UUID      `gorm:"type:uuid;not null" json:"page_id"`
	Type      string         `gorm:"not null;type:varchar(50)" json:"type"` // paragraph, heading, list, image, etc.
	Props     string         `gorm:"type:jsonb" json:"props"`               // Block properties (backgroundColor, textColor, level for headings, etc.)
	Content   string         `gorm:"type:jsonb" json:"content"`            // InlineContent array or TableContent
	ParentID  *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`          // For nested blocks
	Position  float64        `gorm:"not null" json:"position"`            // For ordering blocks
	CreatedBy string         `gorm:"not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for the Block model
func (Block) TableName() string {
	return "blocks"
}

// BlockResponse represents a block with a simplified structure
type BlockResponse struct {
	ID        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	Props     string     `json:"props"`
	Content   string     `json:"content"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Position  float64    `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ToResponse converts a Block model to a BlockResponse
func (b *Block) ToResponse() *BlockResponse {
	return &BlockResponse{
		ID:        b.ID,
		Type:      b.Type,
		Props:     b.Props,
		Content:   b.Content,
		ParentID:  b.ParentID,
		Position:  b.Position,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		DeletedAt: &b.DeletedAt.Time,
	}
}
