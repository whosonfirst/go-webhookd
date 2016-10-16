package dispatchers

import (
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-writer-slackcat"
)

type SlackDispatcher struct {
	webhookd.WebhookDispatcher
	writer *slackcat.Writer
}

func NewSlackDispatcher(slackcat_config string) (*SlackDispatcher, error) {

	writer, err := slackcat.NewWriter(slackcat_config)

	if err != nil {
		return nil, err
	}

	slack := SlackDispatcher{
		writer: writer,
	}

	return &slack, nil
}

func (sl *SlackDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	_, err := sl.writer.Write(body)

	if err != nil {
		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}

	return nil
}
