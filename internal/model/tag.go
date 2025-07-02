package model

import (
	"time"
)

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Posts     []Post    `json:"-" gorm:"many2many:posts_to_tags;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
