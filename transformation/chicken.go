package transformation

import (
	"context"
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

type ChickenTransformation struct {
	webhookd.WebhookTransformation
	chicken *chicken.Chicken
}

func NewChickenTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	lang := u.Host
	str_clucking := q.Get("clucking")

	clucking := false

	if str_clucking != "" {

		v, err := strconv.ParseBool(str_clucking)

		if err != nil {
			return nil, err
		}

		clucking = v
	}

	ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

	if err != nil {
		return nil, err
	}

	tr := ChickenTransformation{
		chicken: ch,
	}

	return &tr, nil
}

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
