package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func accessDB() {
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func insertData(db *sql.DB) {

}
