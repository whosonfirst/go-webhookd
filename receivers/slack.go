package receivers

// This has not been fully tested with an actual Slack message yet
// (20161016/thisisaaronland)

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
	"io/ioutil"
	"net/http"
)

func init() {

	ctx := context.Background()
	err := RegisterReceiver(ctx, "slack", NewSlackReceiver)

	if err != nil {
		panic(err)
	}
}

type SlackReceiver struct {
	webhookd.WebhookReceiver
}

func NewSlackReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	slack := SlackReceiver{}
	return slack, nil
}

func (sl SlackReceiver) Receive(ctx context.Context, req *http.Request) ([]byte, *webhookd.WebhookError) {

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
