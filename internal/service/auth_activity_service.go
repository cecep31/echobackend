package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type AuthActivityService interface {
	LogActivity(ctx context.Context, userID *string, activityType, status, ipAddress, userAgent string, errorMessage *string, metadata map[string]any)
	GetActivityLogs(ctx context.Context, userID, activityType string, limit, offset int) ([]*model.AuthActivityLog, int64, error)
	GetRecentActivity(ctx context.Context, userID string, limit int) ([]*model.AuthActivityLog, error)
	GetFailedLogins(ctx context.Context, since time.Time, limit, offset int) ([]*model.AuthActivityLog, int64, error)
}

type authActivityService struct {
	repo repository.AuthActivityLogRepository
}

func NewAuthActivityService(repo repository.AuthActivityLogRepository) AuthActivityService {
	return &authActivityService{repo: repo}
}

func (s *authActivityService) LogActivity(ctx context.Context, userID *string, activityType, status, ipAddress, userAgent string, errorMessage *string, metadata map[string]any) {
	var metadataStr *string
	if metadata != nil {
		bytes, err := json.Marshal(metadata)
		if err == nil {
			str := string(bytes)
			metadataStr = &str
		}
	}

	entry := &model.AuthActivityLog{
		UserID:       userID,
		ActivityType: activityType,
		Status:        status,
		IPAddress:    &ipAddress,
		UserAgent:    &userAgent,
		ErrorMessage: errorMessage,
		Metadata:     metadataStr,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		slog.Error("failed to log auth activity", "error", err, "activity_type", activityType)
	}
}

func (s *authActivityService) GetActivityLogs(ctx context.Context, userID, activityType string, limit, offset int) ([]*model.AuthActivityLog, int64, error) {
	return s.repo.GetByUserID(ctx, userID, activityType, limit, offset)
}

func (s *authActivityService) GetRecentActivity(ctx context.Context, userID string, limit int) ([]*model.AuthActivityLog, error) {
	return s.repo.GetRecentByUserID(ctx, userID, limit)
}

func (s *authActivityService) GetFailedLogins(ctx context.Context, since time.Time, limit, offset int) ([]*model.AuthActivityLog, int64, error) {
	return s.repo.GetFailedLogins(ctx, since, limit, offset)
}