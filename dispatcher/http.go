package dispatcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/whosonfirst/go-webhookd/v3"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "http", NewHTTPDispatcher)
	if err != nil {
		panic(err)
	}

	err = RegisterDispatcher(ctx, "https", NewHTTPDispatcher)
	if err != nil {
		panic(err)
	}
}

// GET and POST are the only supported HTTP methods
const GET = "GET"
const POST = "POST"

// HTTPDispatcher implements the `webhookd.WebhookDispatcher` interface for dispatching messages to a `log.Logger` instance `http.get` or `http.post`.
type HTTPDispatcher struct {
	webhookd.WebhookDispatcher
	// logger is the `log.Logger` instance associated with the dispatcher.
	logger *log.Logger
	// url to send the message to
	url url.URL
	// method to use when sending the message
	method string
	// client to use when sending the message
	client HTTPClient
}

// HTTPClient is an interface for `http.Client` to allow for mocking in tests.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// HTTPDispatcherOptions is a struct containing the options for `NewHTTPDispatcherWithOptions`.
type HTTPDispatcherOptions struct {
	// Logger is the `log.Logger` instance associated with the dispatcher.
	Logger *log.Logger
	// URL to send the message to
	URL url.URL
	// Client to use when sending the message
	Client HTTPClient
}

// NewHTTPDispatcher returns a new `HTTPDispatcher` instance configured by 'uri' in the form of:
//
//	http://
//	https://
//
// Messages are dispatched to the default `log.Default()` instance along with the uri parsed.
func NewHTTPDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {
	logger := log.Default()

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	opts := HTTPDispatcherOptions{
		Logger: logger,
		URL:    *parsed,
		Client: &http.Client{},
	}

	return NewHTTPDispatcherWithOptions(ctx, &opts)
}

// NewHTTPDispatcher returns a new `HTTPDispatcher` instance that dispatches messages to `http.Get` or `http.Post`.
func NewHTTPDispatcherWithOptions(ctx context.Context, opts *HTTPDispatcherOptions) (webhookd.WebhookDispatcher, error) {

	opts.Logger.Print("Parsed dispatcher URL: ", opts.URL.String())

	// check the method and default to POST
	method := opts.URL.Query().Get("method")
	if method != GET {
		method = POST
	}

	d := HTTPDispatcher{
		logger: opts.Logger,
		url:    opts.URL,
		method: method,
		client: opts.Client,
	}

	return &d, nil
}

// Dispatch sends 'body' to the `log.Logger` and `http.Get`/`http.Post` that 'd' has been instantiated with.
func (d *HTTPDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {
	var resp *http.Response
	var err error

	if d.method == GET {
		d.logger.Println("Dispatching GET:", d.url.String(), "not forwarding body: ", string(body))
		resp, err = d.client.Get(d.url.String())
	} else {
		d.logger.Println("Dispatching POST:", d.url.String(), "forwarding body: ", string(body))
		resp, err = d.client.Post(d.url.String(), "application/json", bytes.NewBuffer(body))
	}

	// if we get a nil response the destination is unreachable
	if resp == nil {
		code := http.StatusRequestTimeout
		message := "Timeout likely destination unreachable"
		whErr := &webhookd.WebhookError{Code: code, Message: message}
		return whErr
	}

	defer resp.Body.Close()

	// if we get any other status code than 200
	if resp.StatusCode != http.StatusOK {
		code := resp.StatusCode
		message := fmt.Sprintf("Failed to dispatch message: %s", resp.Status)
		whErr := &webhookd.WebhookError{Code: code, Message: message}
		return whErr
	}

	if err != nil {
		code := http.StatusInternalServerError
		message := err.Error()
		whErr := &webhookd.WebhookError{Code: code, Message: message}
		return whErr
	}

	return nil
}
