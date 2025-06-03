package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"sdhrt/shortener/internals/helpers"
	"strings"
)

type body struct {
	Url string `json:"Url"`
}

type JsonErr struct {
	Err string `json:"Error"`
}

func (app *application) Hash_route_handler(w http.ResponseWriter, r *http.Request) {
	headers := r.Header.Clone()
	if r.Method != "POST" {
		app.ErrorResponse(w, http.StatusMethodNotAllowed, headers, "Only POST method allowed")
		return
	}

	var user_url body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&user_url)

	if err == io.EOF {
		app.ErrorResponse(w, http.StatusBadRequest, headers, "The format of the body is {\"Url\":\"url_to_be_shortened\"}")
		return
	}

	if strings.HasPrefix(fmt.Sprint(err), "json: unknown field") {
		app.ErrorResponse(w, http.StatusBadRequest, headers, "The format of the body is {\"Url\":\"url_to_be_shortened\"}")
		return
	}

	if err != nil {
		app.ErrorResponse(w, http.StatusInternalServerError, headers, "Internal server error")
		return
	}

	req_url, err := url.ParseRequestURI(user_url.Url)
	if err != nil {
		app.ErrorResponse(w, http.StatusBadRequest, headers, "Please pass valid urls only")
		return
	}

	// Using CRC
	crc32q := crc32.MakeTable(uint32(app.conf.crc_poly))
	var hashed_url uint32 = helpers.Crc_hash(req_url.String(), crc32q)

	err = app.db.Write(hashed_url, req_url.String())
	if strings.HasPrefix(fmt.Sprint(err), "constraint failed: UNIQUE constraint failed") {
		app.ErrorResponse(w, http.StatusConflict, r.Header.Clone(), "link has already been shortened")
		return
	}
	if err != nil {
		app.ErrorResponse(w, http.StatusInternalServerError, headers, "Internal server error, please try again later")
		return
	}
	app.OkResponse(w, http.StatusOK, headers, fmt.Sprintf("%v/%v", r.Host, hashed_url))
}
