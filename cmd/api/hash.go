package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type request_body struct {
	Url string `json:"Url"`
}

func (app *application) Hash_route_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method: %v\n", r.Method)
		fmt.Fprintf(w, "Can hash only using POST")
		return
	}

	var user_url request_body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&user_url)

	if err == io.EOF {
		fmt.Fprintln(w, "Please provide url in the body in json format")
		fmt.Fprintf(w, "%v\n", user_url)
		return
	}

	if strings.HasPrefix(fmt.Sprint(err), "json: unknown field") {
		fmt.Fprintln(w, "Please provide url in the body in json format\n{\n\t\"Url\":\"www.example.com\"\n}")
		return
	}

	if err != nil {
		fmt.Fprintf(w, "Failure, please try again later\n")
		return
	}

	req_url, err := url.ParseRequestURI(user_url.Url)
	if err != nil {
		fmt.Fprintf(w, "Please pass valid urls\n")
		return
	}

	// Using CRC
	crc32q := crc32.MakeTable(uint32(app.conf.crc_poly))
	var hashed_url uint32 = app.crc_hash(req_url.String(), crc32q)

	err = app.write(hashed_url, req_url.String())
	if strings.HasPrefix(fmt.Sprint(err), "constraint failed: UNIQUE constraint failed") {
		fmt.Fprintln(w, "This url has already been hashed")
		return
	}
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "Couldn't create hash, try again later")
		return
	}

	fmt.Fprintf(w, "Your new hash is: %v/%v\n", r.Host, hashed_url)
}
