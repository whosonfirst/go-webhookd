package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const example_config string = "../docs/config/config.json.example"

func newConfigFromURI() (*WebhookConfig, error) {

	ctx := context.Background()

	path_config, err := filepath.Abs(example_config)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive absolute path for %s, %v", example_config, err)
	}

	uri := fmt.Sprintf("file://%s?decoder=string", path_config)

	cfg, err := NewConfigFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new config for %s, %v", path_config, err)
	}

	return cfg, nil
}

func TestNewConfigFromURI(t *testing.T) {

	cfg, err := newConfigFromURI()

	if err != nil {
		t.Fatalf("Failed to create new config from URI, %v", err)
	}

	if cfg.Daemon != "http://localhost:8080" {
		t.Fatalf("Unexpected cfg.Daemon value: %s", cfg.Daemon)
	}
}

func TestNewConfigFromReader(t *testing.T) {

	ctx := context.Background()

	path_config, err := filepath.Abs(example_config)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", example_config, err)
	}

	r, err := os.Open(path_config)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path_config, err)
	}

	defer r.Close()

	cfg, err := NewConfigFromReader(ctx, r)

	if err != nil {
		t.Fatalf("Failed to create new config for %s, %v", path_config, err)
	}

	if cfg.Daemon != "http://localhost:8080" {
		t.Fatalf("Unexpected cfg.Daemon value: %s", cfg.Daemon)
	}
}

func TestGetReceiverConfigByName(t *testing.T) {

	cfg, err := newConfigFromURI()

	if err != nil {
		t.Fatalf("Failed to create new config from URI, %v", err)
	}

	name := "github"
	expected := "github://?secret=s33kret&ref=refs/heads/main"

	uri, err := cfg.GetReceiverConfigByName(name)

	if err != nil {
		t.Fatalf("Failed to get receiver config for %s, %v", name, err)
	}

	if uri != expected {
		t.Fatalf("Unexpected value for %s receiver: %s", name, uri)
	}

}

func TestGetTransformationConfigByName(t *testing.T) {

	cfg, err := newConfigFromURI()

	if err != nil {
		t.Fatalf("Failed to create new config from URI, %v", err)
	}

	name := "chicken"
	expected := "chicken://zxx"

	uri, err := cfg.GetTransformationConfigByName(name)

	if err != nil {
		t.Fatalf("Failed to get transformation config for %s, %v", name, err)
	}

	if uri != expected {
		t.Fatalf("Unexpected value for %s transformation: %s", name, uri)
	}

	t.Skip()
}

func TestGetDispatcherConfigByName(t *testing.T) {

	cfg, err := newConfigFromURI()

	if err != nil {
		t.Fatalf("Failed to create new config from URI, %v", err)
	}

	name := "log"
	expected := "log://"

	uri, err := cfg.GetDispatcherConfigByName(name)

	if err != nil {
		t.Fatalf("Failed to get dispatcher config for %s, %v", name, err)
	}

	if uri != expected {
		t.Fatalf("Unexpected value for %s dispatcher: %s", name, uri)
	}
}
