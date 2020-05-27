package receiver

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
)

var receivers roster.Roster

type ReceiverInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookReceiver, error)

func NewReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	err := ensureReceiverRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := receivers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ReceiverInitializationFunc)
	return init_func(ctx, uri)
}

func RegisterReceiver(ctx context.Context, scheme string, init_func ReceiverInitializationFunc) error {

	err := ensureReceiverRoster()

	if err != nil {
		return err
	}

	return receivers.Register(ctx, scheme, init_func)
}

func ensureReceiverRoster() error {

	if receivers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		receivers = r
	}

	return nil
}

func Receivers() []string {
	ctx := context.Background()
	return receivers.Drivers(ctx)
}
