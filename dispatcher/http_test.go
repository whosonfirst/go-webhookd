package dispatcher

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	Resp  *http.Response
	Error error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Resp, m.Error
}

func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return m.Resp, m.Error
}

func TestNewHTTPDispatcherWithOptions(t *testing.T) {

	ctx := context.Background()

	var buf bytes.Buffer

	logger := log.New(&buf, "testing ", log.Lshortfile)

	parsed, err := url.Parse("http://testing?method=GET")
	if err != nil {
		t.Fatalf("Failed to parse url, %v", err)
	}

	// Create a mock response
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("Mock response body")),
	}

	// Create a mock HTTP client with the desired behavior
	mockClient := &MockHTTPClient{
		Resp:  mockResponse,
		Error: nil,
	}

	d, err := NewHTTPDispatcherWithOptions(ctx, &HTTPDispatcherOptions{logger, *parsed, mockClient})

	if err != nil {
		t.Fatalf("Failed to create new http dispatcher with logger, %v", err)
	}

	err2 := d.Dispatch(ctx, []byte("hello world"))

	if err2 != nil {
		t.Fatalf("Failed to dispatch message, %v", err2)
	}

	expected := "testing http.go:89: Parsed dispatcher URL: http://testing?method=GET\ntesting http.go:113: Dispatching GET: http://testing?method=GET not forwarding body:  hello world"
	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Fatalf("Unexpected output from custom writer: '%s'", output)
	}
}
