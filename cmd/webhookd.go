package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"
	"github.com/whosonfirst/go-whosonfirst-webhookd/receivers"
	"os"
)

func ensure_ok(err error){

	if err != nil {
	   panic(err)
	}

}

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")
	var endpoint = flag.String("endpoint", "", "")

	var pubsub_host = flag.String("pubsub-host", "localhost", "...")
	var pubsub_port = flag.Int("pubsub-port", 6379, "...")
	var pubsub_channel = flag.String("pubsub-channel", "webhookd", "...")

	flag.Parse()

	daemon, err := webhookd.NewWebhookDaemon(*host, *port)
	ensure_ok(err)

	dispatcher, err := dispatchers.NewPubSubDispatcher(*pubsub_host, *pubsub_port, *pubsub_channel)
	ensure_ok(err)

	receiver, err := receivers.NewInsecureReceiver()
	ensure_ok(err)

	webhook, err := webhookd.NewWebhook(*endpoint, receiver, dispatcher)
	ensure_ok(err)

	err = daemon.AddWebhook(webhook)
	ensure_ok(err)

	err = daemon.Start()
	ensure_ok(err)

	os.Exit(0)
}
