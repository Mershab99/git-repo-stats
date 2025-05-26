package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mershab99/git-repo-stats/git"
	"github.com/mershab99/git-repo-stats/models"
	"net/http"
)

// CommitsPost - Get commits from one or more repositories
func (c *Container) CommitsPost(ctx echo.Context) error {
	// This will be auto-wired from the OpenAPI-generated handler interface
	var req models.CommitsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var resp, err = git.GetCommitStats(req.Repositories, int(req.Days))
	if err != nil {
		fmt.Printf("Failed to get Git Commit Stats: %s", err)
	}
	print(resp)

	return ctx.JSON(http.StatusOK, resp)
}
