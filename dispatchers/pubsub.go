package dispatchers

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v2"
	"gopkg.in/redis.v1"
)

type PubSubDispatcher struct {
	webhookd.WebhookDispatcher
	client  *redis.Client
	channel string
}

func NewPubSubDispatcher(ctx context.Context, host string, port int, channel string) (*PubSubDispatcher, error) {

	endpoint := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewTCPClient(&redis.Options{
		Addr: endpoint,
	})

	// defer client.Close()

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	dispatcher := PubSubDispatcher{
		client:  client,
		channel: channel,
	}

	return &dispatcher, nil
}

func (dispatcher *PubSubDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	rsp := dispatcher.client.Publish(dispatcher.channel, string(body))

	_, err := rsp.Result()

	if err != nil {

		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}

	return nil
}
