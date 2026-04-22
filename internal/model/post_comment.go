package model

import (
	"time"

	"gorm.io/gorm"
)

// PostComment represents a comment on a post
type PostComment struct {
	ID               string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt        *time.Time     `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
	Text             string         `json:"text" gorm:"type:text;not null"`
	PostID           string         `json:"post_id" gorm:"type:uuid;not null"`
	ParentCommentID  *string        `json:"parent_comment_id" gorm:"type:uuid;index"`
	CreatedBy        string         `json:"created_by" gorm:"type:uuid;not null"`

	// Relationships
	User           *User        `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	Posts          *Post        `gorm:"foreignKey:PostID" json:"posts,omitempty"`
	ParentComment  *PostComment `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
}

func (PostComment) TableName() string {
	return "post_comments"
}

type PostCommentResponse struct {
	ID              string        `json:"id"`
	PostID          string        `json:"post_id"`
	ParentCommentID *string       `json:"parent_comment_id,omitempty"`
	Text            string        `json:"text"`
	User            *UserResponse `json:"user,omitempty"`
	CreatedAt       *time.Time    `json:"created_at"`
	UpdatedAt       *time.Time    `json:"updated_at"`
}

func (pc *PostComment) ToResponse() *PostCommentResponse {
	var userResp *UserResponse
	if pc.User != nil && pc.User.ID != "" {
		userResp = pc.User.ToResponse()
	}

	return &PostCommentResponse{
		ID:              pc.ID,
		PostID:          pc.PostID,
		ParentCommentID: pc.ParentCommentID,
		Text:            pc.Text,
		User:            userResp,
		CreatedAt:       pc.CreatedAt,
		UpdatedAt:       pc.UpdatedAt,
	}
}

type CreatePostCommentDTO struct {
	Text string `json:"text" validate:"required,min=1,max=1000"`
}
