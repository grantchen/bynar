package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CheckUserExists check is email exist in db
func (r *accountRepositoryHandler) CheckUserExists(email string) error {
	var id int
	err := r.db.QueryRow("SELECT id FROM accounts WHERE email = ? AND status=1 AND archived IS FALSE;", email).Scan(&id)
	if err == nil && id != 0 {
		return fmt.Errorf("account with email: %s has already exist", email)
	}
	return nil
}

func (r *accountRepositoryHandler) CreateOrganization(tx *sql.Tx, description, vat, country string) (int, error) {
	var organizationID int
	err := tx.QueryRow(`SELECT id FROM organizations where vat_number=?;`, vat).Scan(&organizationID)
	if err == nil {
		return organizationID, nil
	}
	if err != sql.ErrNoRows {
		logrus.Errorf("CheckOrganizationExists: error: %v", err)
		return 0, err
	}
	err = tx.QueryRow(
		`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);SELECT LAST_INSERT_ID();`,
		description, vat, country, uuid.New().String(), 1, 1,
	).Scan(&organizationID)
	if err != nil {
		logrus.Errorf("CreateOrganization: error: %v", err)
		tx.Rollback()
		return 0, err
	}
	return 0, err
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
	organizationID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
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
		userID, userID, 1, 1, 1, userID,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	var tanantID int
	var organizations int
	var organizationsAllowed int
	err = r.db.QueryRow("SELECT id, oraginzations, organizations_allowed FROM tenants WHERE region = ?", state).Scan(&tanantID, &organizations, organizationsAllowed)
	if err != nil || organizations >= organizationsAllowed {
		return errors.New("not allowed to insert tenants")
	}
	_, err = tx.Exec(
		`INSERT INTO tenants_management (organization_id, tenant_id, status, suspended) VALUES (?, ?, ?, ?)`,
		organizationID, tanantID, 1, 0,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}
