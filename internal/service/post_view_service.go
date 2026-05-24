package service

import (
	"context"
	"fmt"
	"time"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type PostViewService interface {
	RecordView(ctx context.Context, postID, userID string, ipAddress, userAgent *string) error
	GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*dto.PostViewResponse, int64, error)
	GetViewStats(ctx context.Context, postID string) (*dto.PostViewStats, error)
	HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error)
	GetMyPostsAnalytics(ctx context.Context, userID string, q *dto.MyPostsAnalyticsQuery) (*dto.MyPostsAnalyticsResponse, error)
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

func (s *postViewService) RecordView(ctx context.Context, postID, userID string, ipAddress, userAgent *string) error {
	if postID == "" {
		return apperrors.ErrEmptyPostID
	}

	if _, err := s.postRepo.GetPostByID(ctx, postID); err != nil {
		return fmt.Errorf("failed to verify post existence: %w", err)
	}

	if userID != "" {
		hasViewed, err := s.postViewRepo.HasUserViewedPost(ctx, postID, userID)
		if err != nil {
			return fmt.Errorf("failed to check if user viewed post: %w", err)
		}
		if hasViewed {
			return nil
		}
	}

	now := time.Now()
	view := &model.PostView{
		PostID:    postID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if userID != "" {
		view.UserID = &userID
	}
	if ipAddress != nil && *ipAddress != "" {
		view.IPAddress = ipAddress
	}
	if userAgent != nil && *userAgent != "" {
		view.UserAgent = userAgent
	}

	if err := s.postViewRepo.CreateView(ctx, view); err != nil {
		return fmt.Errorf("failed to create view record: %w", err)
	}

	if err := s.postViewRepo.IncrementPostViewCount(ctx, postID); err != nil {
		return fmt.Errorf("failed to increment post view count: %w", err)
	}

	return nil
}

func (s *postViewService) GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*dto.PostViewResponse, int64, error) {
	if postID == "" {
		return nil, 0, apperrors.ErrEmptyPostID
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	views, total, err := s.postViewRepo.GetViewsByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get views by post ID: %w", err)
	}

	responses := make([]*dto.PostViewResponse, len(views))
	for i, view := range views {
		responses[i] = dto.PostViewToResponse(view)
	}

	return responses, total, nil
}

func (s *postViewService) GetViewStats(ctx context.Context, postID string) (*dto.PostViewStats, error) {
	if postID == "" {
		return nil, apperrors.ErrEmptyPostID
	}

	stats, err := s.postViewRepo.GetViewStats(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get view stats: %w", err)
	}

	return stats, nil
}

func (s *postViewService) HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error) {
	if postID == "" {
		return false, apperrors.ErrEmptyPostID
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

func (s *postViewService) GetMyPostsAnalytics(ctx context.Context, userID string, q *dto.MyPostsAnalyticsQuery) (*dto.MyPostsAnalyticsResponse, error) {
	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()
	if q != nil {
		if q.StartDate != "" {
			if parsed, err := time.Parse("2006-01-02", q.StartDate); err == nil {
				start = parsed
			}
		}
		if q.EndDate != "" {
			if parsed, err := time.Parse("2006-01-02", q.EndDate); err == nil {
				end = parsed
			}
		}
	}
	if start.After(end) {
		start, end = end, start
	}

	startKey := start.Format("2006-01-02")
	endKey := end.Format("2006-01-02")

	summary, err := s.postRepo.GetAuthorPostStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author post stats: %w", err)
	}

	topPosts, err := s.postRepo.GetTopPostsByAuthor(ctx, userID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top posts: %w", err)
	}
	if topPosts == nil {
		topPosts = []dto.MyPostPerformance{}
	}

	rows, err := s.postViewRepo.GetViewTrendByAuthor(ctx, userID, startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get view trend: %w", err)
	}

	rowMap := make(map[string]int64, len(rows))
	for _, row := range rows {
		rowMap[row.Date] = row.Count
	}

	cumulative, err := s.postViewRepo.CountViewsByAuthorBefore(ctx, userID, startKey)
	if err != nil {
		return nil, fmt.Errorf("failed to count views before period: %w", err)
	}

	viewTrend := make([]dto.MyPostsViewTrendPoint, 0)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		views := rowMap[dateKey]
		cumulative += views
		viewTrend = append(viewTrend, dto.MyPostsViewTrendPoint{
			Date:            dateKey,
			Views:           views,
			CumulativeViews: cumulative,
		})
	}

	return &dto.MyPostsAnalyticsResponse{
		Summary:   *summary,
		ViewTrend: viewTrend,
		TopPosts:  topPosts,
	}, nil
}
