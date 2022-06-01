// webhookd-flatten-config is a command line tool for "flattening" a webhookd configuration file in to a string.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/v3/config"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {

	config_path := flag.String("config", "", "The path your webhookd config file")
	constvar := flag.Bool("constvar", false, "A boolean flag indicating flattened config should be encoded as a gocloud.dev/runtimevar 'constant://' URI.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "webhookd-flatten-config is a command line tool for \"flattening\" a webhookd configuration file in to a string.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s /path/to/config.json\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	abs_path, err := filepath.Abs(*config_path)

	if err != nil {
		log.Fatalf("Failed to derive absolute path for '%s', %v", *config_path, err)
	}

	_, err = os.Stat(abs_path)

	if err != nil {
		log.Fatalf("Failed to stat '%s', %v", abs_path, err)
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", abs_path, err)
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		log.Fatalf("Failed to read '%s', %v", abs_path, err)
	}

	var cfg config.WebhookConfig

	err = json.Unmarshal(body, &cfg)

	if err != nil {
		log.Fatalf("Failed to decode config, %v", err)
	}

	body, err = json.Marshal(cfg)

	if err != nil {
		log.Fatalf("Failed to encode config, %v", err)
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
