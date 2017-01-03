package receivers

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/config"
)

func NewReceiverFromConfig(cfg *config.WebhookReceiverConfig) (webhookd.WebhookReceiver, error) {

	switch cfg.Name {
	case "GitHub":
		return NewGitHubReceiver(cfg.Secret, cfg.Ref)
	case "Insecure":
		return NewInsecureReceiver()
	case "Slack":
		return NewSlackReceiver()
	default:
		msg := fmt.Sprintf("Undefined receiver: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}
