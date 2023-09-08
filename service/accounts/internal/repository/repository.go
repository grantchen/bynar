package repository

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

// AccountRepository provides a interface on db level for user
type AccountRepository interface {
	CheckUserExists(email string) error
	CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID string) error
	SelectSignInColumns(email string) (*model.SignIn, error)
}

type accountRepositoryHandler struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepositoryHandler{db}
}
