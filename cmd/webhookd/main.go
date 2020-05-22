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

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")

	flag.Parse()

	ctx := context.Background()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	wh_config, err := config.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	wh_daemon, err := daemon.NewWebhookDaemonFromConfig(ctx, wh_config)

	if err != nil {
		log.Fatal(err)
	}

	err = wh_daemon.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
