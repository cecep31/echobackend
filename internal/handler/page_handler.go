package handler

import (
	"net/http"

	"echobackend/internal/model"
	"echobackend/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PageHandler struct {
	pageService service.PageService
}

func NewPageHandler(pageService service.PageService) *PageHandler {
	return &PageHandler{pageService: pageService}
}

// CreatePage handles the creation of a new page
func (h *PageHandler) CreatePage(c echo.Context) error {
	var page model.Page
	if err := c.Bind(&page); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid request format",
			"success": false,
		})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)
	page.CreatedBy = userID

	if err := h.pageService.CreatePage(c.Request().Context(), &page); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to create page",
			"success": false,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data":    page,
		"message": "Page created successfully",
		"success": true,
	})
}

// GetPage retrieves a page by ID
func (h *PageHandler) GetPage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid page ID",
			"success": false,
		})
	}

	page, err := h.pageService.GetPageByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error":   err.Error(),
			"message": "Page not found",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    page,
		"message": "Page retrieved successfully",
		"success": true,
	})
}

// GetWorkspacePages retrieves all pages in a workspace
func (h *PageHandler) GetWorkspacePages(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid workspace ID",
			"success": false,
		})
	}

	pages, err := h.pageService.GetPagesByWorkspaceID(c.Request().Context(), workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to retrieve pages",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    pages,
		"message": "Pages retrieved successfully",
		"success": true,
	})
}

// GetChildPages retrieves all child pages of a given page
func (h *PageHandler) GetChildPages(c echo.Context) error {
	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid parent page ID",
			"success": false,
		})
	}

	pages, err := h.pageService.GetChildPages(c.Request().Context(), parentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to retrieve child pages",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    pages,
		"message": "Child pages retrieved successfully",
		"success": true,
	})
}

// UpdatePage updates an existing page
func (h *PageHandler) UpdatePage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid page ID",
			"success": false,
		})
	}

	var page model.Page
	if err := c.Bind(&page); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid request payload",
			"success": false,
		})
	}

	page.ID = id
	if err := h.pageService.UpdatePage(c.Request().Context(), &page); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to update page",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    page,
		"message": "Page updated successfully",
		"success": true,
	})
}

// DeletePage deletes a page by ID
func (h *PageHandler) DeletePage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "Invalid page ID",
			"success": false,
		})
	}

	if err := h.pageService.DeletePage(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to delete page",
			"success": false,
		})
	}

	return c.NoContent(http.StatusNoContent)
}
