package sql_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// NewConnection - creates db connection and does ping for checking.
// Keep in mind that's reusing the same connection.
func NewConnection(connString string) (*sql.DB, error) {
	// connection already exist so return it
	if db != nil {
		return db, nil
	}

	return InitializeConnection(connString)
}

// InitializeConnection - creates db connection and does ping for checking.
// It doesn't cache connection.
func InitializeConnection(connString string) (_ *sql.DB, err error) {
	db, err = sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("open db connection: [%w], connString: [%s]", err, connString)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: [%w], connString: [%s]", err, connString)
	}

	return db, nil
}

func Conn() *sql.DB {
	return db
}
