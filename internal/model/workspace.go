package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Workspace represents a container for pages and content
type Workspace struct {
	ID          uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string            `gorm:"not null" json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	CreatedBy   string            `gorm:"not null" json:"created_by"`
	Members     []WorkspaceMember `gorm:"foreignKey:WorkspaceID" json:"members"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
}

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	WorkspaceID uuid.UUID `gorm:"type:uuid" json:"workspace_id"`
	UserID      string    `gorm:"not null" json:"user_id"`
	Role        string    `gorm:"not null;type:varchar(20)" json:"role"` // admin, editor, viewer
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for the WorkspaceMember model
func (WorkspaceMember) TableName() string {
	return "workspace_members"
}

// TableName specifies the table name for the Workspace model
func (Workspace) TableName() string {
	return "workspaces"
}
