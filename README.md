# go-whosonfirst-webhookd

What is the simplest webhook-wrangling server-daemon-thing.

In many ways (at least so far) this is nothing more than a fancy bucket-brigade. By design. Receivers handle the actual webhook side of things, doing auth and basic sanity checking and validation. Assuming everything is as it should be receivers return a bag of bytes (the actual webhook message that may or may not be massaged depending the receiver). That bag is then handed to a dispatcher which does _something_ with those bytes. Those details, including security considerations are left as an exercise to the reader.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

### webhookd

```
./bin/webhookd -h
Usage of ./bin/webhookd:
  -config string
    	Path to a valid webhookd config file
```

`webhookd` is an HTTP daemon for handling webhook requests. Individual webhook endpoints (and how they are processed) are defined in a config file that is read at start-up time.

_In the future there might be dynamic (or runtime) webhook endpoints but today there are not._


### Setting up a `webhookd` server "by hand"

_All error handling has been removed from the examples below for the sake of brevity._

```
import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"
	"github.com/whosonfirst/go-whosonfirst-webhookd/receivers"
)

// imagine flags here but also imagine that *pubsub_channel is "webhookd"
// and *endpoint is "/foo"

flag.Parse()

daemon, _ := webhookd.NewWebhookDaemon(*host, *port)

receiver, _ := receivers.NewInsecureReceiver()
dispatcher, _ := dispatchers.NewPubSubDispatcher(*pubsub_host, *pubsub_port, *pubsub_channel)

webhook, _ := webhookd.NewWebhook(*endpoint, receiver, dispatcher)
daemon.AddWebhook(webhook)

daemon.Start()
```

See the way we're using an `Insecure` receiver and a `PubSub` dispatcher. Both are discussed in detail below.

### Setting up a `webhookd` server with a handy config file

```
import (
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/daemon"
)

config, _ := webhookd.NewConfigFromFile("config.json")

d, _ := daemon.NewWebhookDaemonFromConfig(config)
d.Start()
```

While you can set up a `webhookd` server by hand it's probably easier to all that work with a config file and let code take care of all the details, including registering all the webhooks. Config files are discussed in detail below.

### Sending stuff to webhookd

```
curl -v -X POST http://localhost:8080/foo -d @README.md

* upload completely sent off: 703 out of 703 bytes
< HTTP/1.1 200 OK
< Date: Fri, 20 May 2016 01:31:05 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
```

### Where did it go...

```
./bin/subscribe webhookd
{'pattern': None, 'type': 'subscribe', 'channel': 'webhookd', 'data': 1L}
{'pattern': None, 'type': 'message', 'channel': 'webhookd', 'data': '# go-whosonfirst-webhookd## ImportantYou should not try to use this, yet. No. No, really.## UsageIt _should_ work something like this. If you\'re reading this sentence that means it _doesn\'t_.```import (\t"github.com/whosonfirst/go-whosonfirst-webhookd"\t"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"\t"github.com/whosonfirst/go-whosonfirst-webhookd/receivers")dispatcher := dispatchers.NewPubSubDispatcher("localhost", 6379, "pubsub-channel")receiver := receivers.NewGitHubReceiver("github-webhook-s33kret")endpoint := "/wubwubwub"webhook := webhookd.NewWebhook(endpoint, receiver, dispatcher)daemon := webhookd.NewWebHookDaemon(webhook)daemon.AddWebhook(webhook)daemon.Start()```## See also'}
```

It went to Redis [PubSub](http://redis.io/topics/pubsub) land!

## Config files

Config files for `webhookd` are JSON files consisting of four top-level sections. They are:

### daemon

```
	"daemon": {
		"host": "localhost",
		"port": 8080
	}
```

The `daemon` section is a dictionary defining configuration details for the `webhookd` daemon itself.

### receivers

```
	"receivers": {
		"insecure": {
			"name": "Insecure"
		},
		"github": {
			"name": "GitHub",
			"secret": "s33kret"
		}
	}
```

The `receivers` section is a dictionary of "named" receiver configuations. This allows the actual webhook configurations (described below) to signal their respective receivers using the receiver "name" as a simple short-hand.

### dispatchers

```
	"dispatchers": {
		"pubsub": {
			"name": "PubSub",
			"host": "localhost",
			"port": 6379,
			"channel": "webhookd"
		}
	}
```

The `dispatchers` section is a dictionary of "named" dispatcher configuations. This allows the actual webhook configurations (described below) to signal their respective dispatchers using the dispatcher "name" as a simple short-hand.

### webhooks

```
	"webhooks": [
		{ "endpoint": "/github-test", "receiver": "github", "dispatcher": "pubsub" },
		{ "endpoint": "/insecure-test", "receiver": "insecure", "dispatcher": "pubsub" }		
	]
```

The `webhooks` section is a list of dictionaries. These are the actual webhook endpoints that clients (out there on the internet) will access.

#### endpoint

This is the path that a client will access. It _is_ the webhook.

#### receiver

The named receiver (defined in the `receivers` section) that the webhook will use to process requests.

#### dispatcher

The named dispatcher (defined in the `dispatchers` section) that the webhook will relay a successful request to.

## Receivers

### Insecure

```
	{
		"name": "Insecure"
	}
```

As the name suggests this receiver is completely insecure. It will happily accept anything you send to it and relay it on to the dispatcher defined for that webhook. This receiver exists primarily for debugging purposes and **you should not deploy it in production**.

The `Insecure` receiver has the following properties:

#### name

This is always `Insecure`.

### GitHub

```
	{
		"name": "GitHub",
		"secret": "s33kret"
	}
```

TBW.

The `GitHub` receiver has the following properties:

#### name

This is always `GitHub`.

#### secret

TBW

## Dispatchers

### PubSub

```
	{
		"name": "PubSub",
		"host": "localhost",
		"port": 6379,
		"channel": "webhookd"
	}
```

TBW.

The `PubSub` dispatcher has the following properties:

#### name _string_

This is always `PubSub`.

#### host _string_

The address of the Redis host you want to connect to.

#### port _int_

The port number of the Redis host you want to connect to.

#### channel _string_

The name of the Redis PubSub channel you want to send messages to.

## To do

* More documentation
* Logging

## See also

* https://en.wikipedia.org/wiki/Webhook
* http://redis.io/topics/pubsub
