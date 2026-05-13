package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type ChatConversationHandler struct {
	chatConversationService service.ChatConversationService
}

func NewChatConversationHandler(chatConversationService service.ChatConversationService) *ChatConversationHandler {
	return &ChatConversationHandler{
		chatConversationService: chatConversationService,
	}
}

func (h *ChatConversationHandler) CreateConversation(c *echo.Context) error {
	var conversationReq dto.CreateChatConversationRequest
	if err := c.Bind(&conversationReq); err != nil {
		return response.BadRequest(c, "Failed to create conversation", err)
	}

	if err := c.Validate(conversationReq); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	newConversation, err := h.chatConversationService.CreateConversation(c.Request().Context(), userID, &conversationReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to create conversation", err)
	}

	return response.Created(c, "Successfully created conversation", newConversation)
}

func (h *ChatConversationHandler) CreateConversationStream(c *echo.Context) error {
	var req dto.CreateChatConversationStreamRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Failed to create conversation stream", err)
	}
	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	result, chunks, complete, errCh, err := h.chatConversationService.CreateConversationStream(c.Request().Context(), userID, &req)
	if err != nil {
		return h.respondChatError(c, "Failed to create conversation stream", err)
	}
	if chunks == nil {
		return response.Created(c, "Message created successfully", []any{result.UserMessage})
	}
	return streamChatEvents(c, map[string]any{
		"type": "conversation_created",
		"data": map[string]any{
			"conversation_id": result.ConversationID,
			"user_message":    result.UserMessage,
		},
	}, chunks, complete, errCh)
}

func (h *ChatConversationHandler) GetConversation(c *echo.Context) error {
	id := c.Param("id")

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	conversation, err := h.chatConversationService.GetConversationByID(c.Request().Context(), id, userID)
	if err != nil {
		return h.respondChatError(c, "Failed to get conversation", err)
	}

	return response.Success(c, "Successfully retrieved conversation", conversation)
}

func (h *ChatConversationHandler) GetConversations(c *echo.Context) error {
	limit, offset := ParsePaginationParams(c, 10)

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	conversations, total, err := h.chatConversationService.GetUserConversations(c.Request().Context(), userID, offset, limit)
	if err != nil {
		return h.respondChatError(c, "Failed to get conversations", err)
	}

	meta := response.CalculatePaginationMeta(total, offset, limit)
	return response.SuccessWithMeta(c, "Successfully retrieved conversations", conversations, meta)
}

func (h *ChatConversationHandler) UpdateConversation(c *echo.Context) error {
	id := c.Param("id")
	var updateDTO dto.UpdateChatConversationRequest
	if err := c.Bind(&updateDTO); err != nil {
		return response.BadRequest(c, "Failed to update conversation", err)
	}

	if err := c.Validate(updateDTO); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	updatedConversation, err := h.chatConversationService.UpdateConversation(c.Request().Context(), id, userID, &updateDTO)
	if err != nil {
		return h.respondChatError(c, "Failed to update conversation", err)
	}

	return response.Success(c, "Conversation updated successfully", updatedConversation)
}

func (h *ChatConversationHandler) DeleteConversation(c *echo.Context) error {
	id := c.Param("id")

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	err := h.chatConversationService.DeleteConversation(c.Request().Context(), id, userID)
	if err != nil {
		return h.respondChatError(c, "Failed to delete conversation", err)
	}

	return response.Success(c, "Successfully deleted conversation", nil)
}

func (h *ChatConversationHandler) CreateMessage(c *echo.Context) error {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		return response.BadRequest(c, "Conversation ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateChatMessageRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Failed to create message", err)
	}
	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	messages, err := h.chatConversationService.CreateMessage(c.Request().Context(), userID, conversationID, &req)
	if err != nil {
		return h.respondChatError(c, "Failed to create message", err)
	}
	return response.Created(c, "Messages created successfully", messages)
}

func (h *ChatConversationHandler) CreateMessageStream(c *echo.Context) error {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		return response.BadRequest(c, "Conversation ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateChatMessageRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Failed to create message stream", err)
	}
	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	result, chunks, complete, errCh, err := h.chatConversationService.CreateStreamingMessage(c.Request().Context(), userID, conversationID, &req)
	if err != nil {
		return h.respondChatError(c, "Failed to create message stream", err)
	}
	if chunks == nil {
		return response.Created(c, "Message created successfully", []any{result.UserMessage})
	}
	return streamChatEvents(c, map[string]any{
		"type": "user_message",
		"data": result.UserMessage,
	}, chunks, complete, errCh)
}

func (h *ChatConversationHandler) GetMessages(c *echo.Context) error {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		return response.BadRequest(c, "Conversation ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	messages, err := h.chatConversationService.GetMessages(c.Request().Context(), conversationID, userID)
	if err != nil {
		return h.respondChatError(c, "Failed to get messages", err)
	}
	return response.Success(c, "Messages fetched successfully", messages)
}

func (h *ChatConversationHandler) GetMessage(c *echo.Context) error {
	id := c.Param("messageId")
	if id == "" {
		return response.BadRequest(c, "Message ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	message, err := h.chatConversationService.GetMessage(c.Request().Context(), id, userID)
	if err != nil {
		return h.respondChatError(c, "Failed to get message", err)
	}
	return response.Success(c, "Message fetched successfully", message)
}

func (h *ChatConversationHandler) DeleteMessage(c *echo.Context) error {
	id := c.Param("messageId")
	if id == "" {
		return response.BadRequest(c, "Message ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	message, err := h.chatConversationService.DeleteMessage(c.Request().Context(), id, userID)
	if err != nil {
		return h.respondChatError(c, "Failed to delete message", err)
	}
	return response.Success(c, "Message deleted successfully", message)
}

func (h *ChatConversationHandler) respondChatError(c *echo.Context, message string, err error) error {
	switch {
	case errors.Is(err, apperrors.ErrChatConversationNotFound), errors.Is(err, apperrors.ErrChatMessageNotFound):
		return response.NotFound(c, message, err)
	case errors.Is(err, apperrors.ErrConversationNotOwned):
		return response.Forbidden(c, err.Error())
	default:
		return response.InternalServerError(c, message, err)
	}
}

func streamChatEvents(c *echo.Context, initialEvent map[string]any, chunks <-chan string, complete <-chan dto.ChatMessageResponse, errCh <-chan error) error {
	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.WriteHeader(http.StatusOK)

	flusher, ok := res.(http.Flusher)
	if !ok {
		return response.InternalServerError(c, "Streaming is not supported", nil)
	}

	if initialEvent != nil {
		if err := writeSSE(res, initialEvent); err != nil {
			return err
		}
		flusher.Flush()
	}

	for chunk := range chunks {
		if err := writeSSE(res, map[string]any{
			"type": "ai_chunk",
			"data": chunk,
		}); err != nil {
			return err
		}
		flusher.Flush()
	}

	select {
	case err, ok := <-errCh:
		if ok && err != nil {
			_ = writeSSE(res, map[string]any{
				"type": "error",
				"data": "Failed to generate AI response",
			})
			flusher.Flush()
			return nil
		}
	default:
	}

	if msg, ok := <-complete; ok {
		if err := writeSSE(res, map[string]any{
			"type": "ai_complete",
			"data": msg,
		}); err != nil {
			return err
		}
		flusher.Flush()
	}

	_, _ = res.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
	return nil
}

func writeSSE(w http.ResponseWriter, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
	return err
}
