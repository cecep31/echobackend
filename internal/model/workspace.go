package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm" // Import GORM
)

// Workspace represents a container for pages and content
type Workspace struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string            `gorm:"not null" json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	CreatedBy   string            `gorm:"type:uuid;not null" json:"created_by"`                                                 // Assuming User UUID
	Members     []WorkspaceMember `gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"members"` // Or OnDelete:CASCADE if members should be deleted with workspace

	// Optional: Define relationship to Creator User
	// Creator     User              `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Workspace) TableName() string {
	return "workspaces"
}

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	// No gorm.Model embedding if we use composite primary key and manage CreatedAt/UpdatedAt manually or via hooks
	WorkspaceID uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false" json:"workspace_id"`
	UserID      string    `gorm:"type:uuid;primaryKey;autoIncrement:false" json:"user_id"` // Assuming User ID is UUID
	Role        string    `gorm:"not null" json:"role"`                                    // e.g., admin, editor, viewer

	// Optional: Define relationships to Workspace and User
	// Workspace   Workspace `gorm:"foreignKey:WorkspaceID"`
	// User        User      `gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"` // GORM will manage these if fields exist
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index"` // If WorkspaceMembers can be soft-deleted independently
}

// TableName specifies the table name for GORM
func (WorkspaceMember) TableName() string {
	return "workspace_members"
}
