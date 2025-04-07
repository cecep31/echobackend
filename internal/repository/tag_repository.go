package repository

import (
	"context"
	"database/sql"
	"errors"

	"echobackend/internal/model"

	"github.com/uptrace/bun"
)

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	FindAll(ctx context.Context) ([]model.Tag, error)
	FindByID(ctx context.Context, id uint) (*model.Tag, error)
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id uint) error
}

type tagRepository struct {
	db *bun.DB
}

func NewTagRepository(db *bun.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	_, err := r.db.NewInsert().
		Model(tag).
		Exec(ctx)
	return err
}

func (r *tagRepository) FindAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.NewSelect().
		Model(&tags).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uint) (*model.Tag, error) {
	if id == 0 {
		return nil, errors.New("invalid tag ID")
	}

	var tag model.Tag
	err := r.db.NewSelect().
		Model(&tag).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.ID == 0 {
		return errors.New("invalid tag ID")
	}

	res, err := r.db.NewUpdate().
		Model(tag).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("tag not found")
	}
	return nil
}

func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid tag ID")
	}

	res, err := r.db.NewDelete().
		Model(&model.Tag{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("tag not found")
	}
	return nil
}
