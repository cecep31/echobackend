package service

import (
	"context"
	"strings"
	"time"

	"echobackend/config"
	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type ChatConversationService interface {
	CreateConversation(ctx context.Context, userID string, conversation *dto.CreateChatConversationRequest) (*dto.ChatConversationResponse, error)
	CreateConversationStream(ctx context.Context, userID string, req *dto.CreateChatConversationStreamRequest) (*dto.ChatStreamResult, <-chan string, <-chan dto.ChatMessageResponse, <-chan error, error)
	GetConversationByID(ctx context.Context, id string, userID string) (*dto.ChatConversationResponse, error)
	GetUserConversations(ctx context.Context, userID string, offset int, limit int) ([]*dto.ChatConversationResponse, int64, error)
	UpdateConversation(ctx context.Context, id string, userID string, conversation *dto.UpdateChatConversationRequest) (*dto.ChatConversationResponse, error)
	DeleteConversation(ctx context.Context, id string, userID string) error
	CreateMessage(ctx context.Context, userID, conversationID string, req *dto.CreateChatMessageRequest) ([]*dto.ChatMessageResponse, error)
	CreateStreamingMessage(ctx context.Context, userID, conversationID string, req *dto.CreateChatMessageRequest) (*dto.ChatStreamResult, <-chan string, <-chan dto.ChatMessageResponse, <-chan error, error)
	SaveStreamingMessage(ctx context.Context, conversationID, userID, content string, model *string, usage OpenRouterUsage) (*dto.ChatMessageResponse, error)
	GetMessages(ctx context.Context, conversationID, userID string) ([]*dto.ChatMessageResponse, error)
	GetMessage(ctx context.Context, id, userID string) (*dto.ChatMessageResponse, error)
	DeleteMessage(ctx context.Context, id, userID string) (*dto.ChatMessageResponse, error)
}

type chatConversationService struct {
	conversationRepo repository.ChatConversationRepository
	openRouter       OpenRouterService
	config           *config.Config
}

func NewChatConversationService(conversationRepo repository.ChatConversationRepository, openRouter OpenRouterService, cfg *config.Config) ChatConversationService {
	return &chatConversationService{
		conversationRepo: conversationRepo,
		openRouter:       openRouter,
		config:           cfg,
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

func (s *chatConversationService) CreateConversationStream(ctx context.Context, userID string, req *dto.CreateChatConversationStreamRequest) (*dto.ChatStreamResult, <-chan string, <-chan dto.ChatMessageResponse, <-chan error, error) {
	chatConversation := &model.ChatConversation{
		Title:  buildConversationTitle(req.Title, req.Content),
		UserID: userID,
	}
	createdConversation, err := s.conversationRepo.CreateConversation(ctx, chatConversation)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	msgReq := &dto.CreateChatMessageRequest{
		Content:     req.Content,
		Role:        "user",
		Model:       req.Model,
		Temperature: req.Temperature,
	}
	result, chunks, complete, errCh, err := s.createStreamingMessageInternal(ctx, userID, createdConversation.ID, msgReq)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	result.ConversationID = createdConversation.ID
	return result, chunks, complete, errCh, nil
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
	if conversation.IsPinned != nil {
		updates["is_pinned"] = *conversation.IsPinned
		if *conversation.IsPinned {
			now := time.Now()
			updates["pinned_at"] = &now
		} else {
			updates["pinned_at"] = nil
		}
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

func (s *chatConversationService) CreateMessage(ctx context.Context, userID, conversationID string, req *dto.CreateChatMessageRequest) ([]*dto.ChatMessageResponse, error) {
	conversation, err := s.getOwnedConversation(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	message := &model.ChatMessage{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           normalizedRole(req.Role),
		Content:        req.Content,
		Model:          req.Model,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	userMessage, err := s.conversationRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{"updated_at": now}); err != nil {
		return nil, err
	}

	responses := []*dto.ChatMessageResponse{dto.ChatMessageToResponse(userMessage)}
	if userMessage.Role != "user" || s.openRouter == nil {
		return responses, nil
	}

	contextMessages, err := s.conversationRepo.GetMessagesByConversationIDAsc(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if conversation.Title == "New conversation" && len(contextMessages) == 1 {
		if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{
			"title":      buildConversationTitle(nil, req.Content),
			"updated_at": now,
		}); err != nil {
			return nil, err
		}
	}

	reply, err := s.openRouter.GenerateResponse(ctx, toOpenRouterMessages(contextMessages), req.Model, normalizedTemperature(req.Temperature))
	if err != nil || len(reply.Choices) == 0 {
		return responses, nil
	}

	assistant := &model.ChatMessage{
		ConversationID:   conversationID,
		UserID:           userID,
		Role:             reply.Choices[0].Message.Role,
		Content:          reply.Choices[0].Message.Content,
		Model:            s.effectiveModel(req.Model),
		PromptTokens:     reply.Usage.PromptTokens,
		CompletionTokens: reply.Usage.CompletionTokens,
		TotalTokens:      reply.Usage.TotalTokens,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	createdAssistant, err := s.conversationRepo.CreateMessage(ctx, assistant)
	if err != nil {
		return nil, err
	}
	if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{"updated_at": createdAssistant.UpdatedAt}); err != nil {
		return nil, err
	}
	responses = append(responses, dto.ChatMessageToResponse(createdAssistant))
	return responses, nil
}

func (s *chatConversationService) CreateStreamingMessage(ctx context.Context, userID, conversationID string, req *dto.CreateChatMessageRequest) (*dto.ChatStreamResult, <-chan string, <-chan dto.ChatMessageResponse, <-chan error, error) {
	return s.createStreamingMessageInternal(ctx, userID, conversationID, req)
}

func (s *chatConversationService) createStreamingMessageInternal(ctx context.Context, userID, conversationID string, req *dto.CreateChatMessageRequest) (*dto.ChatStreamResult, <-chan string, <-chan dto.ChatMessageResponse, <-chan error, error) {
	conversation, err := s.getOwnedConversation(ctx, conversationID, userID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	now := time.Now()
	message := &model.ChatMessage{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           normalizedRole(req.Role),
		Content:        req.Content,
		Model:          req.Model,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	userMessage, err := s.conversationRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{"updated_at": now}); err != nil {
		return nil, nil, nil, nil, err
	}

	result := &dto.ChatStreamResult{
		UserMessage: dto.ChatMessageToResponse(userMessage),
	}
	if userMessage.Role != "user" || s.openRouter == nil {
		complete := make(chan dto.ChatMessageResponse)
		close(complete)
		errCh := make(chan error)
		close(errCh)
		return result, nil, complete, errCh, nil
	}

	contextMessages, err := s.conversationRepo.GetMessagesByConversationIDAsc(ctx, conversationID, userID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if conversation.Title == "New conversation" && len(contextMessages) == 1 {
		if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{
			"title":      buildConversationTitle(nil, req.Content),
			"updated_at": now,
		}); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	chunks, usageCh, upstreamErrCh := s.openRouter.GenerateStream(ctx, toOpenRouterMessages(contextMessages), req.Model, normalizedTemperature(req.Temperature))
	complete := make(chan dto.ChatMessageResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(complete)
		defer close(errCh)

		var content strings.Builder
		for chunk := range chunks {
			content.WriteString(chunk)
		}

		select {
		case err, ok := <-upstreamErrCh:
			if ok && err != nil {
				errCh <- err
				return
			}
		default:
		}

		usage := OpenRouterUsage{}
		if u, ok := <-usageCh; ok {
			usage = u
		}
		if content.Len() == 0 {
			return
		}
		saved, err := s.SaveStreamingMessage(ctx, conversationID, userID, content.String(), req.Model, usage)
		if err != nil {
			errCh <- err
			return
		}
		complete <- *saved
	}()

	return result, chunks, complete, errCh, nil
}

func (s *chatConversationService) SaveStreamingMessage(ctx context.Context, conversationID, userID, content string, modelID *string, usage OpenRouterUsage) (*dto.ChatMessageResponse, error) {
	now := time.Now()
	message := &model.ChatMessage{
		ConversationID:   conversationID,
		UserID:           userID,
		Role:             "assistant",
		Content:          content,
		Model:            s.effectiveModel(modelID),
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	created, err := s.conversationRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	if _, err := s.conversationRepo.UpdateConversation(ctx, conversationID, map[string]interface{}{"updated_at": now}); err != nil {
		return nil, err
	}
	return dto.ChatMessageToResponse(created), nil
}

func (s *chatConversationService) GetMessages(ctx context.Context, conversationID, userID string) ([]*dto.ChatMessageResponse, error) {
	if _, err := s.getOwnedConversation(ctx, conversationID, userID); err != nil {
		return nil, err
	}
	messages, err := s.conversationRepo.GetMessagesByConversationID(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	responses := make([]*dto.ChatMessageResponse, 0, len(messages))
	for _, message := range messages {
		responses = append(responses, dto.ChatMessageToResponse(message))
	}
	return responses, nil
}

func (s *chatConversationService) GetMessage(ctx context.Context, id, userID string) (*dto.ChatMessageResponse, error) {
	message, err := s.conversationRepo.GetMessageByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return dto.ChatMessageToResponse(message), nil
}

func (s *chatConversationService) DeleteMessage(ctx context.Context, id, userID string) (*dto.ChatMessageResponse, error) {
	message, err := s.conversationRepo.GetMessageByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := s.conversationRepo.DeleteMessage(ctx, id, userID); err != nil {
		return nil, err
	}
	return dto.ChatMessageToResponse(message), nil
}

func (s *chatConversationService) getOwnedConversation(ctx context.Context, id, userID string) (*model.ChatConversation, error) {
	conversation, err := s.conversationRepo.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if conversation.UserID != userID {
		return nil, apperrors.ErrConversationNotOwned
	}
	return conversation, nil
}

func buildConversationTitle(title *string, fallbackContent string) string {
	if title != nil && strings.TrimSpace(*title) != "" {
		return strings.TrimSpace(*title)
	}
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(fallbackContent)), " ")
	if trimmed == "" {
		return "New conversation"
	}
	if len(trimmed) > 50 {
		return trimmed[:50]
	}
	return trimmed
}

func normalizedRole(role string) string {
	role = strings.TrimSpace(role)
	if role == "" {
		return "user"
	}
	return role
}

func normalizedTemperature(v *float64) float64 {
	if v == nil || *v < 0 || *v > 2 {
		return 0.7
	}
	return *v
}

func toOpenRouterMessages(messages []*model.ChatMessage) []OpenRouterMessage {
	result := make([]OpenRouterMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, OpenRouterMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return result
}

func (s *chatConversationService) effectiveModel(model *string) *string {
	if model != nil && strings.TrimSpace(*model) != "" {
		trimmed := strings.TrimSpace(*model)
		return &trimmed
	}
	if s.config == nil {
		return nil
	}
	defaultModel := strings.TrimSpace(s.config.OpenRouter.DefaultModel)
	if defaultModel == "" {
		return nil
	}
	return &defaultModel
}
