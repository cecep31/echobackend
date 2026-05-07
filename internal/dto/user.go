package dto

import (
	"echobackend/internal/model"
	"time"
)

type UserResponse struct {
	ID             string         `json:"id"`
	Email          string         `json:"email"`
	Name           string         `json:"name"`
	Username       *string        `json:"username"`
	Image          *string        `json:"image"`
	FirstName      *string        `json:"first_name"`
	LastName       *string        `json:"last_name"`
	FollowersCount int64          `json:"followers_count"`
	FollowingCount int64          `json:"following_count"`
	IsFollowing    *bool          `json:"is_following,omitempty"`
	Profile        *model.Profile `json:"profile,omitempty"`
	CreatedAt      *time.Time     `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
}

type PublicUserResponse struct {
	ID             string         `json:"id"`
	Email          string         `json:"email"`
	Name           string         `json:"name"`
	Username       *string        `json:"username"`
	Image          *string        `json:"image"`
	FirstName      *string        `json:"first_name"`
	LastName       *string        `json:"last_name"`
	FollowersCount int64          `json:"followers_count"`
	FollowingCount int64          `json:"following_count"`
	IsFollowing    *bool          `json:"is_following,omitempty"`
	Profile        *model.Profile `json:"profile,omitempty"`
	CreatedAt      *time.Time     `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
}

func UserToResponse(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}
	name := ""
	if u.FirstName != nil && u.LastName != nil {
		name = *u.FirstName + " " + *u.LastName
	}
	return &UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           name,
		Username:       u.Username,
		Image:          u.Image,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		Profile:        u.Profile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func UserToPublicResponse(u *model.User) *PublicUserResponse {
	if u == nil {
		return nil
	}
	name := ""
	if u.FirstName != nil && u.LastName != nil {
		name = *u.FirstName + " " + *u.LastName
	}
	return &PublicUserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           name,
		Username:       u.Username,
		Image:          u.Image,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		Profile:        u.Profile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
