package dispatchers

import (
	"github.com/whosonfirst/go-webhookd"
)

type ECSTaskDispatcher struct {
	webhookd.WebhookDispatcher
	Credentials    string
	Container      string
	Cluster        string
	Task           string
	LaunchType     string
	PublicIP       string
	Subnets        []string
	SecurityGroups []string
}

func NewECSTaskDispatcher() (*ECSTaskDispatcher, error) {

	d := ECSTaskDispatcher{
		Credentials:    "",
		Container:      "",
		Cluster:        "",
		Task:           "",
		LaunchType:     "",
		PublicIP:       "",
		Subnets:        []string{},
		SecurityGroups: []string{},
	}

	return &d, nil
}

func (d *ECSTaskDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	return nil
}
