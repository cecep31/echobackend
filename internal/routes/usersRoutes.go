package routes

import "github.com/labstack/echo/v4"

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users")
	{
		// Public routes
		users.GET("/:id", r.userHandler.GetByID)

		// Authenticated routes
		authUsers := users.Group("", r.authMiddleware.Auth())
		{
			authUsers.GET("/me", r.userHandler.GetMe)
			authUsers.GET("", r.userHandler.GetUsers, r.authMiddleware.AuthAdmin())
			authUsers.DELETE("/:id", r.userHandler.DeleteUser, r.authMiddleware.AuthAdmin())

			// Follow routes
			authUsers.POST("/follow", r.userFollowHandler.FollowUser)
			authUsers.DELETE("/:id/follow", r.userFollowHandler.UnfollowUser)
			authUsers.GET("/:id/follow-status", r.userFollowHandler.CheckFollowStatus)
			authUsers.GET("/:id/mutual-follows", r.userFollowHandler.GetMutualFollows)
		}

		// Follow-related public routes
		users.GET("/:id/followers", r.userFollowHandler.GetFollowers)
		users.GET("/:id/following", r.userFollowHandler.GetFollowing)
		users.GET("/:id/follow-stats", r.userFollowHandler.GetFollowStats)
	}
}
