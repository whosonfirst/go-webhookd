package slackcat

import (
	slack "github.com/whosonfirst/slackcat"
	"strings"
)

type Writer struct {
	Config *slack.Config
}

func NewWriter(path string) (*Writer, error) {

	config, err := slack.ReadConfig(path)

	if err != nil {
		return nil, err
	}

	w := Writer{
		Config: config,
	}

	return &w, nil
}

func (w Writer) WriteString(s string) (n int64, err error) {
	r := strings.NewReader(s)
	return r.WriteTo(w)
}

func (w Writer) Write(p []byte) (int, error) {

	var text string
	text = string(p[:])

	msg := slack.SlackMsg{
		Channel:  w.Config.Channel,
		Username: w.Config.Username,
		Parse:    "full",
		Text:     text,
	}

	// sudo put me in a goroutine?

	err := msg.Post(w.Config.WebhookUrl)

	if err != nil {
		return 0, err
	}

	count := len(text)
	return count, nil
}

func (w Writer) Close() error {
	return nil
}
