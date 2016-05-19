package service

import (
	"gopkg.in/redis.v1"
	"io/ioutil"
	"net/http"
)

type InsecureWebhook struct {
	ps *redis.Client
}

func NewInsecureWebhook(ps *redis.Client) (InsecureWebhook, error) {

	wh := InsecureWebhook{
		ps: ps,
	}

	return wh, nil
}

func (wh InsecureWebhook) Receive(rsp http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
		return
	}

	wh.Dispatch(body)
}

func (wh InsecureWebhook) Dispatch(body []byte) {
	wh.ps.Publish("foo", string(body))
}
