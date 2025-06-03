package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// View Route Handler prints all the available hashes from the db to the stdout //
// Additionally, it also prints the host to the Response Writer
func (app *application) View_route_handler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.View()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintln(w, "index,hash,url,clicks,date created,expiry date")

	for i, row := range rows {
		var date_created, expiry_date string
		date, err := time.Parse(time.RFC1123, row.Date_created)
		if err != nil {
			log.Fatalln(err)
		}
		date_created = date.String()
		date, err = time.Parse(time.RFC1123, row.Expiry_date)
		if err != nil {
			log.Fatalln(err)
		}
		expiry_date = date.String()

		frow := fmt.Sprintf("%v,%v,%v,%v,%v,%v", i, row.Hash, row.Url, row.Clicks, date_created, expiry_date)
		fmt.Fprintln(w, frow)
	}
}
