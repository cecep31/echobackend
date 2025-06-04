package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm" // Import GORM
)

// Comment represents a comment on a post or block
type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	BlockID   uuid.UUID `gorm:"type:uuid;not null;index" json:"block_id"` // Foreign key to Block
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedBy string    `gorm:"type:uuid;not null" json:"created_by"` // Assuming User UUID

	// Optional: Define relationships explicitly if needed
	// Block     Block     `gorm:"foreignKey:BlockID"`
	// Creator   User      `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Use gorm.DeletedAt for soft delete
}

// TableName specifies the table name for GORM
func (Comment) TableName() string {
	return "comments"
}
