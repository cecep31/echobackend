package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatConversation struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Title     string         `gorm:"type:varchar(255);not null"`
	IsPinned  bool           `json:"is_pinned" gorm:"column:is_pinned;default:false;not null"`
	PinnedAt  *time.Time     `json:"pinned_at" gorm:"column:pinned_at"`
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
	ID           string     `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Title        string     `json:"title"`
	IsPinned     bool       `json:"is_pinned"`
	PinnedAt     *time.Time `json:"pinned_at,omitempty"`
	UserID       string     `json:"user_id"`
	MessageCount int        `json:"message_count"`
}

func (c *ChatConversation) ToResponse() *ChatConversationResponse {
	return &ChatConversationResponse{
		ID:           c.ID,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		Title:        c.Title,
		IsPinned:     c.IsPinned,
		PinnedAt:     c.PinnedAt,
		UserID:       c.UserID,
		MessageCount: len(c.Messages),
	}
}
