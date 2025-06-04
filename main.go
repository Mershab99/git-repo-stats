package main

import (
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mershab99/git-repo-stats/handlers"
	"time"
)

func main() {
	e := echo.New()

	//todo: handle the error!
	c, _ := handlers.NewContainer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(cache.New(&cache.Config{
		TTL:        15 * time.Minute,
		Methods:    []string{"GET", "POST"},
		StatusCode: []int{200, 204},
	}, c.InMemCache))

	// CommitsPost - Get commits from one or more repositories
	e.POST("/commits", c.CommitsPost)

	// StatusGet - Liveness probe endpoint
	e.GET("/status", c.StatusGet)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
