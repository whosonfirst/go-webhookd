package main

import (
	"flag"
	"github.com/whosonfirst/go-webhookd/config"
	"github.com/whosonfirst/go-webhookd/daemon"
	"log"
	"os"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	wh_config, err := config.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	wh_daemon, err := daemon.NewWebhookDaemonFromConfig(wh_config)

	if err != nil {
		log.Fatal(err)
	}

	err = wh_daemon.Start()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
