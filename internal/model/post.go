package model

import (
	"time"

	"github.com/uptrace/bun"
)

// Post represents the post model in the database
type Post struct {
	bun.BaseModel `bun:"table:posts,alias:p"`

	ID        string     `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v7()"`
	Title     string     `json:"title" bun:"title,notnull"`
	Photo_url string     `json:"photo_url" bun:"photo_url,type:text"`
	Body      string     `json:"body" bun:"body,notnull"`
	Slug      string     `json:"slug" bun:"slug,unique,notnull"`
	CreatedBy string     `json:"created_by" bun:"created_by"`
	Creator   User       `json:"creator" bun:"rel:belongs-to,join:created_by=id"`
	Tags      []Tag      `json:"tags" bun:"m2m:posts_to_tags,join:Post=Tag"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at,soft_delete"`
}

// PostResponse represents the post data that can be safely sent to clients
type PostResponse struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Photo_url string       `json:"photo_url"`
	Body      string       `json:"body"`
	Slug      string       `json:"slug"`
	Creator   UserResponse `json:"creator"`
	Tags      []Tag        `json:"tags"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty"`
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
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}
}

type CreatePostDTO struct {
	Title     string   `json:"title" validate:"required,min=7"`
	Photo_url string   `json:"photo_url"`
	Slug      string   `json:"slug" validate:"required,min=7"`
	Body      string   `json:"body" validate:"required,min=10"`
	Tags      []string `json:"tags" `
}

type UpdatePostDTO struct {
	Title     string   `json:"title"`
	Photo_url string   `json:"photo_url"`
	Slug      string   `json:"slug"`
	Body      string   `json:"body"`
	Tags      []string `json:"tags"`
}
