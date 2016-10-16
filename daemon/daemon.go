package daemon

import (
	"errors"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/dispatchers"
	"github.com/whosonfirst/go-webhookd/receivers"
	_ "log"
	"net/http"
)

type WebhookDaemon struct {
	host     string
	port     int
	webhooks map[string]webhookd.WebhookHandler
}

func NewWebhookDaemonFromConfig(config *webhookd.WebhookConfig) (*WebhookDaemon, error) {

	d, err := NewWebhookDaemon(config.Daemon.Host, config.Daemon.Port)

	if err != nil {
		return nil, err
	}

	err = d.AddWebhooksFromConfig(config)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func NewWebhookDaemon(host string, port int) (*WebhookDaemon, error) {

	webhooks := make(map[string]webhookd.WebhookHandler)

	d := WebhookDaemon{
		host:     host,
		port:     port,
		webhooks: webhooks,
	}

	return &d, nil
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

		var transformations []webhookd.WebhookTransformation

		webhook, err := webhookd.NewWebhook(hook.Endpoint, receiver, transformations, dispatcher)

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

func (d *WebhookDaemon) HandlerFunc() (http.HandlerFunc, error) {

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

	return http.HandlerFunc(handler), nil
}

func (d *WebhookDaemon) Start() error {

	handler, err := d.HandlerFunc()

	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", d.host, d.port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = gracehttp.Serve(&http.Server{Addr: addr, Handler: mux})

	if err != nil {
		return err
	}

	return nil
}
