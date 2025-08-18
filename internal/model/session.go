package model

import "time"

type Session struct {
	RefreshToken string `gorm:"primaryKey"`
	UserID       string
	CreatedAt    time.Time
}
