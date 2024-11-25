package repository

import (
	"echobackend/internal/domain"

	"gorm.io/gorm"
)

type PostRepository interface {
	GetPosts(limit int, offset int) ([]*domain.Post, error)
	GetPostsRandom(limit int) ([]*domain.Post, error)
	GetTotalPosts() (int64, error)
	GetPostByID(id string) (*domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPosts(limit int, offset int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).
		Error

	return posts, err
}

func (r *postRepository) GetPostsBySlugAndUsername(slug string, username string) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Where("slug = ? AND creator.username = ?", slug, username).
		Find(&posts).
		Error

	return posts, err
}

func (r *postRepository) GetPostByID(id string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Where("id = ?", id).
		First(&post).
		Error

	return &post, err
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

func (r *postRepository) GetTotalPosts() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Post{}).Count(&count).Error

	return count, err
}
