package main

import (
	"database/sql"
	// _ "github.com/go-sql-driver/mysql"
)

func accessDB() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func search(db *sql.DB) {

}
