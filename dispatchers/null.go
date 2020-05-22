package dispatchers

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
)

type NullDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewNullDispatcher(ctx context.Context) (*NullDispatcher, error) {

	n := NullDispatcher{}
	return &n, nil
}

func (n *NullDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	return nil
}
