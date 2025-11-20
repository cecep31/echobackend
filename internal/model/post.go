package model

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

type Post struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt     *time.Time     `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Title         *string        `json:"title" gorm:"type:varchar(255)"`
	CreatedBy     *string        `json:"created_by" gorm:"type:uuid;uniqueIndex:creator_slug_unique"`
	Body          *string        `json:"body"`
	Slug          *string        `json:"slug" gorm:"type:varchar(255);uniqueIndex:creator_slug_unique"`
	Photo_url     *string        `json:"photo_url"`
	Published     *bool          `json:"published" gorm:"default:true"`
	ViewCount     int64          `json:"view_count" gorm:"type:bigint;default:0"`
	LikeCount     int64          `json:"like_count" gorm:"type:bigint;default:0"`
	PostComments  []PostComment  `gorm:"foreignKey:PostID"`
	PostLikes     []PostLike     `gorm:"foreignKey:PostID"`
	PostBookmarks []PostBookmark `gorm:"foreignKey:PostID"`
	Creator       *User          `gorm:"foreignKey:CreatedBy"`
	Tags          []Tag          `gorm:"many2many:posts_to_tags;"`
}

func (Post) TableName() string {
	return "posts"
}

type PostResponse struct {
	ID        string        `json:"id"`
	Title     *string       `json:"title"`
	Photo_url *string       `json:"photo_url"`
	Body      *string       `json:"body"`
	Slug      *string       `json:"slug"`
	ViewCount int64         `json:"view_count"`
	LikeCount int64         `json:"like_count"`
	Published *bool         `json:"published"`
	Creator   *UserResponse `json:"creator,omitempty"`
	Tags      []TagResponse `json:"tags,omitempty"`
	CreatedAt *time.Time    `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at,omitempty"`
}

func (p *Post) ToResponse() *PostResponse {
	var creatorResp *UserResponse
	if p.Creator.ID != "" {
		creatorResp = p.Creator.ToResponse()
	}

	var tagResponses []TagResponse
	if p.Tags != nil {
		tagResponses = make([]TagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			tagResponses[i] = *tag.ToResponse()
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
		ViewCount: p.ViewCount,
		LikeCount: p.LikeCount,
		Published: p.Published,
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
	Published bool     `json:"published"`
	Tags      []string `json:"tags" `
}

type UpdatePostDTO struct {
	Title     string   `json:"title"`
	Photo_url string   `json:"photo_url"`
	Slug      string   `json:"slug"`
	Body      string   `json:"body"`
	Published *bool    `json:"published"`
	Tags      []string `json:"tags"`
}
