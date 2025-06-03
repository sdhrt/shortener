package main

import (
	"fmt"
	"net/http"
	"sdhrt/shortener/internals/database"
)

func (app *application) Catch_all_route_handler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.String()[1:]

	data := &database.Data{}
	data, err := app.db.Read(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error ", err)
		return
	}

	if data.Url == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Couldn't find data associated with hash", hash)
		return
	}

	app.db.Increment(data.Hash)
	http.Redirect(w, r, data.Url, http.StatusFound)
}
