package receivers

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v2"
	"github.com/whosonfirst/go-webhookd/v2/config"
)

func NewReceiverFromConfig(ctx context.Context, cfg *config.WebhookReceiverConfig) (webhookd.WebhookReceiver, error) {

	switch cfg.Name {
	case "GitHub":
		return NewGitHubReceiver(ctx, cfg.Secret, cfg.Ref)
	case "Insecure":
		return NewInsecureReceiver(ctx)
	case "Slack":
		return NewSlackReceiver(ctx)
	default:
		msg := fmt.Sprintf("Undefined receiver: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
