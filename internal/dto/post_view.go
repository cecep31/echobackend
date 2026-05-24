package dto

import (
	"echobackend/internal/model"
	"time"
)

type PostViewStats struct {
	PostID             string `json:"post_id"`
	TotalViews         int64  `json:"total_views"`
	UniqueViews        int64  `json:"unique_views"`
	AnonymousViews     int64  `json:"anonymous_views"`
	AuthenticatedViews int64  `json:"authenticated_views"`
}

type PostViewResponse struct {
	ID        string     `json:"id"`
	PostID    string     `json:"post_id"`
	UserID    *string    `json:"user_id"`
	IPAddress *string    `json:"ip_address"`
	UserAgent *string    `json:"user_agent"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func PostViewToResponse(v *model.PostView) *PostViewResponse {
	if v == nil {
		return nil
	}
	return &PostViewResponse{
		ID:        v.ID,
		PostID:    v.PostID,
		UserID:    v.UserID,
		IPAddress: v.IPAddress,
		UserAgent: v.UserAgent,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
