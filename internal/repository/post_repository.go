package repository

import (
	"context"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error)
	GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error)
	GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error)
	GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error)
	DeletePostByID(ctx context.Context, id string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var posts []*model.Post
	var count int64

	if errcount := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Joins("JOIN users ON users.id = posts.created_by").
		Where("users.username = ?", username).
		Count(&count).Error; errcount != nil {
		return nil, 0, errcount
	}

	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = posts.created_by").
		Preload("Creator").
		Preload("Tags").
		Where("users.username = ?", username).
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Error
	return posts, count, err
}

func (r *postRepository) DeletePostByID(ctx context.Context, id string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	return r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Post{}).
		Error
}

func (r *postRepository) GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var posts []*model.Post

	var count int64
	if errcount := r.db.WithContext(ctx).Model(&model.Post{}).Count(&count).Error; errcount != nil {
		return nil, 0, errcount
	}

	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).
		Error

	return posts, count, err
}

func (r *postRepository) GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var post model.Post
	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = posts.created_by").
		Preload("Creator").
		Preload("Tags").
		Where("posts.slug = ?", slug).Where("users.username = ?", username).
		First(&post).
		Error

	return &post, err
}

func (r *postRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var post model.Post
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Where("id = ?", id).
		First(&post).
		Error

	return &post, err
}

func (r *postRepository) GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var randomPosts []*model.Post
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Order("RANDOM()").
		Limit(limit).
		Find(&randomPosts).
		Error

	return randomPosts, err
}

func (r *postRepository) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var posts []*model.Post
	var count int64

	if errcount := r.db.WithContext(ctx).
		Where("created_by = ?", createdBy).
		Count(&count).Error; errcount != nil {
		return nil, 0, errcount
	}

	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Where("created_by = ?", createdBy).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Error
	return posts, count, err
}
