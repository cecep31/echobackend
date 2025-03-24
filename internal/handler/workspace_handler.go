package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WorkspaceHandler struct {
	workspaceService service.WorkspaceService
}

func NewWorkspaceHandler(workspaceService service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: workspaceService}
}

// CreateWorkspace handles the creation of a new workspace
func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	type CreateWorkspaceRequest struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	var req CreateWorkspaceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid request format",
			"success": false,
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Validation failed",
			"success": false,
		})
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.Get("user_id").(string)

	workspace := &model.Workspace{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.workspaceService.Create(c.Request().Context(), workspace); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to create workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data":    workspace,
		"message": "Workspace created successfully",
		"success": true,
	})
}

// GetWorkspaceByID retrieves a workspace by its ID
func (h *WorkspaceHandler) GetWorkspaceByID(c echo.Context) error {
	workspaceID := c.Param("id")

	workspace, err := h.workspaceService.GetByID(c.Request().Context(), workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to retrieve workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    workspace,
		"success": true,
	})
}

// GetAllWorkspaces retrieves all workspaces with pagination
func (h *WorkspaceHandler) GetAllWorkspaces(c echo.Context) error {
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	workspaces, total, err := h.workspaceService.GetAll(c.Request().Context(), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to retrieve workspaces",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    workspaces,
		"success": true,
		"metadata": map[string]any{
			"totalItems": total,
		},
	})
}

// GetUserWorkspaces retrieves all workspaces a user is a member of
func (h *WorkspaceHandler) GetUserWorkspaces(c echo.Context) error {
	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.Get("user_id").(string)

	workspaces, err := h.workspaceService.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to retrieve user workspaces",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    workspaces,
		"success": true,
	})
}

// UpdateWorkspace updates an existing workspace
func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	workspaceID := c.Param("id")

	type UpdateWorkspaceRequest struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	var req UpdateWorkspaceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Invalid request format",
			"success": false,
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Validation failed",
			"success": false,
		})
	}

	// Get the existing workspace
	wsID, err := uuid.Parse(workspaceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Invalid workspace ID",
			"success": false,
		})
	}

	workspace := &model.Workspace{
		ID:          wsID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		UpdatedAt:   time.Now(),
	}

	if err := h.workspaceService.Update(c.Request().Context(), workspace); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to update workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Workspace updated successfully",
		"success": true,
	})
}

// DeleteWorkspace soft deletes a workspace
func (h *WorkspaceHandler) DeleteWorkspace(c echo.Context) error {
	workspaceID := c.Param("id")

	if err := h.workspaceService.Delete(c.Request().Context(), workspaceID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to delete workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Workspace deleted successfully",
		"success": true,
	})
}

// AddMember adds a new member to a workspace
func (h *WorkspaceHandler) AddMember(c echo.Context) error {
	workspaceID := c.Param("id")

	type AddMemberRequest struct {
		UserID string `json:"user_id" validate:"required"`
		Role   string `json:"role" validate:"required,oneof=admin editor viewer"`
	}

	var req AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Invalid request format",
			"success": false,
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Validation failed",
			"success": false,
		})
	}

	if err := h.workspaceService.AddMember(c.Request().Context(), workspaceID, req.UserID, req.Role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to add member to workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Member added to workspace successfully",
		"success": true,
	})
}

// GetMembers retrieves all members of a workspace
func (h *WorkspaceHandler) GetMembers(c echo.Context) error {
	workspaceID := c.Param("id")

	members, err := h.workspaceService.GetMembers(c.Request().Context(), workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to retrieve workspace members",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    members,
		"success": true,
	})
}

// UpdateMemberRole updates a member's role in a workspace
func (h *WorkspaceHandler) UpdateMemberRole(c echo.Context) error {
	workspaceID := c.Param("id")
	userID := c.Param("user_id")

	type UpdateRoleRequest struct {
		Role string `json:"role" validate:"required,oneof=admin editor viewer"`
	}

	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Invalid request format",
			"success": false,
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Validation failed",
			"success": false,
		})
	}

	if err := h.workspaceService.UpdateMemberRole(c.Request().Context(), workspaceID, userID, req.Role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to update member role",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Member role updated successfully",
		"success": true,
	})
}

// RemoveMember removes a member from a workspace
func (h *WorkspaceHandler) RemoveMember(c echo.Context) error {
	workspaceID := c.Param("id")
	userID := c.Param("user_id")

	if err := h.workspaceService.RemoveMember(c.Request().Context(), workspaceID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to remove member from workspace",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Member removed from workspace successfully",
		"success": true,
	})
}
