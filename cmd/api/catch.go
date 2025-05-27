package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
)

func (app *application) Catch_all_route_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path: ", r.URL.Path)
	hash := r.URL.String()[1:]
	hash_url, err := app.get_url(hash)
	if strings.Contains(err.Error(), "Key not found") {
		msg := fmt.Sprint(hash, " Key not found in the database")
		fmt.Printf(msg)
		fmt.Fprintf(w, msg)
		return
	}
	if err != nil {
		fmt.Println("error occured")
		log.Printf("ERROR: %v\n", err)
		return
	}
	http.Redirect(w, r, hash_url, http.StatusFound)
}

func (app *application) get_url(hash string) (string, error) {
	var hash_url []byte
	db, err := badger.Open(badger.DefaultOptions(app.conf.db_location))
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
		if err != nil {
			return err
		}
		hash_url, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(hash_url), nil
}
