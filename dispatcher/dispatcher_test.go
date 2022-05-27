package dispatcher

import (
	"context"
	"testing"
)

func TestRegisterDispatcher(t *testing.T) {

	ctx := context.Background()

	err := RegisterDispatcher(ctx, "null", NewNullDispatcher)

	if err == nil {
		t.Fatalf("Expected NewNullDispatcher to be registered already")
	}
}

func TestNewDispatcher(t *testing.T) {

	ctx := context.Background()

	uri := "log://"

	_, err := NewDispatcher(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher for '%s', %v", uri, err)
	}
}
