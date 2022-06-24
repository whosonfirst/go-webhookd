// Package webhook provides data structures and methods for definining and configuring individual webhooks.
package webhook

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
)

// type Webhook defines a struct that implements the `webhookd.WebhookHandler` interface for definining and configuring an individual webhook.
type Webhook struct {
	webhookd.WebhookHandler
	// endpoint is the relative URI of the webhook.
	endpoint string
	// receiver is the `webhookd.WebhookReceiver` instance used to process a webhook message on arrival.
	receiver webhookd.WebhookReceiver
	// transformations is a list of zero or more `webhookd.WebhookTransformation` instances that will be applied to a message after receipt.
	transformations []webhookd.WebhookTransformation
	// dispatchers is a list of zero or more `webhookd.WebhookDispatcher` instances which will be to relay the body of a webhook message after it's been transformed.
	dispatchers []webhookd.WebhookDispatcher
}

// NewWebhook return a new `Wehook` instance.
func NewWebhook(ctx context.Context, endpoint string, rc webhookd.WebhookReceiver, tr []webhookd.WebhookTransformation, ds []webhookd.WebhookDispatcher) (Webhook, error) {

	wh := Webhook{
		endpoint:        endpoint,
		receiver:        rc,
		transformations: tr,
		dispatchers:     ds,
	}

	return wh, nil
}

// Endpoint() returns the relative URI	of the webhook.
func (wh Webhook) Endpoint() string {
	return wh.endpoint
}

// Receiver() returns the `webhookd.WebhookReceiver` instance used to process a webhook message on arrival.
func (wh Webhook) Receiver() webhookd.WebhookReceiver {
	return wh.receiver
}

// Transformations() returns the list of zero or more `webhookd.WebhookTransformation` instances that will be applied to a message after receipt.
func (wh Webhook) Transformations() []webhookd.WebhookTransformation {
	return wh.transformations
}

// Dispatchers() returns the zero of zero or more `webhookd.WebhookDispatcher` instances which will be to relay the body of a webhook message after it's been transformed.
func (wh Webhook) Dispatchers() []webhookd.WebhookDispatcher {
	return wh.dispatchers
}
