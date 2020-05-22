package transformations

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	_ "fmt"
	gogithub "github.com/google/go-github/github"
	"github.com/whosonfirst/go-webhookd/v2"
	_ "log"
	"net/url"
	"strconv"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "githubcommits", NewGitHubCommitsTransformation)

	if err != nil {
		panic(err)
	}
}

// see also: https://github.com/whosonfirst/go-whosonfirst-updated/issues/8

type GitHubCommitsTransformation struct {
	webhookd.WebhookTransformation
	ExcludeAdditions     bool
	ExcludeModifications bool
	ExcludeDeletions     bool
}

func NewGitHubCommitsTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	str_additions := q.Get("exclude_additions")
	str_modifications := q.Get("exclude_modifications")
	str_deletions := q.Get("exclude_deletions")

	exclude_additions := false
	exclude_modifications := false
	exclude_deletions := false

	if str_additions != "" {

		v, err := strconv.ParseBool(str_additions)

		if err != nil {
			return nil, err
		}

		exclude_additions = v
	}

	if str_modifications == "" {

		v, err := strconv.ParseBool(str_modifications)

		if err != nil {
			return nil, err
		}

		exclude_modifications = v
	}

	if str_deletions == "" {

		v, err := strconv.ParseBool(str_deletions)

		if err != nil {
			return nil, err
		}

		exclude_deletions = v
	}

	p := GitHubCommitsTransformation{
		ExcludeAdditions:     exclude_additions,
		ExcludeModifications: exclude_modifications,
		ExcludeDeletions:     exclude_deletions,
	}

	return &p, nil
}

func (p *GitHubCommitsTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {

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
