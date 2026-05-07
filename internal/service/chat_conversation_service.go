package service

import (
	"context"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type ChatConversationService interface {
	CreateConversation(ctx context.Context, userID string, conversation *dto.CreateChatConversationRequest) (*dto.ChatConversationResponse, error)
	GetConversationByID(ctx context.Context, id string, userID string) (*dto.ChatConversationResponse, error)
	GetUserConversations(ctx context.Context, userID string, offset int, limit int) ([]*dto.ChatConversationResponse, int64, error)
	UpdateConversation(ctx context.Context, id string, userID string, conversation *dto.UpdateChatConversationRequest) (*dto.ChatConversationResponse, error)
	DeleteConversation(ctx context.Context, id string, userID string) error
}

type chatConversationService struct {
	conversationRepo repository.ChatConversationRepository
}

func NewChatConversationService(conversationRepo repository.ChatConversationRepository) ChatConversationService {
	return &chatConversationService{
		conversationRepo: conversationRepo,
	}
}

func (s *chatConversationService) CreateConversation(ctx context.Context, userID string, conversation *dto.CreateChatConversationRequest) (*dto.ChatConversationResponse, error) {
	chatConversation := &model.ChatConversation{
		Title:  conversation.Title,
		UserID: userID,
	}

	createdConversation, err := s.conversationRepo.CreateConversation(ctx, chatConversation)
	if err != nil {
		return nil, err
	}

	return dto.ChatConversationToResponse(createdConversation), nil
}

func (s *chatConversationService) GetConversationByID(ctx context.Context, id string, userID string) (*dto.ChatConversationResponse, error) {
	conversation, err := s.conversationRepo.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if conversation.UserID != userID {
		return nil, apperrors.ErrConversationNotOwned
	}

	return dto.ChatConversationToResponse(conversation), nil
}

func (s *chatConversationService) GetUserConversations(ctx context.Context, userID string, offset int, limit int) ([]*dto.ChatConversationResponse, int64, error) {
	conversations, total, err := s.conversationRepo.GetUserConversations(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ChatConversationResponse
	for _, conversation := range conversations {
		responses = append(responses, dto.ChatConversationToResponse(conversation))
	}

	return responses, total, nil
}

func (s *chatConversationService) UpdateConversation(ctx context.Context, id string, userID string, conversation *dto.UpdateChatConversationRequest) (*dto.ChatConversationResponse, error) {
	existingConversation, err := s.conversationRepo.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingConversation.UserID != userID {
		return nil, apperrors.ErrConversationNotOwned
	}

	updates := make(map[string]interface{})
	if conversation.Title != "" {
		updates["title"] = conversation.Title
	}

	updatedConversation, err := s.conversationRepo.UpdateConversation(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return dto.ChatConversationToResponse(updatedConversation), nil
}

func (s *chatConversationService) DeleteConversation(ctx context.Context, id string, userID string) error {
	existingConversation, err := s.conversationRepo.GetConversationByID(ctx, id)
	if err != nil {
		return err
	}

	if existingConversation.UserID != userID {
		return apperrors.ErrConversationNotOwned
	}

	return s.conversationRepo.DeleteConversation(ctx, id)
}
