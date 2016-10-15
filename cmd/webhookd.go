package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/daemon"
	"log"
	"os"
)

func ensure_ok(err error) {

	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := webhookd.NewConfigFromFile(*cfg)
	ensure_ok(err)

	d, err := daemon.NewWebhookDaemon(config.Daemon.Host, config.Daemon.Port)
	ensure_ok(err)

	err = d.AddWebhooksFromConfig(config)
	ensure_ok(err)

	err = d.Start()
	ensure_ok(err)

	os.Exit(0)
}
