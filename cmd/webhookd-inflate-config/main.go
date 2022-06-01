// webhookd-inflate-config is a command line tool for parsing one or more webhookd configuration files encoded as gocloud.dev/runtimevar "constant://" URIs and emitting their indented structure ("pretty-printing") to STDOUT.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"log"
	"net/url"
	"os"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "webhookd-inflate-config is a command line tool for parsing one or more webhookd configuration files encoded as gocloud.dev/runtimevar \"constant://val=\" URIs and emitting their indented structure (\"pretty-printing\") to STDOUT.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s string(N) string(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	for _, uri := range flag.Args() {

		u, err := url.Parse(uri)

		if err != nil {
			log.Fatalf("Failed to parse '%s', %v", uri, err)
		}

		q := u.Query()

		val := q.Get("val")

		if val == "" {
			log.Fatalf("Missing ?val parameter")
		}

		var cfg *config.WebhookConfig

		err = json.Unmarshal([]byte(val), &cfg)

		if err != nil {
			log.Fatalf("Failed to parse config, %v", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")

		err = enc.Encode(cfg)

		if err != nil {
			log.Fatalf("Failed to encode config, %v", err)
		}

	}
}
