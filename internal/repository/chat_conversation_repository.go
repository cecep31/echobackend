package repository

import (
	"context"
	"errors"
	"fmt"

	"echobackend/internal/model"

	"gorm.io/gorm"
)

var (
	ErrChatConversationNotFound = errors.New("chat conversation not found")
)

type ChatConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *model.ChatConversation) (*model.ChatConversation, error)
	GetConversationByID(ctx context.Context, id string) (*model.ChatConversation, error)
	GetConversationsByUserID(ctx context.Context, userID string, offset int, limit int) ([]*model.ChatConversation, int64, error)
	UpdateConversation(ctx context.Context, id string, conversation *model.UpdateChatConversationDTO) (*model.ChatConversation, error)
	DeleteConversation(ctx context.Context, id string) error
	GetUserConversations(ctx context.Context, userID string, offset int, limit int) ([]*model.ChatConversation, int64, error)
}

type chatConversationRepository struct {
	db *gorm.DB
}

func NewChatConversationRepository(db *gorm.DB) ChatConversationRepository {
	return &chatConversationRepository{db: db}
}

func (r *chatConversationRepository) CreateConversation(ctx context.Context, conversation *model.ChatConversation) (*model.ChatConversation, error) {
	err := r.db.WithContext(ctx).Create(conversation).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create chat conversation: %w", err)
	}
	return conversation, nil
}

func (r *chatConversationRepository) GetConversationByID(ctx context.Context, id string) (*model.ChatConversation, error) {
	var conversation model.ChatConversation
	err := r.db.WithContext(ctx).
		Preload("Messages").
		First(&conversation, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatConversationNotFound
		}
		return nil, fmt.Errorf("failed to get chat conversation: %w", err)
	}
	return &conversation, nil
}

func (r *chatConversationRepository) GetConversationsByUserID(ctx context.Context, userID string, offset int, limit int) ([]*model.ChatConversation, int64, error) {
	return r.GetUserConversations(ctx, userID, offset, limit)
}

func (r *chatConversationRepository) GetUserConversations(ctx context.Context, userID string, offset int, limit int) ([]*model.ChatConversation, int64, error) {
	var conversations []*model.ChatConversation
	var count int64

	// Count total records
	err := r.db.WithContext(ctx).Model(&model.ChatConversation{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chat conversations: %w", err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Preload("Messages").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chat conversations: %w", err)
	}

	return conversations, count, nil
}

func (r *chatConversationRepository) UpdateConversation(ctx context.Context, id string, conversation *model.UpdateChatConversationDTO) (*model.ChatConversation, error) {
	// Check if the conversation exists first
	var existingConversation model.ChatConversation
	err := r.db.WithContext(ctx).First(&existingConversation, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatConversationNotFound
		}
		return nil, fmt.Errorf("error checking conversation existence: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if conversation.Title != "" {
		updates["title"] = conversation.Title
	}

	if len(updates) == 0 {
		// No fields to update, return current conversation
		return &existingConversation, nil
	}

	// Update the conversation
	result := r.db.WithContext(ctx).Model(&model.ChatConversation{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update chat conversation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrChatConversationNotFound
	}

	// Fetch and return the updated conversation
	err = r.db.WithContext(ctx).
		Preload("Messages").
		First(&existingConversation, "id = ?", id).Error

	if err != nil {
		return nil, fmt.Errorf("conversation updated, but failed to retrieve updated record: %w", err)
	}

	return &existingConversation, nil
}

func (r *chatConversationRepository) DeleteConversation(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ChatConversation{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete chat conversation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrChatConversationNotFound
	}
	return nil
}
