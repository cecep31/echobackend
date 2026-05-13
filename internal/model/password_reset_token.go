package model

import (
	"time"
)

type PasswordResetToken struct {
	ID        string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string     `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string     `json:"token" gorm:"type:text;not null;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:now()"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at"`
	User      *User      `json:"-" gorm:"foreignKey:UserID"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}