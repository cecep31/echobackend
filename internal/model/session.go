package model

import (
	"time"
)

type Session struct {
	RefreshToken string     `gorm:"primaryKey;type:text;not null"`
	UserID       string     `gorm:"type:uuid;not null"`
	CreatedAt    *time.Time `json:"created_at"`
	UserAgent    *string    `json:"user_agent"`
	ExpiresAt    *time.Time `json:"expires_at"`
	User         *User      `gorm:"foreignKey:UserID"`
}

func (Session) TableName() string {
	return "sessions"
}
