package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
)

type tagService struct {
	tagRepo repository.TagRepository
}

type TagService interface {
	CreateTag(ctx context.Context, tag *model.Tag) error
	GetTags(ctx context.Context) ([]model.Tag, error)
	GetTagByID(ctx context.Context, id uint) (*model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, id uint) error
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) CreateTag(ctx context.Context, tag *model.Tag) error {
	if tag.Name == "" {
		return errors.New("tag name is required")
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

func (s *tagService) DeleteTag(ctx context.Context, id uint) error {
	_, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.tagRepo.Delete(ctx, id)
}
