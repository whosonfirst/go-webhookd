// webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.
package main

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"github.com/whosonfirst/go-webhookd/v3/daemon"
	"log"
	"os"
)

func main() {

	fs := flagset.NewFlagSet("webhooks")

	config_uri := fs.String("config-uri", "", "A valid Go Cloud runtimevar URI representing your webhookd config.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "WEBHOOKD", false)

	if err != nil {
		log.Fatalf("Failed to set flags from env vars, %v", err)
	}

	ctx := context.Background()

	cfg, err := config.NewConfigFromURI(ctx, *config_uri)

	if err != nil {
		log.Fatalf("Failed to load config %s, %v", *config_uri, err)
	}

	wh_daemon, err := daemon.NewWebhookDaemonFromConfig(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to create webhook daemon, %v", err)
	}

	err = wh_daemon.Start(ctx)

	if err != nil {
		log.Fatalf("Failed to start webhook daemon, %v", err)
	}

	os.Exit(0)
}
