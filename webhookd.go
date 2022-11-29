// package webhookd implements a bucket-brigrade style webhook server where requests are relayed through a receiver, one or more transformations and one or more dispatchers each of which have interfaces and are defined using a URI-based syntax to allow for custom processing.
package webhookd

import (
	"context"
	"net/http"
)

// type WebhookHandler is an interface for definining and configuring an individual webhooks.
type WebhookHandler interface {
	// Endpoint is the relative URI of the webhook.
	Endpoint() string // sudo make me a net.URL or something
	// Receiver() is the `WebhookReceiver` instance used to process a webhook message on arrival.
	Receiver() WebhookReceiver
	// Transformations() is a list of zero or more `WebhookTransformation` instances that will be applied to a message after receipt.
	Transformations() []WebhookTransformation
	// Dispatchers() is a list of zero or more `WebhookDispatcher` instances which will be to relay the body of a webhook message after it's been transformed.
	Dispatchers() []WebhookDispatcher
}

// WebhookReceiver is an interface that defines methods for processing a webhook message on arrival.
type WebhookReceiver interface {
	// Receive() process the body of an `http.Request` instance (according to rules defined by the package implementing the `WebhookReceiver` interface)..
	Receive(context.Context, *http.Request) ([]byte, *WebhookError)
}

// WebhookTransformation is an interface that defines methods for altering (transforming) the body of a (webhook) message after receipt.
type WebhookTransformation interface {
	// Transforms() alters the body of a (webhook) message (according to rules defined by the package implementing the `WebhookTransformation` interface).
	Transform(context.Context, []byte) ([]byte, *WebhookError)
}

// WebhookDispatcher is an interface that defines methods for relaying the body of a (webhook) message after it has been transformed.
type WebhookDispatcher interface {
	// Dispatch() relays the body of a message (according to rules defined defined by the package implementing the `WebhookDispatcher` interface).
	Dispatch(context.Context, []byte) *WebhookError
}
