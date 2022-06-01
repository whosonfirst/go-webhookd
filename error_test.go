package webhookd

import (
	"testing"
)

func TestWebhookError(t *testing.T) {

	e := WebhookError{
		Code:    999,
		Message: "Self-destruct",
	}

	if e.Error() != "999 Self-destruct" {
		t.Fatalf("Unexpected error string: %s", e.Error())
	}
}
