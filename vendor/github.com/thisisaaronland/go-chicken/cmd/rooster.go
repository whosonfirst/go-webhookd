package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/thisisaaronland/go-chicken"
	"github.com/whosonfirst/go-sanitize"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 1280, "The port number to listen for requests on")	// because ROOSTER is U+1F413 (or 128019)

	flag.Parse()

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		opts := sanitize.DefaultOptions()
		opts.AllowNewlines = true
		
		query := req.URL.Query()

		lang := "zxx"
		clucking := false

		_, has_lang := query["language"]

		if has_lang {

			language, err := sanitize.SanitizeString(query.Get("language"), opts)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			lang = language
		}

		_, clucking = query["clucking"]

		ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		input, err := sanitize.SanitizeString(string(body), opts)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		output := ch.TextToChicken(input)

		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Content-Type", "text/chicken")

		rsp.Write([]byte(output))
	}

	ch_handler := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Content-Type", "text/chicken")
		rsp.Write([]byte("🐔"))
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("%s on %s\n", "🐓", endpoint)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/chicken", ch_handler)

	err := gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
