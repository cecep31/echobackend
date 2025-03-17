package repository

import (
	"context"
	"errors"
	"fmt"

	"echobackend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Errors that can be returned by the repository
var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrWorkspaceExists   = errors.New("workspace already exists with this name")
	ErrMemberNotFound    = errors.New("workspace member not found")
)

// WorkspaceRepository defines the interface for workspace data operations
type WorkspaceRepository interface {
	// Workspace operations
	Create(ctx context.Context, workspace *model.Workspace) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
	GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Workspace, error)
	Update(ctx context.Context, workspace *model.Workspace) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, name string, createdBy string) (bool, error)
	
	// Workspace member operations
	AddMember(ctx context.Context, member *model.WorkspaceMember) error
	GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*model.WorkspaceMember, error)
	UpdateMemberRole(ctx context.Context, workspaceID uuid.UUID, userID string, role string) error
	RemoveMember(ctx context.Context, workspaceID uuid.UUID, userID string) error
	IsMember(ctx context.Context, workspaceID uuid.UUID, userID string) (bool, string, error)
}

type workspaceRepository struct {
	db *gorm.DB
}

// NewWorkspaceRepository creates a new workspace repository instance
func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

// Create adds a new workspace to the database
func (r *workspaceRepository) Create(ctx context.Context, workspace *model.Workspace) error {
	exists, err := r.Exists(ctx, workspace.Name, workspace.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to check workspace existence: %w", err)
	}
	if exists {
		return ErrWorkspaceExists
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create the workspace
	if err := tx.Create(workspace).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Add creator as admin member if members are not provided
	if len(workspace.Members) == 0 {
		member := model.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      workspace.CreatedBy,
			Role:        "admin",
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add creator as member: %w", err)
		}
	}

	return tx.Commit().Error
}

// GetByID retrieves a workspace by its ID
func (r *workspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	var workspace model.Workspace
	if err := r.db.WithContext(ctx).Preload("Members").Where("id = ?", id).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}
	return &workspace, nil
}

// GetAll retrieves all workspaces with pagination
func (r *workspaceRepository) GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error) {
	var workspaces []*model.Workspace
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&model.Workspace{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count workspaces: %w", err)
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).Preload("Members").Offset(offset).Limit(limit).Find(&workspaces).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get workspaces: %w", err)
	}

	return workspaces, total, nil
}

// GetByUserID retrieves all workspaces a user is a member of
func (r *workspaceRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Workspace, error) {
	var workspaces []*model.Workspace

	err := r.db.WithContext(ctx).
		Joins("JOIN workspace_members ON workspaces.id = workspace_members.workspace_id").
		Where("workspace_members.user_id = ?", userID).
		Preload("Members").
		Find(&workspaces).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces by user ID: %w", err)
	}

	return workspaces, nil
}

// Update updates an existing workspace
func (r *workspaceRepository) Update(ctx context.Context, workspace *model.Workspace) error {
	result := r.db.WithContext(ctx).Model(&model.Workspace{}).
		Where("id = ?", workspace.ID).
		Updates(map[string]interface{}{
			"name":        workspace.Name,
			"description": workspace.Description,
			"icon":        workspace.Icon,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update workspace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrWorkspaceNotFound
	}
	return nil
}

// SoftDelete soft deletes a workspace
func (r *workspaceRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Workspace{})
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete workspace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrWorkspaceNotFound
	}
	return nil
}

// HardDelete permanently deletes a workspace
func (r *workspaceRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Delete workspace members first
	if err := tx.Where("workspace_id = ?", id).Delete(&model.WorkspaceMember{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete workspace members: %w", err)
	}

	// Delete workspace
	result := tx.Unscoped().Where("id = ?", id).Delete(&model.Workspace{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to hard delete workspace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return ErrWorkspaceNotFound
	}

	return tx.Commit().Error
}

// Exists checks if a workspace with the given name already exists for the user
func (r *workspaceRepository) Exists(ctx context.Context, name string, createdBy string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Workspace{}).
		Where("name = ? AND created_by = ?", name, createdBy).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
	}
	return count > 0, nil
}

// AddMember adds a new member to a workspace
func (r *workspaceRepository) AddMember(ctx context.Context, member *model.WorkspaceMember) error {
	// Check if workspace exists
	exists, err := r.workspaceExists(ctx, member.WorkspaceID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWorkspaceNotFound
	}

	// Check if member already exists
	var count int64
	err = r.db.WithContext(ctx).Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", member.WorkspaceID, member.UserID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check member existence: %w", err)
	}

	// If member exists, update role
	if count > 0 {
		return r.UpdateMemberRole(ctx, member.WorkspaceID, member.UserID, member.Role)
	}

	// Otherwise add new member
	return r.db.WithContext(ctx).Create(member).Error
}

// GetMembers retrieves all members of a workspace
func (r *workspaceRepository) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*model.WorkspaceMember, error) {
	// Check if workspace exists
	exists, err := r.workspaceExists(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrWorkspaceNotFound
	}

	var members []*model.WorkspaceMember
	if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&members).Error; err != nil {
		return nil, fmt.Errorf("failed to get workspace members: %w", err)
	}
	return members, nil
}

// UpdateMemberRole updates a member's role in a workspace
func (r *workspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID uuid.UUID, userID string, role string) error {
	result := r.db.WithContext(ctx).Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Update("role", role)

	if result.Error != nil {
		return fmt.Errorf("failed to update member role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// RemoveMember removes a member from a workspace
func (r *workspaceRepository) RemoveMember(ctx context.Context, workspaceID uuid.UUID, userID string) error {
	result := r.db.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Delete(&model.WorkspaceMember{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove workspace member: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// IsMember checks if a user is a member of a workspace and returns their role
func (r *workspaceRepository) IsMember(ctx context.Context, workspaceID uuid.UUID, userID string) (bool, string, error) {
	var member model.WorkspaceMember
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check membership: %w", err)
	}

	return true, member.Role, nil
}

// workspaceExists is a helper function to check if a workspace exists
func (r *workspaceRepository) workspaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Workspace{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
	}
	return count > 0, nil
}
