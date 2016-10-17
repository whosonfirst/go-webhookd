package webhookd

import (
	"encoding/json"
	"errors"
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
	Daemon          WebhookDaemonConfig                    `json:"daemon"`
	Receivers       map[string]WebhookReceiverConfig       `json:"receivers"`
	Dispatchers     map[string]WebhookDispatcherConfig     `json:"dispatchers"`
	Transformations map[string]WebhookTransformationConfig `json:"transformations"`
	Webhooks        []WebhookWebhooksConfig                `json:"webhooks"`
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
	Config  string `json:"config,omitempty"`
}

type WebhookTransformationConfig struct {
	Name     string `json:"name"`
	Language string `json:"language,omitempty"`
	Clucking bool   `json:"clucking,omitempty"`
}

type WebhookWebhooksConfig struct {
	Endpoint        string   `json:"endpoint"`
	Receiver        string   `json:"receiver"`
	Transformations []string `json:"transformations"`
	Dispatchers     []string `json:"dispatchers"`
}

type WebhookReceiver interface {
	Receive(*http.Request) ([]byte, *WebhookError)
}

type WebhookDispatcher interface {
	Dispatch([]byte) *WebhookError
}

type WebhookTransformation interface {
	Transform([]byte) ([]byte, *WebhookError)
}

type WebhookHandler interface {
	Endpoint() string // sudo make me a net.URI or something
	Receiver() WebhookReceiver
	Transformations() []WebhookTransformation
	Dispatchers() []WebhookDispatcher
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

func (c *WebhookConfig) GetReceiverConfigByName(name string) (*WebhookReceiverConfig, error) {

	config, ok := c.Receivers[name]

	if !ok {
		return nil, errors.New("Invalid receiver name")
	}

	return &config, nil
}

func (c *WebhookConfig) GetDispatcherConfigByName(name string) (*WebhookDispatcherConfig, error) {

	config, ok := c.Dispatchers[name]

	if !ok {
		return nil, errors.New("Invalid dispatcher name")
	}

	return &config, nil
}

func (c *WebhookConfig) GetTransformationConfigByName(name string) (*WebhookTransformationConfig, error) {

	config, ok := c.Transformations[name]

	if !ok {
		return nil, errors.New("Invalid transformations name")
	}

	return &config, nil
}
