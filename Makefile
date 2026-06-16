GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck -show verbose ./...


bump-version:
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(OLD)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(OLD)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/whosonfirst\/go-webhookd\/$(OLD)/github.com\/whosonfirst\/go-webhookd\/$(NEW)/g'

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/webhookd cmd/webhookd/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/webhookd-generate-hook cmd/webhookd-generate-hook/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/webhookd-flatten-config cmd/webhookd-flatten-config/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/webhookd-inflate-config cmd/webhookd-inflate-config/main.go

local-scan:
	/usr/local/bin/sonar-scanner/bin/sonar-scanner -Dsonar.projectKey=go-webhookd -Dsonar.sources=. -Dsonar.host.url=http://localhost:9000 -Dsonar.login=$(TOKEN)

godoc:
	godoc -http=:6060
