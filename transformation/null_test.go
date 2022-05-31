package transformation

import (
	"bytes"
	"context"
	"testing"
)

func TestNullTransformation(t *testing.T) {

	ctx := context.Background()

	input := []byte("hello world")
	expected := input

	tr, err := NewTransformation(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new null transformation, %v", err)
	}

	output, err2 := tr.Transform(ctx, input)

	if err2 != nil {
		t.Fatalf("Failed to transform body, %v", err2)
	}

	if !bytes.Equal(output, expected) {
		t.Fatalf("Unexpected output '%s'", string(output))
	}

}
