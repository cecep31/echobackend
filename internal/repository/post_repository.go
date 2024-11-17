package repository

import (
	"context"
	"echobackend/internal/domain"

	"gorm.io/gorm"
)

type PostRepository interface {
	GetPosts(ctx context.Context) ([]*domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	var posts []*domain.Post
	return posts, r.db.WithContext(ctx).Preload("Creator").Preload("Tags").Find(&posts).Error
}
