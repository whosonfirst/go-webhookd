CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/thisisaaronland/go-chicken; then rm -rf src/github.com/thisisaaronland/go-chicken; fi
	mkdir -p src/github.com/thisisaaronland/go-chicken
	cp -r strings src/github.com/thisisaaronland/go-chicken/
	cp -r emoji src/github.com/thisisaaronland/go-chicken/
	cp chicken.go src/github.com/thisisaaronland/go-chicken/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/cooperhewitt/go-ucd"
	# find vendor -name '.git' -print -type d -exec rm -rf {} +

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/chicken cmd/chicken.go

fmt:
	@GOPATH=$(GOPATH) go fmt cmd/*.go
	@GOPATH=$(GOPATH) go fmt chicken.go
	@GOPATH=$(GOPATH) go fmt strings/*.go
	@GOPATH=$(GOPATH) go fmt emoji/*.go

alpha:
	@GOPATH=$(GOPATH) go run cmd/build-alpha-codes.go > emoji/emoji.go
