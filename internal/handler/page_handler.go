package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"echobackend/internal/model"
	"echobackend/internal/service"
)

type PageHandler struct {
	pageService *service.PageService
}

func NewPageHandler(pageService *service.PageService) *PageHandler {
	return &PageHandler{pageService: pageService}
}

// CreatePage handles the creation of a new page
func (h *PageHandler) CreatePage(c echo.Context) error {
	var page model.Page
	if err := c.Bind(&page); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)
	page.CreatedBy = userID

	if err := h.pageService.CreatePage(&page); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, page)
}

// GetPage retrieves a page by ID
func (h *PageHandler) GetPage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page ID"})
	}

	page, err := h.pageService.GetPageByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Page not found"})
	}

	return c.JSON(http.StatusOK, page)
}

// GetWorkspacePages retrieves all pages in a workspace
func (h *PageHandler) GetWorkspacePages(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	pages, err := h.pageService.GetPagesByWorkspaceID(workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, pages)
}

// GetChildPages retrieves all child pages of a given page
func (h *PageHandler) GetChildPages(c echo.Context) error {
	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid parent page ID"})
	}

	pages, err := h.pageService.GetChildPages(parentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, pages)
}

// UpdatePage updates an existing page
func (h *PageHandler) UpdatePage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page ID"})
	}

	var page model.Page
	if err := c.Bind(&page); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	page.ID = id
	if err := h.pageService.UpdatePage(&page); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, page)
}

// DeletePage deletes a page by ID
func (h *PageHandler) DeletePage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page ID"})
	}

	if err := h.pageService.DeletePage(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}