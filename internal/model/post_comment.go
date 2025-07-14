package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostComment represents a comment on a post
type PostComment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PostID    string    `gorm:"type:uuid;not null;index" json:"post_id"` // Foreign key to Post
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedBy string    `gorm:"type:uuid;not null" json:"created_by"` // User UUID

	// Relationships
	Post    Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostComment) TableName() string {
	return "post_comments"
}

type PostCommentResponse struct {
	ID        uuid.UUID     `json:"id"`
	PostID    string        `json:"post_id"`
	Content   string        `json:"content"`
	Creator   *UserResponse `json:"creator,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (pc *PostComment) ToResponse() *PostCommentResponse {
	var creatorResp *UserResponse
	if pc.Creator.ID != "" {
		creatorResp = pc.Creator.ToResponse()
	}

	return &PostCommentResponse{
		ID:        pc.ID,
		PostID:    pc.PostID,
		Content:   pc.Content,
		Creator:   creatorResp,
		CreatedAt: pc.CreatedAt,
		UpdatedAt: pc.UpdatedAt,
	}
}

type CreatePostCommentDTO struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}