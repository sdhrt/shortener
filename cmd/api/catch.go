package main

import (
	"fmt"
	"net/http"
)

func (app *application) Catch_all_route_handler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.String()[1:]

	data := &stat{}
	data, err := app.read(hash)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}

	if data.url == "" {
		fmt.Fprintln(w, "Couldn't find data associated with hash ", hash)
		return
	}

	app.increment_click(data.hash)

	http.Redirect(w, r, data.url, http.StatusFound)

}
