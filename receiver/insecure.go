package receiver

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
	"io/ioutil"
	"net/http"
)

func init() {

	ctx := context.Background()
	err := RegisterReceiver(ctx, "insecure", NewInsecureReceiver)

	if err != nil {
		panic(err)
	}
}

type InsecureReceiver struct {
	webhookd.WebhookReceiver
}

func NewInsecureReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	wh := InsecureReceiver{}
	return wh, nil
}

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

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {

		code := http.StatusInternalServerError
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	return body, nil
}
