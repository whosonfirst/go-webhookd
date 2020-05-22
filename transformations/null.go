package transformations

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
)

type NullTransformation struct {
	webhookd.WebhookTransformation
}

func NewNullTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	p := NullTransformation{}
	return &p, nil
}

func (p *NullTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {
	return body, nil
}
