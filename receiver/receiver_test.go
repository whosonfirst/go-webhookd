package receiver

import (
	"context"
	"testing"
)

func TestRegisterReceiver(t *testing.T) {

	ctx := context.Background()

	err := RegisterReceiver(ctx, "insecure", NewInsecureReceiver)

	if err == nil {
		t.Fatalf("Expected NewNullReceiver to be registered already")
	}
}

func TestNewReceiver(t *testing.T) {

	ctx := context.Background()

	uri := "insecure://"

	_, err := NewReceiver(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new receiver for '%s', %v", uri, err)
	}
}
