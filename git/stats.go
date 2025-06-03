package git

import (
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mershab99/git-repo-stats/models"
)

func ExtractCommits(repoPath string, since time.Time) ([]models.CommitInfo, error) {
	// Open the git repository at the specified path
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Retrieve local branches
	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var (
		allCommits []models.CommitInfo
		seen       = make(map[string]bool)
	)

	// Iterate over each branch
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		// Get commit history starting from the branch tip
		logIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return fmt.Errorf("failed to get log for branch %s: %w", ref.Name(), err)
		}

		// Iterate over commits
		return logIter.ForEach(func(commit *object.Commit) error {
			// Skip commits before the specified time
			if commit.Committer.When.Before(since) {
				return nil
			}
			// Skip duplicate commits
			if seen[commit.Hash.String()] {
				return nil
			}
			seen[commit.Hash.String()] = true

			// Append commit info
			allCommits = append(allCommits, models.CommitInfo{
				Hash:        commit.Hash.String(),
				Message:     commit.Message,
				AuthorName:  commit.Author.Name,
				AuthorEmail: commit.Author.Email,
				Timestamp:   commit.Author.When,
			})

			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return allCommits, nil
}
