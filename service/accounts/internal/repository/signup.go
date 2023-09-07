package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CheckUserExists check is email exist in db
func (r *accountRepositoryHandler) CheckUserExists(email string) error {
	var id int
	err := r.db.QueryRow("SELECT id FROM accounts WHERE email = ? AND status=1;", email).Scan(&id)
	if err == nil && id != 0 {
		return errors.New("username already exist")
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
	res, err := tx.Exec(
		`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);`,
		description, vat, country, uuid.New().String(), 1, 1,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	newOrganizationID, _ := res.LastInsertId()
	organizationID = int(newOrganizationID)
	// Add relationship between the account and organization
	_, err = tx.Exec(
		`INSERT INTO oraginzation_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES (?, ?, ?, ?)`,
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
	err := tx.QueryRow("SELECT id, organizations, organizations_allowed FROM tenants WHERE region = ?", region).Scan(&tanantID, &organizations, &organizationsAllowed)
	if err != nil || organizations >= organizationsAllowed {
		logrus.Error("create tenant error: ", err)
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
	_, err = tx.Exec(
		`UPDATE tenants SET organizations = ? WHERE id = ?`,
		organizations+1, tanantID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *accountRepositoryHandler) CreateCard(tx *sql.Tx, customerID, sourceID string, userID int) error {
	var count int
	err := tx.QueryRow(`SELECT COUNT(*) FROM accounts_cards WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		return err
	}
	isDefault := 1
	if count > 0 {
		isDefault = 0
	}
	_, err = tx.Exec(
		`INSERT INTO accounts_cards (user_payment_gateway_id, user_id, status, is_default, source_id, account_id) VALUES (?, ?, ?, ?, ?, ?)`,
		customerID, userID, 1, isDefault, sourceID, userID,
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
	res, err := tx.Exec(
		`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	userID, _ := res.LastInsertId()
	organizationID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry, uid, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	err = r.CreateCard(tx, customerID, sourceID, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	err = r.CreateTenantManagement(tx, state, organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	return tx.Commit()
}
