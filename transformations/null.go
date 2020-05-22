package transformations

import (
	"context"
	"github.com/whosonfirst/go-webhookd"
)

type NullTransformation struct {
	webhookd.WebhookTransformation
}

func NewNullTransformation(ctx context.Context) (*NullTransformation, error) {

	p := NullTransformation{}
	return &p, nil
}

func (p *NullTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {
	return body, nil
}
