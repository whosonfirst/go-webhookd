// Package config provides data structures and methods for configuring a `webhookd` instance.
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"gocloud.dev/runtimevar"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"io"
	_ "log"
	"strings"
)

// type WebhookConfig is a struct containing configuration information for a `webhookd` instance.
type WebhookConfig struct {
	// Daemon is a valid `aaronland/go-http-server` URI. This determines how the `webhookd` server will be
	// instantiated and listen for requests.
	Daemon string `json:"daemon"`
	// Receivers is a dictionary of available receivers where the key is a unique label used to identify the
	// receiver (in `WebhookWebhooksConfig`) and the value is a URI used to instantiate the reciever.
	Receivers map[string]string `json:"receivers"`
	// Dispatchers is a dictionary of available dispatchers where the key is a unique label used to identify the
	// dispatcher (in `WebhookWebhooksConfig`) and the value is a URI used to instantiate the dispatcher.
	Dispatchers map[string]string `json:"dispatchers"`
	// Transformations is a dictionary of available transformations where the key is a unique label used to identify the
	// transformation (in `WebhookWebhooksConfig`) and the value is a URI used to instantiate the transformation.
	Transformations map[string]string `json:"transformations"`
	// Webhooks is a list of `WebhookWebhooksConfig` used to configure the webhooks that a `webhookd` instance will respond to.
	Webhooks []WebhookWebhooksConfig `json:"webhooks"`
}

// type WebhookWebhooksConfig is a struct containing configuration information for an individual webhook.
type WebhookWebhooksConfig struct {
	// Endpoint is the relative URI where the webhook will be installed.
	Endpoint string `json:"endpoint"`
	// Receiver the label for a recievier configured in `WebhookConfig.Receivers` that will be used to process an
	// initial webhook request.
	Receiver string `json:"receiver"`
	// Transformations is a list of transformation labels configured in `WebhookConfig.Transformations`. These transformations
	// will be applied in the order they are listed. The first transformation will be applied to the output of `Receiver` and
	// subsequent transformations will be applied to the output of the previous transformation.
	Transformations []string `json:"transformations"`
	// Dispatchers is a list of dispatcher labels configured in `WebhookConfig.Dispatchers`. Each dispatcher takes the output
	// of the last transformation and relays ("dispatches") it acccording to its internal rules.
	Dispatchers []string `json:"dispatchers"`
}

// NewConfigFromURI returns a new `WebhookConfig` instance derived from 'uri' which is expected to take the form of
// a valid `gocloud.dev/runtimevar` URI. The value of that URI is expected to be a JSON-encoded `WebhookConfig` string.
func NewConfigFromURI(ctx context.Context, uri string) (*WebhookConfig, error) {

	v, err := runtimevar.OpenVariable(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open config URI, %w", err)
	}

	latest, err := v.Latest(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to determine latest value for config URI, %w", err)
	}

	str_cfg := latest.Value.(string)
	cfg_fh := strings.NewReader(str_cfg)

	return NewConfigFromReader(ctx, cfg_fh)
}

// NewConfigFromReader returns a new `WebhookConfig` instance derived from 'r'.The body of 'r' is expected to be a JSON-encoded `WebhookConfig`
// string.
func NewConfigFromReader(ctx context.Context, r io.Reader) (*WebhookConfig, error) {

	var cfg *WebhookConfig

	dec := json.NewDecoder(r)
	err := dec.Decode(&cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode config, %w", err)
	}

	return cfg, nil
}

// GetReceiverConfigByName returns the receiver URI for 'name'.
func (c *WebhookConfig) GetReceiverConfigByName(name string) (string, error) {

	config, ok := c.Receivers[name]

	if !ok {
		return "", fmt.Errorf("Invalid receiver name '%s'", name)
	}

	return config, nil
}

// GetDispatcherConfigByName returns the dispatcher URI for 'name'.
func (c *WebhookConfig) GetDispatcherConfigByName(name string) (string, error) {

	config, ok := c.Dispatchers[name]

	if !ok {
		return "", fmt.Errorf("Invalid dispatcher name '%s'", name)
	}

	return config, nil
}

// GetTransformationConfigByName returns the dispatcher URI for 'name'.
func (c *WebhookConfig) GetTransformationConfigByName(name string) (string, error) {

	config, ok := c.Transformations[name]

	if !ok {
		return "", fmt.Errorf("Invalid transformations name '%s'", name)
	}

	return config, nil
}
