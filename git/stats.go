package git

import (
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mershab99/git-repo-stats/models"
)

func ExtractCommits(repoPath string, since time.Time) ([]models.CommitInfo, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	branches, err := r.Branches()
	if err != nil {
		return nil, err
	}

	var allCommits []models.CommitInfo
	seen := map[string]bool{}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return err
		}

		return cIter.ForEach(func(c *object.Commit) error {
			if c.Committer.When.Before(since) {
				return nil
			}
			if seen[c.Hash.String()] {
				return nil
			}
			seen[c.Hash.String()] = true

			allCommits = append(allCommits, models.CommitInfo{
				Hash:        c.Hash.String(),
				Message:     c.Message,
				AuthorName:  c.Author.Name,
				AuthorEmail: c.Author.Email,
				Timestamp:   c.Author.When,
			})
			return nil
		})
	})

	return allCommits, err
}
