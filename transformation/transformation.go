package transformation

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
)

var transformations roster.Roster

type TransformationInitializationFunc func(ctx context.Context, uri string) (webhookd.WebhookTransformation, error)

func NewTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	err := ensureTransformationRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := transformations.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(TransformationInitializationFunc)
	return init_func(ctx, uri)
}

func RegisterTransformation(ctx context.Context, scheme string, init_func TransformationInitializationFunc) error {

	err := ensureTransformationRoster()

	if err != nil {
		return err
	}

	return transformations.Register(ctx, scheme, init_func)
}

func ensureTransformationRoster() error {

	if transformations == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		transformations = r
	}

	return nil
}

func Transformations() []string {
	ctx := context.Background()
	return transformations.Drivers(ctx)
}
