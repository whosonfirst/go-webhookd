# go-writer-slackcat

A Go package for sending messages to a Slack channel using a standard io.Writer interface

## Usage

### Simple

```
package main

import (
	"flag"
	"github.com/whosonfirst/go-writer"
	"strings"
)

func main() {

	var config = flag.String("config", "", "The path to your Slackcat config file")

	flag.Parse()
	args := flag.Args()

	w, err := slackcat.NewWriter(*config)

	if err != nil {
		panic(err)
	}

	msg := strings.Join(args, " ")

	b := []byte("Write " + msg)

	_, err = w.Write(b)

	if err != nil {
		panic(err)
	}

	// WriteString is provided as a convenience if you don't
	// feel like []byte("-ing") all the things per the default
	// io.Writer interface spec

	_, err = w.WriteString("WriteString " + msg)

	if err != nil {
		panic(err)
	}
}
```

### Fancy

```
package main

import (
	"flag"
	"github.com/whosonfirst/go-writer-slackcat"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var config = flag.String("config", "", "The path to your Slackcat config file")

	flag.Parse()
	args := flag.Args()

	msg := strings.Join(args, " ")

	slack, _ := slackcat.NewWriter(*config)

	writer := io.MultiWriter(os.Stdout, slack)

	logger := log.New(writer, "[example] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(msg)
}
```

## See also

* https://github.com/whosonfirst/slackcat
