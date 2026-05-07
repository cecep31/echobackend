package service

import (
	"context"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/repository"
)

type UserFollowService interface {
	FollowUser(ctx context.Context, followerID, followingID string) (*dto.FollowResponse, error)
	UnfollowUser(ctx context.Context, followerID, followingID string) (*dto.FollowResponse, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*dto.UserResponse, int64, error)
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*dto.UserResponse, int64, error)
	GetFollowStats(ctx context.Context, userID string) (*dto.UserFollowStats, error)
	GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*dto.UserResponse, error)
	GetUserWithFollowStatus(ctx context.Context, userID, currentUserID string) (*dto.UserResponse, error)
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

func (s *userFollowService) FollowUser(ctx context.Context, followerID, followingID string) (*dto.FollowResponse, error) {
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	_, err = s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	err = s.userFollowRepo.Follow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &dto.FollowResponse{
		IsFollowing: true,
		Message:     "Successfully followed user",
	}, nil
}

func (s *userFollowService) UnfollowUser(ctx context.Context, followerID, followingID string) (*dto.FollowResponse, error) {
	err := s.userFollowRepo.Unfollow(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &dto.FollowResponse{
		IsFollowing: false,
		Message:     "Successfully unfollowed user",
	}, nil
}

func (s *userFollowService) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	return s.userFollowRepo.IsFollowing(ctx, followerID, followingID)
}

func (s *userFollowService) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*dto.UserResponse, int64, error) {
	users, total, err := s.userFollowRepo.GetFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserToResponse(user)
	}

	return userResponses, total, nil
}

func (s *userFollowService) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*dto.UserResponse, int64, error) {
	users, total, err := s.userFollowRepo.GetFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserToResponse(user)
	}

	return userResponses, total, nil
}

func (s *userFollowService) GetFollowStats(ctx context.Context, userID string) (*dto.UserFollowStats, error) {
	return s.userFollowRepo.GetFollowStats(ctx, userID)
}

func (s *userFollowService) GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*dto.UserResponse, error) {
	users, err := s.userFollowRepo.GetMutualFollows(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserToResponse(user)
	}

	return userResponses, nil
}

func (s *userFollowService) GetUserWithFollowStatus(ctx context.Context, userID, currentUserID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserToResponse(user)

	if currentUserID != "" && currentUserID != userID {
		isFollowing, err := s.userFollowRepo.IsFollowing(ctx, currentUserID, userID)
		if err != nil {
			return nil, err
		}
		userResponse.IsFollowing = &isFollowing
	}

	return userResponse, nil
}
