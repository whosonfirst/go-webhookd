prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-slackcat-writer; then rm -rf src/github.com/whosonfirst/go-slackcat-writer; fi
	mkdir -p src/github.com/whosonfirst/go-slackcat-writer
	cp slackcat.go src/github.com/whosonfirst/go-slackcat-writer/

deps:   self
	go get -u "github.com/whosonfirst/slackcat"

fmt:
	go fmt cmd/*.go
	go fmt *.go

test:	self fmt
	go build -o bin/test cmd/test.go
	go build -o bin/test_multi cmd/test_multi.go
