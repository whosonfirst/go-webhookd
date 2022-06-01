// package receiver provides an interface for the receipt of webhook messages.
package receiver

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
	"sort"
)

// receiver is a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookReceiver` initialization functions.
var receivers roster.Roster

// ReceiverInitializationFunc is a function used to initialize an implementation of the `webhookd.WebhookReceiver` interface.
type ReceiverInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookReceiver, error)

// NewReceiver() returns a new `webhookd.WebhookReceiver` instance derived from 'uri'. The semantics of and requirements for
// 'uri' as specific to the package implementing the interface.
func NewReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	err := ensureReceiverRoster()

	if err != nil {
		return nil, fmt.Errorf("Failed to ensure receiver roster, %w", err)
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := parsed.Scheme

	i, err := receivers.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to find initialization function for '%s', %w", scheme, err)
	}

	init_func := i.(ReceiverInitializationFunc)
	return init_func(ctx, uri)
}

// RegisterReceiver() associates 'scheme' with 'init_func' in an internal list of avilable `webhookd.WebhookReceiver` implementations.
func RegisterReceiver(ctx context.Context, scheme string, init_func ReceiverInitializationFunc) error {

	err := ensureReceiverRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure receiver roster, %w", err)
	}

	return receivers.Register(ctx, scheme, init_func)
}

// ensureReceiverRoster() ensures that a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookReceiver`
// initialization functions is present
func ensureReceiverRoster() error {

	if receivers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster, %w", err)
		}

		receivers = r
	}

	return nil
}

// Receivers() returns the list of schemes that have been "registered".
// This method is deprecated and you should use `Schemes()` instead.
func Receivers() []string {
	return Schemes()
}

// Schemes() returns the list of schemes that have been "registered".
func Schemes() []string {
	ctx := context.Background()
	drivers := receivers.Drivers(ctx)

	schemes := make([]string, len(drivers))

	for idx, dr := range drivers {
		schemes[idx] = fmt.Sprintf("%s://", dr)
	}

	sort.Strings(schemes)
	return schemes
}
