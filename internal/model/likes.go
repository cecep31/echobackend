package model

import (
	"time"
)

type Like struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time `json:"created_at"`
	PostID    *string    `json:"post_id" gorm:"type:uuid"`
	UserID    *string    `json:"user_id" gorm:"type:uuid"`
	User      *User      `gorm:"foreignKey:UserID"`
	Post      *Post      `gorm:"foreignKey:PostID"`
}

func (Like) TableName() string {
	return "likes"
}
