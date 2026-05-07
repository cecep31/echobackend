package dto

import (
	"echobackend/internal/model"
	"time"
)

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
