package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	service service.TagService
}

func NewTagHandler(service service.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) CreateTag(c echo.Context) error {
	tag := new(model.Tag)
	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request payload"})
	}

	if err := h.service.CreateTag(tag); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": tag})
}

func (h *TagHandler) GetTags(c echo.Context) error {
	tags, err := h.service.GetTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": tags})
}

func (h *TagHandler) GetTagByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid tag ID"})
	}

	tag, err := h.service.GetTagByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": tag})
}

func (h *TagHandler) UpdateTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid tag ID"})
	}

	tag := new(model.Tag)
	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request payload"})
	}
	tag.ID = uint(id)

	if err := h.service.UpdateTag(tag); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "Tag updated successfully", "success": true, "data": tag})
}

func (h *TagHandler) DeleteTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid tag ID"})
	}

	if err := h.service.DeleteTag(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "Tag deleted successfully", "success": true})
}
