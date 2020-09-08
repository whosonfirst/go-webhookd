package main

import (
	_ "context"
	"flag"
	"net/url"
	"log"
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"github.com/whosonfirst/go-webhookd/v3/config"
)

func main (){

	flag.Parse()

	for _, uri := range flag.Args() {

		u, err := url.Parse(uri)

		if err != nil {
			log.Fatal(err)
		}
		
		q := u.Query()

		val := q.Get("val")

		if val == "" {
			log.Fatal("Missing ?val parameter")
		}

		var cfg *config.WebhookConfig

		err = json.Unmarshal([]byte(val), &cfg)

		if err != nil {
			log.Fatal(err)
		}

		body, err := json.Marshal(cfg)

		if err != nil {
			log.Fatal(err)
		}

		body = pretty.Pretty(body)
		fmt.Println(string(body))
	}
}
