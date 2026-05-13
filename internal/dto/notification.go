package dto

import (
	"echobackend/internal/model"
	"encoding/json"
	"time"
)

type CreateNotificationRequest struct {
	UserID  string         `json:"user_id" validate:"required,uuid4"`
	Type    string         `json:"type" validate:"required,max=50"`
	Title   string         `json:"title" validate:"required,max=255"`
	Message *string        `json:"message"`
	Data    map[string]any `json:"data"`
}

type NotificationResponse struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Message   *string        `json:"message"`
	Read      bool           `json:"read"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
}

type NotificationListFilter struct {
	Unread bool
	Limit  int
	Offset int
}

type NotificationUnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}

type NotificationMarkAllReadResponse struct {
	UpdatedCount int64 `json:"updated_count"`
}

func NotificationToResponse(notification *model.Notification) *NotificationResponse {
	if notification == nil {
		return nil
	}

	var payload map[string]any
	if notification.Data != nil && *notification.Data != "" {
		_ = json.Unmarshal([]byte(*notification.Data), &payload)
	}

	return &NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		Read:      notification.Read,
		Data:      payload,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}
}
