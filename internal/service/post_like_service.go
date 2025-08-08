package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
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
	// Check if post exists
	_, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	// Check if user already liked this post
	hasLiked, err := s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return err
	}

	// If user already liked, return error
	if hasLiked {
		return errors.New("user has already liked this post")
	}

	// Create new like
	like := &model.PostLike{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.postLikeRepo.CreateLike(ctx, like)
}

func (s *postLikeService) UnlikePost(ctx context.Context, postID, userID string) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	// Check if user has liked this post
	hasLiked, err := s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return err
	}

	// If user hasn't liked, return error
	if !hasLiked {
		return errors.New("user has not liked this post")
	}

	return s.postLikeRepo.DeleteLike(ctx, postID, userID)
}

func (s *postLikeService) GetLikesByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostLike, int64, error) {
	return s.postLikeRepo.GetLikesByPostID(ctx, postID, limit, offset)
}

func (s *postLikeService) GetLikeStats(ctx context.Context, postID string) (*model.PostLikeStats, error) {
	return s.postLikeRepo.GetLikeStats(ctx, postID)
}

func (s *postLikeService) HasUserLikedPost(ctx context.Context, postID, userID string) (bool, error) {
	return s.postLikeRepo.HasUserLikedPost(ctx, postID, userID)
}