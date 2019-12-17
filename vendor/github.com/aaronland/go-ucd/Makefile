CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/cooperhewitt/go-ucd; then rm -rf src/github.com/cooperhewitt/go-ucd; fi
	mkdir -p src/github.com/cooperhewitt/go-ucd/unicodedata
	mkdir -p src/github.com/cooperhewitt/go-ucd/unihan
	cp ucd.go src/github.com/cooperhewitt/go-ucd/
	cp unicodedata/unicodedata.go src/github.com/cooperhewitt/go-ucd/unicodedata/
	cp unihan/unihan.go src/github.com/cooperhewitt/go-ucd/unihan/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:	rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/pretty"

vendor-deps: rmdeps deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt unicodedata/*.go
	go fmt unihan/*.go
	go fmt cmd/*.go

data:
	@GOPATH=$(GOPATH) go run cmd/ucd-build-unicodedata.go -data https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt > unicodedata/unicodedata.go
	@GOPATH=$(GOPATH) go run cmd/ucd-build-unihan.go -data https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip > unihan/unihan.go

build:	bin

bin:	fmt self
	@GOPATH=$(GOPATH) go build -o bin/ucd cmd/ucd.go
	@GOPATH=$(GOPATH) go build -o bin/ucd-dump cmd/ucd-dump.go
	@GOPATH=$(GOPATH) go build -o bin/ucd-server cmd/ucd-server.go
