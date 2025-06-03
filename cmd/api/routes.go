package main

import "net/http"

func (app *application) setupRoutes() {
	http.HandleFunc("/create", app.Hash_route_handler)
	http.HandleFunc("/view", app.View_route_handler)
	http.HandleFunc("/", app.Catch_all_route_handler)
}
