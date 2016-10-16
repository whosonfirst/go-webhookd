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

### Setting up webhookd "by hand"

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

### Setting up webhookd with a handy config file

```
import (
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/daemon"
)

config, _ := webhookd.NewConfigFromFile("config.json")

d, _ := daemon.NewWebhookDaemon(config.Daemon.Host, config.Daemon.Port)

d.AddWebhooksFromConfig(config)
d.Start()
```

### Sending stuff to webhookd

```
curl -v -X POST http://localhost:8080/foo -d @README.md
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /foo HTTP/1.1
> User-Agent: curl/7.35.0
> Host: localhost:8080
> Accept: */*
> Content-Length: 703
> Content-Type: application/x-www-form-urlencoded
> 
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

## Utilities

### webhookd

```
./bin/webhookd -h
Usage of ./bin/webhookd:
  -config string
    	Path to a valid webhookd config file
```

## Config files

### daemon

```
	"daemon": {
		"host": "localhost",
		"port": 8080
	}
```

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

### webhooks

```
	"webhooks": [
		{ "endpoint": "/github-test", "receiver": "github", "dispatcher": "pubsub" },
		{ "endpoint": "/insecure-test", "receiver": "insecure", "dispatcher": "pubsub" }		
	]
```

## Receivers

### Insecure

### GitHub

## Dispatchers

### PubSub

## To do

* Documentation
* Logging

## See also

* http://redis.io/topics/pubsub
