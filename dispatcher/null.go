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

// NullDispatcher implements the `webhookd.WebhookDispatcher` interface for dispatching messages to nowhere.
type NullDispatcher struct {
	webhookd.WebhookDispatcher
}

// NewNullDispatcher returns a new `NullDispatcher` instance that dispatches messages to nowhere
// configured by 'uri' in the form of:
//
//	null://
func NewNullDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	d := NullDispatcher{}
	return &d, nil
}

// Dispatch sends 'body' to nowhere.
func (d *NullDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	return nil
}
