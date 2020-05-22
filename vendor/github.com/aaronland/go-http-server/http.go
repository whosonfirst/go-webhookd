package server

import (
	"context"
	_ "log"
	"net/http"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "http", NewHTTPServer)
}

type HTTPServer struct {
	Server
	url *url.URL
}

func NewHTTPServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	u.Scheme = "http"

	server := HTTPServer{
		url: u,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {
	return http.ListenAndServe(s.url.Host, mux)
}
