package model

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Post represents the post model in the database
type Post struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"` // Assuming uuid_generate_v7() is a DB func or handled by app
	Title     string         `json:"title" gorm:"not null"`
	Photo_url string         `json:"photo_url" gorm:"type:text"`
	Body      string         `json:"body" gorm:"type:text;not null"` // Changed to type:text for potentially long content
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null"`
	CreatedBy string         `json:"created_by" gorm:"type:uuid"`          // Foreign key for User
	Creator   User           `json:"creator" gorm:"foreignKey:CreatedBy"`  // Belongs to User
	Tags      []Tag          `json:"tags" gorm:"many2many:posts_to_tags;"` // Many to many with Tag, posts_tags is the join table
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // GORM soft delete
}

// TableName specifies the table name for GORM
func (Post) TableName() string {
	return "posts"
}

// PostResponse represents the post data that can be safely sent to clients
type PostResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Photo_url string        `json:"photo_url"`
	Body      string        `json:"body"`
	Slug      string        `json:"slug"`
	Creator   *UserResponse `json:"creator,omitempty"` // Made pointer to handle potential nil Creator
	Tags      []TagResponse `json:"tags,omitempty"`    // Assuming TagResponse exists or will be created
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at,omitempty"`
}

// ToResponse converts a Post model to a PostResponse
func (p *Post) ToResponse() *PostResponse {
	var creatorResp *UserResponse
	// Check if Creator object is valid (e.g., ID is not zero value for User's PK type)
	// This check depends on how User struct is defined and if Creator is always preloaded.
	// A simple check for non-zero ID (if User ID is string, check for non-empty string)
	if p.Creator.ID != "" { // Assuming User ID is string and non-empty means valid
		creatorResp = p.Creator.ToResponse()
	}

	var tagResponses []TagResponse
	if p.Tags != nil {
		tagResponses = make([]TagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			tagResponses[i] = *tag.ToResponse() // Assuming Tag has a ToResponse method returning *TagResponse
		}
	}

	var deletedAtTime *time.Time
	if p.DeletedAt.Valid {
		deletedAtTime = &p.DeletedAt.Time
	}

	return &PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Photo_url: p.Photo_url,
		Body:      p.Body,
		Slug:      p.Slug,
		Creator:   creatorResp,
		Tags:      tagResponses,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: deletedAtTime,
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
