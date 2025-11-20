package model

import (
	"time"
)

type Tag struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      *string    `json:"name" gorm:"uniqueIndex;type:varchar(30)"`
	CreatedAt *time.Time `json:"created_at"`
	Posts     []Post     `gorm:"many2many:posts_to_tags;"`
}

// TableName specifies the table name for GORM
func (Tag) TableName() string {
	return "tags"
}

// TagResponse represents the tag data that can be safely sent to clients
type TagResponse struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
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
