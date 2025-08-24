package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

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
