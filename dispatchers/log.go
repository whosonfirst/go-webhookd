package dispatchers

// PLEASE MAKE ME MORE SOPHISTICATED (20161016/thisisaaronland)

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
	"log"
)

type LogDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewLogDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	n := LogDispatcher{}
	return &n, nil
}

func (n *LogDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	log.Println(string(body))
	return nil
}
