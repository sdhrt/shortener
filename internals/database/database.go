package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Data struct {
	Hash         int
	Url          string
	Clicks       int
	Date_created string
	Expiry_date  string
}

type Database struct {
	db *sql.DB
}

// NewDatabase is a method of application that is used to create
// a sqlite database and attach the db connection to application.db
func NewDatabase(location string) (*Database, func() error, error) {
	var err error
	db, err := sql.Open("sqlite", location)
	if err != nil {
		return nil, nil, err
	}
	// err = db.Ping()
	// if err != nil {
	// 	return nil, nil, err
	// }

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
		return nil, nil, err
	}

	return &Database{db: db}, db.Close, nil
}

func (database *Database) Write(key uint32, value string) error {
	date_now := time.Now().UTC().Format(time.RFC1123)
	date_later := time.Now().UTC().Add(time.Hour * 24 * 30 * 2).Format(time.RFC1123) // set expiry to two months
	result, err := database.db.ExecContext(
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

func (database *Database) Read(key string) (*Data, error) {
	data := Data{}
	rows, err := database.db.QueryContext(
		context.Background(),
		`SELECT * FROM urls WHERE hash=(?)`, key,
	)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(&data.Hash, &data.Url, &data.Clicks, &data.Date_created, &data.Expiry_date); err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func (database *Database) View() ([]Data, error) {
	data := []Data{}
	row_stat := Data{}
	rows, err := database.db.QueryContext(
		context.Background(),
		`SELECT * FROM urls`,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&row_stat.Hash, &row_stat.Url, &row_stat.Clicks, &row_stat.Date_created, &row_stat.Expiry_date); err != nil {
			return nil, err
		}
		data = append(data, row_stat)
	}
	return data, nil
}

// Increment the clicks field in the database using primary key, which is the hash
func (database *Database) Increment(hash int) error {
	result, err := database.db.ExecContext(
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

func (database *Database) Test() {
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
