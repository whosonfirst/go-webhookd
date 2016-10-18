package transformations

import (
	"encoding/json"
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

	commits := make([]string, 0)

	for _, c := range event.Commits {

		for _, path := range c.Added {
			commits = append(commits, path)
		}

		for _, path := range c.Modified {
			commits = append(commits, path)
		}

		for _, path := range c.Removed {
			commits = append(commits, path)
		}
	}

	body, err = json.Marshal(commits)

	if err != nil {
		err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
		return nil, err
	}

	return body, nil
}
