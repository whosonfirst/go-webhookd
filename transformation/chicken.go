package transformation

import (
	"context"
	"fmt"
	"github.com/aaronland/go-chicken"
	"github.com/whosonfirst/go-webhookd/v3"
	"net/url"
	"strconv"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "chicken", NewChickenTransformation)

	if err != nil {
		panic(err)
	}
}

// ChickenTransformation implements the `webhookd.WebhookTransformation` interface for transforming messages using the `aaronland/go-chicken` package.
// the output of the `Transform` method is the same as its input.
type ChickenTransformation struct {
	webhookd.WebhookTransformation
	chicken *chicken.Chicken
}

// NewInsecureTransformation returns a new `ChickenTransformation` instance configured by 'uri' in the form of:
//
// 	chicken://{LANGUAGE_TAG}?{PARAMETERS}
//
// Where {LANGUAGE_TAG} is any valid language tag supported by the `aaronland/go-chicken` package. Valid {PARAMETERS} are:
// * `clucking={BOOLEAN}` A boolean flag to indicate whether messages should be transformed in the form of chicken noises.
func NewChickenTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	lang := u.Host
	str_clucking := q.Get("clucking")

	clucking := false

	if str_clucking != "" {

		v, err := strconv.ParseBool(str_clucking)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?clucking= parameter, %v", err)
		}

		clucking = v
	}

	ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

	if err != nil {
		return nil, fmt.Errorf("Failed to get chicken for language tag '%s', %v", lang, err)
	}

	tr := ChickenTransformation{
		chicken: ch,
	}

	return &tr, nil
}

// Transform returns 'body' translated in to "chicken".
func (tr *ChickenTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	txt := tr.chicken.TextToChicken(string(body))
	return []byte(txt), nil
}
