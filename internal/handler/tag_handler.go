package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/response"
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
		return response.BadRequest(c, "Invalid request payload", err)
	}

	if err := h.service.CreateTag(c.Request().Context(), tag); err != nil {
		return response.InternalServerError(c, "Failed to create tag", err)
	}

	return response.Created(c, "Tag created successfully", tag)
}

func (h *TagHandler) GetTags(c echo.Context) error {
	tags, err := h.service.GetTags(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get tags", err)
	}

	return response.Success(c, "Successfully retrieved tags", tags)
}

func (h *TagHandler) GetTagByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid tag ID", err)
	}

	tag, err := h.service.GetTagByID(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Tag not found", err)
	}

	return response.Success(c, "Successfully retrieved tag", tag)
}

func (h *TagHandler) UpdateTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid tag ID", err)
	}

	tag := new(model.Tag)
	if err := c.Bind(tag); err != nil {
		return response.BadRequest(c, "Invalid request payload", err)
	}
	tag.ID = uint(id)

	if err := h.service.UpdateTag(c.Request().Context(), tag); err != nil {
		return response.InternalServerError(c, "Failed to update tag", err)
	}

	return response.Success(c, "Tag updated successfully", tag)
}

func (h *TagHandler) DeleteTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid tag ID", err)
	}

	if err := h.service.DeleteTag(c.Request().Context(), uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete tag", err)
	}

	return response.Success(c, "Tag deleted successfully", nil)
}
