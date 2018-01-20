CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-webhookd; then rm -rf src/github.com/whosonfirst/go-webhookd; fi
	mkdir -p src/github.com/whosonfirst/go-webhookd/config
	mkdir -p src/github.com/whosonfirst/go-webhookd/daemon
	mkdir -p src/github.com/whosonfirst/go-webhookd/dispatchers
	mkdir -p src/github.com/whosonfirst/go-webhookd/github
	mkdir -p src/github.com/whosonfirst/go-webhookd/receivers
	mkdir -p src/github.com/whosonfirst/go-webhookd/transformations
	mkdir -p src/github.com/whosonfirst/go-webhookd/webhook
	cp webhookd.go src/github.com/whosonfirst/go-webhookd/
	cp config/*.go src/github.com/whosonfirst/go-webhookd/config/
	cp daemon/*.go src/github.com/whosonfirst/go-webhookd/daemon/
	cp dispatchers/*.go src/github.com/whosonfirst/go-webhookd/dispatchers/
	cp github/*.go src/github.com/whosonfirst/go-webhookd/github/
	cp receivers/*.go src/github.com/whosonfirst/go-webhookd/receivers/
	cp transformations/*.go src/github.com/whosonfirst/go-webhookd/transformations/
	cp webhook/*.go src/github.com/whosonfirst/go-webhookd/webhook/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "gopkg.in/redis.v1"
	@GOPATH=$(GOPATH) go get -u "github.com/google/go-github/github"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-writer-slackcat"
	@GOPATH=$(GOPATH) go get -u "github.com/thisisaaronland/go-chicken"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt config/*.go
	go fmt daemon/*.go
	go fmt dispatchers/*.go
	go fmt github/*.go
	go fmt receivers/*.go
	go fmt transformations/*.go
	go fmt webhook/*.go
	go fmt *.go

bin: 	rmdeps self
	@GOPATH=$(shell pwd) go build -o bin/webhookd cmd/webhookd.go
	@GOPATH=$(shell pwd) go build -o bin/webhookd-test-github cmd/webhookd-test-github.go
	@GOPATH=$(shell pwd) go build -o bin/webhookd-generate-hook cmd/webhookd-generate-hook.go

debug:
	bin/webhookd -config ./config.json

hook:
	./bin/webhookd-generate-hook
