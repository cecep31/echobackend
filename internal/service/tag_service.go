package service

import (
	"context"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type tagService struct {
	tagRepo repository.TagRepository
}

type TagService interface {
	CreateTag(ctx context.Context, tag *model.Tag) error
	GetTags(ctx context.Context) ([]model.Tag, error)
	GetTagByID(ctx context.Context, id uint) (*model.Tag, error)
	GetTagByName(ctx context.Context, name string) (*model.Tag, error)
	GetTagsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapTag, error)
	FindOrCreateByName(ctx context.Context, name string) (*model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, id uint) error
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) CreateTag(ctx context.Context, tag *model.Tag) error {
	if tag.Name == "" {
		return apperrors.ErrTagNameRequired
	}
	return s.tagRepo.Create(ctx, tag)
}

func (s *tagService) GetTags(ctx context.Context) ([]model.Tag, error) {
	return s.tagRepo.FindAll(ctx)
}

func (s *tagService) GetTagByID(ctx context.Context, id uint) (*model.Tag, error) {
	return s.tagRepo.FindByID(ctx, id)
}

func (s *tagService) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return s.tagRepo.Update(ctx, tag)
}

func (s *tagService) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	return s.tagRepo.FindByName(ctx, name)
}

func (s *tagService) GetTagsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapTag, error) {
	return s.tagRepo.GetTagsForSitemap(ctx, limit)
}

func (s *tagService) FindOrCreateByName(ctx context.Context, name string) (*model.Tag, error) {
	if name == "" {
		return nil, apperrors.ErrTagNameEmpty
	}

	tag, err := s.tagRepo.FindByName(ctx, name)
	if err == nil {
		return tag, nil
	}

	newTag := &model.Tag{
		Name: name,
	}

	err = s.tagRepo.Create(ctx, newTag)
	if err != nil {
		return nil, err
	}

	return newTag, nil
}

func (s *tagService) DeleteTag(ctx context.Context, id uint) error {
	_, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.tagRepo.Delete(ctx, id)
}
