tools:
	go build -mod vendor -o bin/webhookd cmd/webhookd/main.go
	go build -mod vendor -o bin/webhookd-test-github cmd/webhookd-test-github/main.go
	go build -mod vendor -o bin/webhookd-generate-hook cmd/webhookd-generate-hook/main.go

debug:
	bin/webhookd -config ./config.json

hook:
	./bin/webhookd-generate-hook
