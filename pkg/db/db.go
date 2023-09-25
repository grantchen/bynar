package sql_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// NewConnection - creates db connection and does ping for checking.
// Keep in mind that's reusing the same connection.
func NewConnection(connString string) (_ *sql.DB, err error) {
	// connection already exist so return it
	if db != nil {
		return db, nil
	}

	db, err = InitializeConnection(connString)

	return db, err
}

// InitializeConnection - creates db connection and does ping for checking.
// It doesn't cache connection.
func InitializeConnection(connString string) (_ *sql.DB, err error) {
	uncachedDb, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("open db connection: [%w], connString: [%s]", err, connString)
	}

	if err = uncachedDb.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: [%w], connString: [%s]", err, connString)
	}

	return uncachedDb, nil
}

func Conn() *sql.DB {
	return db
}
