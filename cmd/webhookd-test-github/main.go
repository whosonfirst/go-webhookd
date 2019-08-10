package main

// A simple CLI tool for testing the GitHub receiver

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/config"
	"github.com/whosonfirst/go-webhookd/github"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")
	var receiver_name = flag.String("receiver", "", "A valid webhookd config receiver name to test")
	var endpoint = flag.String("endpoint", "", "A valid webhookd (relative) endpoint")
	var file = flag.String("file", "", "The path to a file to test the endpoint with. If empty the webhookd-test-github tool will concatenate arguments passed on the command line.")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	if *receiver_name == "" {
		log.Fatal("Missing receiver name")
	}

	wh_config, err := config.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	receiver_config, err := wh_config.GetReceiverConfigByName(*receiver_name)

	if err != nil {
		log.Fatal(err)
	}

	var body string

	if *file != "" {
		stuff, err := ioutil.ReadFile(*file)

		if err != nil {
			log.Fatal(err)
		}

		body = string(stuff)
	} else {
		body = strings.Join(flag.Args(), " ")
	}

	secret := receiver_config.Secret

	sig, _ := github.GenerateSignature(body, secret)

	client := &http.Client{}

	uri := fmt.Sprintf("http://%s:%d%s", wh_config.Daemon.Host, wh_config.Daemon.Port, *endpoint)

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
