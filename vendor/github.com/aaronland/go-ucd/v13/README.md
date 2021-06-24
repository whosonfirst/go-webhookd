# go-ucd

Go libraries and utilities for working with Unicode character data.

_This package was forked from the [cooperhewitt/ucd](https://github.com/cooperhewitt/go-ucd) package but builds everything under its own `aaronland` namespace._

## Example

	package main

	import(
		"fmt"
		"flag"
		"github.com/aaronland/go-ucd/v13"
	)

	func main(){

	     flag.Parse()
	     char := flag.Arg(0)

	     name := ucd.Name(char)
	     fmt.Println(name)
	}

## Tools

The following tools are included in the `cmd` directory. Note however that you will need to [compile them](https://golang.org/cmd/go/#hdr-Compile_and_run_Go_program) yourself. You can do this (and all the steps in-between using the handy [Makefile](Makefile) and the `build` target included in this repository. Like this:

```
$> make tools
```

This will build the `ucd` and `ucd-server` applications and place them in the `bin` directory.

### ucd

	$> bin/ucd A
	LATIN CAPITAL LETTER A

	$> ucd THIS → WAY
	LATIN CAPITAL LETTER T
	LATIN CAPITAL LETTER H
	LATIN CAPITAL LETTER I
	LATIN CAPITAL LETTER S
	SPACE
	RIGHTWARDS ARROW
	SPACE
	LATIN CAPITAL LETTER W
	LATIN CAPITAL LETTER A
	LATIN CAPITAL LETTER Y

`ucd` supports the Unicode Han Data character set, or at least endeavours to. There may still be bugs.

	$> bin/ucd 䍕
	NET; WEB; NETWORK, NET FOR CATCHING RABBIT

### ucd-server

#### Usage

	$> bin/ucd-server --help
	Usage of ./ucd-server:
	  -host="localhost": host
	  -port=8080: port

#### Install as a service

To install as an init.d script, copy the example provided, replace the values of UCD_USER, UCD_DAEMON and UCD_PORT, and start the service.

        $> sudo cp init.d/ucd-server.sh.example /etc/init.d/ucd-server.sh
        $> sudo service ucd-server start

#### as JSON

	$> curl -X GET -s 'http://localhost:8080/?text=♕%20HAT' | python -mjson.tool
	{
	    "Chars": [
	        {
	            "Char": "\u2655",
	            "Hex": "2655",
	            "Name": "WHITE CHESS QUEEN"
	        },
	        {
	            "Char": " ",
	            "Hex": "0020",
	            "Name": "SPACE"
	        },
	        {
	            "Char": "H",
	            "Hex": "0048",
	            "Name": "LATIN CAPITAL LETTER H"
	        },
	        {
	            "Char": "A",
	            "Hex": "0041",
	            "Name": "LATIN CAPITAL LETTER A"
	        },
	        {
	            "Char": "T",
	            "Hex": "0054",
	            "Name": "LATIN CAPITAL LETTER T"
	        }
	    ]
	}

#### As plain text

	$> curl -H 'Accept: text/plain' -s 'http://localhost:8080/?text=♕%20HAT%20WITH%20😸'
	WHITE CHESS QUEEN
	SPACE
	LATIN CAPITAL LETTER H
	LATIN CAPITAL LETTER A
	LATIN CAPITAL LETTER T
	SPACE
	LATIN CAPITAL LETTER W
	LATIN CAPITAL LETTER I
	LATIN CAPITAL LETTER T
	LATIN CAPITAL LETTER H
	SPACE
	GRINNING CAT FACE WITH SMILING EYES

## Versions

`go-ucd` supports Unicode 13.0 as of February 16, 2021 and requires Go [1.16](https://golang.org/doc/go1.16) or higher to compile.

This package exports data defined in the `UnicodeData.txt` and the `Unihan.zip`. Both are available from
http://unicode.org/Public/UCD/latest/ucd/.

If the Unicode consortium releases newer data files and you want or need to
updated your version of `go-ucd` before we do you do so manually by using the
`ucd-build-unicodedata` and `ucd-build-unihan` tools included in the [bin
directory](https://github.com/aaronland/go-ucd/tree/master/bin). For example:

```
go run ./cmd/ucd-build-unicodedata.go > ./unicodedata/unicodedata.go
go run ./cmd/ucd-build-unihan.go > ./unihan/unihan.go
```

_Note: You will need to recompile your `ucd` and `ucd-server` binaries manually._

## Shout outs

Many thanks to friend and Go-friend [Richard Crowley](https://github.com/rcrowley) who is always kind and patient answering my Go-related questions. Go is lovely but Go is weird.

## See also

* http://unicode.org/Public/UCD/latest/ucd/
* http://www.washingtonpost.com/news/the-intersect/wp/2015/02/23/the-surprisingly-complex-reason-you-never-see-emoji-urls/
* https://modelviewculture.com/pieces/i-can-text-you-a-pile-of-poo-but-i-cant-write-my-name
