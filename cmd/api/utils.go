package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"strings"
)

type envelope map[string]any

func (app *application) OkResponse(
	w http.ResponseWriter, status int, headers http.Header, msg string,
) error {
	maps.Copy(headers, w.Header())

	data := envelope{"message": msg}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)

	return nil
}

func (app *application) ErrorResponse(
	w http.ResponseWriter, status int, headers http.Header, msg string,
) {
	maps.Copy(headers, w.Header())

	data := envelope{"error": msg}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (app *application) MiddlewareCors(next http.Handler, methods []string) http.Handler {
	if methods == nil {
		methods = []string{"GET", "POST"}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.conf.frontend)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
