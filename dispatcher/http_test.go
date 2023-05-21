package dispatcher

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestNewHTTPDispatcher(t *testing.T) {

	ctx := context.Background()

	d, err := NewDispatcher(ctx, "http://")

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	err2 := d.Dispatch(ctx, []byte("hello world"))

	if err2 != nil {
		t.Fatalf("Failed to dispatch message, %v", err2)
	}
}

func TestNewHTTPDispatcherWithLogger(t *testing.T) {

	ctx := context.Background()

	var buf bytes.Buffer

	logger := log.New(&buf, "testing ", log.Lshortfile)

	d, err := NewLogDispatcherWithLogger(ctx, logger)
	fmt.Print("http run")

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
