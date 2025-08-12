package model

import (
	"time"

	"gorm.io/gorm"
)

// PostView represents a view record for a post
type PostView struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	PostID    string         `json:"post_id" gorm:"type:uuid;not null;index"`
	UserID    string         `json:"user_id" gorm:"type:uuid;index"`     // nullable for anonymous views
	IPAddress string         `json:"ip_address" gorm:"type:varchar(45)"` // for tracking anonymous views
	UserAgent string         `json:"user_agent" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Post Post  `json:"post" gorm:"foreignKey:PostID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (PostView) TableName() string {
	return "post_views"
}

// PostViewStats represents aggregated view statistics for a post
type PostViewStats struct {
	PostID             string `json:"post_id"`
	TotalViews         int64  `json:"total_views"`
	UniqueViews        int64  `json:"unique_views"`
	AnonymousViews     int64  `json:"anonymous_views"`
	AuthenticatedViews int64  `json:"authenticated_views"`
}
