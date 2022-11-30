# go-webhookd

![](docs/images/webhookd-arch.png)

What is the simplest webhook-wrangling server-daemon-thing?

`go-webhook` is a Go package that implements a bucket-brigrade style webhook server where requests are relayed through a receiver, one or more transformations and one or more dispatchers each of which have interfaces and are defined using a URI-based syntax to allow for custom processing.

Receivers handle the actual webhook side of things, doing auth and basic sanity checking and validation. Assuming everything is as it should be receivers return a bag of bytes (the actual webhook message that may or may not be massaged depending the receiver). That bag is then handed to one or more dispatchers which do _something_ with those bytes. Those details, including security considerations are left as an exercise to the reader.

In between (receivers and dispatchers) are an optional chain of transformations which accept bytes as their input, do _something_ with those bytes, and then return bytes.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-webhookd.svg)](https://pkg.go.dev/github.com/whosonfirst/go-webhookd)

## Install

You will need to have both `Go` (specifically version [1.18](https://golang.org/dl) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
$> make cli
go build -mod vendor -o bin/webhookd cmd/webhookd/main.go
go build -mod vendor -o bin/webhookd-generate-hook cmd/webhookd-generate-hook/main.go
go build -mod vendor -o bin/webhookd-flatten-config cmd/webhookd-flatten-config/main.go
go build -mod vendor -o bin/webhookd-inflate-config cmd/webhookd-inflate-config/main.go
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

### webhookd

```
./bin/webhookd -h
Usage of ./bin/webhookd:
  -config-uri string
    	A valid Go Cloud runtimevar URI representing your webhookd config file
```

`webhookd` is an HTTP daemon for handling webhook requests. Individual webhook endpoints (and how they are processed) are defined in a [config file](#config-files) that is read at start-up time.

#### Config URIs

The following [Go Cloud runtimevar URL schemes](https://gocloud.dev/concepts/urls/) are supported, by default, for defining config URIs:

* [constvar](https://godoc.org/gocloud.dev/runtimevar/constantvar)
* [file://](https://godoc.org/gocloud.dev/runtimevar/filevar)

#### Example

This is a deliberately juvenile example, just to keep things simple. 

Let's assume an insecure receiver with debugging enabled that reads input,
transforms it using the [go-chicken](https://github.com/aaronland/go-chicken)
`clucking` method and drops the results on the floor.

Here are the relevant settings in the config file:

```
{
	"daemon": "http://localhost:8080",
	...
	"webhooks": [
		{
			"endpoint": "/insecure-test",
	 		"receiver": "insecure://",
			"transformations": [ "clucking" ],
			"dispatchers": [ "null" ]
		}
	]
}
```

First we start `webhookd`:

```
./bin/webhookd -config-uri 'file:///usr/local/webhookd/config.json?decoder=string'
2018/07/21 08:43:37 webhookd listening for requests on http://localhost:8080
```

Then we pass `webhookd` a file along with a `debug=1` query parameter so that we
can see the output:

```
curl -v 'http://localhost:8080/insecure-test?debug=1' -d @README.md
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /insecure-test?debug=1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 12790
> Content-Type: application/x-www-form-urlencoded
> Expect: 100-continue
> 
< HTTP/1.1 100 Continue
* We are completely uploaded and fine
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Content-Type: text/plain
< X-Webhookd-Time-To-Dispatch: 16.907µs
< X-Webhookd-Time-To-Process: 13.033089ms
< X-Webhookd-Time-To-Receive: 209.332µs
< X-Webhookd-Time-To-Transform: 12.802186ms
< Date: Sat, 21 Jul 2018 15:43:40 GMT
< Transfer-Encoding: chunked
< 
# bok bok b'gawk-cluck cluck![](bok bok b'gawk/bok bok b'gawk-bok bok bok.cluck
cluck)bok bok bok bok bok b'gawk bok bok bok cluck cluck bok bok b'gawk-bok bok
bok cluck cluck-bok bok b'gawk-bok bok bok.bok bok b'gawk cluck cluck bok bok
b'gawk bok bok b'gawk bok bok b'gawk bok bok b'gawk cluck cluck bok bok b'gawk
bok bok bok cluck cluck bok bok bok-bok bok bok. bok bok
... and so on
```

#### Caveats

##### Dynamic endpoints

At some point there might be dynamic (or runtime) webhook endpoints but today there are not.

In the meantime you can gracefully restart `webhookd` by sending its PID a `USR2` signal which will cause the config file (and all the endpoints it defines) to be re-read. It's not elegant but it works. For example:

```
$> ./bin/webhookd -config-uri 'file:///usr/local/webhookd/config.json?decoder=string'
2016/10/16 00:19:47 Serving 127.0.0.1:8080 with pid 2723

$> kill -USR2 2723
2016/10/16 00:19:59 Graceful handoff of 127.0.0.1:8080 with new pid 2724 and old pid 2723
2016/10/16 00:19:59 Exiting pid 2723.
```

### Setting up a `webhookd` server

While you can set up a `webhookd` server by hand it's probably easier to all that work with a config file and let code take care of all the details, including registering all the webhooks. [Config files](#config-files) are discussed in detail below.

_All error handling in the examples below have been removed for the sake of brevity._

#### Setting up a `webhookd` server with a handy config file

```
import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"github.com/whosonfirst/go-webhookd/v3/daemon"
)

ctx := context.Background()
cfg, _ := config.NewConfigFromURI(ctx, "file:///usr/local/webhookd/config.json?decoder=string")

wh_daemon, _ := daemon.NewWebhookDaemonFromConfig(ctx, cfg)
wh_daemon.Start()
```

_You can also just grab the HTTP handler func with `wh_daemon.HandlerFunc()` if you need or want to start a webhookd daemon in your own way._

#### Setting up a `webhookd` server "by hand"

```
import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"		
	"github.com/whosonfirst/go-webhookd/v3/daemon"	
	"github.com/whosonfirst/go-webhookd/v3/dispatchers"
	"github.com/whosonfirst/go-webhookd/v3/receivers"
	"github.com/whosonfirst/go-webhookd/v3/transformations"
	"github.com/whosonfirst/go-webhookd/v3/webhook"
	_ "github.com/whosonfirst/go-webhookd-pubsub"		
)

ctx := context.Background()

wh_receiver, _ := receivers.NewReceiver(ctx, "insecure://")
null, _ := transfromations.NewTransformation(ctx, "null://")
pubsub, _ := dispatchers.NewDispatcher(ctx, "pubsub://localhost:6379/websocketd")

wh_transformations := []webhookd.WebhookTransformation{ null }
wh_dispatchers, _ := []webhookd.WebhookDispatcher{ pubsub }

wh, _ := webhook.NewWebhook("/foo", wh_receiver, wh_transformations, wh_dispatchers)

wh_daemon, _ := daemon.NewWebhookDaemon(ctx, "http://localhost:8080")
wh_daemon.AddWebhook(ctx, wh)
wh_daemon.Start(ctx)
```

Two important things to note:

* We're using an `Insecure` receiver with a `Null` transformation? These are included with the base `go-webhookd` package and are discussed in detail below.
* We're using a `PubSub` dispatcher which is made available by importing the [go-webhookd-pubsub](https://github.com/whosonfirst/go-webhookd-pubsub) package.

## Sending stuff to webhookd

```
curl -v http://localhost:8080/foo -d @README.md

* upload completely sent off: 703 out of 703 bytes
< HTTP/1.1 200 OK
< Date: Fri, 20 May 2016 01:31:05 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
```

## Where did it go...

```
./bin/subscribe webhookd
{'pattern': None, 'type': 'subscribe', 'channel': 'webhookd', 'data': 1L}
{'pattern': None, 'type': 'message', 'channel': 'webhookd', 'data': '# go-webhookd## ImportantYou should not try to use this, yet. No. No, really.## UsageIt _should_ work something like this. If you\'re reading this sentence that means it _doesn\'t_.```import (\t"github.com/whosonfirst/go-webhookd"\t"github.com/whosonfirst/go-webhookd/dispatchers"\t"github.com/whosonfirst/go-webhookd/receivers")dispatcher := dispatchers.NewPubSubDispatcher("localhost", 6379, "pubsub-channel")receiver := receivers.NewGitHubReceiver("github-webhook-s33kret")endpoint := "/wubwubwub"webhook := webhookd.NewWebhook(endpoint, receiver, dispatcher)daemon := webhookd.NewWebHookDaemon(webhook)daemon.AddWebhook(webhook)daemon.Start()```## See also'}
```

In this case, it went to Redis [PubSub](http://redis.io/topics/pubsub) land! Where things go depend on how you've configured your [dispatchers](#dispatchers-1).

## Config files

Config files for `webhookd` are JSON files consisting of five top-level sections. An [example config file](docs/config/config.json.example) is included with this repository. The five top-level sections are:

### daemon

```
	"daemon": "http://localhost:8080"
```

The `daemon` section is a dictionary defining configuration details for the `webhookd` daemon itself.

Valid daemon URI strings can be anything supported by the [aaronland/go-http-server](https://github.com/aaronland/go-http-server#server-schemes) package.

### receivers

```
	"receivers": {
		"insecure": "insecure://"
		"github": "github://?secret=s33kret"
	}
```

The `receivers` section is a dictionary of "named" receiver configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective receivers using the dictionary "name" as a simple short-hand.

_Note: This example includes a `github://` receiver which assumes you've imported the [go-webhookd-github](https://github.com/whosonfirst/go-webhookd-github) package in your code._

### transformations

```
	"transformations": {
		"chicken": "chicken://zxx?clucking=false"
	}
```

The `transformations` section is a dictionary of "named" tranformation configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective transformations using the dictionary "name" as a simple short-hand.

### dispatchers

```
	"dispatchers": {
		"pubsub": "pubsub://localhost:6379/webhookd"
	}
```

The `dispatchers` section is a dictionary of "named" dispatcher configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective dispatchers using the dictionary "name" as a simple short-hand.

_Note: This example includes a `pubsub://` receiver which assumes you've imported the [go-webhookd-github](https://github.com/whosonfirst/go-webhookd-pubsub) package in your code._

### webhooks

```
	"webhooks": [
		{
			"endpoint": "/github-test",
			"receiver": "github",
			"dispatchers": [ "pubsub" ]
		},
		{
			"endpoint": "/insecure-test",
		 	"receiver": "insecure",
			"dispatchers": [ "pubsub" ]
		}
	]
```

The `webhooks` section is a list of dictionaries. These are the actual webhook endpoints that clients (out there on the internet) will access.

* **endpoint** This is the path that a client will access. It _is_ the webhook URI that clients will send requests to.
* **receiver** The named receiver (defined in the `receivers` section) that the webhook will use to process requests.
* **transformations** An optional list of named transformations (defined in the `transformations` section) that the webhook process the message body with.
* **dispatchers** The list of named dispatchers (defined in the `dispatchers` section) that the webhook will relay a successful request to.

## Receivers

### Insecure

As the name suggests the `Insecure` receiver is completely insecure. It will happily accept anything you send to it and relay it on to the dispatcher defined for that webhook. It is defined as a URI string in the form of:

```
insecure://
```

This receiver exists primarily for debugging purposes and **you should not deploy it in production**.

## Transformations

### Chicken

The `Chicken` transformation will convert every word in your message to 🐔 using the [go-chicken](https://github.com/thisisaaronland/go-chicken) package. It is defined as a URI string in the form of:

```
chicken://{LANGUAGE}?clucking={CLUCKING}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| language | string | A three-letter language code specifying which language `go-chicken` should use. | yes |
| clucking | boolean | A boolean flag indicating whether or not to [cluck](https://github.com/thisisaaronland/go-chicken#clucking) when generating results. | no |

If this seems silly that's because it is. It's also more fun that yet-another boring _"make all the words upper-cased"_ example.

### Null

The `Null` transformation will not do _anything_. It's not clear why you would ever use this outside of debugging but that's your business. It is defined as a URI string in the form of:

```
null://
```

## Dispatchers

### Log

The `Log` dispatcher will send messages to Go's logging facility. As of this writing that means everything is logged to STDOUT but eventually it will be more sophisticated. It is defined as a URI string in the form of:

```
log://
```

### Null

The `Null` dispatcher will send messages in to the vortex, never to be seen again. This can be useful for debugging. It is defined as a URI string in the form of:

```
null://
```

## Halting a `webhookd` processing flow

As of `go-webhookd` v3.2.0 it is possible to "halt" a processing flow in mid-stream. This occurs is a receiver or transformation returns a `webhookd.WebhookError` with `Code` property whose value is `webhookd.HaltEvent`. These errors are treated as non-fatal but are treated as a signal to end processing and return immediately. Support for `webhookd.HaltEvent` in dispatchers is also enabled but they do not stop processing since dispatchers are invoked asynchronously.

## Testing

In advance of proper tests. In a terminal start `webhookd` like this:

```
go run -mod vendor cmd/webhookd/main.go \
	-config-uri 'file:///usr/local/go-webhookd/docs/config/config.json.example?decoder=string'

2020/05/22 17:18:49 webhookd listening for requests on http://localhost:8080
```

In another terminal, run the `webhookd-test` command like this:

```
go run cmd/webhookd-test-github/main.go \
	-config-uri 'file:///usr/local/whosonfirst/go-webhookd/config.json.example?decoder=string' \
	-endpoint insecure-test \
	-receiver insecure \
	-file docs/events/flights.json

2020/05/22 17:18:53 200 OK
```

In the first terminal you should see the following:

```
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/221/7/1713162217.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/221/9/1713162219.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/222/1/1713162221.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/222/3/1713162223.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/222/5/1713162225.geojson
... and so on
```

## To do

* [Add a general purpose "shared-secret/signed-message" receiver](https://github.com/whosonfirst/go-webhookd/issues/5)
* [Restrict access to receivers by host/IP](https://github.com/whosonfirst/go-webhookd/issues/6)
* Better logging

## Upgrading from `whosonfirst/go-webhookd/v2` 

`whosonfirst/go-webhookd/v3` does not introduce any _new_ functionality relative to `whosonfirst/go-webhookd/v2` but no longer comes with support for external platforms (GitHub, Slack, etc.) enabled by default. This functionality has been moved in to a number of separate `go-webhookd-{PLATFORM}` packages. This was done to make developing and adding custom receivers, transformations and dispatchers easier and modular.

You will need to add the relevant packages to your `cmd/webhookd/main.go` program. For example if your `webhookd` config file defines a GitHub receiver, a GitHub transformation and an AWS dispatcher you would need to import the [go-webhookd-github](https://github.com/whosonfirst/go-webhookd-github) and [go-webhookd-aws](https://github.com/whosonfirst/go-webhookd-aws) packages. Here's an abbreviated example in code, with error handling removed for the sake of brevity:

```
package main

import (
	"context"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"github.com/whosonfirst/go-webhookd/v3/daemon"
	_ "github.com/whosonfirst/go-webhookd-aws"
	_ "github.com/whosonfirst/go-webhookd-github"		
	"log"
	"os"
)

func main() {

	fs := flagset.NewFlagSet("webhooks")
	config_uri := fs.String("config-uri", "", "A valid Go Cloud runtimevar URI representing your webhookd config.")

	flagset.Parse(fs)

	ctx := context.Background()
	cfg, _ := config.NewConfigFromURI(ctx, *config_uri)

	wh_daemon, _ := daemon.NewWebhookDaemonFromConfig(ctx, cfg)
	wh_daemon.Start(ctx)
}

```

## See also

* https://github.com/whosonfirst/go-webhookd-aws
* https://github.com/whosonfirst/go-webhookd-github
* https://github.com/whosonfirst/go-webhookd-pubsub
* https://github.com/whosonfirst/go-webhookd-slack

## Related

* https://github.com/aaronland/go-http-server
* https://gocloud.dev/howto/runtimevar