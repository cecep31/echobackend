package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

// WorkspaceService defines the interface for workspace business logic
type WorkspaceService interface {
	// Workspace operations
	Create(ctx context.Context, workspace *model.Workspace) error
	GetByID(ctx context.Context, id string) (*model.Workspace, error)
	GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Workspace, error)
	Update(ctx context.Context, workspace *model.Workspace) error
	Delete(ctx context.Context, id string) error

	// Workspace member operations
	AddMember(ctx context.Context, workspaceID string, userID string, role string) error
	GetMembers(ctx context.Context, workspaceID string) ([]*model.WorkspaceMember, error)
	UpdateMemberRole(ctx context.Context, workspaceID string, userID string, role string) error
	RemoveMember(ctx context.Context, workspaceID string, userID string) error
	CheckMembership(ctx context.Context, workspaceID string, userID string) (bool, string, error)
}

type workspaceService struct {
	workspaceRepo repository.WorkspaceRepository
}

// NewWorkspaceService creates a new workspace service instance
func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{workspaceRepo: workspaceRepo}
}

// Create adds a new workspace
func (s *workspaceService) Create(ctx context.Context, workspace *model.Workspace) error {
	return s.workspaceRepo.Create(ctx, workspace)
}

// GetByID retrieves a workspace by its ID
func (s *workspaceService) GetByID(ctx context.Context, id string) (*model.Workspace, error) {
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.GetByID(ctx, workspaceID)
}

// GetAll retrieves all workspaces with pagination
func (s *workspaceService) GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error) {
	return s.workspaceRepo.GetAll(ctx, offset, limit)
}

// GetByUserID retrieves all workspaces a user is a member of
func (s *workspaceService) GetByUserID(ctx context.Context, userID string) ([]*model.Workspace, error) {
	return s.workspaceRepo.GetByUserID(ctx, userID)
}

// Update updates an existing workspace
func (s *workspaceService) Update(ctx context.Context, workspace *model.Workspace) error {
	return s.workspaceRepo.Update(ctx, workspace)
}

// Delete soft deletes a workspace
func (s *workspaceService) Delete(ctx context.Context, id string) error {
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.SoftDelete(ctx, workspaceID)
}

// AddMember adds a new member to a workspace
func (s *workspaceService) AddMember(ctx context.Context, workspaceID string, userID string, role string) error {
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return fmt.Errorf("invalid workspace ID: %w", err)
	}

	member := &model.WorkspaceMember{
		WorkspaceID: wsID,
		UserID:      userID,
		Role:        role,
	}

	return s.workspaceRepo.AddMember(ctx, member)
}

// GetMembers retrieves all members of a workspace
func (s *workspaceService) GetMembers(ctx context.Context, workspaceID string) ([]*model.WorkspaceMember, error) {
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.GetMembers(ctx, wsID)
}

// UpdateMemberRole updates a member's role in a workspace
func (s *workspaceService) UpdateMemberRole(ctx context.Context, workspaceID string, userID string, role string) error {
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.UpdateMemberRole(ctx, wsID, userID, role)
}

// RemoveMember removes a member from a workspace
func (s *workspaceService) RemoveMember(ctx context.Context, workspaceID string, userID string) error {
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.RemoveMember(ctx, wsID, userID)
}

// CheckMembership checks if a user is a member of a workspace and returns their role
func (s *workspaceService) CheckMembership(ctx context.Context, workspaceID string, userID string) (bool, string, error) {
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return false, "", fmt.Errorf("invalid workspace ID: %w", err)
	}

	return s.workspaceRepo.IsMember(ctx, wsID, userID)
}
