CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-webhookd; then rm -rf src/github.com/whosonfirst/go-whosonfirst-webhookd; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-webhookd/receivers
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers
	cp webhookd.go src/github.com/whosonfirst/go-whosonfirst-webhookd/
	cp receivers/*.go src/github.com/whosonfirst/go-whosonfirst-webhookd/receivers/
	cp dispatchers/*.go src/github.com/whosonfirst/go-whosonfirst-webhookd/dispatchers/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(shell pwd) go get -u "gopkg.in/redis.v1"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt receivers/*.go
	go fmt dispatchers/*.go
	go fmt *.go

bin: 	self
	@GOPATH=$(shell pwd) go build -o bin/webhookd cmd/webhookd.go
