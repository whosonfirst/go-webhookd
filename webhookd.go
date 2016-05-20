package webhookd

import (
	"errors"
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
	Dispatch([]byte) *WebhookError // sudo make me an io.Reader or something
}

type WebhookHandler interface {
	Endpoint() string // sudo make me a net.URI or something
	Receiver() WebhookReceiver
	Dispatcher() WebhookDispatcher
}

type Webhook struct {
	WebhookHandler
	endpoint   string
	receiver   WebhookReceiver
	dispatcher WebhookDispatcher
}

func NewWebhook(endpoint string, receiver WebhookReceiver, dispatcher WebhookDispatcher) (Webhook, error) {

	wh := Webhook{
		endpoint:   endpoint,
		receiver:   receiver,
		dispatcher: dispatcher,
	}

	return wh, nil
}

func (wh Webhook) Endpoint() string {
	return wh.endpoint
}

func (wh Webhook) Receiver() WebhookReceiver {
	return wh.receiver
}

func (wh Webhook) Dispatcher() WebhookDispatcher {
	return wh.dispatcher
}

type WebhookDaemon struct {
	host     string
	port     int
	webhooks map[string]WebhookHandler
}

func NewWebhookDaemon(host string, port int) (WebhookDaemon, error) {

	webhooks := make(map[string]WebhookHandler)

	d := WebhookDaemon{
		host:     host,
		port:     port,
		webhooks: webhooks,
	}

	return d, nil
}

func (d *WebhookDaemon) AddWebhook(wh Webhook) error {

	endpoint := wh.Endpoint()
	_, ok := d.webhooks[endpoint]

	if ok {
		return errors.New("endpoint already configured")
	}

	d.webhooks[endpoint] = wh
	return nil
}

func (d *WebhookDaemon) Start() error {

	if len(d.webhooks) == 0 {
		return errors.New("no webhooks configured")
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		endpoint := req.URL.Path

		wh, ok := d.webhooks[endpoint]

		if !ok {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		rcvr := wh.Receiver()
		dspt := wh.Dispatcher()

		body, err := rcvr.Receive(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		err = dspt.Dispatch(body)

		if err != nil {
			http.Error(rsp, err.Error(), 500)
			return
		}

		return
	}

	endpoint := fmt.Sprintf("%s:%d", d.host, d.port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(endpoint, nil)

	if err != nil {
		return err
	}

	return nil
}
