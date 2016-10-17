package transformations

import (
	"github.com/thisisaaronland/go-chicken"
	"github.com/whosonfirst/go-webhookd"
)

type ChickenTransformation struct {
	webhookd.WebhookTransformation
	chicken *chicken.Chicken
}

func NewChickenTransformation(lang string, clucking bool) (*ChickenTransformation, error) {

	ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

	if err != nil {
		return nil, err
	}

	tr := ChickenTransformation{
		chicken: ch,
	}

	return &tr, nil
}

func (tr *ChickenTransformation) Transform(body []byte) ([]byte, *webhookd.WebhookError) {

	txt := tr.chicken.TextToChicken(string(body))
	return []byte(txt), nil
}
