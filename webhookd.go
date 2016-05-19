package webhookd

import (
	"net/http"
)

type Receiver interface {
	Receive(rsp http.ResponseWriter, req *http.Request)
}

type Dispatcher interface {
	Dispatch(body []byte) // sudo make an io.Reader or something
}

type Webhook interface {
	Receiver
	Dispatcher
}

type WebhookDaemon struct {
	Endpoint string
	Webhook  Webhook
}

func NewWebhookDaemon(endpoint string, webhook Webhook) (WebhookDaemon, error) {

	d := WebhookDaemon{
		Endpoint: endpoint,
		Webhook:  webhook,
	}

	return d, nil
}

func (d *WebhookDaemon) Start() error {

	handler := func(rsp http.ResponseWriter, req *http.Request) {
		d.Webhook.Receive(rsp, req)
	}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(d.Endpoint, nil)

	if err != nil {
		return err
	}

	return nil
}
