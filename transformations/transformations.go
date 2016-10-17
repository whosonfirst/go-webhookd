package transformations

import (
	"errors"
	"github.com/whosonfirst/go-webhookd"
)

func NewTransformationFromConfig(config *webhookd.WebhookTransformationConfig) (webhookd.WebhookTransformation, error) {

	switch config.Name {
	case "Chicken":
		return NewChickenTransformation(config.Language, config.Clucking)
	case "Null":
		return NewNullTransformation()
	case "Slack":
		return NewSlackTransformation()
	default:
		return nil, errors.New("Undefined transformation")
	}
}
