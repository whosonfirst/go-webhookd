package receivers

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
)

func NewReceiverFromConfig(config *webhookd.WebhookConfig) (webhookd.WebhookReceiver, error) {

	if config.Receiver.Name == "Insecure" {
		return NewInsecureReceiver()
	} else if config.Receiver.Name == "GitHub" {
		return NewGitHubReceiver(config.Receiver.Secret)
	} else {
		return nil, errors.New("Invalid receiver")
	}
}
