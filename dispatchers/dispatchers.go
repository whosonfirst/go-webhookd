package dispatchers

import (
	"errors"
	"github.com/whosonfirst/go-webhookd"
)

func NewDispatcherFromConfig(config *webhookd.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	switch config.Name {
	case "Log":
		return NewLogDispatcher()
	case "Null":
		return NewNullDispatcher()
	case "PubSub":
		return NewPubSubDispatcher(config.Host, config.Port, config.Channel)
	case "Slack":
		return NewSlackDispatcher(config.Config)
	default:
		return nil, errors.New("Invalid dispatcher")
	}
}
