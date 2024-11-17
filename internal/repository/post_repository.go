package repository

import (
	"context"
	"echobackend/internal/domain"

	"gorm.io/gorm"
)

type PostRepository interface {
	GetPosts(ctx context.Context) ([]*domain.Post, error)
	GetPostsRandom(limit int) ([]*domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Find(&posts).
		Error

	return posts, err
}

func (r *postRepository) GetPostsRandom(limit int) ([]*domain.Post, error) {
	var randomPosts []*domain.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Order("RANDOM()").
		Limit(limit).
		Find(&randomPosts).
		Error

	return randomPosts, err
}
