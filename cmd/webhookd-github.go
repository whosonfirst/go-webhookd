package main

// A simple CLI tool for testing the GitHub receiver

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-webhookd"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")
	var endpoint = flag.String("endpoint", "", "...")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := webhookd.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	body := strings.Join(flag.Args(), " ")

	secret := config.Receiver.Secret

	// sudo make me a function

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(body))

	sum := mac.Sum(nil)
	enc := hex.EncodeToString(sum)

	sig := fmt.Sprintf("sha1=%s", enc)

	// end of sudo make me a function

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

	b, err := ioutil.ReadAll(rsp.Body)
	log.Println(string(b))
}
