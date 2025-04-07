package service

import (
	"context"
	"errors"

	"echobackend/internal/model"
	"echobackend/internal/repository"

	"github.com/google/uuid"
)

type PageService struct {
	pagesRepo repository.PageRepository
}

func NewPageService(pagesRepo repository.PageRepository) PageService {
	return PageService{pagesRepo: pagesRepo}
}

// CreatePage creates a new page in the workspace
func (s *PageService) CreatePage(ctx context.Context, page *model.Page) error {
	if page.Title == "" {
		return errors.New("page title is required")
	}

	if page.WorkspaceID == uuid.Nil {
		return errors.New("workspace ID is required")
	}

	if page.CreatedBy == "" {
		return errors.New("creator information is required")
	}

	return s.pagesRepo.CreatePage(ctx, page)
}

// GetPageByID retrieves a page by its ID
func (s *PageService) GetPageByID(ctx context.Context, id uuid.UUID) (*model.Page, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid page ID")
	}

	return s.pagesRepo.GetPageByID(ctx, id)
}

// GetPagesByWorkspaceID retrieves all pages in a workspace
func (s *PageService) GetPagesByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]model.Page, error) {
	if workspaceID == uuid.Nil {
		return nil, errors.New("invalid workspace ID")
	}

	return s.pagesRepo.GetPagesByWorkspaceID(ctx, workspaceID)
}

// GetChildPages retrieves all child pages of a given page
func (s *PageService) GetChildPages(ctx context.Context, parentID uuid.UUID) ([]model.Page, error) {
	if parentID == uuid.Nil {
		return nil, errors.New("invalid parent page ID")
	}

	return s.pagesRepo.GetChildPages(ctx, parentID)
}

// UpdatePage updates an existing page
func (s *PageService) UpdatePage(ctx context.Context, page *model.Page) error {
	if page.ID == uuid.Nil {
		return errors.New("invalid page ID")
	}

	if page.Title == "" {
		return errors.New("page title is required")
	}

	existing, err := s.pagesRepo.GetPageByID(ctx, page.ID)
	if err != nil {
		return err
	}

	// Preserve certain fields from the existing page
	page.CreatedAt = existing.CreatedAt
	page.CreatedBy = existing.CreatedBy

	return s.pagesRepo.UpdatePage(ctx, page)
}

// DeletePage deletes a page by its ID
func (s *PageService) DeletePage(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid page ID")
	}

	return s.pagesRepo.DeletePage(ctx, id)
}
