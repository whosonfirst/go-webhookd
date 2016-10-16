package main

import (
	"flag"
	"fmt"
	ucd "github.com/cooperhewitt/go-ucd"
	"strings"
)

func main() {

	flag.Parse()

	args := flag.Args()
	str := strings.Join(args, " ")

	chars := strings.Split(str, "")

	for _, char := range chars {
		n := ucd.Name(char)
		fmt.Println(n)
	}
}
