package service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"gopkg.in/redis.v1"
	"io/ioutil"
	"net/http"
)

type GitHubWebhook struct {
	secret string
	ps     *redis.Client
}

func NewGitHubWebhook(secret string, ps *redis.Client) (GitHubWebhook, error) {

	wh := GitHubWebhook{
		secret: secret,
		ps:     ps,
	}

	return wh, nil
}

func (wh GitHubWebhook) Receive(rsp http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Error(rsp, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventType := req.Header.Get("X-GitHub-Event")

	if eventType == "" {
		http.Error(rsp, "400 Bad Request - Missing X-GitHub-Event Header", http.StatusBadRequest)
		return
	}

	sig := req.Header.Get("X-Hub-Signature")

	if sig == "" {
		http.Error(rsp, "403 Forbidden - Missing X-Hub-Signature required for HMAC verification", http.StatusForbidden)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
		return
	}

	mac := hmac.New(sha1.New, []byte(wh.secret))
	mac.Write(body)

	expectedMAC := mac.Sum(nil)
	expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
		http.Error(rsp, "403 Forbidden - HMAC verification failed", http.StatusForbidden)
		return
	}

	wh.Dispatch(body)
}

func (wh GitHubWebhook) Dispatch(body []byte) {
	wh.ps.Publish("foo", string(body))
}
