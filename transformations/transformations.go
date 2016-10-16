package transformations

import (
	"errors"
	"github.com/whosonfirst/go-webhookd"
)

func NewTransformationFromConfig(config *webhookd.WebhookTransformationConfig) (webhookd.WebhookTransformation, error) {

	if config.Name == "Chicken" {
		return NewChickenTransformation(config.Language, config.Clucking)
	} else if config.Name == "Null" {
		return NewNullTransformation()
	} else {
		return nil, errors.New("Invalid receiver")
	}
}
