package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
)

type UserFollowService interface {
	FollowUser(ctx context.Context, followerID, followingID string) (*model.FollowResponse, error)
	UnfollowUser(ctx context.Context, followerID, followingID string) (*model.FollowResponse, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*model.UserResponse, int64, error)
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*model.UserResponse, int64, error)
	GetFollowStats(ctx context.Context, userID string) (*model.UserFollowStats, error)
	GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*model.UserResponse, error)
	GetUserWithFollowStatus(ctx context.Context, userID, currentUserID string) (*model.UserResponse, error)
}

type userFollowService struct {
	userFollowRepo repository.UserFollowRepository
	userRepo       repository.UserRepository
}

func NewUserFollowService(
	userFollowRepo repository.UserFollowRepository,
	userRepo repository.UserRepository,
) UserFollowService {
	return &userFollowService{
		userFollowRepo: userFollowRepo,
		userRepo:       userRepo,
	}
}

func (s *userFollowService) FollowUser(ctx context.Context, followerID, followingID string) (*model.FollowResponse, error) {
	// Check if both users exist
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return nil, errors.New("follower user not found")
	}

	_, err = s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return nil, errors.New("user to follow not found")
	}

	// Follow the user
	err = s.userFollowRepo.Follow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &model.FollowResponse{
		IsFollowing: true,
		Message:     "Successfully followed user",
	}, nil
}

func (s *userFollowService) UnfollowUser(ctx context.Context, followerID, followingID string) (*model.FollowResponse, error) {
	// Unfollow the user
	err := s.userFollowRepo.Unfollow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &model.FollowResponse{
		IsFollowing: false,
		Message:     "Successfully unfollowed user",
	}, nil
}

func (s *userFollowService) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	return s.userFollowRepo.IsFollowing(ctx, followerID, followingID)
}

func (s *userFollowService) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userFollowRepo.GetFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response format
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, total, nil
}

func (s *userFollowService) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userFollowRepo.GetFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response format
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, total, nil
}

func (s *userFollowService) GetFollowStats(ctx context.Context, userID string) (*model.UserFollowStats, error) {
	return s.userFollowRepo.GetFollowStats(ctx, userID)
}

func (s *userFollowService) GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*model.UserResponse, error) {
	users, err := s.userFollowRepo.GetMutualFollows(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, nil
}

func (s *userFollowService) GetUserWithFollowStatus(ctx context.Context, userID, currentUserID string) (*model.UserResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToResponse()

	// If current user is provided and different from target user, check follow status
	if currentUserID != "" && currentUserID != userID {
		isFollowing, err := s.userFollowRepo.IsFollowing(ctx, currentUserID, userID)
		if err != nil {
			return nil, err
		}
		userResponse.IsFollowing = &isFollowing
	}

	return userResponse, nil
}
