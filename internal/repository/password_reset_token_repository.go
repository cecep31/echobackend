package repository

import (
	"context"
	"time"

	"echobackend/internal/model"

	"gorm.io/gorm"
)

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *model.PasswordResetToken) error
	FindByToken(ctx context.Context, token string) (*model.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *model.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *passwordResetTokenRepository) FindByToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	var prt model.PasswordResetToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&prt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &prt, nil
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *passwordResetTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.PasswordResetToken{}).Error
}