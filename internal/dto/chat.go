package dto

import (
	"echobackend/internal/model"
	"time"
)

type CreateChatConversationRequest struct {
	Title string `json:"title" validate:"required,max=255"`
}

type UpdateChatConversationRequest struct {
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

func ChatConversationToResponse(c *model.ChatConversation) *ChatConversationResponse {
	if c == nil {
		return nil
	}
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
