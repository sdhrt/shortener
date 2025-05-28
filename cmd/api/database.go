package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type stat struct {
	hash         int
	url          string
	clicks       int
	date_created string
	expiry_date  string
}

// create_database is a method of application that is used to create a sqlite database and attach the db connection to application.db
func (app *application) create_database() error {
	var err error
	db, err := sql.Open("sqlite", app.conf.db_location)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS urls(
			hash INTEGER PRIMARY KEY,
			url TEXT NOT NULL,
			clicks INTEGER NOT NULL,
			date_created TEXT NOT NULL,
			expiry_date TEXT NOT NULL
		);`,
	)
	if err != nil {
		return err
	}

	app.db = db
	return nil
}

func (app *application) write(key uint32, value string) error {
	date_now := time.Now().UTC().Format(time.RFC1123)
	date_later := time.Now().UTC().Add(time.Hour * 24 * 30 * 2).Format(time.RFC1123) // set expiry to two months
	result, err := app.db.ExecContext(
		context.Background(),
		`INSERT into urls(hash, url, clicks, date_created, expiry_date) 
		VALUES
		(?,?,?,?,?);`, key, value, 0, date_now, date_later,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 1 {
		fmt.Println("Inserted data")
	}
	return nil
}

func (app *application) read(key string) (*stat, error) {
	data := stat{}
	rows, err := app.db.QueryContext(
		context.Background(),
		`SELECT * FROM urls WHERE hash=(?)`, key,
	)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(&data.hash, &data.url, &data.clicks, &data.date_created, &data.expiry_date); err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func (app *application) view() ([]stat, error) {
	data := []stat{}
	row_stat := stat{}
	rows, err := app.db.QueryContext(
		context.Background(),
		`SELECT * FROM urls`,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&row_stat.hash, &row_stat.url, &row_stat.clicks, &row_stat.date_created, &row_stat.expiry_date); err != nil {
			return nil, err
		}
		data = append(data, row_stat)
	}
	return data, nil
}

// Increment the clicks field in the database using primary key, which is the hash
func (app *application) increment_click(hash int) error {
	result, err := app.db.ExecContext(
		context.Background(),
		"UPDATE urls SET clicks = clicks + 1 WHERE hash=(?)", hash,
	)
	if err != nil {
		fmt.Println("Couldn't increment click count")
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 1 {
		return nil
	}
	return nil
}

func (app *application) test() {
	var err error
	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	res, err := db.ExecContext(
		context.Background(),
		"INSERT INTO test VALUES (3)",
	)
	if err != nil {
		log.Fatalln(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("rows affected: ", rows)
}
