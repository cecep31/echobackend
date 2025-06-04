package repository

import (
	"context"
	"echobackend/internal/model"
	"errors" // For custom errors and gorm.ErrRecordNotFound
	"fmt"    // For error wrapping

	"github.com/google/uuid"
	"gorm.io/gorm" // For gorm.DB
)

var (
	ErrPageNotFound = errors.New("page not found")
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
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

// CreatePage creates a new page in the database
func (r *pageRepository) CreatePage(ctx context.Context, page *model.Page) error {
	result := r.db.WithContext(ctx).Create(page)
	if result.Error != nil {
		return fmt.Errorf("failed to create page: %w", result.Error)
	}
	return nil
}

// GetPageByID retrieves a page by its ID
func (r *pageRepository) GetPageByID(ctx context.Context, id uuid.UUID) (*model.Page, error) {
	var page model.Page
	// Assuming model.Page has a "Blocks" relation defined for GORM
	err := r.db.WithContext(ctx).Preload("Blocks").First(&page, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPageNotFound
		}
		return nil, fmt.Errorf("failed to get page by ID %s: %w", id.String(), err)
	}
	return &page, nil
}

// GetPagesByWorkspaceID retrieves all pages in a workspace
func (r *pageRepository) GetPagesByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&pages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pages by workspace ID %s: %w", workspaceID.String(), err)
	}
	return pages, nil
}

// GetChildPages retrieves all child pages of a given page
func (r *pageRepository) GetChildPages(ctx context.Context, parentID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&pages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get child pages for parent ID %s: %w", parentID.String(), err)
	}
	return pages, nil
}

// UpdatePage updates an existing page
func (r *pageRepository) UpdatePage(ctx context.Context, page *model.Page) error {
	if page.ID == uuid.Nil {
		return errors.New("invalid page ID for update")
	}
	result := r.db.WithContext(ctx).Save(page)
	if result.Error != nil {
		return fmt.Errorf("failed to update page: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrPageNotFound // Or handle as no change / record not found
	}
	return nil
}

// DeletePage soft deletes a page
func (r *pageRepository) DeletePage(ctx context.Context, id uuid.UUID) error {
	// Assumes model.Page has gorm.DeletedAt for soft delete
	result := r.db.WithContext(ctx).Delete(&model.Page{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete page: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrPageNotFound
	}
	return nil
}

// HardDeletePage permanently deletes a page
func (r *pageRepository) HardDeletePage(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.Page{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to hard delete page: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrPageNotFound
	}
	return nil
}
