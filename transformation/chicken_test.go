package transformation

import (
	"bytes"
	"context"
	"testing"
)

func TestChickenTransformation(t *testing.T) {

	ctx := context.Background()

	input := []byte("hello world")
	expected := []byte("🐔 🐔")

	tr, err := NewTransformation(ctx, "chicken://zxx?clucking=false")

	if err != nil {
		t.Fatalf("Failed to create new chicken transformation, %v", err)
	}

	output, err2 := tr.Transform(ctx, input)

	if err2 != nil {
		t.Fatalf("Failed to transform body, %v", err2)
	}

	if !bytes.Equal(output, expected) {
		t.Fatalf("Unexpected output '%s'", string(output))
	}

}
