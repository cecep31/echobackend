package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatConversation struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Title     string         `gorm:"type:varchar(255);not null"`
	UserID    string         `gorm:"type:uuid;not null"`
	Messages  []ChatMessage  `gorm:"foreignKey:ConversationID"`
}

func (ChatConversation) TableName() string {
	return "chat_conversations"
}
