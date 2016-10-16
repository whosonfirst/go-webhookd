package receivers

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
)

func NewReceiverFromConfig(config *webhookd.WebhookReceiverConfig) (webhookd.WebhookReceiver, error) {

	if config.Name == "Insecure" {
		return NewInsecureReceiver()
	} else if config.Name == "GitHub" {
		return NewGitHubReceiver(config.Secret)
	} else {
		return nil, errors.New("Invalid receiver")
	}
}
