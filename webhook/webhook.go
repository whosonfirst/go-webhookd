package webhook

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
)

type Webhook struct {
	webhookd.WebhookHandler
	endpoint        string
	receiver        webhookd.WebhookReceiver
	transformations []webhookd.WebhookTransformation
	dispatchers     []webhookd.WebhookDispatcher
}

func NewWebhook(ctx context.Context, endpoint string, rc webhookd.WebhookReceiver, tr []webhookd.WebhookTransformation, ds []webhookd.WebhookDispatcher) (Webhook, error) {

	wh := Webhook{
		endpoint:        endpoint,
		receiver:        rc,
		transformations: tr,
		dispatchers:     ds,
	}

	return wh, nil
}

func (wh Webhook) Endpoint() string {
	return wh.endpoint
}

func (wh Webhook) Receiver() webhookd.WebhookReceiver {
	return wh.receiver
}

func (wh Webhook) Dispatchers() []webhookd.WebhookDispatcher {
	return wh.dispatchers
}

func (wh Webhook) Transformations() []webhookd.WebhookTransformation {
	return wh.transformations
}
