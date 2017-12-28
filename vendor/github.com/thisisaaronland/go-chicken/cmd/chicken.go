package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-chicken"
	"os"
)

func main() {

	var lang = flag.String("language", "zxx", "A valid ISO-639-3 language code.")
	var clucking = flag.Bool("clucking", false, "Make chicken noises")

	flag.Parse()

	ch, err := chicken.GetChickenForLanguageTag(*lang, *clucking)

	if err != nil {
		panic(err)
	}

	for _, path := range flag.Args() {

		var buf *bufio.Scanner

		if path == "-" {
			buf = bufio.NewScanner(os.Stdin)
		} else {

			fh, err := os.Open(path)

			if err != nil {
				panic(err)
			}

			buf = bufio.NewScanner(fh)
		}

		for buf.Scan() {
			txt := buf.Text()
			fmt.Println(ch.TextToChicken(txt))
		}
	}
}
