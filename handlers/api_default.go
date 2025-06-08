package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mershab99/git-repo-stats/git"
	"github.com/mershab99/git-repo-stats/models"
	"net/http"
)

// CommitsHeatmapPost - Get commits from one or more repositories and return in a cal-heatmap compatible format
func (c *Container) CommitsHeatmapPost(ctx echo.Context) error {
	var req models.CommitsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Create a hash of the request body to use as a cache key
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to encode request")
	}
	hash := sha256.Sum256(bodyBytes)
	cacheKey := hex.EncodeToString(hash[:])

	// Try to get cached response
	if cached, err := c.InMemCache.Get([]byte(cacheKey)); err == nil {
		var cachedResp []models.HeatmapData
		if err := json.Unmarshal(cached, &cachedResp); err == nil {
			return ctx.JSON(http.StatusOK, cachedResp)
		}
	}

	// Get fresh commit data
	commitsByRepo, err := git.GetCommitStats(req.Repositories, int(req.Days))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to get commit stats")
	}

	// Count commits per day (in YYYY-MM-DD format)
	commitCountByDay := make(map[string]int)
	for _, commits := range commitsByRepo {
		for _, commit := range commits {
			day := commit.Timestamp.UTC().Format("2006-01-02")
			commitCountByDay[day]++
		}
	}

	// Convert map to array of HeatmapData
	var heatmapData []models.HeatmapData
	for day, count := range commitCountByDay {
		heatmapData = append(heatmapData, models.HeatmapData{
			Date:  day,
			Value: int32(count),
		})
	}

	// Store in cache
	if encoded, err := json.Marshal(heatmapData); err == nil {
		_ = c.InMemCache.Set([]byte(cacheKey), encoded, 15*60) // 15-minute TTL
	}

	return ctx.JSON(http.StatusOK, heatmapData)
}

// CommitsPost - Get commits from one or more repositories
func (c *Container) CommitsPost(ctx echo.Context) error {
	var req models.CommitsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var resp, err = git.GetCommitStats(req.Repositories, int(req.Days))
	if err != nil {
		fmt.Printf("Failed to get Git Commit Stats: %s", err)
		return ctx.JSON(http.StatusInternalServerError, "Failed To Get Commit Stats")
	}

	return ctx.JSON(http.StatusOK, resp)
}

// StatusGet - Liveness probe endpoint
func (c *Container) StatusGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}
