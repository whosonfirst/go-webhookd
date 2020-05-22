package config

import (
	"context"
	"encoding/json"
	"errors"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	"io"
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

func NewConfigFromURI(ctx context.Context, uri string) (*WebhookConfig, error) {

	config_root := filepath.Dir(uri)
	config_name := filepath.Base(uri)

	bucket, err := blob.OpenBucket(ctx, config_root)

	if err != nil {
		return nil, err
	}

	cfg_fh, err := bucket.NewReader(ctx, config_name, nil)

	if err != nil {
		return nil, err
	}

	defer cfg_fh.Close()

	return NewConfigFromReader(ctx, cfg_fh)
}

func NewConfigFromReader(ctx context.Context, fh io.Reader) (*WebhookConfig, error) {

	var cfg *WebhookConfig

	dec := json.NewDecoder(fh)
	err := dec.Decode(&cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
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
