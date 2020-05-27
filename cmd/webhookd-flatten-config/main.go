package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {

	config_path := flag.String("config", "", "The path your webhookd config file")
	constvar := flag.Bool("constvar", false, "...")

	flag.Parse()

	abs_path, err := filepath.Abs(*config_path)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(abs_path)

	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		log.Fatal(err)
	}

	var cfg config.WebhookConfig

	err = json.Unmarshal(body, &cfg)

	if err != nil {
		log.Fatal(err)
	}

	body, err = json.Marshal(cfg)

	if err != nil {
		log.Fatal(err)
	}

	str_body := string(body)

	if *constvar {

		q := url.Values{}
		q.Set("decoder", "string")
		q.Set("val", str_body)

		str_body = fmt.Sprintf("constant://?%s", q.Encode())
	}

	fmt.Println(str_body)
}
