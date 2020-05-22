package main

// A simple CLI tool for testing the GitHub receiver

import (
	"bytes"
	"context"
	"flag"
	"github.com/whosonfirst/go-webhookd/v2/config"
	"github.com/whosonfirst/go-webhookd/v2/github"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	config_uri := flag.String("config-uri", "", "A valid Go Cloud blob URI where your webhookd config file lives")

	var receiver_name = flag.String("receiver", "", "A valid webhookd config receiver name to test")
	var endpoint = flag.String("endpoint", "", "A valid webhookd (relative) endpoint")
	var file = flag.String("file", "", "The path to a file to test the endpoint with. If empty the webhookd-test-github tool will concatenate arguments passed on the command line.")

	flag.Parse()

	ctx := context.Background()

	if *receiver_name == "" {
		log.Fatal("Missing receiver name")
	}

	cfg, err := config.NewConfigFromURI(ctx, *config_uri)

	if err != nil {
		log.Fatalf("Failed to load config %s, %v", *config_uri, err)
	}

	receiver_uri, err := cfg.GetReceiverConfigByName(*receiver_name)

	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(receiver_uri)

	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	secret := q.Get("secret")

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

	sig, _ := github.GenerateSignature(body, secret)

	client := &http.Client{}

	d_uri, _ := url.Parse(cfg.Daemon)
	d_uri.Path = filepath.Join(d_uri.Path, *endpoint)

	req, err := http.NewRequest("POST", d_uri.String(), bytes.NewBufferString(body))

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
