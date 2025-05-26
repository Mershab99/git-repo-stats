package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mershab99/git-repo-stats/models"
	"os"
)

func CloneWithAuth(repo models.RepoConfig) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cloneOptions := &git.CloneOptions{
		URL:      repo.Url,
		Progress: os.Stdout,
		Depth:    0, // full history
	}

	if repo.Auth.Token != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: "token", // can be anything for GitHub
			Password: repo.Auth.Token,
		}
	} else if repo.Auth.Username != "" && repo.Auth.Password != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: repo.Auth.Username,
			Password: repo.Auth.Password,
		}
	}

	_, err = git.PlainClone(tmpDir, true, cloneOptions)
	if err != nil {
		return "", nil, fmt.Errorf("clone failed: %w", err)
	}

	cleanup := func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			return
		}
	}

	return tmpDir, cleanup, nil
}
