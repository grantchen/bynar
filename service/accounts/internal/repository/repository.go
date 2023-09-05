package repository

import (
	"database/sql"
)

// AccountRepository provides a interface on db level for user
type AccountRepository interface {
	CreateUser(email string) error
}

type accountRepositoryHandler struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepositoryHandler{db}
}
