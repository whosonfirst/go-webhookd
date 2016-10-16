package daemon

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"
	"github.com/whosonfirst/go-whosonfirst-webhookd/receivers"
	_ "log"
	"net/http"
)

type WebhookDaemon struct {
	host     string
	port     int
	webhooks map[string]webhookd.WebhookHandler
}

func NewWebhookDaemon(host string, port int) (WebhookDaemon, error) {

	webhooks := make(map[string]webhookd.WebhookHandler)

	d := WebhookDaemon{
		host:     host,
		port:     port,
		webhooks: webhooks,
	}

	return d, nil
}

func (d *WebhookDaemon) AddWebhooksFromConfig(config *webhookd.WebhookConfig) error {

	if len(config.Webhooks) == 0 {
		return errors.New("No webhooks defined")
	}

	for i, hook := range config.Webhooks {

		if hook.Endpoint == "" {
			msg := fmt.Sprintf("Missing endpoint at offset %d", i+1)
			return errors.New(msg)
		}

		if hook.Receiver == "" {
			msg := fmt.Sprintf("Missing receiver at offset %d", i+1)
			return errors.New(msg)
		}

		if hook.Dispatcher == "" {
			msg := fmt.Sprintf("Missing dispatcher at offset %d", i+1)
			return errors.New(msg)
		}

		receiver_config, err := config.GetReceiverConfigByName(hook.Receiver)

		if err != nil {
			return err
		}

		receiver, err := receivers.NewReceiverFromConfig(receiver_config)

		if err != nil {
			return err
		}

		dispatcher_config, err := config.GetDispatcherConfigByName(hook.Dispatcher)

		if err != nil {
			return err
		}

		dispatcher, err := dispatchers.NewDispatcherFromConfig(dispatcher_config)

		if err != nil {
			return err
		}

		webhook, err := webhookd.NewWebhook(hook.Endpoint, receiver, dispatcher)

		if err != nil {
			return err
		}

		err = d.AddWebhook(webhook)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *WebhookDaemon) AddWebhook(wh webhookd.Webhook) error {

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
			http.Error(rsp, err.Error(), err.Code)
			return
		}

		err = dspt.Dispatch(body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
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
