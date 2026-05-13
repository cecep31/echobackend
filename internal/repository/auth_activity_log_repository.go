package repository

import (
	"context"
	"time"

	"echobackend/internal/model"

	"gorm.io/gorm"
)

type AuthActivityLogRepository interface {
	Create(ctx context.Context, log *model.AuthActivityLog) error
	GetByUserID(ctx context.Context, userID string, activityType string, limit, offset int) ([]*model.AuthActivityLog, int64, error)
	GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.AuthActivityLog, error)
	GetFailedLogins(ctx context.Context, since time.Time, limit, offset int) ([]*model.AuthActivityLog, int64, error)
}

type authActivityLogRepository struct {
	db *gorm.DB
}

func NewAuthActivityLogRepository(db *gorm.DB) AuthActivityLogRepository {
	return &authActivityLogRepository{db: db}
}

func (r *authActivityLogRepository) Create(ctx context.Context, log *model.AuthActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *authActivityLogRepository) GetByUserID(ctx context.Context, userID string, activityType string, limit, offset int) ([]*model.AuthActivityLog, int64, error) {
	var logs []*model.AuthActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&model.AuthActivityLog{}).Where("user_id = ?", userID)
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}

func (r *authActivityLogRepository) GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.AuthActivityLog, error) {
	var logs []*model.AuthActivityLog
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *authActivityLogRepository) GetFailedLogins(ctx context.Context, since time.Time, limit, offset int) ([]*model.AuthActivityLog, int64, error) {
	var logs []*model.AuthActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&model.AuthActivityLog{}).
		Where("activity_type IN ? AND created_at >= ?", []string{model.ActivityLoginFailed}, since)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}