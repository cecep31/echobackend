package repository

import (
	"context"

	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"

	"echobackend/internal/dto"
	"gorm.io/gorm"
)

type UserFollowRepository interface {
	Follow(ctx context.Context, followerID, followingID string) error
	Unfollow(ctx context.Context, followerID, followingID string) error
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*model.User, int64, error)
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*model.User, int64, error)
	GetFollowStats(ctx context.Context, userID string) (*dto.UserFollowStats, error)
	UpdateFollowCounts(ctx context.Context, userID string) error
	GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*model.User, error)
}

type userFollowRepository struct {
	db *gorm.DB
}

func NewUserFollowRepository(db *gorm.DB) UserFollowRepository {
	return &userFollowRepository{db: db}
}

func (r *userFollowRepository) Follow(ctx context.Context, followerID, followingID string) error {
	isFollowing, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return apperrors.ErrAlreadyFollowing
	}

	if followerID == followingID {
		return apperrors.ErrCannotFollowSelf
	}

	follow := &model.UserFollow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(follow).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", followerID).
			Update("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", followingID).
			Update("followers_count", gorm.Expr("followers_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userFollowRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	isFollowing, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return apperrors.ErrNotFollowing
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).
			Delete(&model.UserFollow{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", followerID).
			Update("following_count", gorm.Expr("following_count - 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", followingID).
			Update("followers_count", gorm.Expr("followers_count - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userFollowRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserFollow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *userFollowRepository) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.UserFollow{}).
		Where("following_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("JOIN user_follows ON users.id = user_follows.follower_id").
		Where("user_follows.following_id = ?", userID).
		Order("user_follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *userFollowRepository) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.UserFollow{}).
		Where("follower_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("JOIN user_follows ON users.id = user_follows.following_id").
		Where("user_follows.follower_id = ?", userID).
		Order("user_follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *userFollowRepository) GetFollowStats(ctx context.Context, userID string) (*dto.UserFollowStats, error) {
	stats := &dto.UserFollowStats{UserID: userID}

	if err := r.db.WithContext(ctx).Model(&model.UserFollow{}).
		Where("following_id = ?", userID).Count(&stats.FollowersCount).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&model.UserFollow{}).
		Where("follower_id = ?", userID).Count(&stats.FollowingCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *userFollowRepository) UpdateFollowCounts(ctx context.Context, userID string) error {
	stats, err := r.GetFollowStats(ctx, userID)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"followers_count": stats.FollowersCount,
			"following_count": stats.FollowingCount,
		}).Error
}

func (r *userFollowRepository) GetMutualFollows(ctx context.Context, userID1, userID2 string) ([]*model.User, error) {
	var users []*model.User

	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("JOIN user_follows uf1 ON users.id = uf1.following_id").
		Joins("JOIN user_follows uf2 ON users.id = uf2.following_id").
		Where("uf1.follower_id = ? AND uf2.follower_id = ?", userID1, userID2).
		Find(&users).Error

	return users, err
}
