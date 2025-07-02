package repository

import (
	"context"
	"errors"
	"fmt" // For error wrapping

	"echobackend/internal/model"

	"gorm.io/gorm" // For gorm.DB and gorm.ErrRecordNotFound
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	FindAll(ctx context.Context) ([]model.Tag, error)
	FindByID(ctx context.Context, id uint) (*model.Tag, error)
	FindByName(ctx context.Context, name string) (*model.Tag, error)
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
		return errors.New("tag cannot be nil")
	}
	// GORM's Create will also update CreatedAt/UpdatedAt if they exist in model.Tag
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
		return nil, errors.New("invalid tag ID")
	}
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to find tag by ID %d: %w", id, err)
	}
	return &tag, nil
}

func (r *tagRepository) FindByName(ctx context.Context, name string) (*model.Tag, error) {
	if name == "" {
		return nil, errors.New("tag name cannot be empty")
	}
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, fmt.Errorf("failed to find tag by name %s: %w", name, err)
	}
	return &tag, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.ID == 0 { // Check for zero value of primary key
		return errors.New("invalid tag ID for update")
	}

	// GORM's Save updates all fields or creates if record not found (based on PK).
	// If you only want to update, ensure the record exists or use Updates for specific fields.
	// For updating based on a full model, Save is common.
	// It will also update UpdatedAt field if present in model.Tag.
	result := r.db.WithContext(ctx).Save(tag)
	if result.Error != nil {
		return fmt.Errorf("failed to update tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		// This case might mean the record didn't exist, or no fields were changed
		// depending on GORM version and specific behavior of Save.
		// For Save, if PK exists, it updates. If PK doesn't exist, it inserts.
		// If we want to ensure it's an update only, a prior check or using Updates might be better.
		// However, if Save is used and RowsAffected is 0 after an update attempt on existing PK,
		// it could mean no actual data changed.
		// Let's assume for now that if Save doesn't error, it's fine.
		// A more robust check for "not found" would be to query first or use `Updates` and check RowsAffected.
		// For simplicity, if Save doesn't error, we assume success.
		// If an update must happen, check RowsAffected.
		return ErrTagNotFound // If we strictly expect an update
	}
	return nil
}

func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid tag ID")
	}
	// This will perform a soft delete if model.Tag has gorm.DeletedAt.
	// Otherwise, it's a hard delete.
	result := r.db.WithContext(ctx).Delete(&model.Tag{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrTagNotFound
	}
	return nil
}
