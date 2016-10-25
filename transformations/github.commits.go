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

type GitHubCommitsTransformation struct {
	webhookd.WebhookTransformation
}

func NewGitHubCommitsTransformation() (*GitHubCommitsTransformation, error) {

	p := GitHubCommitsTransformation{}
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

		for _, path := range c.Added {
			commit := []string{commit_hash, repo_name, path}
			wr.Write(commit)
		}

		for _, path := range c.Modified {
			commit := []string{commit_hash, repo_name, path}
			wr.Write(commit)
		}

		for _, path := range c.Removed {
			commit := []string{commit_hash, repo_name, path}
			wr.Write(commit)
		}
	}

	wr.Flush()

	return buf.Bytes(), nil
}
