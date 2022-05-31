package receiver

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
	"io"
	"net/http"
)

func init() {

	ctx := context.Background()
	err := RegisterReceiver(ctx, "insecure", NewInsecureReceiver)

	if err != nil {
		panic(err)
	}
}

// LogReceiver implements the `webhookd.WebhookReceiver` interface for receiving webhook messages in an insecure fashion.
type InsecureReceiver struct {
	webhookd.WebhookReceiver
}

// NewInsecureReceiver returns a new `InsecureReceiver` instance configured by 'uri' in the form of:
//
// 	insecure://
func NewInsecureReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	wh := InsecureReceiver{}
	return wh, nil
}

// Receive returns the body of the message in 'req'. It does not check its provenance or validate the message body in any way. You should not use this in production.
func (wh InsecureReceiver) Receive(ctx context.Context, req *http.Request) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	if req.Method != "POST" {

		code := http.StatusMethodNotAllowed
		message := "Method not allowed"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {

		code := http.StatusInternalServerError
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	return body, nil
}
