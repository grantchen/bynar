package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Signup check is email exist in db
func (r *accountRepositoryHandler) Signup(email string) error {
	var id int
	err := r.db.QueryRow("SELECT id FROM accounts WHERE email = ?", email).Scan(&id)
	if err == nil && id != 0 {
		return fmt.Errorf("account with email: %s has already exist", email)
	}
	return nil
}

// CreateUser create a new account in db
func (r *accountRepositoryHandler) CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry string) error {
	_, err := r.db.Exec(
		`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);`,
		email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	return nil
}
