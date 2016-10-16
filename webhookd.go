package webhookd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebhookError struct {
	Code    int
	Message string
}

func (e WebhookError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type WebhookConfig struct {
	Daemon      WebhookDaemonConfig                `json:"daemon"`
	Receivers   map[string]WebhookReceiverConfig   `json:"receivers"`
	Dispatchers map[string]WebhookDispatcherConfig `json:"dispatchers"`
	Webhooks    []WebhookWebhooksConfig            `json:"webhooks"`
}

type WebhookDaemonConfig struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type WebhookReceiverConfig struct {
	Name   string `json:"name"`
	Secret string `json:"secret,omitempty"`
}

type WebhookDispatcherConfig struct {
	Name    string `json:"name"`
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
	Channel string `json:"channel,omitempty"`
}

type WebhookWebhooksConfig struct {
	Endpoint   string `json:"endpoint"`
	Dispatcher string `json:"dispatcher"`
	Receiver   string `json:"receiver"`
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

func NewConfigFromFile(file string) (*WebhookConfig, error) {

	body, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	c := WebhookConfig{}
	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
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
