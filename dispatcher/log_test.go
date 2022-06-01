package dispatcher

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func TestNewLogDispatcher(t *testing.T) {

	ctx := context.Background()

	d, err := NewDispatcher(ctx, "log://")

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	err2 := d.Dispatch(ctx, []byte("hello world"))

	if err2 != nil {
		t.Fatalf("Failed to dispatch message, %v", err2)
	}
}

func TestNewLogDispatcherWithLogger(t *testing.T) {

	ctx := context.Background()

	var buf bytes.Buffer

	logger := log.New(&buf, "testing ", log.Lshortfile)

	d, err := NewLogDispatcherWithLogger(ctx, logger)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher with logger, %v", err)
	}

	err2 := d.Dispatch(ctx, []byte("hello world"))

	if err2 != nil {
		t.Fatalf("Failed to dispatch message, %v", err2)
	}

	expected := "testing log.go:50: hello world"
	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Fatalf("Unexpected output from custom writer: '%s'", output)
	}
}
