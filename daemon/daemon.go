package daemon

import (
	"errors"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/config"
	"github.com/whosonfirst/go-webhookd/dispatchers"
	"github.com/whosonfirst/go-webhookd/receivers"
	"github.com/whosonfirst/go-webhookd/transformations"
	"github.com/whosonfirst/go-webhookd/webhook"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebhookDaemon struct {
	host     string
	port     int
	webhooks map[string]webhookd.WebhookHandler
}

func NewWebhookDaemonFromConfig(cfg *config.WebhookConfig) (*WebhookDaemon, error) {

	d, err := NewWebhookDaemon(cfg.Daemon.Host, cfg.Daemon.Port)

	if err != nil {
		return nil, err
	}

	err = d.AddWebhooksFromConfig(cfg)

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

func (d *WebhookDaemon) AddWebhooksFromConfig(cfg *config.WebhookConfig) error {

	if len(cfg.Webhooks) == 0 {
		return errors.New("No webhooks defined")
	}

	for i, hook := range cfg.Webhooks {

		if hook.Endpoint == "" {
			msg := fmt.Sprintf("Missing endpoint at offset %d", i+1)
			return errors.New(msg)
		}

		if hook.Receiver == "" {
			msg := fmt.Sprintf("Missing receiver at offset %d", i+1)
			return errors.New(msg)
		}

		if len(hook.Dispatchers) == 0 {
			msg := fmt.Sprintf("Missing dispatchers at offset %d", i+1)
			return errors.New(msg)
		}

		receiver_config, err := cfg.GetReceiverConfigByName(hook.Receiver)

		if err != nil {
			return err
		}

		receiver, err := receivers.NewReceiverFromConfig(receiver_config)

		if err != nil {
			return err
		}

		var steps []webhookd.WebhookTransformation

		for _, name := range hook.Transformations {

			transformation_config, err := cfg.GetTransformationConfigByName(name)

			if err != nil {
				return err
			}

			step, err := transformations.NewTransformationFromConfig(transformation_config)

			if err != nil {
				return err
			}

			steps = append(steps, step)
		}

		var sendto []webhookd.WebhookDispatcher

		for _, name := range hook.Dispatchers {

			dispatcher_config, err := cfg.GetDispatcherConfigByName(name)

			if err != nil {
				return err
			}

			dispatcher, err := dispatchers.NewDispatcherFromConfig(dispatcher_config)

			if err != nil {
				return err
			}

			sendto = append(sendto, dispatcher)
		}

		wh, err := webhook.NewWebhook(hook.Endpoint, receiver, steps, sendto)

		if err != nil {
			return err
		}

		err = d.AddWebhook(wh)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *WebhookDaemon) AddWebhook(wh webhook.Webhook) error {

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

		t1 := time.Now()

		var ta time.Time
		var tb time.Duration

		var ttr time.Duration // time to receive
		var ttt time.Duration // time to transform
		var ttd time.Duration // time to dispatch

		ta = time.Now()

		rcvr := wh.Receiver()

		body, err := rcvr.Receive(req)

		if err != nil {
			http.Error(rsp, err.Error(), err.Code)
			return
		}

		tb = time.Since(ta)

		ttr = tb

		ta = time.Now()

		for _, step := range wh.Transformations() {

			body, err = step.Transform(body)

			if err != nil {
				http.Error(rsp, err.Error(), err.Code)
				return
			}

			// check to see if there is anything left the transformation
			// https://github.com/whosonfirst/go-webhookd/issues/7
		}

		tb = time.Since(ta)
		ttt = tb

		// check to see if there is anything to dispatch
		// https://github.com/whosonfirst/go-webhookd/issues/7

		ta = time.Now()

		wg := new(sync.WaitGroup)
		ch := make(chan *webhookd.WebhookError)

		for _, d := range wh.Dispatchers() {

			wg.Add(1)

			go func(d webhookd.WebhookDispatcher, body []byte) {

				defer wg.Done()

				err = d.Dispatch(body)

				if err != nil {
					ch <- err
				}

				// err = &webhookd.WebhookError{Code: 000, Message: "o_O"}
				// ch <- err

			}(d, body)
		}

		errors := make([]string, 0)

		go func() {

			for e := range ch {
				errors = append(errors, e.Error())
			}
		}()

		wg.Wait()

		if len(errors) > 0 {

			msg := strings.Join(errors, "\n\n")
			http.Error(rsp, msg, http.StatusInternalServerError)
			return
		}

		tb = time.Since(ta)
		ttd = tb

		t2 := time.Since(t1)

		rsp.Header().Set("X-Webhookd-Time-To-Receive", fmt.Sprintf("%v", ttr))
		rsp.Header().Set("X-Webhookd-Time-To-Transform", fmt.Sprintf("%v", ttt))
		rsp.Header().Set("X-Webhookd-Time-To-Dispatch", fmt.Sprintf("%v", ttd))
		rsp.Header().Set("X-Webhookd-Time-To-Process", fmt.Sprintf("%v", t2))

		query := req.URL.Query()

		if query.Get("debug") != "" {
			rsp.Header().Set("Content-Type", "text/plain")
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
			rsp.Write(body)
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
	log.Printf("webhookd listening for requests on %s\n", addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = gracehttp.Serve(&http.Server{Addr: addr, Handler: mux})

	if err != nil {
		return err
	}

	return nil
}
