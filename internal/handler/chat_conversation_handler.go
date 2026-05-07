package handler

import (
	"echobackend/internal/dto"
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

func (h *ChatConversationHandler) GetConversation(c *echo.Context) error {
	id := c.Param("id")

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	conversation, err := h.chatConversationService.GetConversationByID(c.Request().Context(), id, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get conversation", err)
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
		return response.InternalServerError(c, "Failed to get conversations", err)
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offset,
		Limit:      limit,
		TotalPages: int(total)/limit + 1,
	}
	if int(total)%limit == 0 {
		meta.TotalPages = int(total) / limit
	}

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
		return response.InternalServerError(c, "Failed to update conversation", err)
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
		return response.InternalServerError(c, "Failed to delete conversation", err)
	}

	return response.Success(c, "Successfully deleted conversation", nil)
}
