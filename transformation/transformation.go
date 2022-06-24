// Package transformation provides an interface to altering (transforming) webhook messages.
package transformation

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
	"sort"
)

// transformation is a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookTransformation` initialization functions.
var transformations roster.Roster

// TransformationInitializationFunc is a function used to initialize an implementation of the `webhookd.WebhookTransformation` interface.
type TransformationInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookTransformation, error)

// NewTransformation() returns a new `webhookd.WebhookTransformation` instance derived from 'uri'. The semantics of and requirements for
// 'uri' as specific to the package implementing the interface.
func NewTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	err := ensureTransformationRoster()

	if err != nil {
		return nil, fmt.Errorf("Failed to ensure transformation roster, %w", err)
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := parsed.Scheme

	i, err := transformations.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to find initialization function for '%s', %w", scheme, err)
	}

	init_func := i.(TransformationInitializationFunc)
	return init_func(ctx, uri)
}

// RegisterTransformation() associates 'scheme' with 'init_func' in an internal list of avilable `webhookd.WebhookTransformation` implementations.
func RegisterTransformation(ctx context.Context, scheme string, init_func TransformationInitializationFunc) error {

	err := ensureTransformationRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure transformation roster, %w", err)
	}

	return transformations.Register(ctx, scheme, init_func)
}

// ensureTransformationRoster() ensures that a `aaronland/go-roster.Roster` instance used to maintain a list of registered `webhookd.WebhookTransformation`
// initialization functions is present
func ensureTransformationRoster() error {

	if transformations == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster, %w", err)
		}

		transformations = r
	}

	return nil
}

// Transformations() returns the list of schemes that have been "registered".
// This method is deprecated and you should use `Schemes()` instead.
func Transformations() []string {
	return Schemes()
}

// Schemes() returns the list of schemes that have been "registered".
func Schemes() []string {
	ctx := context.Background()
	drivers := transformations.Drivers(ctx)

	schemes := make([]string, len(drivers))

	for idx, dr := range drivers {
		schemes[idx] = fmt.Sprintf("%s://", dr)
	}

	sort.Strings(schemes)
	return schemes
}
