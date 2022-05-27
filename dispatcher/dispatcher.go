// package dispatcher provides methods for relaying webhook messages after they have been transformed.
package dispatcher

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
	"sort"
)

// dispatcher is a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookDispatcher` initialization functions.
var dispatchers roster.Roster

// DispatcherInitializationFunc is a function used to initialize an implementation of the `webhookd.WebhookDispatcher` interface.
type DispatcherInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error)

// NewDispatcher() returns a new `webhookd.WebhookDispatcher` instance derived from 'uri'. The semantics of and requirements for
// 'uri' as specific to the package implementing the interface.
func NewDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	err := ensureDispatcherRoster()

	if err != nil {
		return nil, fmt.Errorf("Failed to ensure dispatcher roster, %w", err)
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := parsed.Scheme

	i, err := dispatchers.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to find initialization function for '%s', %w", scheme, err)
	}

	init_func := i.(DispatcherInitializationFunc)
	return init_func(ctx, uri)
}

// RegisterDispatcher() associates 'scheme' with 'init_func' in an internal list of avilable `webhookd.WebhookDispatcher` implementations.
func RegisterDispatcher(ctx context.Context, scheme string, init_func DispatcherInitializationFunc) error {

	err := ensureDispatcherRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure dispatcher roster, %w", err)
	}

	return dispatchers.Register(ctx, scheme, init_func)
}

// ensureDispatcherRoster() ensures that a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookDispatcher`
// initialization functions is present
func ensureDispatcherRoster() error {

	if dispatchers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster, %w", err)
		}

		dispatchers = r
	}

	return nil
}

// Dispatchers() returns the list of schemes that have been "registered".
// This method is deprecated and you should use `Schemes()` instead.
func Dispatchers() []string {
	return Schemes()
}

// Schemes() returns the list of schemes that have been "registered".
func Schemes() []string {
	ctx := context.Background()
	drivers := dispatchers.Drivers(ctx)

	schemes := make([]string, len(drivers))

	for idx, dr := range drivers {
		schemes[idx] = fmt.Sprintf("%s://", dr)
	}

	sort.Strings(schemes)
	return schemes
}
