package model

import (
	"time" // For CreatedAt, UpdatedAt if we add gorm.Model or similar
	// "gorm.io/gorm" // No longer needed if gorm.DeletedAt is not used
)

type Tag struct {
	// gorm.Model // Embed if you want ID, CreatedAt, UpdatedAt, DeletedAt by default
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"uniqueIndex;not null"` // Assuming tag names are unique and not null
	Posts []Post `json:"-" gorm:"many2many:posts_tags;"`   // Many to many with Post, posts_tags is the join table

	// Optional: Add CreatedAt, UpdatedAt if not embedding gorm.Model
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // If soft delete is needed for tags
}

// TableName specifies the table name for GORM
func (Tag) TableName() string {
	return "tags"
}

// TagResponse represents the tag data that can be safely sent to clients
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Add other fields if needed in response
}

// ToResponse converts a Tag model to a TagResponse
func (t *Tag) ToResponse() *TagResponse {
	if t == nil {
		return nil
	}
	return &TagResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}
