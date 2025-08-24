package routes

import "github.com/labstack/echo/v4"

func (r *Routes) setupWorkspaceRoutes(v1 *echo.Group) {
	workspaces := v1.Group("/workspaces", r.authMiddleware.Auth())
	{
		workspaces.POST("", r.workspaceHandler.CreateWorkspace)
		workspaces.GET("", r.workspaceHandler.GetAllWorkspaces)
		workspaces.GET("/me", r.workspaceHandler.GetUserWorkspaces)
		workspaces.GET("/:id", r.workspaceHandler.GetWorkspaceByID)
		workspaces.PUT("/:id", r.workspaceHandler.UpdateWorkspace)
		workspaces.DELETE("/:id", r.workspaceHandler.DeleteWorkspace)

		// Workspace members
		workspaces.POST("/:id/members", r.workspaceHandler.AddMember)
		workspaces.GET("/:id/members", r.workspaceHandler.GetMembers)
		workspaces.PUT("/:id/members/:user_id", r.workspaceHandler.UpdateMemberRole)
		workspaces.DELETE("/:id/members/:user_id", r.workspaceHandler.RemoveMember)
	}
}
