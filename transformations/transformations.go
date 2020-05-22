package transformations

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v2"
	"github.com/whosonfirst/go-webhookd/v2/config"
)

func NewTransformationFromConfig(ctx context.Context, cfg *config.WebhookTransformationConfig) (webhookd.WebhookTransformation, error) {

	switch cfg.Name {
	case "Chicken":
		return NewChickenTransformation(ctx, cfg.Language, cfg.Clucking)
	case "GitHubCommits":
		return NewGitHubCommitsTransformation(ctx, cfg.ExcludeAdditions, cfg.ExcludeModifications, cfg.ExcludeDeletions)
	case "GitHubRepo":
		return NewGitHubRepoTransformation(ctx, cfg.ExcludeAdditions, cfg.ExcludeModifications, cfg.ExcludeDeletions)
	case "Null":
		return NewNullTransformation(ctx)
	case "SlackText":
		return NewSlackTextTransformation(ctx)
	default:
		msg := fmt.Sprintf("Undefined transformation: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
