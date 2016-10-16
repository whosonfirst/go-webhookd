package transformations

import (
	"github.com/whosonfirst/go-webhookd"
)

type NullTransformation struct {
	webhookd.WebhookTransformation
}

func NewNullTransformation() (*NullTransformation, error) {

	p := NullTransformation{}
	return &p, nil
}

func (p *NullTransformation) Transform(body []byte) ([]byte, *webhookd.WebhookError) {
	return body, nil
}
