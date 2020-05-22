package dispatchers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v2"
	"github.com/whosonfirst/go-webhookd/v2/config"
	"net/url"
)

var dispatchers roster.Roster

type DispatcherInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error)

func NewDispatcherFromConfig(ctx context.Context, cfg *config.WebhookDispatcherConfig) (webhookd.WebhookDispatcher, error) {

	switch cfg.Name {
	case "Lambda":
		uri := fmt.Sprintf("lambda://dsn=%s&function=%s")
		return NewLambdaDispatcher(ctx, uri)
	case "Log":
		return NewLogDispatcher(ctx, "log://")
	case "Null":
		return NewNullDispatcher(ctx, "null://")
	case "PubSub":
		uri := fmt.Sprintf("pubsub:%s//%s/%s", cfg.Port, cfg.Host, cfg.Channel)
		return NewPubSubDispatcher(ctx, uri)
	case "Slack":
		uri := fmt.Sprintf("slack://")
		return NewSlackDispatcher(ctx, uri)
	default:
		msg := fmt.Sprintf("Undefined dispatcher: '%s'", cfg.Name)
		return nil, errors.New(msg)
	}
}

func NewDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	err := ensureDispatcherRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := dispatchers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(DispatcherInitializationFunc)
	return init_func(ctx, uri)
}

func RegisterDispatcher(ctx context.Context, scheme string, init_func DispatcherInitializationFunc) error {

	err := ensureDispatcherRoster()

	if err != nil {
		return err
	}

	return dispatchers.Register(ctx, scheme, init_func)
}

func ensureDispatcherRoster() error {

	if dispatchers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		dispatchers = r
	}

	return nil
}

func Dispatchers() []string {
	ctx := context.Background()
	return dispatchers.Drivers(ctx)
}
