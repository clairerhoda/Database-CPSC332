package main

import (
	"database/sql"
	"fmt"
	// "log"
	// "net/http"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "chese21"
	dbname   = "Rational_Room_Reservations"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	// http.HandleFunc("/", department.GETHandler)
	// http.HandleFunc("/insert", department.POSTHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}