package repository

import (
	"database/sql"
)

// AccountRepository provides a interface on db level for user
type AccountRepository interface {
	CheckUserExists(email string) error
	CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry string) error
}

type accountRepositoryHandler struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepositoryHandler{db}
}
