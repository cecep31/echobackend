package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
)

type tagService struct {
	repo repository.TagRepository
}

type TagService interface {
	CreateTag(ctx context.Context, tag *model.Tag) error
	GetTags(ctx context.Context) ([]model.Tag, error)
	GetTagByID(ctx context.Context, id uint) (*model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, id uint) error
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) CreateTag(ctx context.Context, tag *model.Tag) error {
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	return s.repo.Create(ctx, tag)
}

func (s *tagService) GetTags(ctx context.Context) ([]model.Tag, error) {
	return s.repo.FindAll(ctx)
}

func (s *tagService) GetTagByID(ctx context.Context, id uint) (*model.Tag, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *tagService) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return s.repo.Update(ctx, tag)
}

func (s *tagService) DeleteTag(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
