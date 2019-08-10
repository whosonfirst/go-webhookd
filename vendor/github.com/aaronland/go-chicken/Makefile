tools:
	go build -o bin/chicken cmd/chicken/main.go
	go build -o bin/rooster cmd/rooster/main.go

wasm:	
	GOARCH=wasm GOOS=js go build -o www/chicken.wasm cmd/chicken-wasm/main.go

alpha:
	go run cmd/build-alpha-codes/main.go > emoji/emoji.go

docker-build:
	docker build -t rooster .

docker-debug: docker-build
	docker run -it -p 1280:1280 rooster
