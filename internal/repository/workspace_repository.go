package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"echobackend/internal/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
	db *bun.DB
}

// NewWorkspaceRepository creates a new workspace repository instance
func NewWorkspaceRepository(db *bun.DB) WorkspaceRepository {
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
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	exists, err := r.Exists(ctx, workspace.Name, workspace.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to check workspace existence: %w", err)
	}
	if exists {
		return ErrWorkspaceExists
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Use defer with named return value to ensure proper rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the workspace
	_, err = tx.NewInsert().
		Model(workspace).
		Exec(ctx)
	if err != nil {
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
		_, err = tx.NewInsert().
			Model(&member).
			Exec(ctx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add creator as member: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID retrieves a workspace by its ID
func (r *workspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var workspace model.Workspace
	err := r.db.NewSelect().
		Model(&workspace).
		Relation("Members").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}
	return &workspace, nil
}

// GetAll retrieves all workspaces with pagination
func (r *workspaceRepository) GetAll(ctx context.Context, offset int, limit int) ([]*model.Workspace, int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var workspaces []*model.Workspace

	// Count total records
	totalCount, err := r.db.NewSelect().
		Model((*model.Workspace)(nil)).
		Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count workspaces: %w", err)
	}
	total := int64(totalCount)

	// Get paginated records
	err = r.db.NewSelect().
		Model(&workspaces).
		Relation("Members").
		Offset(offset).
		Limit(limit).
		Scan(ctx)
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

	err := r.db.NewSelect().
		Model(&workspaces).
		Join("JOIN workspace_members ON workspaces.id = workspace_members.workspace_id").
		Where("workspace_members.user_id = ?", userID).
		Relation("Members").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces by user ID: %w", err)
	}

	return workspaces, nil
}

// Update updates an existing workspace
func (r *workspaceRepository) Update(ctx context.Context, workspace *model.Workspace) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	res, err := r.db.NewUpdate().
		Model(&model.Workspace{}).
		Set("name = ?", workspace.Name).
		Set("description = ?", workspace.Description).
		Set("icon = ?", workspace.Icon).
		Where("id = ?", workspace.ID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}
	return nil
}

// SoftDelete soft deletes a workspace
func (r *workspaceRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	res, err := r.db.NewDelete().
		Model(&model.Workspace{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to soft delete workspace: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}
	return nil
}

// HardDelete permanently deletes a workspace
func (r *workspaceRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Use defer with named return value to ensure proper rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete workspace members first
	_, err = tx.NewDelete().
		Model(&model.WorkspaceMember{}).
		Where("workspace_id = ?", id).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete workspace members: %w", err)
	}

	// Delete workspace
	res, err := tx.NewDelete().
		Model(&model.Workspace{}).
		Where("id = ?", id).
		ForceDelete().
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to hard delete workspace: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return ErrWorkspaceNotFound
	}

	return tx.Commit()
}

// Exists checks if a workspace with the given name already exists for the user
func (r *workspaceRepository) Exists(ctx context.Context, name string, createdBy string) (bool, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	count, err := r.db.NewSelect().
		Model(&model.Workspace{}).
		Where("name = ? AND created_by = ?", name, createdBy).
		Count(ctx)
	if err != nil {
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
	count, err := r.db.NewSelect().
		Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", member.WorkspaceID, member.UserID).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check member existence: %w", err)
	}
	if count > 0 {
		// Update role instead of creating a new member
		return r.UpdateMemberRole(ctx, member.WorkspaceID, member.UserID, member.Role)
	}

	// Create new member
	_, err = r.db.NewInsert().
		Model(member).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add workspace member: %w", err)
	}
	return nil
}

// GetMembers retrieves all members of a workspace
func (r *workspaceRepository) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*model.WorkspaceMember, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var members []*model.WorkspaceMember
	err := r.db.NewSelect().
		Model(&members).
		Where("workspace_id = ?", workspaceID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace members: %w", err)
	}
	return members, nil
}

// UpdateMemberRole updates a member's role in a workspace
func (r *workspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID uuid.UUID, userID string, role string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	res, err := r.db.NewUpdate().
		Model(&model.WorkspaceMember{}).
		Set("role = ?", role).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// RemoveMember removes a member from a workspace
func (r *workspaceRepository) RemoveMember(ctx context.Context, workspaceID uuid.UUID, userID string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	res, err := r.db.NewDelete().
		Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove workspace member: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// IsMember checks if a user is a member of a workspace and returns their role
func (r *workspaceRepository) IsMember(ctx context.Context, workspaceID uuid.UUID, userID string) (bool, string, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var member model.WorkspaceMember
	err := r.db.NewSelect().
		Model(&member).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check membership: %w", err)
	}
	return true, member.Role, nil
}

// workspaceExists is a helper function to check if a workspace exists
func (r *workspaceRepository) workspaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.db.NewSelect().
		Model(&model.Workspace{}).
		Where("id = ?", id).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check workspace existence: %w", err)
	}
	return count > 0, nil
}
