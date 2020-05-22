package server

import (
	"context"
	"github.com/whosonfirst/algnhsa"
	_ "log"
	"net/http"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "lambda", NewLambdaServer)
}

type LambdaServer struct {
	Server
	url *url.URL
}

func NewLambdaServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	server := LambdaServer{
		url: u,
	}

	return &server, nil
}

func (s *LambdaServer) Address() string {
	return s.url.String()
}

func (s *LambdaServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {

	// this cr^H^H^H stuff is important (20180713/thisisaaronland)
	// go-rasterzen/README.md#lambda-api-gateway-and-images#lambda-api-gateway-and-images

	lambda_opts := new(algnhsa.Options)
	lambda_opts.BinaryContentTypes = []string{"application/zip"}

	algnhsa.ListenAndServe(mux, lambda_opts)
	return nil
}
