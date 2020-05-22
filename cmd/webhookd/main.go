package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-webhookd/v2/config"
	"github.com/whosonfirst/go-webhookd/v2/daemon"
	"log"
	"os"
)

func main() {

	config_uri := flag.String("config-uri", "", "A valid Go Cloud blob URI where your webhookd config file lives")

	flag.Parse()

	ctx := context.Background()

	cfg, err := config.NewConfigFromURI(ctx, *config_uri)

	if err != nil {
		log.Fatalf("Failed to load config %s, %v", *config_uri, err)
	}

	wh_daemon, err := daemon.NewWebhookDaemonFromConfig(ctx, cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = wh_daemon.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
