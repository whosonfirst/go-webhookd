package dispatchers

import (
	"errors"
	"github.com/whosonfirst/go-webhookd"
)

func NewDispatcherFromConfig(config *webhookd.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	if config.Name == "PubSub" {
		return NewPubSubDispatcher(config.Host, config.Port, config.Channel)
	} else if config.Name == "Slack" {
		return NewSlackDispatcher(config.Config)
	} else {
		return nil, errors.New("Invalid dispatcher")
	}
}
