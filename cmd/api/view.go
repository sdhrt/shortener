package main

import (
	"fmt"
	"log"
	"net/http"
	"text/tabwriter"
	"time"
)

// View Route Handler prints all the available hashes from the db to the stdout //
// Additionally, it also prints the host to the Response Writer
func (app *application) View_route_handler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.view()
	if err != nil {
		log.Fatalln(err)
	}

	tabwrite := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	fmt.Fprintln(tabwrite, "index\thash\turl\tclicks\tdate created\texpiry date")

	for i, row := range rows {
		var date_created, expiry_date string
		date, err := time.Parse(time.RFC1123, row.date_created)
		if err != nil {
			log.Fatalln(err)
		}
		date_created = date.String()
		date, err = time.Parse(time.RFC1123, row.expiry_date)
		if err != nil {
			log.Fatalln(err)
		}
		expiry_date = date.String()

		frow := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v", i, row.hash, row.url, row.clicks, date_created, expiry_date)
		fmt.Fprintln(tabwrite, frow)
	}

	tabwrite.Flush()
}
