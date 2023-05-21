package dispatcher

import (
	"context"
	"log"
	"net/http"

	"github.com/whosonfirst/go-webhookd/v3"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "http", NewHTTPDispatcher)

	if err != nil {
		panic(err)
	}
}

// HTTPDispatcher implements the `webhookd.WebhookDispatcher` interface for dispatching messages to a `log.Logger` instance.
type HTTPDispatcher struct {
	webhookd.WebhookDispatcher
	// logger is the `log.Logger` instance associated with the dispatcher.
	logger *log.Logger
	// url to send the message to
	url string
}

// NewHTTPDispatcher returns a new `HTTPDispatcher` instance configured by 'uri' in the form of:
//
//	http://
//
// Messasges are dispatched to the default `HTTP.Default()` instance.
func NewHTTPDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {
	logger := log.Default()
	return NewHTTPDispatcherWithLogger(ctx, logger)
}

// NewHTTPDispatcher returns a new `HTTPDispatcher` instance that dispatches messages to 'HTTPger'.
func NewHTTPDispatcherWithLogger(ctx context.Context, logger *log.Logger) (webhookd.WebhookDispatcher, error) {

	d := HTTPDispatcher{
		logger: logger,
		url:    ctx.Value(ctxUrl{}).(string),
	}

	return &d, nil
}

// Dispatch sends 'body' to the `log.Logger` that 'd' has been instantiated with.
func (d *HTTPDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	d.logger.Println("GET:", d.url, "forwarding body: ", string(body))

	resp, err := http.Get(d.url)
	if err != nil {
		d.logger.Println(err)
	}
	defer resp.Body.Close()

	return nil
}
