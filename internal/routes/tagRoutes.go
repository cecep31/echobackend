package routes

import "github.com/labstack/echo/v4"

func (r *Routes) setupTagRoutes(v1 *echo.Group) {
	tags := v1.Group("/tags")
	{
		tags.POST("", r.tagHandler.CreateTag, r.authMiddleware.Auth())
		tags.GET("", r.tagHandler.GetTags)
		tags.GET("/:id", r.tagHandler.GetTagByID)
		tags.PUT("/:id", r.tagHandler.UpdateTag, r.authMiddleware.Auth(), r.authMiddleware.AuthAdmin())
		tags.DELETE("/:id", r.tagHandler.DeleteTag, r.authMiddleware.Auth(), r.authMiddleware.AuthAdmin())
	}
}
