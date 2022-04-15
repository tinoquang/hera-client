package main

import (
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("hera", "localhost:10001")
	if err != nil {
		log.Fatal(err)
	}

}
