package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
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
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Post ID is required",
			"success": false,
		})
	}

	var dto model.CreatePostCommentDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"success": false,
		})
	}

	if err := c.Validate(dto); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"message": "Validation failed",
				"errors":  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  []string{err.Error()},
		})
	}

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	comment, err := h.commentService.CreateComment(c.Request().Context(), postID, &dto, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data":    comment,
		"success": true,
		"message": "Comment created successfully",
	})
}

// GetCommentsByPostID handles getting all comments for a specific post
func (h *CommentHandler) GetCommentsByPostID(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Post ID is required",
			"success": false,
			"message": "Post ID is required",
		})
	}

	comments, err := h.commentService.GetCommentsByPostID(c.Request().Context(), postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"success": false,
			"message": "Comments fetched failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    comments,
		"success": true,
		"message": "Comments fetched successfully",
	})
}

// UpdateComment handles updating a comment
func (h *CommentHandler) UpdateComment(c echo.Context) error {
	commentID := c.Param("comment_id")
	if commentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Comment ID is required",
			"success": false,
		})
	}

	var dto model.CreatePostCommentDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"success": false,
		})
	}

	if err := c.Validate(dto); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"message": "Validation failed",
				"errors":  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  []string{err.Error()},
		})
	}

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	comment, err := h.commentService.UpdateComment(c.Request().Context(), commentID, dto.Content, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    comment,
		"success": true,
		"message": "Comment updated successfully",
	})
}

// DeleteComment handles deleting a comment
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	commentID := c.Param("comment_id")
	if commentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Comment ID is required",
			"success": false,
		})
	}

	// Get user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	if err := h.commentService.DeleteComment(c.Request().Context(), commentID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Comment deleted successfully",
	})
}
