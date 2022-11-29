package dispatcher

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
	"log"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "log", NewLogDispatcher)

	if err != nil {
		panic(err)
	}
}

// LogDispatcher implements the `webhookd.WebhookDispatcher` interface for dispatching messages to a `log.Logger` instance.
type LogDispatcher struct {
	webhookd.WebhookDispatcher
	// logger is the `log.Logger` instance associated with the dispatcher.
	logger *log.Logger
}

// NewLogDispatcher returns a new `LogDispatcher` instance configured by 'uri' in the form of:
//
//	log://
//
// Messasges are dispatched to the default `log.Default()` instance.
func NewLogDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	logger := log.Default()
	return NewLogDispatcherWithLogger(ctx, logger)
}

// NewLogDispatcher returns a new `LogDispatcher` instance that dispatches messages to 'logger'.
func NewLogDispatcherWithLogger(ctx context.Context, logger *log.Logger) (webhookd.WebhookDispatcher, error) {

	d := LogDispatcher{
		logger: logger,
	}

	return &d, nil
}

// Dispatch sends 'body' to the `log.Logger` that 'd' has been instantiated with.
func (d *LogDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	d.logger.Println(string(body))
	return nil
}
