package lambda_handler

import (
	"database/sql"
)

type ConnectionResolver interface {
	Get(string) (*sql.DB, error)
}
