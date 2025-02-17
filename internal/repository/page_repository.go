package repository

import (
	"echobackend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{db: db}
}

// CreatePage creates a new page in the database
func (r *PageRepository) CreatePage(page *model.Page) error {
	return r.db.Create(page).Error
}

// GetPageByID retrieves a page by its ID
func (r *PageRepository) GetPageByID(id uuid.UUID) (*model.Page, error) {
	var page model.Page
	err := r.db.Preload("Blocks").First(&page, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetPagesByWorkspaceID retrieves all pages in a workspace
func (r *PageRepository) GetPagesByWorkspaceID(workspaceID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.Where("workspace_id = ?", workspaceID).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetChildPages retrieves all child pages of a given page
func (r *PageRepository) GetChildPages(parentID uuid.UUID) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.Where("parent_id = ?", parentID).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// UpdatePage updates an existing page
func (r *PageRepository) UpdatePage(page *model.Page) error {
	return r.db.Save(page).Error
}

// DeletePage soft deletes a page
func (r *PageRepository) DeletePage(id uuid.UUID) error {
	return r.db.Delete(&model.Page{}, "id = ?", id).Error
}

// HardDeletePage permanently deletes a page
func (r *PageRepository) HardDeletePage(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&model.Page{}, "id = ?", id).Error
}
