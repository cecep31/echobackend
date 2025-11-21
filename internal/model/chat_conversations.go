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

// DTOs for chat conversations
type CreateChatConversationDTO struct {
	Title string `json:"title" validate:"required,max=255"`
}

type UpdateChatConversationDTO struct {
	Title string `json:"title" validate:"max=255"`
}

type ChatConversationResponse struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
	UserID       string    `json:"user_id"`
	MessageCount int       `json:"message_count"`
}

func (c *ChatConversation) ToResponse() *ChatConversationResponse {
	return &ChatConversationResponse{
		ID:           c.ID,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		Title:        c.Title,
		UserID:       c.UserID,
		MessageCount: len(c.Messages),
	}
}
