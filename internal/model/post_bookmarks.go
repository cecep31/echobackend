package model

import (
	"time"

	"gorm.io/gorm"
)

type PostBookmark struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
	PostID    string         `json:"post_id" gorm:"type:uuid;not null"`
	UserID    string         `json:"user_id" gorm:"type:uuid;not null"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Post      Post           `gorm:"foreignKey:PostID"`
	User      User           `gorm:"foreignKey:UserID"`
}

func (PostBookmark) TableName() string {
	return "post_bookmarks"
}
