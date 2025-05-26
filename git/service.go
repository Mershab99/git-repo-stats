package git

import (
	"fmt"
	"time"

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

	for _, repo := range repos {

		localPath, _, err := CloneWithAuth(repo)
		if err != nil {
			return nil, fmt.Errorf("failed to clone %s: %w", repo.Url, err)
		}
		//defer cleanup()

		commits, err := ExtractCommits(localPath, since)
		if err != nil {
			return nil, fmt.Errorf("failed to extract commits: %w", err)
		}

		results[repo.Url] = append(results[repo.Url], commits...)
	}

	return results, nil
}
