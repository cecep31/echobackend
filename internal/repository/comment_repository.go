package repository

import (
	"context"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.PostComment) error
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostComment, int64, error)
	GetCommentByID(ctx context.Context, id string) (*model.PostComment, error)
	UpdateComment(ctx context.Context, comment *model.PostComment) error
	DeleteComment(ctx context.Context, id string) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *model.PostComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostComment, int64, error) {
	var comments []*model.PostComment
	var total int64

	// Count total comments for the post
	if err := r.db.WithContext(ctx).Model(&model.PostComment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get comments with creator information
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id string) (*model.PostComment, error) {
	var comment model.PostComment
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Post").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, comment *model.PostComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) DeleteComment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.PostComment{}).Error
}