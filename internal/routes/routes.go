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
	config             *config.Config
	userHandler        *handler.UserHandler
	postHandler        *handler.PostHandler
	authHandler        *handler.AuthHandler
	authMiddleware     *middleware.AuthMiddleware
	tagHandler         *handler.TagHandler
	pageHandler        *handler.PageHandler
	workspaceHandler   *handler.WorkspaceHandler
	commentHandler     *handler.CommentHandler
	postViewHandler    *handler.PostViewHandler
	userFollowHandler  *handler.UserFollowHandler
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
	commentHandler *handler.CommentHandler,
	postViewHandler *handler.PostViewHandler,
	userFollowHandler *handler.UserFollowHandler,
) *Routes {
	return &Routes{
		config:            config,
		userHandler:       userHandler,
		postHandler:       postHandler,
		authHandler:       authHandler,
		authMiddleware:    authMiddleware,
		tagHandler:        tagHandler,
		pageHandler:       pageHandler,
		workspaceHandler:  workspaceHandler,
		commentHandler:    commentHandler,
		postViewHandler:   postViewHandler,
		userFollowHandler: userFollowHandler,
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

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")
	{
		posts.POST("", r.postHandler.CreatePost, r.authMiddleware.Auth())
		posts.GET("/username/:username", r.postHandler.GetPostsByUsername)
		posts.GET("/u/:username/:slug", r.postHandler.GetPostBySlugAndUsername)
		posts.GET("/tag/:tag", r.postHandler.GetPostsByTag)
		posts.GET("", r.postHandler.GetPosts)
		posts.PUT("/:id", r.postHandler.UpdatePost, r.authMiddleware.Auth())
		posts.DELETE("/:id", r.postHandler.DeletePost, r.authMiddleware.Auth())
		posts.GET("/random", r.postHandler.GetPostsRandom)
		posts.GET("/:id", r.postHandler.GetPost)
		posts.GET("/mine", r.postHandler.GetMyPosts, r.authMiddleware.Auth())
		posts.POST("/image", r.postHandler.UploadImagePosts, r.authMiddleware.Auth())

		// Comment routes
		posts.GET("/:id/comments", r.commentHandler.GetCommentsByPostID)
		posts.POST("/:id/comments", r.commentHandler.CreateComment, r.authMiddleware.Auth())
		posts.PUT("/:id/comments/:comment_id", r.commentHandler.UpdateComment, r.authMiddleware.Auth())
		posts.DELETE("/:id/comments/:comment_id", r.commentHandler.DeleteComment, r.authMiddleware.Auth())
		
		// View routes
		posts.POST("/:id/view", r.postViewHandler.RecordView) // Can be called by anonymous users
		posts.GET("/:id/views", r.postViewHandler.GetPostViews, r.authMiddleware.Auth())
		posts.GET("/:id/view-stats", r.postViewHandler.GetPostViewStats)
		posts.GET("/:id/viewed", r.postViewHandler.CheckUserViewed, r.authMiddleware.Auth())
	}
}

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	confratelimit := echomidleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(5), ExpiresIn: 5 * time.Minute, Burst: 5}
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
		auth.POST("/check-username", r.authHandler.CheckUsername)
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
	tags := v1.Group("/tags")
	{
		tags.POST("", r.tagHandler.CreateTag, r.authMiddleware.Auth())
		tags.GET("", r.tagHandler.GetTags)
		tags.GET("/:id", r.tagHandler.GetTagByID)
		tags.PUT("/:id", r.tagHandler.UpdateTag, r.authMiddleware.Auth(), r.authMiddleware.AuthAdmin())
		tags.DELETE("/:id", r.tagHandler.DeleteTag, r.authMiddleware.Auth(), r.authMiddleware.AuthAdmin())
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
