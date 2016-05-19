prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-webhookd; then rm -rf src/github.com/whosonfirst/go-whosonfirst-webhookd; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-webhookd/service
	cp webhookd.go src/github.com/whosonfirst/go-whosonfirst-webhookd/
	cp service/*.go src/github.com/whosonfirst/go-whosonfirst-webhookd/service/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps deps fmt bin

deps:   self
	@GOPATH=$(shell pwd) get -u "gopkg.in/redis.v1"

fmt:
	go fmt cmd/*.go
	go fmt service/*.go
	go fmt *.go

bin: 	self
	@GOPATH=$(shell pwd) go build -o bin/webhookd cmd/webhookd.go
