package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/service"
	"gopkg.in/redis.v1"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	pubsub := redis.NewTCPClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer pubsub.Close()
	pubsub.Publish("foo", "starting up")

	webhook, _ := service.NewInsecureWebhook(pubsub)

	fmt.Println(endpoint)

	daemon, _ := webhookd.NewWebhookDaemon(endpoint, webhook)
	daemon.Start()

	os.Exit(0)
}
