package service

import (
	"context"
	"encoding/json"
	"time"

	"echobackend/internal/dto"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	GetNotifications(ctx context.Context, userID string, filter *dto.NotificationListFilter) ([]*dto.NotificationResponse, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (*dto.NotificationUnreadCountResponse, error)
	MarkAsRead(ctx context.Context, id, userID string) (*dto.NotificationResponse, error)
	MarkAllAsRead(ctx context.Context, userID string) (*dto.NotificationMarkAllReadResponse, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) CreateNotification(ctx context.Context, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	var encodedData *string
	if len(req.Data) > 0 {
		payload, err := json.Marshal(req.Data)
		if err != nil {
			return nil, err
		}
		payloadStr := string(payload)
		encodedData = &payloadStr
	}

	notification := &model.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    encodedData,
	}
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}
	created, err := s.notificationRepo.GetByID(ctx, notification.ID, req.UserID)
	if err != nil {
		return nil, err
	}
	return dto.NotificationToResponse(created), nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID string, filter *dto.NotificationListFilter) ([]*dto.NotificationResponse, int64, error) {
	notifications, total, err := s.notificationRepo.GetByUser(ctx, userID, filter.Unread, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*dto.NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, dto.NotificationToResponse(notification))
	}
	return result, total, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID string) (*dto.NotificationUnreadCountResponse, error) {
	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.NotificationUnreadCountResponse{UnreadCount: count}, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, id, userID string) (*dto.NotificationResponse, error) {
	notification, err := s.notificationRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	notification.Read = true
	now := time.Now()
	notification.UpdatedAt = &now
	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		return nil, err
	}
	updated, err := s.notificationRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return dto.NotificationToResponse(updated), nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) (*dto.NotificationMarkAllReadResponse, error) {
	updated, err := s.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.NotificationMarkAllReadResponse{UpdatedCount: updated}, nil
}
