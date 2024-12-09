package repository

import (
	"echobackend/internal/model"
	"errors"

	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *model.Tag) error
	FindAll() ([]model.Tag, error)
	FindByID(id uint) (*model.Tag, error)
	Update(tag *model.Tag) error
	Delete(id uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *model.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	return r.db.Create(tag).Error
}

func (r *tagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	result := r.db.Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (r *tagRepository) FindByID(id uint) (*model.Tag, error) {
	if id == 0 {
		return nil, errors.New("invalid tag ID")
	}

	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(tag *model.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.ID == 0 {
		return errors.New("invalid tag ID")
	}

	result := r.db.Save(tag)
	if result.RowsAffected == 0 {
		return errors.New("tag not found")
	}
	return result.Error
}

func (r *tagRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid tag ID")
	}

	result := r.db.Delete(&model.Tag{}, id)
	if result.RowsAffected == 0 {
		return errors.New("tag not found")
	}
	return result.Error
}
