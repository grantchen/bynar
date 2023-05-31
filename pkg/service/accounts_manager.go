package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/scope"
)

type accountManagerRepository struct {
	db *sql.DB
	accountManagerRepository repository.AccountManagerRepository
}

func NewAccountManagerService(
	db *sql.DB,
	accaccountManagerRepository repository.AccountManagerRepository,
	) AccountManagerService {
	return &accountManagerRepository{
		db: db,
		accountManagerRepository: accaccountManagerRepository,
	}
}

// CheckPermission implements AccountManagerService
func (*accountManagerRepository) CheckPermission(token string) (bool, error) {
	requestScope, err := scope.ResolveFromToken(token)
	if err != nil {
		return false, fmt.Errorf("check permission error: [%w]", err)
	}

	requestScope.
}

// GetNewStringConnection implements AccountManagerService
func (*accountManagerRepository) GetNewStringConnection(token string) (string, error) {
	return "", nil
}
