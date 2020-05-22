package server

// https://github.com/FiloSottile/mkcert/issues/154
// https://smallstep.com/blog/private-acme-server/

import (
	"context"
	"fmt"
	_ "log"
	"net/http"
	"net/url"
	"os"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "https", NewTLSServer)
	RegisterServer(ctx, "tls", NewTLSServer)
}

type TLSServer struct {
	Server
	url  *url.URL
	cert string
	key  string
}

func NewTLSServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	tls_cert := q.Get("cert")
	tls_key := q.Get("key")

	_, err = os.Stat(tls_cert)

	if err != nil {
		return nil, err
	}

	_, err = os.Stat(tls_key)

	if err != nil {
		return nil, err
	}

	server_uri := fmt.Sprintf("https://%s", u.Host)
	server_u, err := url.Parse(server_uri)

	if err != nil {
		return nil, err
	}

	server := TLSServer{
		url:  server_u,
		cert: tls_cert,
		key:  tls_key,
	}

	return &server, nil
}

func (s *TLSServer) Address() string {
	return s.url.String()
}

func (s *TLSServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {
	return http.ListenAndServeTLS(s.url.Host, s.cert, s.key, mux)
}
