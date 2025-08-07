package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"time"
)

type PostViewService interface {
	RecordView(ctx context.Context, postID, userID string) error
	GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error)
	GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error)
	HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error)
}

type postViewService struct {
	postViewRepo repository.PostViewRepository
	postRepo     repository.PostRepository
}

func NewPostViewService(
	postViewRepo repository.PostViewRepository,
	postRepo repository.PostRepository,
) PostViewService {
	return &postViewService{
		postViewRepo: postViewRepo,
		postRepo:     postRepo,
	}
}

func (s *postViewService) RecordView(ctx context.Context, postID, userID string) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	// For authenticated users, check if they already viewed this post
	if userID != "" {
		hasViewed, err := s.postViewRepo.HasUserViewedPost(ctx, postID, userID)
		if err != nil {
			return err
		}
		// If user already viewed, don't record another view
		if hasViewed {
			return nil
		}
	}

	// Create view record
	view := &model.PostView{
		PostID:    postID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set user ID if authenticated
	if userID != "" {
		view.UserID = userID
	}

	// Record the view
	if err := s.postViewRepo.CreateView(ctx, view); err != nil {
		return err
	}

	// Increment post view count
	return s.postViewRepo.IncrementPostViewCount(ctx, postID)
}

func (s *postViewService) GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error) {
	return s.postViewRepo.GetViewsByPostID(ctx, postID, limit, offset)
}

func (s *postViewService) GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error) {
	return s.postViewRepo.GetViewStats(ctx, postID)
}

func (s *postViewService) HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error) {
	if userID == "" {
		return false, nil
	}
	return s.postViewRepo.HasUserViewedPost(ctx, postID, userID)
}
