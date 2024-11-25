package model

import (
	"time"

	"gorm.io/gorm"
)

// Post represents the post model in the database
type Post struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v7()"`
	Title     string         `json:"title" gorm:"not null"`
	Photo_url string         `json:"photo_url" gorm:"type:text"`
	Body      string         `json:"body" gorm:"not null"`
	Slug      string         `json:"slug" gorm:"unique;not null"`
	CreatedBy string         `json:"created_by"`
	Creator   User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	Tags      []Tag          `json:"tags" gorm:"many2many:posts_to_tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}

// PostResponse represents the post data that can be safely sent to clients
type PostResponse struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Photo_url string       `json:"photo_url"`
	Body      string       `json:"body"`
	Slug      string       `json:"slug"`
	Creator   UserResponse `json:"creator"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ToResponse converts a Post model to a PostResponse
func (p *Post) ToResponse() *PostResponse {
	return &PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Photo_url: p.Photo_url,
		Body:      p.Body,
		Slug:      p.Slug,
		Creator:   *p.Creator.ToResponse(),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
