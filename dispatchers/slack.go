package dispatchers

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v2"
	"github.com/whosonfirst/go-writer-slackcat"
	"net/url"
)

type SlackDispatcher struct {
	webhookd.WebhookDispatcher
	writer *slackcat.Writer
}

func NewSlackDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	slackcat_config := u.Path

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
