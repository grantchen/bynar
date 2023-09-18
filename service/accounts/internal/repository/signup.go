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

// CreateOrganization create the organization when creating user
func (r *accountRepositoryHandler) CreateOrganization(tx *sql.Tx, description, vat, country, uid string, accountID int) (int, string, error) {
	var organizationID int
	organizationUUID := uuid.New().String()
	// check if the organization exists
	err := tx.QueryRow(`SELECT id, organization_uuid FROM organizations where vat_number=?;`, vat).Scan(&organizationID, &organizationUUID)
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("select organization error ", err.Error())
		return 0, "", fmt.Errorf("organization if vat_number=%s not exist", vat)
	}
	if err == sql.ErrNoRows {
		// Insert organization info to db
		res, err := tx.Exec(
			`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?);`,
			description, vat, country, organizationUUID, 1, 1,
		)
		if err != nil {
			tx.Rollback()
			logrus.Error("insert organization error ", err.Error())
			return 0, "", errors.New("insert organization failed")
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
		logrus.Error("insert organization_accounts error ", err.Error())
		return 0, "", errors.New("insert organization_accounts failed")
	}
	return organizationID, organizationUUID, nil
}

// CreateTenantManagement create the tanant managemant when creating user
func (r *accountRepositoryHandler) CreateTenantManagement(tx *sql.Tx, region string, organizationID int) (string, int, error) {
	// Check if is allowed to insert
	var tenantID int
	var organizations int
	var organizationsAllowed int
	var tenantUUID string
	var status bool
	// check the tanant exists
	err := tx.QueryRow("SELECT id, organizations, organizations_allowed, tenant_uuid, status FROM tenants WHERE region = ?", region).Scan(&tenantID, &organizations, &organizationsAllowed, &tenantUUID, &status)
	if err != nil {
		logrus.Error("select tenant error: ", err)
		return "", 0, fmt.Errorf("tenants of region %s not exist", region)
	}
	if !status {
		return "", 0, fmt.Errorf("tenant %s cannot be selected ", region)
	}
	if organizations >= organizationsAllowed {
		return "", 0, fmt.Errorf("tenant %s is full", region)
	}

	var tenantManagentID int
	err = tx.QueryRow("SELECT id FROM tenants_management WHERE organization_id = ? AND tenant_id = ?", organizationID, tenantID).Scan(&tenantManagentID)
	if err != nil {
		logrus.Error("select tenants_management error: ", err)
		return "", 0, errors.New("tenants_management count error")
	}
	if tenantManagentID == 0 {
		// insert managemant
		res, err := tx.Exec(
			`INSERT INTO tenants_management (organization_id, tenant_id, status, suspended) VALUES (?, ?, ?, ?)`,
			organizationID, tenantID, 1, 0,
		)
		if err != nil {
			tx.Rollback()
			logrus.Error("insert tenants_management error ", err.Error())
			return "", 0, errors.New("insert tenants_management failed")
		}
		newTenantManagentID, _ := res.LastInsertId()
		tenantManagentID = int(newTenantManagentID)
	}

	// update the organizations count in tanants
	_, err = tx.Exec(
		`UPDATE tenants SET organizations = ? WHERE id = ?`,
		organizations+1, tenantID,
	)
	if err != nil {
		tx.Rollback()
		logrus.Error("update tenants error ", err.Error())
		return "", 0, errors.New("update tenants failed")
	}
	// update the tenant_id in organization
	_, err = tx.Exec(
		`UPDATE organizations SET tenant_id = ? WHERE id = ?`,
		tenantID, organizationID,
	)
	if err != nil {
		tx.Rollback()
		logrus.Error("update organizations error ", err.Error())
		return "", 0, errors.New("update organizations failed")
	}
	return tenantUUID, int(tenantManagentID), nil
}

func (r *accountRepositoryHandler) CreateCard(tx *sql.Tx, customerID, sourceID string, userID int) error {
	var count int
	// check the card is default
	err := tx.QueryRow(`SELECT COUNT(*) FROM accounts_cards WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		logrus.Error("select cards error ", err.Error())
		return errors.New("select cards failed")
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
		logrus.Error("insert cards error ", err.Error())
		return errors.New("insert cards failed")
	}
	return nil
}

// CreateUser create a new account in db
func (r *accountRepositoryHandler) CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	// insert user base info
	res, err := tx.Exec(
		`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1,
	)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, errors.New("insert user failed")
	}
	userID, _ := res.LastInsertId()
	organizationID, organizationUUID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry, uid, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	err = r.CreateCard(tx, customerID, sourceID, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	tenantUUID, tenantManagentID, err := r.CreateTenantManagement(tx, country, organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	// create environment
	return 1, r.CreateEnvironment(tenantUUID, organizationUUID, tenantManagentID, int(userID), email, fullName, phoneNumber)
}

func (r *accountRepositoryHandler) SetStatusToZeroIfEnvFailed(userID, tenantManagentID int) {
	if _, err := r.db.Exec(`UPDATE accounts SET status=? WHERE id=?`, 0, userID); err != nil {
		logrus.Error("create environment failed, update account status to 0 error: ", err.Error())
	}
	if _, err := r.db.Exec(`UPDATE tenants_management SET status=? WHERE id=?`, 0, tenantManagentID); err != nil {
		logrus.Error("create environment failed, update tenants_management status to 0 error: ", err.Error())
	}
}

// CreateEnvironment create a new schema in db
func (r *accountRepositoryHandler) CreateEnvironment(tenantUUID, organizationUUID string, tenantManagentID, userID int, email, fullName, phoneNumber string) error {
	//get tanant mysql connstr from environment
	connStr := os.Getenv(tenantUUID)
	if len(connStr) == 0 {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return errors.New("the tenant mysql connstr is not set")
	}
	if strings.Contains(connStr, "?") {
		connStr += "&multiStatements=true"
	} else {
		connStr += "?multiStatements=true"
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return errors.New("open mysql failed")
	}
	var name string
	err = db.QueryRow(`SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '` + organizationUUID + `'`).Scan(&name)
	if err == nil {
		logrus.Info(organizationUUID, " schema exists")
	}
	if name == "" {
		// create database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", organizationUUID))
		if err != nil {
			r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
			return errors.New("create database failed")
		}
	}
	_, err = db.Exec(fmt.Sprintf("USE `%s`", organizationUUID))
	if err != nil {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return errors.New("use database failed")
	}
	if name == "" {
		// create tables
		_, err = db.Exec(model.SQL_TEMPLATE)
		if err != nil {
			r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
			return errors.New("create tables failed")
		}
	}
	// create user
	_, err = db.Exec(
		`INSERT INTO users (email, full_name, phone, status, language_preference, policy_id, theme) VALUES (?,?,?,?,?,?,?)`,
		email, fullName, phoneNumber, 1, "en", 1, "system",
	)
	if err != nil {
		logrus.Error("insert user error ", err.Error())
		return errors.New("insert users failed")
	}
	return nil
}
