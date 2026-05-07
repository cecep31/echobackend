package repository

import (
	"context"
	"fmt"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	FindAll(ctx context.Context) ([]model.Tag, error)
	FindByID(ctx context.Context, id uint) (*model.Tag, error)
	FindByName(ctx context.Context, name string) (*model.Tag, error)
	GetTagsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapTag, error)
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return apperrors.ErrTagNameRequired
	}
	result := r.db.WithContext(ctx).Create(tag)
	if result.Error != nil {
		return fmt.Errorf("failed to create tag: %w", result.Error)
	}
	return nil
}

func (r *tagRepository) FindAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all tags: %w", err)
	}
	return tags, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uint) (*model.Tag, error) {
	if id == 0 {
		return nil, apperrors.ErrInvalidTagID
	}
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to find tag by ID %d: %w", id, err)
	}
	return &tag, nil
}

func (r *tagRepository) FindByName(ctx context.Context, name string) (*model.Tag, error) {
	if name == "" {
		return nil, apperrors.ErrTagNameEmpty
	}
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to find tag by name %s: %w", name, err)
	}
	return &tag, nil
}

func (r *tagRepository) GetTagsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapTag, error) {
	var sitemapTags []*dto.SitemapTag

	err := r.db.WithContext(ctx).
		Table("tags").
		Select("tags.name, tags.created_at").
		Joins("INNER JOIN posts_to_tags ON posts_to_tags.tag_id = tags.id").
		Joins("INNER JOIN posts ON posts.id = posts_to_tags.post_id").
		Where("posts.published = ?", true).
		Group("tags.id, tags.name, tags.created_at").
		Order("tags.name ASC").
		Limit(limit).
		Find(&sitemapTags).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tags for sitemap: %w", err)
	}

	return sitemapTags, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return apperrors.ErrTagNameRequired
	}
	if tag.ID == 0 {
		return apperrors.ErrInvalidTagID
	}

	result := r.db.WithContext(ctx).Save(tag)
	if result.Error != nil {
		return fmt.Errorf("failed to update tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrTagNotFound
	}
	return nil
}

func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperrors.ErrInvalidTagID
	}
	result := r.db.WithContext(ctx).Delete(&model.Tag{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrTagNotFound
	}
	return nil
}
