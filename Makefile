bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g'

tools:
	go build -mod vendor -o bin/webhookd cmd/webhookd/main.go
	go build -mod vendor -o bin/webhookd-test-github cmd/webhookd-test-github/main.go
	go build -mod vendor -o bin/webhookd-generate-hook cmd/webhookd-generate-hook/main.go

debug:
	bin/webhookd -config ./config.json

hook:
	./bin/webhookd-generate-hook

lambda-config:
	go run cmd/webhookd-flatten-config/main.go -config $(CONFIG) -constvar | pbcopy

lambda: lambda-webhookd

lambda-webhookd:
	if test -f main; then rm -f main; fi
	if test -f webhookd.zip; then rm -f webhookd.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/webhookd/main.go
	zip webhookd.zip main
	rm -f main

