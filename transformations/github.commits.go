package transformations

import (
	"encoding/json"
	"fmt"
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

	repo := event.Repo
	repo_name := *repo.Name

	commits := make([]string, 0)

	for _, c := range event.Commits {

		for _, path := range c.Added {
			commits = append(commits, fmt.Sprintf("%s#%s", repo_name, path))
		}

		for _, path := range c.Modified {
			commits = append(commits, fmt.Sprintf("%s#%s", repo_name, path))
		}

		for _, path := range c.Removed {
			commits = append(commits, fmt.Sprintf("%s#%s", repo_name, path))
		}
	}

	body, err = json.Marshal(commits)

	if err != nil {
		err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
		return nil, err
	}

	return body, nil
}
