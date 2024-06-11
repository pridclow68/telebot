package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	// debug mode
	e.Debug = true
	// Set log level
	if !e.Debug {
		log.SetLevel(log.INFO)
	} else {
		log.SetLevel(log.DEBUG)
	}
	// Body Dump
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if c.Echo().Debug {
			requestDump, _ := httputil.DumpRequest(c.Request(), true)
			fmt.Printf("request: %s\n\n", requestDump)

			reqContentType := http.DetectContentType(reqBody)
			if strings.Contains(reqContentType, "text/") {
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, reqBody)
			} else {
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, reqContentType, len(reqBody)))
			}
			resContentType := http.DetectContentType(resBody)
			if strings.Contains(resContentType, "text/") {
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, resBody)
			} else {
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, resContentType, len(resBody)))
			}
		}
	}))

	// debug
	PProf(e.Group("/debug"))
	// tgbot
	TGBot(e.Group("/tgbot"))

	// Start server
	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
