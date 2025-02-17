package routes

import (
	"echobackend/internal/handler"
	"echobackend/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	tagHandler     *handler.TagHandler
	pageHandler    *handler.PageHandler
}

func NewRoutes(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	tagHandler *handler.TagHandler,
	pageHandler *handler.PageHandler,
) *Routes {
	return &Routes{
		userHandler:    userHandler,
		postHandler:    postHandler,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
		tagHandler:     tagHandler,
		pageHandler:    pageHandler,
	}
}

func (r *Routes) Setup(e *echo.Echo) {
	// API Group
	v1 := e.Group("/v1")
	r.setupV1Routes(v1)
}

func (r *Routes) setupV1Routes(v1 *echo.Group) {
	r.setupUserRoutes(v1)
	r.setupPostRoutes(v1)
	r.setupAuthRoutes(v1)
	r.setupTagRoutes(v1)
	r.setupPageRoutes(v1)
}

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users", r.authMiddleware.Auth())
	{
		users.GET("/:id", r.userHandler.GetByID)
		users.GET("", r.userHandler.GetUsers)
		users.DELETE("/:id", r.userHandler.DeleteUser, r.authMiddleware.AuthAdmin())
	}
}

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")
	{
		posts.GET("/u/:username/:slug", r.postHandler.GetPostBySlugAndUsername)
		posts.GET("", r.postHandler.GetPosts)
		posts.DELETE("/:id", r.postHandler.DeletePost, r.authMiddleware.Auth())
		posts.GET("/random", r.postHandler.GetPostsRandom)
		posts.GET("/:id", r.postHandler.GetPost)
		posts.GET("/mine", r.postHandler.GetMyPosts, r.authMiddleware.Auth())
		posts.GET("/username/:username", r.postHandler.GetPostsByUsername)
		posts.POST("/image", r.postHandler.UploadImagePosts, r.authMiddleware.Auth())
	}
}

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}
}

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

func (r *Routes) setupTagRoutes(v1 *echo.Group) {
	tags := v1.Group("/tags", r.authMiddleware.Auth())
	{
		tags.POST("", r.tagHandler.CreateTag)
		tags.GET("", r.tagHandler.GetTags)
		tags.GET("/:id", r.tagHandler.GetTagByID)
		tags.PUT("/:id", r.tagHandler.UpdateTag)
		tags.DELETE("/:id", r.tagHandler.DeleteTag)
	}
}
