package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"echobackend/pkg/validator"
	"errors"
	"fmt"
	"time"
)

type PostLikeService interface {
	LikePost(ctx context.Context, postID, userID string) error
	UnlikePost(ctx context.Context, postID, userID string) error
	GetLikesByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostLike, int64, error)
	GetLikeStats(ctx context.Context, postID string) (*model.PostLikeStats, error)
	HasUserLikedPost(ctx context.Context, postID, userID string) (bool, error)
}

type postLikeService struct {
	postLikeRepo repository.PostLikeRepository
	postRepo     repository.PostRepository
}

func NewPostLikeService(
	postLikeRepo repository.PostLikeRepository,
	postRepo repository.PostRepository,
) PostLikeService {
	return &postLikeService{
		postLikeRepo: postLikeRepo,
		postRepo:     postRepo,
	}
}

func (s *postLikeService) LikePost(ctx context.Context, postID, userID string) error {
	// Validate input parameters
	if err := validator.ValidatePostLikeInput(postID, userID); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check if post exists
	_, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to check post existence: %w", err)
	}

	// Check if user already liked this post
	hasLiked, err := s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}

	// If user already liked, return error
	if hasLiked {
		return errors.New("user has already liked this post")
	}

	// Create new like
	now := time.Now()
	like := &model.PostLike{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	return s.postLikeRepo.CreateLike(ctx, like)
}

func (s *postLikeService) UnlikePost(ctx context.Context, postID, userID string) error {
	// Validate input parameters
	if err := validator.ValidatePostLikeInput(postID, userID); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check if post exists
	_, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to check post existence: %w", err)
	}

	// Check if user has liked this post
	hasLiked, err := s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}

	// If user hasn't liked, return error
	if !hasLiked {
		return errors.New("user has not liked this post")
	}

	return s.postLikeRepo.DeleteLike(ctx, postID, userID)
}

func (s *postLikeService) GetLikesByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostLike, int64, error) {
	// Validate post ID
	if !validator.IsValidUUID(postID) {
		return nil, 0, fmt.Errorf("invalid post ID format")
	}

	// Validate pagination parameters
	if err := validator.ValidatePagination(limit, offset); err != nil {
		return nil, 0, fmt.Errorf("validation error: %w", err)
	}

	return s.postLikeRepo.GetLikesByPostID(ctx, postID, limit, offset)
}

func (s *postLikeService) GetLikeStats(ctx context.Context, postID string) (*model.PostLikeStats, error) {
	// Validate post ID
	if !validator.IsValidUUID(postID) {
		return nil, fmt.Errorf("invalid post ID format")
	}

	return s.postLikeRepo.GetLikeStats(ctx, postID)
}

func (s *postLikeService) HasUserLikedPost(ctx context.Context, postID, userID string) (bool, error) {
	// Validate input parameters
	if err := validator.ValidatePostLikeInput(postID, userID); err != nil {
		return false, fmt.Errorf("validation error: %w", err)
	}

	return s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
}
