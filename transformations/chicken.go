package transformations

import (
	"context"
	"github.com/aaronland/go-chicken"
	"github.com/whosonfirst/go-webhookd"
)

type ChickenTransformation struct {
	webhookd.WebhookTransformation
	chicken *chicken.Chicken
}

func NewChickenTransformation(ctx context.Context, lang string, clucking bool) (*ChickenTransformation, error) {

	ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

	if err != nil {
		return nil, err
	}

	tr := ChickenTransformation{
		chicken: ch,
	}

	return &tr, nil
}

func (tr *ChickenTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	txt := tr.chicken.TextToChicken(string(body))
	return []byte(txt), nil
}
