package transformation

import (
	"context"
	"testing"
)

func TestRegisterTransformation(t *testing.T) {

	ctx := context.Background()

	err := RegisterTransformation(ctx, "chicken", NewChickenTransformation)

	if err == nil {
		t.Fatalf("Expected NewNullTransformation to be registered already")
	}
}

func TestNewTransformation(t *testing.T) {

	ctx := context.Background()

	uri := "chicken://zxx?clucking=false"

	_, err := NewTransformation(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new transformation for '%s', %v", uri, err)
	}
}
