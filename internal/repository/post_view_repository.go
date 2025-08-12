package repository

import (
	"context"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type PostViewRepository interface {
	CreateView(ctx context.Context, view *model.PostView) error
	GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error)
	GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error)
	HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error)
	GetViewByUserAndPost(ctx context.Context, postID, userID string) (*model.PostView, error)
	IncrementPostViewCount(ctx context.Context, postID string) error
}

type postViewRepository struct {
	db *gorm.DB
}

func NewPostViewRepository(db *gorm.DB) PostViewRepository {
	return &postViewRepository{db: db}
}

func (r *postViewRepository) CreateView(ctx context.Context, view *model.PostView) error {
	return r.db.WithContext(ctx).Create(view).Error
}

func (r *postViewRepository) GetViewsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.PostView, int64, error) {
	var views []*model.PostView
	var total int64

	// Count total views
	if err := r.db.WithContext(ctx).Model(&model.PostView{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated views with user information
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&views).Error

	return views, total, err
}

func (r *postViewRepository) GetViewStats(ctx context.Context, postID string) (*model.PostViewStats, error) {
	stats := &model.PostViewStats{PostID: postID}

	// Total views
	if err := r.db.WithContext(ctx).Model(&model.PostView{}).Where("post_id = ?", postID).Count(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	// Unique views (distinct user_id where user_id is not null)
	if err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Where("post_id = ? AND user_id IS NOT NULL", postID).
		Distinct("user_id").
		Count(&stats.UniqueViews).Error; err != nil {
		return nil, err
	}

	// Anonymous views
	if err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Where("post_id = ? AND user_id IS NULL", postID).
		Count(&stats.AnonymousViews).Error; err != nil {
		return nil, err
	}

	// Authenticated views
	if err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Where("post_id = ? AND user_id IS NOT NULL", postID).
		Count(&stats.AuthenticatedViews).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *postViewRepository) HasUserViewedPost(ctx context.Context, postID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postViewRepository) GetViewByUserAndPost(ctx context.Context, postID, userID string) (*model.PostView, error) {
	var view model.PostView
	err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		First(&view).Error
	return &view, err
}

func (r *postViewRepository) IncrementPostViewCount(ctx context.Context, postID string) error {
	return r.db.WithContext(ctx).Model(&model.Post{}).
		Where("id = ?", postID).
		Update("view_count", r.db.Raw("view_count + 1")).Error
}
