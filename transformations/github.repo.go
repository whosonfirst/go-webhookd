package transformations

import (
	"bytes"
	"context"
	"encoding/json"
	gogithub "github.com/google/go-github/github"
	"github.com/whosonfirst/go-webhookd/v2"
	_ "log"
	"net/url"
	"strconv"
)

// see also: https://github.com/whosonfirst/go-whosonfirst-updated/issues/8

type GitHubRepoTransformation struct {
	webhookd.WebhookTransformation
	ExcludeAdditions     bool
	ExcludeModifications bool
	ExcludeDeletions     bool
}

func NewGitHubRepoTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	str_additions := q.Get("exclude_additions")
	str_modifications := q.Get("exclude_modifications")
	str_deletions := q.Get("exclude_deletions")

	exclude_additions, err := strconv.ParseBool(str_additions)

	if err != nil {
		return nil, err
	}

	exclude_modifications, err := strconv.ParseBool(str_modifications)

	if err != nil {
		return nil, err
	}

	exclude_deletions, err := strconv.ParseBool(str_deletions)

	if err != nil {
		return nil, err
	}

	p := GitHubRepoTransformation{
		ExcludeAdditions:     exclude_additions,
		ExcludeModifications: exclude_modifications,
		ExcludeDeletions:     exclude_deletions,
	}

	return &p, nil
}

func (p *GitHubRepoTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

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
