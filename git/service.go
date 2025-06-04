package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/labstack/gommon/log"

	"github.com/mershab99/git-repo-stats/models"
)

func GetCommitStats(repos []models.RepoConfig, sinceDays int) (map[string][]models.CommitInfo, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	results := make(map[string][]models.CommitInfo)
	since := time.Now().AddDate(0, 0, -sinceDays)

	for _, repoConfig := range repos {
		auth := getAuthMethod(repoConfig.Auth)

		commits, err := getRemoteCommits(repoConfig.Url, auth, since)
		if err != nil {
			return nil, fmt.Errorf("failed to get commits from %s: %w", repoConfig.Url, err)
		}

		results[repoConfig.Url] = commits
	}

	//log.Debugf("%v", results)

	return results, nil
}

func getRemoteCommits(repoURL string, auth transport.AuthMethod, since time.Time) ([]models.CommitInfo, error) {
	memStore := memory.NewStorage()
	fs := memfs.New()

	repo, err := git.Init(memStore, fs)
	if err != nil {
		return nil, fmt.Errorf("failed to init in-memory repo: %w", err)
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create remote: %w", err)
	}

	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return nil, fmt.Errorf("failed to list remote refs: %w", err)
	}

	var allCommits []models.CommitInfo
	seen := make(map[string]bool)

	for _, ref := range refs {
		if !ref.Name().IsBranch() {
			continue
		}

		// Fetch shallow copy of branch
		err := repo.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{
				config.RefSpec(fmt.Sprintf("+%s:%s", ref.Name(), ref.Name())),
			},
			Depth: 1000,
			Auth:  auth,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, fmt.Errorf("failed to fetch branch %s: %w", ref.Name(), err)
		}

		// Get the branch reference
		branchRef, err := repo.Reference(ref.Name(), true)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve branch %s: %w", ref.Name(), err)
		}

		logIter, err := repo.Log(&git.LogOptions{From: branchRef.Hash()})
		if err != nil {
			// Log and continue to the next branch instead of failing
			log.Warnf("Could not get log for branch %s: %v", ref.Name(), err)
			continue
		}

		var count int
		err = logIter.ForEach(func(commit *object.Commit) error {
			if commit.Committer.When.Before(since) {
				return nil // Too old, skip
			}
			if seen[commit.Hash.String()] {
				return nil // Already counted
			}

			seen[commit.Hash.String()] = true
			allCommits = append(allCommits, models.CommitInfo{
				Hash:        commit.Hash.String(),
				Message:     commit.Message,
				AuthorName:  commit.Author.Name,
				AuthorEmail: commit.Author.Email,
				Timestamp:   commit.Author.When,
			})

			count++

			return nil
		})

		if err != nil {
			// Only log, do not return the error
			continue
		}

		log.Printf("Finished iterating commits in branch %s", ref.Name())
		log.Printf("Branch %s: collected %d commits", ref.Name(), count)
	}

	return allCommits, nil
}

func getAuthMethod(auth models.Auth) transport.AuthMethod {
	if auth.Token != "" {
		return &http.BasicAuth{
			Username: "token",
			Password: auth.Token,
		}
	} else if auth.Username != "" && auth.Password != "" {
		return &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Password,
		}
	}
	return nil
}
