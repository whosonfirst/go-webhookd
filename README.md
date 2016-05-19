# go-whosonfirst-webhookd

## Important

You should not try to use this, yet. No. No, really.

## Usage

It _should_ work something like this. If you're reading this sentence that means it _doesn't_.

```
import (
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers"
	"github.com/whosonfirst/go-whosonfirst-webhookd/receivers"
)

dispatcher := dispatchers.NewPubSubDispatcher("localhost", 6379, "pubsub-channel")
receiver := receivers.NewGitHubReceiver("github-webhook-s33kret")

endpoint := "/wubwubwub"
webhook := webhookd.NewWebhook(endpoint, receiver, dispatcher)

daemon := webhookd.NewWebHookDaemon(webhook)
daemon.AddWebhook(webhook)

daemon.Start()
```

## See also