package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

// CheckUserExists check is email exist in db
func (r *accountRepositoryHandler) CheckUserExists(email string) error {
	var id int
	err := r.db.QueryRow("SELECT id FROM accounts WHERE email = ? AND status=1;", email).Scan(&id)
	if err == nil && id != 0 {
		return errors.New("email already exist")
	}
	return nil
}

// CreateOrganization create the organization when creating user
func (r *accountRepositoryHandler) CreateOrganization(tx *sql.Tx, description, vat, country, uid string, accountID int) (int, string, error) {
	var organizationID int
	organizationUUID := uuid.New().String()
	// check if the organization exists
	err := tx.QueryRow(`SELECT id FROM organizations where vat_number=?;`, vat).Scan(&organizationID)
	if err != nil && err != sql.ErrNoRows {
		return 0, "", err
	}
	if err == sql.ErrNoRows {
		// Insert organization info to db
		res, err := tx.Exec(
			`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);`,
			description, vat, country, organizationUUID, 1, 1,
		)
		if err != nil {
			tx.Rollback()
			return 0, "", err
		}
		newOrganizationID, _ := res.LastInsertId()
		organizationID = int(newOrganizationID)
	}
	// Add relationship between the account and organization
	_, err = tx.Exec(
		`INSERT INTO oraginzation_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES (?, ?, ?, ?)`,
		organizationID, uid, accountID, 1,
	)
	if err != nil {
		tx.Rollback()
		return 0, "", err
	}
	return organizationID, organizationUUID, nil
}

// CreateTenantManagement create the tanant managemant when creating user
func (r *accountRepositoryHandler) CreateTenantManagement(tx *sql.Tx, region string, organizationID int) (string, error) {
	// Check if is allowed to insert
	var tenantID int
	var organizations int
	var organizationsAllowed int
	var tenantUUID string
	// check the tanant exists
	err := tx.QueryRow("SELECT id, organizations, organizations_allowed, tenant_uuid FROM tenants WHERE region = ?", region).Scan(&tenantID, &organizations, &organizationsAllowed, &tenantUUID)
	if err != nil || organizations >= organizationsAllowed {
		logrus.Error("create tenant error: ", err)
		return "", errors.New("not allowed to insert tenants")
	}
	// insert managemant
	_, err = tx.Exec(
		`INSERT INTO tenants_management (organization_id, tenant_id, status, suspended) VALUES (?, ?, ?, ?)`,
		organizationID, tenantID, 1, 0,
	)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	// update the organizations count in tanants
	_, err = tx.Exec(
		`UPDATE tenants SET organizations = ? WHERE id = ?`,
		organizations+1, tenantID,
	)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	// update the tenant_id in organization
	_, err = tx.Exec(
		`UPDATE organizations SET tenant_id = ? WHERE id = ?`,
		tenantID, organizationID,
	)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return tenantUUID, nil
}

func (r *accountRepositoryHandler) CreateCard(tx *sql.Tx, customerID, sourceID string, userID int) error {
	var count int
	// check the card is default
	err := tx.QueryRow(`SELECT COUNT(*) FROM accounts_cards WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		return err
	}
	isDefault := 1
	if count > 0 {
		isDefault = 0
	}
	// insert account card info
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
	// insert user base info
	res, err := tx.Exec(
		`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	userID, _ := res.LastInsertId()
	organizationID, organizationUUID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry, uid, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	err = r.CreateCard(tx, customerID, sourceID, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	tenantUUID, err := r.CreateTenantManagement(tx, state, organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	// create environment
	return r.CreateEnvironment(tenantUUID, organizationUUID)
}

// CreateEnvironment create a new schema in db
func (r *accountRepositoryHandler) CreateEnvironment(tenantUUID, organizationUUID string) error {
	//get tanant mysql connstr from environment
	connStr := os.Getenv(tenantUUID)
	if strings.Contains(connStr, "?") {
		connStr += "&multiStatements=true"
	} else {
		connStr += "?multiStatements=true"
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	// create database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", organizationUUID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("USE `%s`", organizationUUID))
	if err != nil {
		return err
	}
	// create tables
	_, err = db.Exec(model.SQL_TEMPLATE)
	return err
}
