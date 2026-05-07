package dto

type FollowRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type FollowResponse struct {
	IsFollowing bool   `json:"is_following"`
	Message     string `json:"message"`
}

type UserFollowStats struct {
	UserID         string `json:"user_id"`
	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
}
