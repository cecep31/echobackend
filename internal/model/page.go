package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Page represents a document/page in the workspace
type Page struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;not null" json:"workspace_id"`
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parent_id"` // For nested pages
	Title       string         `gorm:"not null" json:"title"`
	Icon        string         `json:"icon"`
	Blocks      []Block        `gorm:"foreignKey:PageID" json:"blocks"`
	CreatedBy   string         `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for the Page model
func (Page) TableName() string {
	return "pages"
}
