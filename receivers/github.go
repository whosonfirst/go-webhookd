package receivers

// https://developer.github.com/webhooks/
// https://developer.github.com/webhooks/#payloads
// https://developer.github.com/v3/activity/events/types/#pushevent
// https://developer.github.com/v3/repos/hooks/#ping-a-hook

import (
	"crypto/hmac"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/github"
	"io/ioutil"
	_ "log"
	"net/http"
)

type GitHubReceiver struct {
	webhookd.WebhookReceiver
	secret string
}

func NewGitHubReceiver(secret string) (GitHubReceiver, error) {

	wh := GitHubReceiver{
		secret: secret,
	}

	return wh, nil
}

func (wh GitHubReceiver) Receive(req *http.Request) ([]byte, *webhookd.WebhookError) {

	if req.Method != "POST" {

		code := http.StatusMethodNotAllowed
		message := "Method not allowed"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	event_type := req.Header.Get("X-GitHub-Event")

	if event_type == "" {

		code := http.StatusBadRequest
		message := "Bad Request - Missing X-GitHub-Event Header"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	sig := req.Header.Get("X-Hub-Signature")

	if sig == "" {

		code := http.StatusForbidden
		message := "Missing X-Hub-Signature required for HMAC verification"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {

		code := http.StatusInternalServerError
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	expectedSig, _ := github.GenerateSignature(string(body), wh.secret)

	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {

		code := http.StatusForbidden
		message := "HMAC verification failed"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	return body, nil
}
