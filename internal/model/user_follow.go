package model

import (
	"time"

	"gorm.io/gorm"
)

// UserFollow represents a following relationship between users
type UserFollow struct {
	ID         string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	FollowerID string         `json:"follower_id" gorm:"type:uuid;not null;index"`
	FollowingID string        `json:"following_id" gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Follower  User `json:"follower" gorm:"foreignKey:FollowerID"`
	Following User `json:"following" gorm:"foreignKey:FollowingID"`
}

func (UserFollow) TableName() string {
	return "user_follows"
}

// UserFollowStats represents follow statistics for a user
type UserFollowStats struct {
	UserID        string `json:"user_id"`
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}

// FollowRequest represents a request to follow/unfollow a user
type FollowRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

// FollowResponse represents the response after following/unfollowing
type FollowResponse struct {
	IsFollowing bool `json:"is_following"`
	Message     string `json:"message"`
}