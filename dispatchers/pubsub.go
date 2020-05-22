package dispatchers

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
	"gopkg.in/redis.v1"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "pubsub", NewPubSubDispatcher)

	if err != nil {
		panic(err)
	}
}

type PubSubDispatcher struct {
	webhookd.WebhookDispatcher
	client  *redis.Client
	channel string
}

func NewPubSubDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	endpoint := u.Host
	channel := u.Path
	
	client := redis.NewTCPClient(&redis.Options{
		Addr: endpoint,
	})

	// defer client.Close()

	_, err = client.Ping().Result()

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
