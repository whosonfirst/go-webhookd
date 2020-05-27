package dispatcher

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "null", NewNullDispatcher)

	if err != nil {
		panic(err)
	}
}

type NullDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewNullDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	n := NullDispatcher{}
	return &n, nil
}

func (n *NullDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	return nil
}
