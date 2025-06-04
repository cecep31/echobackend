package repository

import (
	"context"
	// "database/sql" // No longer needed with GORM's error handling
	"errors"
	"fmt"
	"time"

	"echobackend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm" // Changed from bun to gorm
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
	db *gorm.DB // Changed from bun.DB to gorm.DB
}

// NewWorkspaceRepository creates a new workspace repository instance
func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository { // Changed parameter type
	return &workspaceRepository{db: db}
}

// withTimeout adds a timeout to the context if one doesn't exist
func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	// Check if context already has a deadline
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	// Add a reasonable timeout for database operations
	return context.WithTimeout(ctx, 5*time.Second)
}

// Create adds a new workspace to the database
func (r *workspaceRepository) Create(ctx context.Context, workspace *model.Workspace) error {
	timedCtx, cancel := withTimeout(ctx)
	defer cancel()

	exists, err := r.Exists(timedCtx, workspace.Name, workspace.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to check workspace existence: %w", err)
	}
	if exists {
		return ErrWorkspaceExists
	}

	return r.db.WithContext(timedCtx).Transaction(func(tx *gorm.DB) error {
		// Create the workspace
		if err := tx.Create(workspace).Error; err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		// Add creator as admin member if members are not provided
		// This assumes workspace.Members is not populated by the Create DTO,
		// and workspace.ID is populated by GORM after tx.Create(workspace).
		if len(workspace.Members) == 0 {
			member := model.WorkspaceMember{
				WorkspaceID: workspace.ID, // workspace.ID should be populated by GORM after Create
				UserID:      workspace.CreatedBy,
				Role:        "admin",
			}
			if err := tx.Create(&member).Error; err != nil {
				return fmt.Errorf("failed to add creator as member: %w", err)
			}
		}
		return nil // Commit
	})
}

// GetByID retrieves a workspace by its ID
func (r *workspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var workspace model.Workspace
	// Assuming model.Workspace has a "Members" relation defined for GORM
	err := r.db.WithContext(ctx).Preload("Members").First(&workspace, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("failed to get workspace by ID %s: %w", id.String(), err)
	}
	return &workspace, nil
}

// GetAll retrieves all workspaces with pagination
func (r *workspaceRepository) GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var workspaces []*model.Workspace

	// Count total records
	var total int64
	err := r.db.WithContext(ctx).Model((*model.Workspace)(nil)).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count workspaces: %w", err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Preload("Members"). // Assuming Members relation is defined
		Offset(offset).
		Limit(limit).
		Find(&workspaces).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get workspaces: %w", err)
	}

	return workspaces, total, nil
}

// GetByUserID retrieves all workspaces a user is a member of
func (r *workspaceRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Workspace, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var workspaces []*model.Workspace
	// Assuming model.Workspace has a "Members" relation and WorkspaceMember has UserID and WorkspaceID
	// The join condition implies a direct link or through WorkspaceMember table.
	// GORM's Preload("Members") should handle fetching members.
	// To filter workspaces by a user being a member:
	err := r.db.WithContext(ctx).
		Joins("JOIN workspace_members ON workspaces.id = workspace_members.workspace_id").
		Where("workspace_members.user_id = ?", userID).
		Preload("Members"). // Preload members for each workspace found
		Find(&workspaces).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces by user ID %s: %w", userID, err)
	}

	return workspaces, nil
}

// Update updates an existing workspace
func (r *workspaceRepository) Update(ctx context.Context, workspace *model.Workspace) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	if workspace.ID == uuid.Nil {
		return errors.New("invalid workspace ID for update")
	}

	// Using Updates with a map to update specific fields.
	// Alternatively, if workspace struct contains only the fields to update, db.Save(workspace) can be used.
	// For selective updates based on DTO or partial model, Updates is better.
	updates := map[string]interface{}{
		"name":        workspace.Name,
		"description": workspace.Description,
		"icon":        workspace.Icon,
		// Add other updatable fields from model.Workspace here
	}

	result := r.db.WithContext(ctx).Model(&model.Workspace{}).Where("id = ?", workspace.ID).Updates(updates)
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
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	// Assumes model.Workspace has gorm.DeletedAt for soft delete
	result := r.db.WithContext(ctx).Delete(&model.Workspace{}, "id = ?", id)
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
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete workspace members first
		// Unscoped() ensures hard delete even if WorkspaceMember has gorm.DeletedAt
		if err := tx.Unscoped().Where("workspace_id = ?", id).Delete(&model.WorkspaceMember{}).Error; err != nil {
			return fmt.Errorf("failed to delete workspace members: %w", err)
		}

		// Delete workspace
		// Unscoped() ensures hard delete even if Workspace has gorm.DeletedAt
		result := tx.Unscoped().Where("id = ?", id).Delete(&model.Workspace{})
		if result.Error != nil {
			return fmt.Errorf("failed to hard delete workspace: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrWorkspaceNotFound // Workspace itself not found to delete
		}
		return nil // Commit
	})
}

// Exists checks if a workspace with the given name already exists for the user
func (r *workspaceRepository) Exists(ctx context.Context, name string, createdBy string) (bool, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var count int64
	err := r.db.WithContext(ctx).Model(&model.Workspace{}).
		Where("name = ? AND created_by = ?", name, createdBy).
		Count(&count).Error
	if err != nil {
		// If error is record not found, it means no such workspace, so count is 0.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
	}
	return count > 0, nil
}

// AddMember adds a new member to a workspace
func (r *workspaceRepository) AddMember(ctx context.Context, member *model.WorkspaceMember) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	// Check if workspace exists
	exists, err := r.workspaceExists(ctx, member.WorkspaceID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWorkspaceNotFound
	}

	// Check if member already exists
	var existingMember model.WorkspaceMember
	err = r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", member.WorkspaceID, member.UserID).
		First(&existingMember).Error

	if err == nil { // Member found
		// Update role instead of creating a new member, if role is different
		if existingMember.Role != member.Role {
			return r.UpdateMemberRole(ctx, member.WorkspaceID, member.UserID, member.Role)
		}
		return nil // Member exists with the same role, do nothing
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) { // An actual error occurred
		return fmt.Errorf("failed to check member existence: %w", err)
	}
	// Member not found (gorm.ErrRecordNotFound), proceed to create

	// Create new member
	if err := r.db.WithContext(ctx).Create(member).Error; err != nil {
		return fmt.Errorf("failed to add workspace member: %w", err)
	}
	return nil
}

// GetMembers retrieves all members of a workspace
func (r *workspaceRepository) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*model.WorkspaceMember, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var members []*model.WorkspaceMember
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace members for workspace ID %s: %w", workspaceID.String(), err)
	}
	return members, nil
}

// UpdateMemberRole updates a member's role in a workspace
func (r *workspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID uuid.UUID, userID string, role string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Update("role", role) // Update single column

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
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	// Assumes model.WorkspaceMember might have gorm.DeletedAt for soft delete,
	// otherwise it's a hard delete. If soft delete, use .Delete()
	// If hard delete is always intended for members, use .Unscoped().Delete()
	result := r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
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
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var member model.WorkspaceMember
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		First(&member).Error // First will return error if not found

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil // Not a member, not an error
		}
		return false, "", fmt.Errorf("failed to check membership for user %s in workspace %s: %w", userID, workspaceID.String(), err)
	}
	return true, member.Role, nil
}

// workspaceExists is a helper function to check if a workspace exists
func (r *workspaceRepository) workspaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Workspace{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check workspace existence for ID %s: %w", id.String(), err)
	}
	return count > 0, nil
}
