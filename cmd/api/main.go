package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sdhrt/shortener/internals/database"
)

type config struct {
	port        int
	db_location string
	crc_poly    uint
	frontend    string
}

type application struct {
	conf *config
	db   *database.Database
}

func NewApplication(db *database.Database, conf *config) *application {
	return &application{
		db:   db,
		conf: conf,
	}
}

func main() {

	conf := &config{}

	flag.IntVar(&conf.port, "port", 8080, "port for the server")
	flag.StringVar(&conf.db_location, "db_location", "./db.sqlite", "database location for the server")
	flag.UintVar(&conf.crc_poly, "poly", 0xD5828281, "poly for crc")
	flag.StringVar(&conf.frontend, "frontend", "http://localhost:3000", "url of the frontend app, used for cors")
	flag.Parse()

	db, close, err := database.NewDatabase(conf.db_location)
	if err != nil {
		log.Fatalln(err)
	}
	defer close()
	app := NewApplication(db, conf)

	mux := http.NewServeMux()

	mux.HandleFunc("/create", app.Hash_route_handler)
	mux.HandleFunc("/view", app.View_route_handler)
	mux.HandleFunc("/healthcheck", app.Health_check_route_handler)
	mux.HandleFunc("/", app.Catch_all_route_handler)

	handler := app.MiddlewareCors(mux, nil)

	log.Println("LISTENING")
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", app.conf.port), handler))
}
