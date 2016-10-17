package dispatchers

// PLEASE MAKE ME MORE SOPHISTICATED (20161016/thisisaaronland)

import (
	"github.com/whosonfirst/go-webhookd"
	"log"
)

type LogDispatcher struct {
	webhookd.WebhookDispatcher
}

func NewLogDispatcher() (*LogDispatcher, error) {

	n := LogDispatcher{}
	return &n, nil
}

func (n *LogDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	log.Println(string(body))
	return nil
}
