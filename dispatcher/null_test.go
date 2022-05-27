package dispatcher

import (
	"context"
	"testing"
)

func TestNullDispatcher(t *testing.T) {

	ctx := context.Background()

	d, err := NewDispatcher(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	err2 := d.Dispatch(ctx, []byte("hello world"))

	if err2 != nil {
		t.Fatalf("Failed to dispatch message, %v", err2)
	}

}
