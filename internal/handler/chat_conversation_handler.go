package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"echobackend/pkg/validator"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ChatConversationHandler struct {
	chatConversationService service.ChatConversationService
}

func NewChatConversationHandler(chatConversationService service.ChatConversationService) *ChatConversationHandler {
	return &ChatConversationHandler{
		chatConversationService: chatConversationService,
	}
}

func (h *ChatConversationHandler) CreateConversation(c echo.Context) error {
	var conversationReq model.CreateChatConversationDTO
	if err := c.Bind(&conversationReq); err != nil {
		return response.BadRequest(c, "Failed to create conversation", err)
	}

	if err := c.Validate(conversationReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	// Get the user ID from the JWT token
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	newConversation, err := h.chatConversationService.CreateConversation(c.Request().Context(), userID, &conversationReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to create conversation", err)
	}

	return response.Created(c, "Successfully created conversation", newConversation)
}

func (h *ChatConversationHandler) GetConversation(c echo.Context) error {
	id := c.Param("id")

	// Get the user ID from the JWT token
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	conversation, err := h.chatConversationService.GetConversationByID(c.Request().Context(), id, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get conversation", err)
	}

	return response.Success(c, "Successfully retrieved conversation", conversation)
}

func (h *ChatConversationHandler) GetConversations(c echo.Context) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0 // Default offset if not provided or invalid
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10 // Default limit if not provided or invalid
	}

	// Get the user ID from the JWT token
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	conversations, total, err := h.chatConversationService.GetUserConversations(c.Request().Context(), userID, offsetInt, limitInt)
	if err != nil {
		return response.InternalServerError(c, "Failed to get conversations", err)
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offsetInt,
		Limit:      limitInt,
		TotalPages: int(total)/limitInt + 1,
	}
	if int(total)%limitInt == 0 {
		meta.TotalPages = int(total) / limitInt
	}

	return response.SuccessWithMeta(c, "Successfully retrieved conversations", conversations, meta)
}

func (h *ChatConversationHandler) UpdateConversation(c echo.Context) error {
	id := c.Param("id")
	var updateDTO model.UpdateChatConversationDTO
	if err := c.Bind(&updateDTO); err != nil {
		return response.BadRequest(c, "Failed to update conversation", err)
	}

	if err := c.Validate(updateDTO); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	// Get the user ID from the JWT token
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	updatedConversation, err := h.chatConversationService.UpdateConversation(c.Request().Context(), id, userID, &updateDTO)
	if err != nil {
		return response.InternalServerError(c, "Failed to update conversation", err)
	}

	return response.Success(c, "Conversation updated successfully", updatedConversation)
}

func (h *ChatConversationHandler) DeleteConversation(c echo.Context) error {
	id := c.Param("id")

	// Get the user ID from the JWT token
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	err := h.chatConversationService.DeleteConversation(c.Request().Context(), id, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete conversation", err)
	}

	return response.Success(c, "Successfully deleted conversation", nil)
}
