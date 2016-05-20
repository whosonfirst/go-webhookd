# go-whosonfirst-webhookd

What is the simplest webhook-wrangling server-daemon-thing.

## Important

This should be considered "wet paint". It has not been tested much and may yet change in significant ways.

## Usage

### Setting up webhookd


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

// error handling removed for brevity

daemon, _ := webhookd.NewWebhookDaemon(*host, *port)

dispatcher, _ := dispatchers.NewPubSubDispatcher(*pubsub_host, *pubsub_port, *pubsub_channel)
receiver, _ := receivers.NewInsecureReceiver()

webhook, _ := webhookd.NewWebhook(*endpoint, receiver, dispatcher)

daemon.AddWebhook(webhook)
daemon.Start()
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

## Receivers

_Please write me_

## Dispatchers

_Please write me_

## To do

* Documentation
* Logging
* Defining webhooks with a config file or something
