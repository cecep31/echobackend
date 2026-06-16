package service

import (
	"context"
	"math"
	"time"

	"echobackend/internal/dto"
	"echobackend/internal/repository"
)

type ReportService interface {
	GetOverviewStats(ctx context.Context) (*dto.OverviewStatsResponse, error)
	GetUserReport(ctx context.Context, q dto.DateRangeQuery, limit int) (*dto.UserReportResponse, error)
	GetPostReport(ctx context.Context, q dto.DateRangeQuery, limit int, tagID *int) (*dto.PostReportResponse, error)
	GetEngagementMetrics(ctx context.Context, q dto.DateRangeQuery) (*dto.EngagementMetricsResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GetOverviewStats(ctx context.Context) (*dto.OverviewStatsResponse, error) {
	totalUsers, totalPosts, totalViews, totalLikes, totalComments, newUsersToday, newPostsToday, activeUsersThisWeek, err := s.reportRepo.GetOverviewCounts(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.OverviewStatsResponse{
		TotalUsers:          totalUsers,
		TotalPosts:          totalPosts,
		TotalViews:          totalViews,
		TotalLikes:          totalLikes,
		TotalComments:       totalComments,
		NewUsersToday:       newUsersToday,
		NewPostsToday:       newPostsToday,
		ActiveUsersThisWeek: activeUsersThisWeek,
	}, nil
}

func (s *reportService) GetUserReport(ctx context.Context, q dto.DateRangeQuery, limit int) (*dto.UserReportResponse, error) {
	totalUsers, newUsers, activeUsers, err := s.reportRepo.GetUserCounts(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return nil, err
	}

	topContributors, err := s.reportRepo.GetTopContributors(ctx, limit)
	if err != nil {
		return nil, err
	}

	growthTrend, err := s.getUserGrowthTrend(ctx, q)
	if err != nil {
		return nil, err
	}

	return &dto.UserReportResponse{
		TotalUsers:         totalUsers,
		NewUsersThisPeriod: newUsers,
		ActiveUsers:        activeUsers,
		TopContributors:    topContributors,
		GrowthTrend:        growthTrend,
	}, nil
}

func (s *reportService) getUserGrowthTrend(ctx context.Context, q dto.DateRangeQuery) ([]dto.UserGrowthData, error) {
	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()
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

	cumulative, dailyCounts, err := s.reportRepo.GetUserGrowthTrendData(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var result []dto.UserGrowthData
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		newUsers := dailyCounts[dateKey]
		cumulative += newUsers
		result = append(result, dto.UserGrowthData{
			Date:            dateKey,
			NewUsers:        newUsers,
			CumulativeUsers: cumulative,
		})
	}
	return result, nil
}

func (s *reportService) GetPostReport(ctx context.Context, q dto.DateRangeQuery, limit int, tagID *int) (*dto.PostReportResponse, error) {
	totalPosts, newPosts, totalComments, totalViews, totalLikes, err := s.reportRepo.GetPostCounts(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return nil, err
	}

	topPosts, err := s.reportRepo.GetTopPosts(ctx, limit, tagID)
	if err != nil {
		return nil, err
	}

	tagPerformance, err := s.reportRepo.GetTagPerformance(ctx, 10)
	if err != nil {
		return nil, err
	}

	// Calculate engagement rates in Service Layer
	for i := range topPosts {
		if topPosts[i].Views > 0 {
			topPosts[i].EngagementRate = math.Round((float64(topPosts[i].Likes+topPosts[i].Comments)/float64(topPosts[i].Views))*10000) / 100
		}
	}

	avgEngagementRate := 0.0
	if totalViews > 0 {
		avgEngagementRate = math.Round((float64(totalLikes+totalComments)/float64(totalViews))*10000) / 100
	}

	return &dto.PostReportResponse{
		TotalPosts:         totalPosts,
		NewPostsThisPeriod: newPosts,
		TotalViews:         totalViews,
		TotalLikes:         totalLikes,
		TotalComments:      totalComments,
		AvgEngagementRate:  avgEngagementRate,
		TopPosts:           topPosts,
		TagPerformance:     tagPerformance,
	}, nil
}

func (s *reportService) GetEngagementMetrics(ctx context.Context, q dto.DateRangeQuery) (*dto.EngagementMetricsResponse, error) {
	prevPeriodStart := time.Now().AddDate(0, 0, -60)
	prevPeriodEnd := time.Now().AddDate(0, 0, -30)
	if q.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", q.StartDate)
		if err == nil {
			endDate := time.Now()
			if q.EndDate != "" {
				if parsedEnd, parseErr := time.Parse("2006-01-02", q.EndDate); parseErr == nil {
					endDate = parsedEnd
				}
			}
			duration := endDate.Sub(startDate)
			prevPeriodEnd = startDate
			prevPeriodStart = startDate.Add(-duration)
		}
	}

	currentLikes, currentComments, totalPosts, totalViews, prevLikes, err := s.reportRepo.GetEngagementCounts(ctx, q.StartDate, q.EndDate, prevPeriodStart, prevPeriodEnd)
	if err != nil {
		return nil, err
	}

	changePercent := 0.0
	if prevLikes > 0 {
		changePercent = math.Round(((float64(currentLikes-prevLikes) / float64(prevLikes)) * 100 * 100)) / 100
	}

	avgLikes := 0.0
	avgComments := 0.0
	avgViews := 0.0
	if totalPosts > 0 {
		avgLikes = math.Round((float64(currentLikes)/float64(totalPosts))*100) / 100
		avgComments = math.Round((float64(currentComments)/float64(totalPosts))*100) / 100
		avgViews = math.Round((float64(totalViews)/float64(totalPosts))*100) / 100
	}

	return &dto.EngagementMetricsResponse{
		TotalEngagements:   currentLikes + currentComments,
		AvgLikesPerPost:    avgLikes,
		AvgCommentsPerPost: avgComments,
		AvgViewsPerPost:    avgViews,
		PeriodComparison: dto.PeriodComparison{
			Current:       currentLikes,
			Previous:      prevLikes,
			ChangePercent: changePercent,
		},
	}, nil
}
