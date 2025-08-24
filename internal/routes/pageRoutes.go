package routes

import "github.com/labstack/echo/v4"

func (r *Routes) setupPageRoutes(v1 *echo.Group) {
	pages := v1.Group("/pages", r.authMiddleware.Auth())
	{
		pages.POST("", r.pageHandler.CreatePage)
		pages.GET("/:id", r.pageHandler.GetPage)
		pages.PUT("/:id", r.pageHandler.UpdatePage)
		pages.DELETE("/:id", r.pageHandler.DeletePage)
		pages.GET("/workspace/:workspace_id", r.pageHandler.GetWorkspacePages)
		pages.GET("/children/:parent_id", r.pageHandler.GetChildPages)
	}
}
