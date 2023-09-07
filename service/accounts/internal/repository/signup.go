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
	err := r.db.QueryRow("SELECT id FROM accounts WHERE email = ? AND status=1;", email).Scan(&id)
	if err == nil && id != 0 {
		return fmt.Errorf("account with email: %s has already exist", email)
	}
	return nil
}

func (r *accountRepositoryHandler) CreateOrganization(tx *sql.Tx, description, vat, country, uid string, accountID int) (int, error) {
	var organizationID int
	err := tx.QueryRow(`SELECT id FROM organizations where vat_number=?;`, vat).Scan(&organizationID)
	if err == nil {
		return organizationID, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	// Insert organization info to db
	err = tx.QueryRow(
		`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);SELECT LAST_INSERT_ID();`,
		description, vat, country, uuid.New().String(), 1, 1,
	).Scan(&organizationID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Add relationship between the account and organization
	_, err = tx.Exec(
		`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES (?, ?, ?, ?)`,
		organizationID, uid, accountID, 1,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return organizationID, nil
}

func (r *accountRepositoryHandler) CreateTenantManagement(tx *sql.Tx, region string, organizationID int) error {
	// Check if is allowed to insert
	var tanantID int
	var organizations int
	var organizationsAllowed int
	err := r.db.QueryRow("SELECT id, oraginzations, organizations_allowed FROM tenants WHERE region = ?", region).Scan(&tanantID, &organizations, organizationsAllowed)
	if err != nil || organizations >= organizationsAllowed {
		return errors.New("not allowed to insert tenants")
	}
	_, err = tx.Exec(
		`INSERT INTO tenants_management (organization_id, tenant_id, status, suspended) VALUES (?, ?, ?, ?)`,
		organizationID, tanantID, 1, 0,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// CreateUser create a new account in db
func (r *accountRepositoryHandler) CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID string) error {
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
	organizationID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry, uid, userID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO account_cards (user_payment_gateway_id, user_id, status, is_default, source_id, account_id) VALUES (?, ?, ?, ?, ?, ?)`,
		customerID, sourceID, 1, 1, 1, userID,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		tx.Rollback()
		return err
	}
	err = r.CreateTenantManagement(tx, state, organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	return nil
}
