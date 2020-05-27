package transformation

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "null", NewNullTransformation)

	if err != nil {
		panic(err)
	}
}

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
