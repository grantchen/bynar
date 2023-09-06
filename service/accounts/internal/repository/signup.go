package repository

import (
	"fmt"

	"github.com/google/uuid"
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
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var userID int
	err = tx.QueryRow(
		`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);SELECT LAST_INSERT_ID();`,
		email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1).Scan(&userID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	var organizationID int
	err = tx.QueryRow(
		`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);SELECT LAST_INSERT_ID();`,
		organizationName, vat, organisationCountry, uuid.New().String(), 1, 1,
	).Scan(&organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES (?, ?, ?, ?)`,
		organizationID, uid, userID, 1,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO account_cards (user_payment_gateway_id, user_id, status, is_default, source_id, account_id) VALUES (?, ?, ?, ?, ?, ?)`,
		"", userID, 1, 1, 1, userID,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}
