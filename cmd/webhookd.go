package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"
	"github.com/whosonfirst/go-whosonfirst-webhookd/receivers"
	"log"
	"os"
)

func ensure_ok(err error) {

	if err != nil {
		panic(err)
	}

}

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")
	var endpoint = flag.String("endpoint", "", "")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := webhookd.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	daemon, err := webhookd.NewWebhookDaemon(config.Daemon.Host, config.Daemon.Port)
	ensure_ok(err)

	dispatcher, err := dispatchers.NewDispatcherFromConfig(config)
	ensure_ok(err)

	receiver, err := receivers.NewReceiverFromConfig(config)
	ensure_ok(err)

	webhook, err := webhookd.NewWebhook(*endpoint, receiver, dispatcher)
	ensure_ok(err)

	err = daemon.AddWebhook(webhook)
	ensure_ok(err)

	err = daemon.Start()
	ensure_ok(err)

	os.Exit(0)
}
