package repository

import (
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	GetPosts(limit int, offset int) ([]*model.Post, error)
	GetPostsRandom(limit int) ([]*model.Post, error)
	GetTotalPosts() (int64, error)
	GetPostByID(id string) (*model.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetPosts(limit int, offset int) ([]*model.Post, error) {
	var posts []*model.Post
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

func (r *postRepository) GetPostBySlugAndUsername(slug string, username string) (*model.Post, error) {
	var post model.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Where("slug = ? AND username = ?", slug, username).
		First(&post).
		Error

	return &post, err
}

func (r *postRepository) GetPostByID(id string) (*model.Post, error) {
	var post model.Post
	err := r.db.
		Preload("Creator").
		Preload("Tags").
		Where("id = ?", id).
		First(&post).
		Error

	return &post, err
}

func (r *postRepository) GetPostsRandom(limit int) ([]*model.Post, error) {
	var randomPosts []*model.Post
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
	err := r.db.Model(&model.Post{}).Count(&count).Error

	return count, err
}
