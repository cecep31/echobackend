package model

import (
	"time"

	"gorm.io/gorm"
)

// PostLike represents a like record for a post
type PostLike struct {
	ID     string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	PostID string `json:"post_id" gorm:"type:uuid;not null;index"`
	UserID string `json:"user_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Post Post `json:"post" gorm:"foreignKey:PostID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (PostLike) TableName() string {
	return "post_likes"
}

// PostLikeStats represents aggregated like statistics for a post
type PostLikeStats struct {
	PostID     string `json:"post_id"`
	TotalLikes int64  `json:"total_likes"`
}

// PostLikeResponse represents the response structure for post likes
type PostLikeResponse struct {
	ID        string        `json:"id"`
	PostID    string        `json:"post_id"`
	UserID    string        `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

func (pl *PostLike) ToResponse() *PostLikeResponse {
	var userResp *UserResponse
	if pl.User.ID != "" {
		userResp = pl.User.ToResponse()
	}

	return &PostLikeResponse{
		ID:        pl.ID,
		PostID:    pl.PostID,
		UserID:    pl.UserID,
		User:      userResp,
		CreatedAt: pl.CreatedAt,
	}
}