package main

import (
	"encoding/json"
	"flag"
	"fmt"
	ucd "github.com/cooperhewitt/go-ucd"
	"net/http"
	"strings"
)

type UCDResponse struct {
	Chars []ucd.UCDName
}

func string(w http.ResponseWriter, r *http.Request) {

	txt := r.FormValue("text")
	txt = strings.Trim(txt, " ")

	chars := ucd.NamesForString(txt)

	rsp := UCDResponse{chars}
	send(w, r, rsp)
}

/*
func char(w http.ResponseWriter, r *http.Request) {

	txt := r.FormValue("text")
	char := ucd.Name(txt)

	chars := make([]ucd.UCDName, 1)
	chars[0] = char

	rsp := UCDResponse{chars}
	send(w, r, rsp)
}
*/

func send(w http.ResponseWriter, r *http.Request, rsp UCDResponse) {

	accept := r.Header.Get("Accept")

	if accept == "text/plain" {
		send_text(w, rsp)
	} else {
		send_json(w, rsp)
	}

}

func send_text(w http.ResponseWriter, rsp UCDResponse) {

	w.Header().Set("Content-Type", "text/plain")

	for _, char := range rsp.Chars {
		fmt.Fprintln(w, char.String())
	}
}

func send_json(w http.ResponseWriter, rsp UCDResponse) {

	js, err := json.Marshal(rsp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	host := flag.String("host", "localhost", "host")
	port := flag.Int("port", 8080, "port")

	flag.Parse()

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	fmt.Printf("listening on %s\n", endpoint)

	http.HandleFunc("/", string)
	http.ListenAndServe(endpoint, nil)
}
