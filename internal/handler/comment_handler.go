package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"echobackend/pkg/validator"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment handles creating a new comment on a post
func (h *CommentHandler) CreateComment(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	var dto model.CreatePostCommentDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err)
	}

	if err := c.Validate(dto); err != nil {
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

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	comment, err := h.commentService.CreateComment(c.Request().Context(), postID, &dto, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to create comment", err)
	}

	return response.Created(c, "Comment created successfully", comment)
}

// GetCommentsByPostID handles getting all comments for a specific post
func (h *CommentHandler) GetCommentsByPostID(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	comments, err := h.commentService.GetCommentsByPostID(c.Request().Context(), postID)
	if err != nil {
		return response.InternalServerError(c, "Comments fetched failed", err)
	}

	return response.Success(c, "Comments fetched successfully", comments)
}

// UpdateComment handles updating a comment
func (h *CommentHandler) UpdateComment(c echo.Context) error {
	commentID := c.Param("comment_id")
	if commentID == "" {
		return response.BadRequest(c, "Comment ID is required", nil)
	}

	var dto model.CreatePostCommentDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err)
	}

	if err := c.Validate(dto); err != nil {
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

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	comment, err := h.commentService.UpdateComment(c.Request().Context(), commentID, dto.Content, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to update comment", err)
	}

	return response.Success(c, "Comment updated successfully", comment)
}

// DeleteComment handles deleting a comment
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	commentID := c.Param("comment_id")
	if commentID == "" {
		return response.BadRequest(c, "Comment ID is required", nil)
	}

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	if err := h.commentService.DeleteComment(c.Request().Context(), commentID, userID); err != nil {
		return response.InternalServerError(c, "Failed to delete comment", err)
	}

	return response.Success(c, "Comment deleted successfully", nil)
}
