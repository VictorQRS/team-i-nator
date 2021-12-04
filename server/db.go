package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
    dbUser = "team"
    dbPass = "team"
    dbName = "team"
)

func SetupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
    db, err := sql.Open("postgres", dbinfo)

    CheckErr(err)

    return db
}