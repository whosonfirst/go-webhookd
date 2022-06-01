package webhookd

import (
	"fmt"
)

// WebhookError implements the `error` interface for wrapping webhookd error codes and messages.
type WebhookError struct {
	error
	// A numeric status code identifying the error.
	Code int
	// A long-form string describing the error.
	Message string
}

// Error() returns a string containing both the status code and descriptive message associated with an error.
func (e WebhookError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}
