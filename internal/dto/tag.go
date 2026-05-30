package dto

import (
	"echobackend/internal/model"
	"time"
)

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TrendingTagResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	TotalViews    int64  `json:"total_views"`
	TotalLikes    int64  `json:"total_likes"`
	TrendingScore int64  `json:"trending_score"`
}

func TagToResponse(t *model.Tag) *TagResponse {
	if t == nil {
		return nil
	}
	return &TagResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}

type SitemapTag struct {
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}
