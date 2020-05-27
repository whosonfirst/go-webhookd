# go-webhookd

![](docs/images/webhookd-arch.png)

What is the simplest webhook-wrangling server-daemon-thing.

In many ways this is nothing more than a fancy bucket-brigade. By design.

Receivers handle the actual webhook side of things, doing auth and basic sanity checking and validation. Assuming everything is as it should be receivers return a bag of bytes (the actual webhook message that may or may not be massaged depending the receiver). That bag is then handed to one or more dispatchers which do _something_ with those bytes. Those details, including security considerations are left as an exercise to the reader.

In between (receivers and dispatchers) are an optional chain of transformations which accept bytes as their input, do _something_ with those bytes, and then return bytes.

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

`whosonfirst/go-webhookd/v2` does not introduce any new functionality relative to `whosonfirst/go-webhookd` "v1" but does make substantial changes to the package's interfaces and config file definitions.

Once all the kinks have been worked out of `whosonfirst/go-webhookd/v2` it will quickly be superseded by `whosonfirst/go-webhookd/v3` which will move most of the platform or vendor specific functionalities in to their own packages.

## Upgrading from `whosonfirst/go-webhookd` "v1"

* The `-config` flag has been replaced by a `-config-uri` flag which is a fully-qualified [Go Cloud runtimevar URI](https://gocloud.dev/concepts/urls/).
* The config file itself has been simplified. Daemon, receiver, dispatcher and tranformation settings are now defined as URI strings rather than dictionaries. `whosonfirst/go-webhookd` "v1" config files will need to be updated manually.

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

See the way we're using an `Insecure` receiver and a `PubSub` dispatcher with a `Null` transformation? All are these are discussed in detail below.

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
		"insecure": "insecure://",
		"github": "github://?secret=s33kret"
	}
```

The `receivers` section is a dictionary of "named" receiver configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective receivers using the dictionary "name" as a simple short-hand.

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
		},
		{
			"endpoint": "/slack-test",
			"receiver": "slack",
			"transformations": [ "slack", "chicken" ],
			"dispatchers": [ "slack", "log"]
		}
	]
```

The `webhooks` section is a list of dictionaries. These are the actual webhook endpoints that clients (out there on the internet) will access.

* **endpoint** This is the path that a client will access. It _is_ the webhook URI that clients will send requests to.
* **receiver** The named receiver (defined in the `receivers` section) that the webhook will use to process requests.
* **transformations** An optional list of named transformations (defined in the `transformations` section) that the webhook process the message body with.
* **dispatchers** The list of named dispatchers (defined in the `dispatchers` section) that the webhook will relay a successful request to.

## Receivers

### GitHub

The `GitHub` receiver handles Webhooks sent from [GitHub](https://developer.github.com/webhooks/). It validates that the message sent is actually from GitHub (by way of the `X-Hub-Signature` header) but performs no other processing. It is defined as a URI string in the form of:

```
github://?secret={SECRET}&ref={REF}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| secret | string | The secret used to generate [the HMAC hex digest](https://developer.github.com/webhooks/#delivery-headers) of the message payload. | yes |
| ref | string | An optional Git `ref` to filter by. If present and a WebHook is sent with a different ref then the daemon will return a `666` error response. | no |

### Insecure

As the name suggests the `Insecure` receiver is completely insecure. It will happily accept anything you send to it and relay it on to the dispatcher defined for that webhook. It is defined as a URI string in the form of:

```
insecure://
```

This receiver exists primarily for debugging purposes and **you should not deploy it in production**.

### Slack

The `Slack` receiver handles Webhooks sent from [Slack](https://api.slack.com/outgoing-webhooks). It does not process the message at all. It is defined as a URI string in the form of:

```
slack://
```

_This receiver has not been fully tested yet so proceed with caution._

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

### GitHubCommits

The `GitHubCommits` transformation will extract all the commits (added, modified, removed) from a `push` event and return a CSV encoded list of rows consisting of: commit hash, repository name, path. For example:

```
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/450/9/1713164509.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/451/9/1713164519.geojson
e3a18d4de60a5e50ca78ca1733238735ddfaef4c,sfomuseum-data-flights-2020-05,data/171/316/483/5/1713164835.geojson
````

It is defined as a URI string in the form of:

```
githubcommits://?exclude_additions={EXCLUDE_ADDITIONS}&exclude_modification={EXCLUDE_MODIFICATIONS}&exclude_deletions={EXCLUDE_DELETIONS}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| exclude_additions| boolean | A flag to indicate that new additions in a commit should be ignored. | no |
| exclude_modifications| boolean | A flag to indicate that modifications in a commit should be ignored. | no |
| exclude_deletions | boolean | A flag to indicate that deletions in a commit should be ignored. | no |

### GitHubRepo

The `GitHubRepo` transformation will extract the reporsitory name for all the commits matching (added, modified, removed) criteria. It is defined as a URI string in the form of:

```
githubrepo://?exclude_additions={EXCLUDE_ADDITIONS}&exclude_modification={EXCLUDE_MODIFICATIONS}&exclude_deletions={EXCLUDE_DELETIONS}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| exclude_additions| boolean | A flag to indicate that new additions in a commit should be ignored. | no |
| exclude_modifications| boolean | A flag to indicate that modifications in a commit should be ignored. | no |
| exclude_deletions | boolean | A flag to indicate that deletions in a commit should be ignored. | no |

### Null

The `Null` transformation will not do _anything_. It's not clear why you would ever use this outside of debugging but that's your business. It is defined as a URI string in the form of:

```
null://
```

### SlackText

The `SlackText` transformation will extract and return [the `text` property](https://api.slack.com/outgoing-webhooks) from a Webhook sent by Slack. It is defined as URI string in the form of:

```
slacktext://
```

## Dispatchers

### Lambda

The `Lambda` dispatcher will send messages to an Amazon Web Services (ASW) [Lambda function](#). It is defined as a URI string in the form of:

```
lambda://{FUNCTION}?dsn={DSN}&invocation_type={INVOCATION_TYPE}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| dsn | string | A valid `aaronland/go-aws-session` DSN string. | yes |
| function | string | The name of your Lambda function. | yes |
| invocation_type | string | A valid AWS Lambda `Invocation Type` string. | no |

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

### PubSub

The `PubSub` dispatcher will send messages to a Redis PubSub channel. It is defined as a URI string in the form of:

```
pubsub://{REDIS_HOST}:{REDIS_PORT}/{REDIS_CHANNEL}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| redis_host | string | The path to a valid [slackcat](https://github.com/whosonfirst/slackcat#configuring) config file. | yes |

_Eventually you will be able to specify a plain-vanilla Slack Webhook URL but not today._

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

## See also

* https://en.wikipedia.org/wiki/Webhook

## Related

* https://github.com/aaronland/go-http-server
* https://github.com/whosonfirst/go-pubssed
* https://gocloud.dev/howto/runtimevar