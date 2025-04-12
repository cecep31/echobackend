package routes

import (
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"net/http"
	"net/http/pprof"
	"time"

	echomidleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	config           *config.Config
	userHandler      *handler.UserHandler
	postHandler      *handler.PostHandler
	authHandler      *handler.AuthHandler
	authMiddleware   *middleware.AuthMiddleware
	tagHandler       *handler.TagHandler
	pageHandler      *handler.PageHandler
	workspaceHandler *handler.WorkspaceHandler
}

func NewRoutes(
	config *config.Config,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	tagHandler *handler.TagHandler,
	pageHandler *handler.PageHandler,
	workspaceHandler *handler.WorkspaceHandler,
) *Routes {
	return &Routes{
		config:           config,
		userHandler:      userHandler,
		postHandler:      postHandler,
		authHandler:      authHandler,
		authMiddleware:   authMiddleware,
		tagHandler:       tagHandler,
		pageHandler:      pageHandler,
		workspaceHandler: workspaceHandler,
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
	r.setupWorkspaceRoutes(v1)
	if r.config.DEBUG {
		r.setupDebugRoutes(v1)
	}
}

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users", r.authMiddleware.Auth())
	{
		users.GET("/:id", r.userHandler.GetByID)
		users.GET("", r.userHandler.GetUsers, r.authMiddleware.AuthAdmin())
		users.DELETE("/:id", r.userHandler.DeleteUser, r.authMiddleware.AuthAdmin())
	}
}

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")
	{
		posts.GET("/username/:username", r.postHandler.GetPostsByUsername)
		posts.GET("/u/:username/:slug", r.postHandler.GetPostBySlugAndUsername)
		posts.GET("", r.postHandler.GetPosts)
		posts.PUT("/:id", r.postHandler.UpdatePost, r.authMiddleware.Auth())
		posts.DELETE("/:id", r.postHandler.DeletePost, r.authMiddleware.Auth())
		posts.GET("/random", r.postHandler.GetPostsRandom)
		posts.GET("/:id", r.postHandler.GetPost)
		posts.GET("/mine", r.postHandler.GetMyPosts, r.authMiddleware.Auth())
		posts.POST("/image", r.postHandler.UploadImagePosts, r.authMiddleware.Auth())
	}
}

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	confratelimit := echomidleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(5), ExpiresIn: 5 * time.Minute}
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
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

func (r *Routes) setupDebugRoutes(v1 *echo.Group) {
	debug := v1.Group("/debug")
	debug.GET("/pprof/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	debug.GET("/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	debug.GET("/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	debug.GET("/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	debug.GET("/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	debug.GET("/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
	debug.GET("/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
	debug.GET("/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Handler("allocs").ServeHTTP)))
	debug.GET("/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
	debug.GET("/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
}
