package repository

import (
	"context"
	"fmt"
	"time"

	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, id, userID string) (*model.Notification, error)
	GetByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*model.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	Update(ctx context.Context, notification *model.Notification) error
	MarkAllAsRead(ctx context.Context, userID string) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id, userID string) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotificationNotFound
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return &notification, nil
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("\"read\" = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, total, nil
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ? AND \"read\" = ?", userID, false).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *model.Notification) error {
	result := r.db.WithContext(ctx).Save(notification)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotificationNotFound
	}
	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND \"read\" = ?", userID, false).
		Updates(map[string]any{
			"read":       true,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to mark all notifications as read: %w", result.Error)
	}
	return result.RowsAffected, nil
}
