package repository

import (
	"context"
	"echobackend/internal/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PageRepository interface {
	CreatePage(ctx context.Context, page *model.Page) error
	GetPageByID(ctx context.Context, id uuid.UUID) (*model.Page, error)
	GetPagesByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]model.Page, error)
	GetChildPages(ctx context.Context, parentID uuid.UUID) ([]model.Page, error)
	UpdatePage(ctx context.Context, page *model.Page) error
	DeletePage(ctx context.Context, id uuid.UUID) error
	HardDeletePage(ctx context.Context, id uuid.UUID) error
}

type pageRepository struct {
	db *bun.DB
}

func NewPageRepository(db *bun.DB) PageRepository {
	return &pageRepository{db: db}
}

// CreatePage creates a new page in the database
func (r *pageRepository) CreatePage(ctx context.Context, page *model.Page) error {
	_, err := r.db.NewInsert().
		Model(page).
		Exec(ctx)
	return err
}

// GetPageByID retrieves a page by its ID
func (r *pageRepository) GetPageByID(ctx context.Context, id uuid.UUID) (*model.Page, error) {
	var page model.Page
	err := r.db.NewSelect().
		Model(&page).
		Relation("Blocks").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetPagesByWorkspaceID retrieves all pages in a workspace
func (r *pageRepository) GetPagesByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.NewSelect().
		Model(&pages).
		Where("workspace_id = ?", workspaceID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetChildPages retrieves all child pages of a given page
func (r *pageRepository) GetChildPages(ctx context.Context, parentID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.NewSelect().
		Model(&pages).
		Where("parent_id = ?", parentID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// UpdatePage updates an existing page
func (r *pageRepository) UpdatePage(ctx context.Context, page *model.Page) error {
	_, err := r.db.NewUpdate().
		Model(page).
		WherePK().
		Exec(ctx)
	return err
}

// DeletePage soft deletes a page
func (r *pageRepository) DeletePage(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model(&model.Page{}).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// HardDeletePage permanently deletes a page
func (r *pageRepository) HardDeletePage(ctx context.Context, id uuid.UUID) error {
	// In Bun, we need to use ForceDelete() for hard delete
	_, err := r.db.NewDelete().
		Model(&model.Page{}).
		Where("id = ?", id).
		ForceDelete().
		Exec(ctx)
	return err
}
