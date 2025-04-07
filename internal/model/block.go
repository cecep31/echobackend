package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Block represents a content block within a page
type Block struct {
	bun.BaseModel `bun:"table:blocks,alias:b"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	PageID    uuid.UUID  `bun:"page_id,type:uuid,notnull" json:"page_id"`
	Type      string     `bun:"type,notnull" json:"type"`             // paragraph, heading, list, image, etc.
	Props     string     `bun:"props,type:jsonb" json:"props"`        // Block properties (backgroundColor, textColor, level for headings, etc.)
	Content   string     `bun:"content,type:jsonb" json:"content"`    // InlineContent array or TableContent
	ParentID  *uuid.UUID `bun:"parent_id,type:uuid" json:"parent_id"` // For nested blocks
	Position  float64    `bun:"position,notnull" json:"position"`     // For ordering blocks
	CreatedBy string     `bun:"created_by,notnull" json:"created_by"`
	CreatedAt time.Time  `bun:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

// No need for TableName with Bun as it's specified in the struct tag
