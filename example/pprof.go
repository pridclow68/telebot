package main

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PProf is pprof
func PProf(e *echo.Group) {
	e.GET("/pprof/", func(c echo.Context) error {
		pprof.Index(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/allocs", func(c echo.Context) error {
		pprof.Handler("allocs").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/heap", func(c echo.Context) error {
		pprof.Handler("heap").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/goroutine", func(c echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/block", func(c echo.Context) error {
		pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/threadcreate", func(c echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/cmdline", func(c echo.Context) error {
		pprof.Cmdline(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/profile", func(c echo.Context) error {
		pprof.Profile(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/symbol", func(c echo.Context) error {
		pprof.Symbol(c.Response().Writer, c.Request())
		return nil
	})
	e.POST("/pprof/symbol", func(c echo.Context) error {
		pprof.Symbol(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/trace", func(c echo.Context) error {
		pprof.Trace(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/pprof/mutex", func(c echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
}
