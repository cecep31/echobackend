package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Page represents a document/page in the workspace
type Page struct {
	bun.BaseModel `bun:"table:pages,alias:pg"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	WorkspaceID uuid.UUID  `bun:"workspace_id,type:uuid,notnull" json:"workspace_id"`
	ParentID    *uuid.UUID `bun:"parent_id,type:uuid" json:"parent_id"` // For nested pages
	Title       string     `bun:"title,notnull" json:"title"`
	Icon        string     `bun:"icon" json:"icon"`
	Blocks      []Block    `bun:"rel:has-many,join:id=page_id" json:"blocks"`
	CreatedBy   string     `bun:"created_by,notnull" json:"created_by"`
	CreatedAt   time.Time  `bun:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bun:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

// No need for TableName with Bun as it's specified in the struct tag
