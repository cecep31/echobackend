package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Workspace represents a container for pages and content
type Workspace struct {
	bun.BaseModel `bun:"table:workspaces,alias:ws"`

	ID          uuid.UUID         `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Name        string            `bun:"name,notnull" json:"name"`
	Description string            `bun:"description" json:"description"`
	Icon        string            `bun:"icon" json:"icon"`
	CreatedBy   string            `bun:"created_by,notnull" json:"created_by"`
	Members     []WorkspaceMember `bun:"rel:has-many,join:id=workspace_id" json:"members"`
	CreatedAt   time.Time         `bun:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bun:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	bun.BaseModel `bun:"table:workspace_members,alias:wm"`

	WorkspaceID uuid.UUID `bun:"workspace_id,type:uuid" json:"workspace_id"`
	UserID      string    `bun:"user_id,notnull" json:"user_id"`
	Role        string    `bun:"role,notnull" json:"role"` // admin, editor, viewer
	CreatedAt   time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at" json:"updated_at"`
}

// No need for TableName with Bun as it's specified in the struct tag
