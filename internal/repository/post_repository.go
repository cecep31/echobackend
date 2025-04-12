package repository

import (
	"context"

	"echobackend/internal/model"

	"github.com/uptrace/bun"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.CreatePostDTO, creator_id string) (*model.Post, error)
	GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error)
	GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error)
	GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error)
	GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error)
	DeletePostByID(ctx context.Context, id string) error
}

type postRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *model.CreatePostDTO, creator_id string) (*model.Post, error) {
	newpost := &model.Post{
		Title: post.Title,
		Slug:  post.Slug,
		Body:  post.Body,
		// Tags:        post.Tags,
		CreatedBy: creator_id,
	}
	_, err := r.db.NewInsert().
		Model(newpost).
		Exec(ctx)

	return newpost, err
}

func (r *postRepository) GetPostByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	countInt, err := r.db.NewSelect().
		Model(&model.Post{}).
		Join("JOIN users ON users.id = p.created_by").
		Where("users.username = ?", username).
		Count(ctx)

	count = int64(countInt)

	if err != nil {
		return nil, 0, err
	}

	err = r.db.NewSelect().
		Model(&posts).
		Join("JOIN users ON users.id = p.created_by").
		Relation("Creator").
		Relation("Tags").
		Where("users.username = ?", username).
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (r *postRepository) DeletePostByID(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model(&model.Post{}).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (r *postRepository) GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	countInt, err := r.db.NewSelect().
		Model(&model.Post{}).
		Count(ctx)

	count = int64(countInt)

	if err != nil {
		return nil, 0, err
	}

	err = r.db.NewSelect().
		Model(&posts).
		Relation("Creator").
		Relation("Tags").
		OrderExpr("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (r *postRepository) GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.Post, error) {
	var post model.Post

	err := r.db.NewSelect().
		Model(&post).
		Relation("Creator", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("username = ?", username)
		}).
		Relation("Tags").
		Where("p.slug = ?", slug).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post

	err := r.db.NewSelect().
		Model(&post).
		Relation("Creator").
		Relation("Tags").
		Where("p.id = ?", id).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) GetPostsRandom(ctx context.Context, limit int) ([]*model.Post, error) {
	var randomPosts []*model.Post

	err := r.db.NewSelect().
		Model(&randomPosts).
		Relation("Creator").
		Relation("Tags").
		OrderExpr("RANDOM()").
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return randomPosts, nil
}

func (r *postRepository) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var count int64

	countInt, err := r.db.NewSelect().
		Model(&model.Post{}).
		Where("created_by = ?", createdBy).
		Count(ctx)

	count = int64(countInt)

	if err != nil {
		return nil, 0, err
	}

	err = r.db.NewSelect().
		Model(&posts).
		Relation("Creator").
		Relation("Tags").
		Where("created_by = ?", createdBy).
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}
