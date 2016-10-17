package dispatchers

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/config"
)

func NewDispatcherFromConfig(cfg *config.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	switch cfg.Name {
	case "Log":
		return NewLogDispatcher()
	case "Null":
		return NewNullDispatcher()
	case "PubSub":
		return NewPubSubDispatcher(cfg.Host, cfg.Port, cfg.Channel)
	case "Slack":
		return NewSlackDispatcher(cfg.Config)
	default:
		msg := fmt.Sprintf("Undefined dispatcher: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
