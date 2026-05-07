package dto

import (
	"echobackend/internal/model"
	"time"
)

type PostLikeResponse struct {
	ID        string        `json:"id"`
	PostID    string        `json:"post_id"`
	UserID    string        `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt *time.Time    `json:"created_at"`
}

func PostLikeToResponse(pl *model.PostLike) *PostLikeResponse {
	if pl == nil {
		return nil
	}
	var userResp *UserResponse
	if pl.User != nil && pl.User.ID != "" {
		userResp = UserToResponse(pl.User)
	}
	return &PostLikeResponse{
		ID:        pl.ID,
		PostID:    pl.PostID,
		UserID:    pl.UserID,
		User:      userResp,
		CreatedAt: pl.CreatedAt,
	}
}

type PostLikeStats struct {
	PostID     string `json:"post_id"`
	TotalLikes int64  `json:"total_likes"`
}
