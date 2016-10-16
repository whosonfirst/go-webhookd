package main

// A simple CLI tool for testing the GitHub receiver

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"github.com/whosonfirst/go-whosonfirst-webhookd/github"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")
	var receiver_name = flag.String("receiver", "", "...")
	var endpoint = flag.String("endpoint", "", "...")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	if *receiver_name == "" {
		log.Fatal("Missing receiver name")
	}

	config, err := webhookd.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	receiver_config, err := config.GetReceiverConfigByName(*receiver_name)

	if err != nil {
		log.Fatal(err)
	}

	body := strings.Join(flag.Args(), " ")

	secret := receiver_config.Secret

	sig, _ := github.GenerateSignature(body, secret)

	client := &http.Client{}

	uri := fmt.Sprintf("http://%s:%d%s", config.Daemon.Host, config.Daemon.Port, *endpoint)

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(body))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-GitHub-Event", "debug")
	req.Header.Set("X-Hub-Signature", sig)

	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	rsp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if rsp.StatusCode != 200 {
		log.Fatal(rsp.Status)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	log.Println(rsp.Status, string(b))
}
