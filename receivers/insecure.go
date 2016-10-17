package receivers

import (
	"github.com/whosonfirst/go-webhookd"
	"io/ioutil"
	"net/http"
)

type InsecureReceiver struct {
	webhookd.WebhookReceiver
}

func NewInsecureReceiver() (InsecureReceiver, error) {

	wh := InsecureReceiver{}
	return wh, nil
}

func (wh InsecureReceiver) Receive(req *http.Request) ([]byte, *webhookd.WebhookError) {

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
