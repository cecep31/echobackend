package model

import (
	"time"

	"gorm.io/gorm"
)

type UserFollow struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	FollowerID  string         `json:"follower_id" gorm:"type:uuid;not null;index"`
	FollowingID string         `json:"following_id" gorm:"type:uuid;not null;index"`
	CreatedAt   *time.Time     `json:"created_at" gorm:"index"`
	UpdatedAt   *time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Follower  *User `json:"follower" gorm:"foreignKey:FollowerID"`
	Following *User `json:"following" gorm:"foreignKey:FollowingID"`
}

func (UserFollow) TableName() string {
	return "user_follows"
}
