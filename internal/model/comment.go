package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Comment represents comments on blocks
type Comment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BlockID   uuid.UUID      `gorm:"type:uuid;not null" json:"block_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedBy string         `gorm:"not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
