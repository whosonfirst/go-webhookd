package main

import (
	"flag"
	"github.com/whosonfirst/go-slackcat-writer"
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
