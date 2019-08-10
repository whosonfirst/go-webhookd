package transformations

import (
	"bytes"
	"encoding/json"
	gogithub "github.com/google/go-github/github"
	"github.com/whosonfirst/go-webhookd"
	_ "log"
)

// see also: https://github.com/whosonfirst/go-whosonfirst-updated/issues/8

type GitHubRepoTransformation struct {
	webhookd.WebhookTransformation
	ExcludeAdditions     bool
	ExcludeModifications bool
	ExcludeDeletions     bool
}

func NewGitHubRepoTransformation(exclude_additions bool, exclude_modifications bool, exclude_deletions bool) (*GitHubRepoTransformation, error) {

	p := GitHubRepoTransformation{
		ExcludeAdditions:     exclude_additions,
		ExcludeModifications: exclude_modifications,
		ExcludeDeletions:     exclude_deletions,
	}

	return &p, nil
}

func (p *GitHubRepoTransformation) Transform(body []byte) ([]byte, *webhookd.WebhookError) {

	var event gogithub.PushEvent

	err := json.Unmarshal(body, &event)

	if err != nil {
		err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
		return nil, err
	}

	buf := new(bytes.Buffer)

	repo := event.Repo
	repo_name := *repo.Name

	has_updates := false

	for _, c := range event.Commits {

		if !p.ExcludeAdditions {

			if len(c.Added) > 0 {
				has_updates = true
			}
		}

		if !p.ExcludeModifications {

			if len(c.Modified) > 0 {
				has_updates = true
			}
		}

		if !p.ExcludeDeletions {

			if len(c.Removed) > 0 {
				has_updates = true
			}
		}
	}

	if has_updates {
		buf.WriteString(repo_name)
	}

	return buf.Bytes(), nil
}
