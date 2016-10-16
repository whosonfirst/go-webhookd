CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-webhookd; then rm -rf src/github.com/whosonfirst/go-webhookd; fi
	mkdir -p src/github.com/whosonfirst/go-webhookd/daemon
	mkdir -p src/github.com/whosonfirst/go-webhookd/dispatchers
	mkdir -p src/github.com/whosonfirst/go-webhookd/github
	mkdir -p src/github.com/whosonfirst/go-webhookd/receivers
	mkdir -p src/github.com/whosonfirst/go-webhookd/transformations
	cp webhookd.go src/github.com/whosonfirst/go-webhookd/
	cp daemon/*.go src/github.com/whosonfirst/go-webhookd/daemon/
	cp dispatchers/*.go src/github.com/whosonfirst/go-webhookd/dispatchers/
	cp github/*.go src/github.com/whosonfirst/go-webhookd/github/
	cp receivers/*.go src/github.com/whosonfirst/go-webhookd/receivers/
	cp transformations/*.go src/github.com/whosonfirst/go-webhookd/transformations/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "gopkg.in/redis.v1"
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/grace/gracehttp"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-writer-slackcat"
	@GOPATH=$(GOPATH) go get -u "github.com/thisisaaronland/go-chicken"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt daemon/*.go
	go fmt dispatchers/*.go
	go fmt github/*.go
	go fmt receivers/*.go
	go fmt transformations/*.go
	go fmt *.go

bin: 	rmdeps self
	@GOPATH=$(shell pwd) go build -o bin/webhookd cmd/webhookd.go
