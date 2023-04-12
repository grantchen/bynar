package lambda_handler

import (
	"context"
	"database/sql"
)

type ConnectionResolver interface {
	Get(context.Context, string) (*sql.DB, error)
}
