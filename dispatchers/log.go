package dispatchers

// PLEASE MAKE ME MORE SOPHISTICATED (20161016/thisisaaronland)

import (
	"context"
	"github.com/whosonfirst/go-webhookd"
	"log"
)

type LogDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewLogDispatcher(ctx context.Context) (*LogDispatcher, error) {

	n := LogDispatcher{}
	return &n, nil
}

func (n *LogDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	log.Println(string(body))
	return nil
}
