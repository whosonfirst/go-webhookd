package dispatchers

import (
	"bytes"
	"log"
	"net/http"

	"github.com/whosonfirst/go-webhookd"
)

type ForwarderDispatcher struct {
	webhookd.WebhookDispatcher
	// writer *slackcat.Writer
	destinationURL string
}

func NewForwarderDispatcher(destinationURL string) (*ForwarderDispatcher, error) {

	fdr := ForwarderDispatcher{
		destinationURL: destinationURL,
	}

	return &fdr, nil
}

func (fdr *ForwarderDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	req, err := http.NewRequest("POST", fdr.destinationURL, bytes.NewReader(body))
	req.Header.Set("X-Webhookd-Dispatcher-Type", "forwarder")
	// req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		code := 999
		message := err.Error()

		log.Println("err when forwarding request: ", message)

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}
	defer resp.Body.Close()

	log.Println("forwarded request to detination url: ", fdr.destinationURL)
	// log.Println("response Status:", resp.Status)
	// log.Println("response Headers:", resp.Header)
	// body, _ = ioutil.ReadAll(resp.Body)
	// log.Println("response Body:", string(body))

	return nil
}
