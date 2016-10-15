package receivers

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
)

func NewDispatcherFromConfig(config *webhookd.WebhookConfig) (webhookd.WebhookReceiver, error) {

	if config.Dispatcher.Name == "Insecure" {
		return NewInsecureReceiver()
	} else if config.Dispatcher.Name == "GitHub" {
		return NewGitHubReceiver(config.Receiver.Secret)
	} else {
		return nil, errors.New("Invalid receiver")
	}
}
