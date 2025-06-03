package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) Health_check_route_handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: "OK"})
}
