package service

import (
	"context"
	"testing"

	"echobackend/internal/dto"
)

func TestPostViewService_GetMyPostsAnalytics(t *testing.T) {
	postRepo := &mockPostRepo{
		getAuthorPostStatsFn: func(ctx context.Context, userID string) (*dto.MyPostsAnalyticsSummary, error) {
			return &dto.MyPostsAnalyticsSummary{
				TotalPosts:     3,
				PublishedPosts: 2,
				TotalViews:     100,
				TotalLikes:     15,
			}, nil
		},
		getTopPostsByAuthorFn: func(ctx context.Context, userID string, limit int) ([]dto.MyPostPerformance, error) {
			title := "Top post"
			slug := "top-post"
			return []dto.MyPostPerformance{{
				ID:        validPostID,
				Title:     &title,
				Slug:      &slug,
				ViewCount: 50,
				LikeCount: 10,
			}}, nil
		},
	}
	viewRepo := &mockPostViewRepo{
		countViewsByAuthorBeforeFn: func(ctx context.Context, userID, beforeDate string) (int64, error) {
			return 20, nil
		},
		getViewTrendByAuthorFn: func(ctx context.Context, userID, startDate, endDate string) ([]struct {
			Date  string
			Count int64
		}, error) {
			return []struct {
				Date  string
				Count int64
			}{
				{Date: startDate, Count: 5},
				{Date: endDate, Count: 3},
			}, nil
		},
	}

	svc := NewPostViewService(viewRepo, postRepo)
	got, err := svc.GetMyPostsAnalytics(context.Background(), validUserID, &dto.MyPostsAnalyticsQuery{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-03",
	})
	if err != nil {
		t.Fatalf("GetMyPostsAnalytics returned error: %v", err)
	}

	if got.Summary.TotalPosts != 3 || got.Summary.TotalViews != 100 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if len(got.TopPosts) != 1 || got.TopPosts[0].ViewCount != 50 {
		t.Fatalf("unexpected top posts: %+v", got.TopPosts)
	}
	if len(got.ViewTrend) != 3 {
		t.Fatalf("expected 3 trend points, got %d", len(got.ViewTrend))
	}
	if got.ViewTrend[0].Views != 5 || got.ViewTrend[0].CumulativeViews != 25 {
		t.Fatalf("unexpected first trend point: %+v", got.ViewTrend[0])
	}
	if got.ViewTrend[2].Views != 3 || got.ViewTrend[2].CumulativeViews != 28 {
		t.Fatalf("unexpected last trend point: %+v", got.ViewTrend[2])
	}
}
