package server

import (
	"errors"
	_ "log"
	"net/http"
	"net/url"
	"strings"
)

type Server interface {
	ListenAndServe(*http.ServeMux) error
	Address() string
}

func NewServerFromString(addr string, args ...interface{}) (Server, error) {

	u, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return NewServer(u)
}

func NewServer(addr *url.URL, args ...interface{}) (Server, error) {

	var svr Server
	var err error

	switch strings.ToUpper(addr.Scheme) {

	case "HTTP":

		svr, err = NewHTTPServer(addr, args...)

	case "LAMBDA":

		svr, err = NewLambdaServer(addr, args...)

	default:
		return nil, errors.New("Invalid server protocol")
	}

	if err != nil {
		return nil, err
	}

	return svr, nil
}
