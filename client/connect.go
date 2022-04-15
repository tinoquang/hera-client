package client

import (
	"database/sql"
	"time"
)

type connOpts func(*sql.DB)

func setDefaultConnOpts(db *sql.DB) {
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(75)
	db.SetConnMaxIdleTime(30 * time.Second)
}

// Connect creates and returns a new database connection pool
// using hera driver
// Default options:
// 	- 1. max open connections: 100
// 	- 2. max idle connections: 75
// 	- 3. max connection idle time: 30 seconds
func Connect(heraURL string, optsFunc ...connOpts) (DBConn, error) {
	conn, err := sql.Open("hera", heraURL)
	if err != nil {
		return nil, err
	}

	// configure connection options
	setDefaultConnOpts(conn)
	for _, o := range optsFunc {
		o(conn)
	}

	return newClient(conn), nil
}
