package repository

import (
	"context"

	"echobackend/internal/model"

	"gorm.io/gorm"
)

// SessionRepository defines operations for managing user sessions (refresh tokens).
type SessionRepository interface {
	CreateSession(ctx context.Context, s *model.Session) error
	GetByRefreshToken(ctx context.Context, token string) (*model.Session, error)
	DeleteByRefreshToken(ctx context.Context, token string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, s *model.Session) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, token string) (*model.Session, error) {
	var sess model.Session
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&sess).Error; err != nil {
		return nil, err
	}
	return &sess, nil
}

func (r *sessionRepository) DeleteByRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("refresh_token = ?", token).Delete(&model.Session{}).Error
}
