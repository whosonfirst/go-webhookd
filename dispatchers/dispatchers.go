package dispatchers

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
)

func NewDispatcherFromConfig(config *webhookd.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	if config.Name == "PubSub" {
		return NewPubSubDispatcher(config.Host, config.Port, config.Channel)
	} else {
		return nil, errors.New("Invalid dispatcher")
	}
}
