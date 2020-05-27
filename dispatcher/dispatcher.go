package dispatcher

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
)

var dispatchers roster.Roster

type DispatcherInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error)

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
