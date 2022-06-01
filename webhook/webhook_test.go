package webhook

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
	"github.com/whosonfirst/go-webhookd/v3/dispatcher"
	"github.com/whosonfirst/go-webhookd/v3/receiver"
	"github.com/whosonfirst/go-webhookd/v3/transformation"
	"testing"
)

func TestWebhook(t *testing.T) {

	ctx := context.Background()

	r, err := receiver.NewReceiver(ctx, "insecure://")

	if err != nil {
		t.Fatalf("Failed to create new receiver, %v", err)
	}

	tr, err := transformation.NewTransformation(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new transformation, %v", err)
	}

	d, err := dispatcher.NewDispatcher(ctx, "log://")

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	endpoint := "/insecure"

	wh, err := NewWebhook(ctx, endpoint, r, []webhookd.WebhookTransformation{tr}, []webhookd.WebhookDispatcher{d})

	if err != nil {
		t.Fatalf("Failed to create new webhook, %v", err)
	}

	if wh.Endpoint() != endpoint {
		t.Fatalf("Invalid endpoint: %s", wh.Endpoint())
	}

	if wh.Receiver() != r {
		t.Fatalf("Unexpected receiver")
	}

}
