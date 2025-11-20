package model

import (
	"time"
)

type Session struct {
	RefreshToken string     `gorm:"primaryKey;type:varchar(200)"`
	UserID       string     `gorm:"type:uuid"`
	CreatedAt    *time.Time `json:"created_at"`
	UserAgent    *string    `json:"user_agent"`
	ExpiresAt    *time.Time `json:"expires_at"`
	User         *User      `gorm:"foreignKey:UserID"`
}
