package dispatchers

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/config"
)

func NewDispatcherFromConfig(ctx context.Context, cfg *config.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	switch cfg.Name {
	case "Lambda":
		return NewLambdaDispatcher(ctx, cfg.DSN, cfg.Function)
	case "Log":
		return NewLogDispatcher(ctx)
	case "Null":
		return NewNullDispatcher(ctx)
	case "PubSub":
		return NewPubSubDispatcher(ctx, cfg.Host, cfg.Port, cfg.Channel)
	case "Slack":
		return NewSlackDispatcher(ctx, cfg.Config)
	default:
		msg := fmt.Sprintf("Undefined dispatcher: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
