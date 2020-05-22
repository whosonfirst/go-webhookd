package dispatchers

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
	"github.com/whosonfirst/go-writer-slackcat"
)

type SlackDispatcher struct {
	webhookd.WebhookDispatcher
	writer *slackcat.Writer
}

func NewSlackDispatcher(ctx context.Context, slackcat_config string) (*SlackDispatcher, error) {

	writer, err := slackcat.NewWriter(slackcat_config)

	if err != nil {
		return nil, err
	}

	slack := SlackDispatcher{
		writer: writer,
	}

	return &slack, nil
}

func (sl *SlackDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	_, err := sl.writer.Write(body)

	if err != nil {
		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}

	return nil
}
