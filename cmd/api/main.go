package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type config struct {
	port        int
	db_location string
	crc_poly    uint
}

type application struct {
	conf config
	db   *sql.DB
}

func main() {
	app := application{
		conf: config{},
	}

	flag.IntVar(&app.conf.port, "port", 8080, "port for the server")
	flag.StringVar(&app.conf.db_location, "db_location", "./db.sqlite", "database location for the server")
	flag.UintVar(&app.conf.crc_poly, "poly", 0xD5828281, "poly for crc")

	flag.Parse()

	err := app.create_database()
	if err != nil {
		log.Fatalf("Error while creating database \n%v\n", err)
	}
	defer app.db.Close()

	// app.test()

	http.HandleFunc("/create", app.Hash_route_handler)
	http.HandleFunc("/view", app.View_route_handler)
	http.HandleFunc("/", app.Catch_all_route_handler)

	fmt.Println("Listening on port", app.conf.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.conf.port), nil))
}
