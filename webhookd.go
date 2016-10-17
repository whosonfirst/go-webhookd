package webhookd

import (
	"fmt"
	"net/http"
)

type WebhookError struct {
	Code    int
	Message string
}

func (e WebhookError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type WebhookReceiver interface {
	Receive(*http.Request) ([]byte, *WebhookError)
}

type WebhookDispatcher interface {
	Dispatch([]byte) *WebhookError
}

type WebhookTransformation interface {
	Transform([]byte) ([]byte, *WebhookError)
}

type WebhookHandler interface {
	Endpoint() string // sudo make me a net.URI or something
	Receiver() WebhookReceiver
	Transformations() []WebhookTransformation
	Dispatchers() []WebhookDispatcher
}
