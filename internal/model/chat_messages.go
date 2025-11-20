package model

import "time"

type ChatMessage struct {
	ID               string           `gorm:"type:uuid;primaryKey"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	ConversationID   string           `gorm:"type:uuid;not null"`
	UserID           string           `gorm:"type:uuid;not null"`
	Role             string           `gorm:"type:varchar(20);not null"`
	Content          string           `gorm:"not null"`
	Model            *string          `gorm:"type:varchar(100)"`
	PromptTokens     *int             `json:"prompt_tokens"`
	CompletionTokens *int             `json:"completion_tokens"`
	TotalTokens      *int             `json:"total_tokens"`
	Conversation     ChatConversation `gorm:"foreignKey:ConversationID"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
