package receivers

// This has not been fully tested with an actual GitHub message yet
// (20161016/thisisaaronland)

// https://developer.github.com/webhooks/
// https://developer.github.com/webhooks/#payloads
// https://developer.github.com/v3/activity/events/types/#pushevent
// https://developer.github.com/v3/repos/hooks/#ping-a-hook

import (
	"crypto/hmac"
	"encoding/json"
	gogithub "github.com/google/go-github/github"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-webhookd/github"
	"io/ioutil"
	_ "log"
	"net/http"
)

type GitHubReceiver struct {
	webhookd.WebhookReceiver
	secret string
	ref    string
}

func NewGitHubReceiver(secret string, ref string) (GitHubReceiver, error) {

	wh := GitHubReceiver{
		secret: secret,
		ref:    ref,
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

	if wh.ref != "" {

		var event gogithub.PushEvent

		err := json.Unmarshal(body, &event)

		if err != nil {
			err := &webhookd.WebhookError{Code: 999, Message: err.Error()}
			return nil, err
		}

		if wh.ref != *event.Ref {

			msg := "Invalid ref for commit"
			err := &webhookd.WebhookError{Code: 666, Message: msg}
			return nil, err
		}
	}

	/*

		So here's a thing that's not awesome: the event_type is passed in the header
		rather than anywhere in the payload body. So I don't know... maybe we need to
		change the signature of Receive method to be something like this:
		       { Payload: []byte, Extras: map[string]string }

		Which is not something that makes me "happy"... (20161016/thisisaaronland)

	*/

	return body, nil
}
