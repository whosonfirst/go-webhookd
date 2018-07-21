package server

import (
	"github.com/whosonfirst/algnhsa"
	_ "log"
	"net/http"
	"net/url"
)

type LambdaServer struct {
	Server
	url *url.URL
}

func NewLambdaServer(u *url.URL, args ...interface{}) (Server, error) {

	server := LambdaServer{
		url: u,
	}

	return &server, nil
}

func (s *LambdaServer) Address() string {
	return s.url.String()
}

func (s *LambdaServer) ListenAndServe(mux *http.ServeMux) error {

	lambda_opts := new(algnhsa.Options)

	algnhsa.ListenAndServe(mux, lambda_opts)
	return nil
}
