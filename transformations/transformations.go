package transformations

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/config"
)

func NewTransformationFromConfig(cfg *config.WebhookTransformationConfig) (webhookd.WebhookTransformation, error) {

	switch cfg.Name {
	case "Chicken":
		return NewChickenTransformation(cfg.Language, cfg.Clucking)
	case "GitHubCommits":
		return NewGitHubCommitsTransformation(cfg.ExcludeAdditions, cfg.ExcludeModifications, cfg.ExcludeDeletions)
	case "Null":
		return NewNullTransformation()
	case "SlackText":
		return NewSlackTextTransformation()
	default:
		msg := fmt.Sprintf("Undefined transformation: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
