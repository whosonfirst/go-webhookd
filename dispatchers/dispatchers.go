package dispatchers

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
)

func NewDispatcherFromConfig(config *webhookd.WebhookConfig) (webhookd.WebhookDispatcher, error) {

	if config.Dispatcher.Name == "PubSub" {
		return NewPubSubDispatcher(config.Dispatcher.Host, config.Dispatcher.Port, config.Dispatcher.Channel)
	} else {
		return nil, errors.New("Invalid dispatcher")
	}
}
