package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// AccountService is a interface which provide helper methods to access account related operations
type AccountService interface {
	Signup(email string) error
	ConfirmEmail(email, code string) (int, error)
	VerifyCard(id, token, email, name string) error
	CreateUser(email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry string) (string, error)
}

type accountServiceHandler struct {
	ar              repository.AccountRepository
	authProvider    gip.AuthProvider
	paymentProvider checkout.PaymentClient
}

// NewAccountService initiates the account service object
func NewAccountService(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient) AccountService {
	ar := repository.NewAccountRepository(db)
	return &accountServiceHandler{ar, authProvider, paymentProvider}
}
