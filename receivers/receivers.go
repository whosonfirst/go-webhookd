package receivers

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd"
)

func NewReceiverFromConfig(config *webhookd.WebhookReceiverConfig) (webhookd.WebhookReceiver, error) {

	if config.Name == "Insecure" {
		return NewInsecureReceiver()
	} else if config.Name == "GitHub" {
		return NewGitHubReceiver(config.Secret)
	} else if config.Name == "Slack" {
		return NewSlackReceiver()
	} else {
		msg := fmt.Sprintf("Invalid receiver: '%s'", config.Name)
		return nil, errors.New(msg)
	}
}
