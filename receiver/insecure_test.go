package receiver

import (
	"bytes"
	"context"
	"net/http"
	"testing"
)

func TestNullReceiver(t *testing.T) {

	ctx := context.Background()

	r, err := NewReceiver(ctx, "insecure://")

	if err != nil {
		t.Fatalf("Failed to create new receiver, %v", err)
	}

	expected := []byte("hello world")

	req, err := http.NewRequest("POST", "http://localhost:8080/insecure", bytes.NewReader(expected))

	if err != nil {
		t.Fatalf("Failed to create new request, %v", err)
	}

	body, err2 := r.Receive(ctx, req)

	if err2 != nil {
		t.Fatalf("Failed to receive message, %v", err)
	}

	if !bytes.Equal(body, expected) {
		t.Fatalf("Unexpected output '%s'", string(body))
	}
}
