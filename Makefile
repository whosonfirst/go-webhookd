bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(PREVIOUS)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g'

cli:
	go build -mod vendor -o bin/webhookd cmd/webhookd/main.go
	go build -mod vendor -o bin/webhookd-generate-hook cmd/webhookd-generate-hook/main.go
	go build -mod vendor -o bin/webhookd-flatten-config cmd/webhookd-flatten-config/main.go
	go build -mod vendor -o bin/webhookd-inflate-config cmd/webhookd-inflate-config/main.go
