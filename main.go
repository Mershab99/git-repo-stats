package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mershab99/git-repo-stats/handlers"
)

func main() {
	e := echo.New()

	//todo: handle the error!
	c, _ := handlers.NewContainer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CommitsPost - Get commits from one or more repositories
	e.POST("/commits", c.CommitsPost)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
