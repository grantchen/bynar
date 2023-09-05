package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
)

// AccountService is a interface which provide helper methods to access account related operations
type AccountService interface {
	CreateUser(email string) (string, error)
}

type accountServiceHandler struct {
	ar repository.AccountRepository
}

// NewAccountService initiates the account service object
func NewUserService(db *sql.DB) AccountService {
	ar := repository.NewAccountRepository(db)
	return &accountServiceHandler{ar}
}
