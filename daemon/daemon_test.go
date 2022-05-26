package daemon

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

const example_config string = "../docs/config/config.json.example"

func TestNewWebhookDaemonFromConfig(t *testing.T) {

	ctx := context.Background()

	path_config, err := filepath.Abs(example_config)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", example_config, err)
	}

	uri := fmt.Sprintf("file://%s?decoder=string", path_config)

	cfg, err := config.NewConfigFromURI(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new config for %s, %v", path_config, err)
	}

	d, err := NewWebhookDaemonFromConfig(ctx, cfg)

	if err != nil {
		t.Fatalf("Failed to create new daemon from config, %v", err)
	}

	go func() {

		err := d.Start(ctx)

		if err != nil {
			log.Fatalf("Failed to start server, %v", err)
		}
	}()

	body := strings.NewReader("hello world")

	rsp, err := http.Post("http://localhost:8080/insecure-test", "text/plain", body)

	if err != nil {
		t.Fatalf("Failed to issue webhook request, %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected HTTP status: %s", rsp.Status)
	}
}
