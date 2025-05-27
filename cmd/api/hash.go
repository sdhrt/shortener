package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
)

const Size int = 8

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
	if strings.Contains(fmt.Sprint(err), "json: unknown field") {
		fmt.Fprintln(w, "Please provide url in the body in json format\n{\n\t\"Url\":\"www.facebook.com\"\n}")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "failure, please try again later\n")
		log.Fatalf("ERROR: %v", err)
		return
	}

	req_url, err := url.ParseRequestURI(user_url.Url)
	if err != nil {
		fmt.Fprintf(w, "Please pass valid urls\n")
		log.Fatalf("ERROR: %v", err)
		return
	}

	// Omit queries when hashing
	to_hash_url := []byte{}
	to_hash_url = (fmt.Appendf(to_hash_url, "%s://%s%s", req_url.Scheme, req_url.Host, req_url.Path))
	fmt.Printf("to_hash_url: %v\n", string(to_hash_url))

	// Using CRC
	crc32q := crc32.MakeTable(uint32(app.conf.crc_poly))
	hashed_url := app.get_hash(to_hash_url, crc32q)
	fmt.Printf("hashed_url: %v\n", hashed_url)

	db, err := badger.Open(badger.DefaultOptions(app.conf.db_location))
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	defer db.Close()

	err = app.store_url(to_hash_url, []byte(strconv.FormatUint(uint64(hashed_url), 10)), db)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	fmt.Fprintf(w, "%v/%v\n", r.Host, hashed_url)
}

func (app *application) store_url(to_hash_url []byte, hash []byte, db *badger.DB) error {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	fmt.Printf("STORING %v to %v\n", string(to_hash_url), string(hash))

	err := txn.Set(hash, to_hash_url)
	if err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (app *application) get_hash(req_url []byte, table *crc32.Table) uint32 {
	hashed_url := crc32.Checksum(req_url, table)
	return hashed_url
}
