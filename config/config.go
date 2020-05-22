package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	_ "log"
	"path/filepath"
)

type WebhookConfig struct {
	Daemon          WebhookDaemonConfig                    `json:"daemon"`
	Receivers       map[string]WebhookReceiverConfig       `json:"receivers"`
	Dispatchers     map[string]WebhookDispatcherConfig     `json:"dispatchers"`
	Transformations map[string]WebhookTransformationConfig `json:"transformations"`
	Webhooks        []WebhookWebhooksConfig                `json:"webhooks"`
}

type WebhookDaemonConfig string

type WebhookReceiverConfig string

type WebhookDispatcherConfig string

type WebhookTransformationConfig string

type WebhookWebhooksConfig struct {
	Endpoint        string   `json:"endpoint"`
	Receiver        string   `json:"receiver"`
	Transformations []string `json:"transformations"`
	Dispatchers     []string `json:"dispatchers"`
}

func NewConfigFromFile(path string) (*WebhookConfig, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadFile(abs_path)

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
