package transformations

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	_ "fmt"
	gogithub "github.com/google/go-github/github"
	"github.com/whosonfirst/go-webhookd"
	_ "log"
)

// see also: https://github.com/whosonfirst/go-whosonfirst-updated/issues/8

type GitHubCommitsTransformation struct {
	webhookd.WebhookTransformation
	ExcludeAdditions     bool
	ExcludeModifications bool
	ExcludeDeletions     bool
}

func NewGitHubCommitsTransformation(exclude_additions bool, exclude_modifications bool, exclude_deletions bool) (*GitHubCommitsTransformation, error) {

	p := GitHubCommitsTransformation{
		ExcludeAdditions:     exclude_additions,
		ExcludeModifications: exclude_modifications,
		ExcludeDeletions:     exclude_deletions,
	}

	return &p, nil
}

func (p *GitHubCommitsTransformation) Transform(body []byte) ([]byte, *webhookd.WebhookError) {

	var event gogithub.PushEvent

	err := json.Unmarshal(body, &event)

	if err != nil {
		err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
		return nil, err
	}

	buf := new(bytes.Buffer)
	wr := csv.NewWriter(buf)

	repo := event.Repo
	repo_name := *repo.Name
	commit_hash := *event.HeadCommit.ID

	for _, c := range event.Commits {

		if !p.ExcludeAdditions {
			for _, path := range c.Added {
				commit := []string{commit_hash, repo_name, path}
				wr.Write(commit)
			}
		}

		if !p.ExcludeModifications {
			for _, path := range c.Modified {
				commit := []string{commit_hash, repo_name, path}
				wr.Write(commit)
			}
		}

		if !p.ExcludeDeletions {
			for _, path := range c.Removed {
				commit := []string{commit_hash, repo_name, path}
				wr.Write(commit)
			}
		}
	}

	wr.Flush()

	return buf.Bytes(), nil
}
