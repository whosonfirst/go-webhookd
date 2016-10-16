CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/cooperhewitt/go-ucd; then rm -rf src/github.com/cooperhewitt/go-ucd; fi
	mkdir -p src/github.com/cooperhewitt/go-ucd/unicodedata
	mkdir -p src/github.com/cooperhewitt/go-ucd/unihan
	cp ucd.go src/github.com/cooperhewitt/go-ucd/
	cp unicodedata/unicodedata.go src/github.com/cooperhewitt/go-ucd/unicodedata/
	cp unihan/unihan.go src/github.com/cooperhewitt/go-ucd/unihan/

fmt:
	go fmt *.go
	go fmt unicodedata/*.go
	go fmt unihan/*.go
	go fmt cmd/*.go

data: 	
	@GOPATH=$(GOPATH) go run cmd/ucd-build-unicodedata.go > unicodedata/unicodedata.go
	@GOPATH=$(GOPATH) go run cmd/ucd-build-unihan.go > unihan/unihan.go

build:	fmt self
	@GOPATH=$(GOPATH) go build -o bin/ucd cmd/ucd.go
	@GOPATH=$(GOPATH) go build -o bin/ucd-server cmd/ucd-server.go
