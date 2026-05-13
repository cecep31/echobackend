package service

import (
	"context"
	"math"
	"time"

	"echobackend/internal/dto"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type ReportService interface {
	GetOverviewStats(ctx context.Context) (*dto.OverviewStatsResponse, error)
	GetUserReport(ctx context.Context, q dto.DateRangeQuery, limit int) (*dto.UserReportResponse, error)
	GetPostReport(ctx context.Context, q dto.DateRangeQuery, limit int, tagID *int) (*dto.PostReportResponse, error)
	GetEngagementMetrics(ctx context.Context, q dto.DateRangeQuery) (*dto.EngagementMetricsResponse, error)
}

type reportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportService{db: db}
}

func (s *reportService) GetOverviewStats(ctx context.Context) (*dto.OverviewStatsResponse, error) {
	today := time.Now().Format("2006-01-02")
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	type result struct{ Count int64 }
	var totalUsers, totalPosts, totalViews, totalLikes, totalComments result
	var newUsersToday, newPostsToday, activeUsersThisWeek result

	if err := s.db.WithContext(ctx).Model(&model.User{}).Count(&totalUsers.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Where("published = ?", true).Count(&totalPosts.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.PostView{}).Count(&totalViews.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.PostLike{}).Count(&totalLikes.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.PostComment{}).Count(&totalComments.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("DATE(created_at) >= ?", today).Count(&newUsersToday.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Where("published = ? AND DATE(created_at) >= ?", true, today).Count(&newPostsToday.Count).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Table("post_views").Select("COUNT(DISTINCT user_id) AS count").Where("user_id IS NOT NULL AND DATE(created_at) >= ?", weekAgo).Scan(&activeUsersThisWeek).Error; err != nil {
		return nil, err
	}

	return &dto.OverviewStatsResponse{
		TotalUsers:          totalUsers.Count,
		TotalPosts:          totalPosts.Count,
		TotalViews:          totalViews.Count,
		TotalLikes:          totalLikes.Count,
		TotalComments:       totalComments.Count,
		NewUsersToday:       newUsersToday.Count,
		NewPostsToday:       newPostsToday.Count,
		ActiveUsersThisWeek: activeUsersThisWeek.Count,
	}, nil
}

func (s *reportService) GetUserReport(ctx context.Context, q dto.DateRangeQuery, limit int) (*dto.UserReportResponse, error) {
	var totalUsers, newUsers, activeUsers int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	newUsersQuery := s.db.WithContext(ctx).Model(&model.User{})
	if q.StartDate != "" {
		newUsersQuery = newUsersQuery.Where("DATE(created_at) >= ?", q.StartDate)
	}
	if q.EndDate != "" {
		newUsersQuery = newUsersQuery.Where("DATE(created_at) <= ?", q.EndDate)
	}
	if err := newUsersQuery.Count(&newUsers).Error; err != nil {
		return nil, err
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := s.db.WithContext(ctx).Table("post_views").Select("COUNT(DISTINCT user_id)").Where("user_id IS NOT NULL AND created_at >= ?", thirtyDaysAgo).Scan(&activeUsers).Error; err != nil {
		return nil, err
	}

	type contributorRow struct {
		ID         string
		Username   *string
		FirstName  *string
		LastName   *string
		PostCount  int64
		TotalViews int64
		TotalLikes int64
	}
	var topRows []contributorRow
	if err := s.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.username, users.first_name, users.last_name, COUNT(posts.id) AS post_count, COALESCE(SUM(posts.view_count), 0) AS total_views, COALESCE(SUM(posts.like_count), 0) AS total_likes").
		Joins("LEFT JOIN posts ON users.id = posts.created_by AND posts.deleted_at IS NULL").
		Group("users.id, users.username, users.first_name, users.last_name").
		Order("COUNT(posts.id) DESC").
		Limit(limit).
		Scan(&topRows).Error; err != nil {
		return nil, err
	}

	growthTrend, err := s.getUserGrowthTrend(ctx, q)
	if err != nil {
		return nil, err
	}

	topContributors := make([]dto.TopContributor, 0, len(topRows))
	for _, row := range topRows {
		topContributors = append(topContributors, dto.TopContributor{
			ID:         row.ID,
			Username:   row.Username,
			FirstName:  row.FirstName,
			LastName:   row.LastName,
			PostCount:  row.PostCount,
			TotalViews: row.TotalViews,
			TotalLikes: row.TotalLikes,
		})
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

	type dayCount struct {
		Date  string
		Count int64
	}
	var rows []dayCount
	if err := s.db.WithContext(ctx).Table("users").
		Select("DATE(created_at) AS date, COUNT(*) AS count").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", start.Format("2006-01-02"), end.Format("2006-01-02")).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	rowMap := make(map[string]int64, len(rows))
	for _, row := range rows {
		rowMap[row.Date] = row.Count
	}

	var cumulative int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("DATE(created_at) < ?", start.Format("2006-01-02")).Count(&cumulative).Error; err != nil {
		return nil, err
	}

	var result []dto.UserGrowthData
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		newUsers := rowMap[dateKey]
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
	var totalPosts, newPosts, totalComments int64
	type sumResult struct{ Total int64 }
	var totalViews, totalLikes sumResult

	if err := s.db.WithContext(ctx).Model(&model.Post{}).Where("published = ?", true).Count(&totalPosts).Error; err != nil {
		return nil, err
	}

	newPostsQuery := s.db.WithContext(ctx).Model(&model.Post{}).Where("published = ?", true)
	if q.StartDate != "" {
		newPostsQuery = newPostsQuery.Where("DATE(created_at) >= ?", q.StartDate)
	}
	if q.EndDate != "" {
		newPostsQuery = newPostsQuery.Where("DATE(created_at) <= ?", q.EndDate)
	}
	if err := newPostsQuery.Count(&newPosts).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Select("COALESCE(SUM(view_count), 0) AS total").Where("published = ?", true).Scan(&totalViews).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Select("COALESCE(SUM(like_count), 0) AS total").Where("published = ?", true).Scan(&totalLikes).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.PostComment{}).Count(&totalComments).Error; err != nil {
		return nil, err
	}

	type postRow struct {
		ID              string
		Title           *string
		Slug            *string
		Views           int64
		Likes           int64
		CreatedAt       *string
		AuthorID        string
		AuthorUsername  *string
		AuthorFirstName *string
		AuthorLastName  *string
	}
	query := s.db.WithContext(ctx).
		Table("posts").
		Select("posts.id, posts.title, posts.slug, posts.view_count AS views, posts.like_count AS likes, posts.created_at, users.id AS author_id, users.username AS author_username, users.first_name AS author_first_name, users.last_name AS author_last_name").
		Joins("INNER JOIN users ON posts.created_by = users.id").
		Where("posts.published = ?", true).
		Order("posts.view_count DESC").
		Limit(limit)
	if tagID != nil {
		query = query.Joins("INNER JOIN posts_to_tags ON posts.id = posts_to_tags.post_id").Where("posts_to_tags.tag_id = ?", *tagID)
	}
	var topRows []postRow
	if err := query.Scan(&topRows).Error; err != nil {
		return nil, err
	}

	postIDs := make([]string, 0, len(topRows))
	for _, row := range topRows {
		postIDs = append(postIDs, row.ID)
	}

	type commentCountRow struct {
		PostID string
		Count  int64
	}
	commentCountMap := map[string]int64{}
	if len(postIDs) > 0 {
		var commentCounts []commentCountRow
		if err := s.db.WithContext(ctx).
			Table("post_comments").
			Select("post_id, COUNT(*) AS count").
			Where("post_id IN ?", postIDs).
			Group("post_id").
			Scan(&commentCounts).Error; err != nil {
			return nil, err
		}
		for _, row := range commentCounts {
			commentCountMap[row.PostID] = row.Count
		}
	}

	topPosts := make([]dto.PostPerformanceData, 0, len(topRows))
	for _, row := range topRows {
		comments := commentCountMap[row.ID]
		engagementRate := 0.0
		if row.Views > 0 {
			engagementRate = math.Round((float64(row.Likes+comments)/float64(row.Views))*10000) / 100
		}
		topPosts = append(topPosts, dto.PostPerformanceData{
			ID:             row.ID,
			Title:          row.Title,
			Slug:           row.Slug,
			Views:          row.Views,
			Likes:          row.Likes,
			Comments:       comments,
			EngagementRate: engagementRate,
			Author: dto.PostPerformanceAuthor{
				ID:        row.AuthorID,
				Username:  row.AuthorUsername,
				FirstName: row.AuthorFirstName,
				LastName:  row.AuthorLastName,
			},
			CreatedAt: row.CreatedAt,
		})
	}

	type tagRow struct {
		ID         int
		Name       string
		PostCount  int64
		TotalViews int64
		TotalLikes int64
	}
	var tagRows []tagRow
	if err := s.db.WithContext(ctx).
		Table("tags").
		Select("tags.id, tags.name, COUNT(posts_to_tags.post_id) AS post_count, COALESCE(SUM(posts.view_count), 0) AS total_views, COALESCE(SUM(posts.like_count), 0) AS total_likes").
		Joins("INNER JOIN posts_to_tags ON tags.id = posts_to_tags.tag_id").
		Joins("INNER JOIN posts ON posts_to_tags.post_id = posts.id").
		Where("posts.published = ?", true).
		Group("tags.id, tags.name").
		Order("COUNT(posts_to_tags.post_id) DESC").
		Limit(10).
		Scan(&tagRows).Error; err != nil {
		return nil, err
	}

	tagPerformance := make([]dto.TagPerformance, 0, len(tagRows))
	for _, row := range tagRows {
		tagPerformance = append(tagPerformance, dto.TagPerformance(row))
	}

	avgEngagementRate := 0.0
	if totalViews.Total > 0 {
		avgEngagementRate = math.Round((float64(totalLikes.Total+totalComments)/float64(totalViews.Total))*10000) / 100
	}

	return &dto.PostReportResponse{
		TotalPosts:         totalPosts,
		NewPostsThisPeriod: newPosts,
		TotalViews:         totalViews.Total,
		TotalLikes:         totalLikes.Total,
		TotalComments:      totalComments,
		AvgEngagementRate:  avgEngagementRate,
		TopPosts:           topPosts,
		TagPerformance:     tagPerformance,
	}, nil
}

func (s *reportService) GetEngagementMetrics(ctx context.Context, q dto.DateRangeQuery) (*dto.EngagementMetricsResponse, error) {
	var currentLikes, currentComments, totalPosts, totalViews, prevLikes int64

	likesQuery := s.db.WithContext(ctx).Model(&model.PostLike{})
	commentsQuery := s.db.WithContext(ctx).Model(&model.PostComment{})
	if q.StartDate != "" {
		likesQuery = likesQuery.Where("created_at >= ?", q.StartDate)
		commentsQuery = commentsQuery.Where("created_at >= ?", q.StartDate)
	}
	if q.EndDate != "" {
		endDateTime, _ := time.Parse("2006-01-02", q.EndDate)
		if !endDateTime.IsZero() {
			likesQuery = likesQuery.Where("created_at <= ?", endDateTime.Add(24*time.Hour))
			commentsQuery = commentsQuery.Where("created_at <= ?", endDateTime.Add(24*time.Hour))
		}
	}
	if err := likesQuery.Count(&currentLikes).Error; err != nil {
		return nil, err
	}
	if err := commentsQuery.Count(&currentComments).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Where("published = ?", true).Count(&totalPosts).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&model.Post{}).Select("COALESCE(SUM(view_count), 0)").Where("published = ?", true).Scan(&totalViews).Error; err != nil {
		return nil, err
	}

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
	if err := s.db.WithContext(ctx).Model(&model.PostLike{}).Where("created_at >= ? AND created_at <= ?", prevPeriodStart, prevPeriodEnd).Count(&prevLikes).Error; err != nil {
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
