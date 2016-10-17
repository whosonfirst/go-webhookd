package receivers

// This has not been fully tested with an actual Slack message yet
// (20161016/thisisaaronland)

import (
	"github.com/whosonfirst/go-webhookd"
	"io/ioutil"
	"net/http"
)

type SlackReceiver struct {
	webhookd.WebhookReceiver
}

func NewSlackReceiver() (SlackReceiver, error) {

	slack := SlackReceiver{}
	return slack, nil
}

func (sl SlackReceiver) Receive(req *http.Request) ([]byte, *webhookd.WebhookError) {

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
