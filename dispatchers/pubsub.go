package dispatchers

import (
	"fmt"
	"os"

	"github.com/whosonfirst/go-webhookd"
	"gopkg.in/redis.v1"
)

type PubSubDispatcher struct {
	webhookd.WebhookDispatcher
	client  *redis.Client
	channel string
}

func NewPubSubDispatcher(host string, port int, channel string) (*PubSubDispatcher, error) {

	if host == "" || port == "" {
		host = os.Getenv("WEBHOOKD_REDIS_HOST")
		port = os.Getenv("WEBHOOKD_REDIS_PORT")
	}
	password := os.Getenv("WEBHOOKD_REDIS_PASSWORD")

	endpoint := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewTCPClient(&redis.Options{
		Addr:     endpoint,
		Password: password, // read password from env, if your redis no password set, no need to set env WEBHOOKD_REDIS_PASSWORD
		DB:       0,        // use default DB
	})

	// defer client.Close()

	// with redis.v1 https://godoc.org/gopkg.in/redis.v1#NewTCPClient
	// you must run auth first before any commands, even ping
	if password != "" {
		client.Auth(os.Getenv("WEBHOOKD_REDIS_PASSWORD"))
	}

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

func (dispatcher *PubSubDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

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
