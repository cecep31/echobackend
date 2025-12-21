package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
	"fmt"
	"time"
)

// PostViewService defines the interface for post view operations
type PostViewService interface {
	// RecordView records a view for a post with optional IP address and user agent
	RecordView(ctx context.Context, postID, userID string, ipAddress, userAgent *string) error
	// GetViewsByPostID retrieves paginated views for a specific post
	GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error)
	// GetViewStats retrieves aggregated view statistics for a post
	GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error)
	// HasUserViewedPost checks if a user has viewed a specific post
	HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error)
}

type postViewService struct {
	postViewRepo repository.PostViewRepository
	postRepo     repository.PostRepository
}

// NewPostViewService creates a new instance of PostViewService
func NewPostViewService(
	postViewRepo repository.PostViewRepository,
	postRepo repository.PostRepository,
) PostViewService {
	return &postViewService{
		postViewRepo: postViewRepo,
		postRepo:     postRepo,
	}
}

// RecordView records a view for a post, preventing duplicate views from authenticated users
func (s *postViewService) RecordView(ctx context.Context, postID, userID string, ipAddress, userAgent *string) error {
	// Validate input
	if postID == "" {
		return errors.New("post ID cannot be empty")
	}

	// Check if post exists
	if _, err := s.postRepo.GetPostByID(ctx, postID); err != nil {
		return fmt.Errorf("failed to verify post existence: %w", err)
	}

	// For authenticated users, check if they already viewed this post
	if userID != "" {
		hasViewed, err := s.postViewRepo.HasUserViewedPost(ctx, postID, userID)
		if err != nil {
			return fmt.Errorf("failed to check if user viewed post: %w", err)
		}
		// If user already viewed, don't record another view (idempotent operation)
		if hasViewed {
			return nil
		}
	}

	// Create view record
	now := time.Now()
	view := &model.PostView{
		PostID:    postID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Set optional fields
	if userID != "" {
		view.UserID = &userID
	}
	if ipAddress != nil && *ipAddress != "" {
		view.IPAddress = ipAddress
	}
	if userAgent != nil && *userAgent != "" {
		view.UserAgent = userAgent
	}

	// Record the view
	if err := s.postViewRepo.CreateView(ctx, view); err != nil {
		return fmt.Errorf("failed to create view record: %w", err)
	}

	// Increment post view count
	if err := s.postViewRepo.IncrementPostViewCount(ctx, postID); err != nil {
		return fmt.Errorf("failed to increment post view count: %w", err)
	}

	return nil
}

// GetViewsByPostID retrieves paginated views for a specific post
func (s *postViewService) GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error) {
	if postID == "" {
		return nil, 0, errors.New("post ID cannot be empty")
	}

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent excessive data retrieval
	}
	if offset < 0 {
		offset = 0
	}

	views, total, err := s.postViewRepo.GetViewsByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get views by post ID: %w", err)
	}

	return views, total, nil
}

// GetViewStats retrieves aggregated view statistics for a post
func (s *postViewService) GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error) {
	if postID == "" {
		return nil, errors.New("post ID cannot be empty")
	}

	stats, err := s.postViewRepo.GetViewStats(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get view stats: %w", err)
	}

	return stats, nil
}

// HasUserViewedPost checks if a user has viewed a specific post
func (s *postViewService) HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error) {
	if postID == "" {
		return false, errors.New("post ID cannot be empty")
	}
	if userID == "" {
		return false, nil
	}

	hasViewed, err := s.postViewRepo.HasUserViewedPost(ctx, postID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user viewed post: %w", err)
	}

	return hasViewed, nil
}
