package service

import (
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
)

type tagService struct {
	repo repository.TagRepository
}

type TagService interface {
	CreateTag(tag *model.Tag) error
	GetTags() ([]model.Tag, error)
	GetTagByID(id uint) (*model.Tag, error)
	UpdateTag(tag *model.Tag) error
	DeleteTag(id uint) error
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) CreateTag(tag *model.Tag) error {
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	return s.repo.Create(tag)
}

func (s *tagService) GetTags() ([]model.Tag, error) {
	return s.repo.FindAll()
}

func (s *tagService) GetTagByID(id uint) (*model.Tag, error) {
	return s.repo.FindByID(id)
}

func (s *tagService) UpdateTag(tag *model.Tag) error {
	return s.repo.Update(tag)
}

func (s *tagService) DeleteTag(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
