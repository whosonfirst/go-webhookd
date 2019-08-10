fmt:
	go fmt *.go
	go fmt unicodedata/*.go
	go fmt unihan/*.go
	go fmt cmd/*.go

data:
	go run cmd/ucd-build-unicodedata/main.go -data https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt > unicodedata/unicodedata.go
	go run cmd/ucd-build-unihan/main.go -data https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip > unihan/unihan.go

tools:	
	go build -o bin/ucd cmd/ucd/main.go
	go build -o bin/ucd-dump cmd/ucd-dump/main.go
	go build -o bin/ucd-server cmd/ucd-server/main.go
