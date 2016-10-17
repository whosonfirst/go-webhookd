package dispatchers

import (
	"github.com/whosonfirst/go-webhookd"
)

type NullDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewNullDispatcher() (*NullDispatcher, error) {

	n := NullDispatcher{}
	return &n, nil
}

func (n *NullDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	return nil
}
