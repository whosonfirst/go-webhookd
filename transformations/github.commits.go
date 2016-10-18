package transformations

import (
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/github"
	"log"
)

type GitHubCommitsTransformation struct {
	webhookd.WebhookTransformation
}

func NewGitHubCommitsTransformation() (*GitHubCommitsTransformation, error) {

	p := GitHubCommitsTransformation{}
	return &p, nil
}

func (p *GitHubCommitsTransformation) Transform(body []byte) ([]byte, *webhookd.WebhookError) {

	event, err := github.UnmarshalEvent("push", body)

	if err != nil {

		err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
		return nil, err
	}

	log.Println(event)
	return body, nil
}
