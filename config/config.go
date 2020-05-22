package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	_ "log"
	"path/filepath"
)

type WebhookConfig struct {
	Daemon          string                  `json:"daemon"`
	Receivers       map[string]string       `json:"receivers"`
	Dispatchers     map[string]string       `json:"dispatchers"`
	Transformations map[string]string       `json:"transformations"`
	Webhooks        []WebhookWebhooksConfig `json:"webhooks"`
}

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

func (c *WebhookConfig) GetReceiverConfigByName(name string) (string, error) {

	config, ok := c.Receivers[name]

	if !ok {
		return "", errors.New("Invalid receiver name")
	}

	return config, nil
}

func (c *WebhookConfig) GetDispatcherConfigByName(name string) (string, error) {

	config, ok := c.Dispatchers[name]

	if !ok {
		return "", errors.New("Invalid dispatcher name")
	}

	return config, nil
}

func (c *WebhookConfig) GetTransformationConfigByName(name string) (string, error) {

	config, ok := c.Transformations[name]

	if !ok {
		return "", errors.New("Invalid transformations name")
	}

	return config, nil
}
