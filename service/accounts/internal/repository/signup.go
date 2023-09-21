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
func (r *accountRepositoryHandler) CreateOrganization(tx *sql.Tx, description, vat, country, uid string, accountID int) (int, string, int, error) {
	var organizationID int
	organizationUUID := uuid.New().String()
	// check if the organization exists
	stmt, err := tx.Prepare(`SELECT id, organization_uuid FROM organizations where vat_number=?`)
	if err != nil {
		return 0, "", 0, err
	}
	err = stmt.QueryRow(vat).Scan(&organizationID, &organizationUUID)
	stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("select organization error ", err.Error())
		return 0, "", 0, fmt.Errorf("select organization of vat_number=%s error", vat)
	}
	if err == sql.ErrNoRows {
		// Insert organization info to db
		stmt, err := tx.Prepare(`INSERT INTO organizations (description, vat_number, country, organization_uuid, status, verified) VALUES (?, ?, ?, ?, ?, ?)`)
		if err != nil {
			return 0, "", 0, err
		}
		res, err := stmt.Exec(description, vat, country, organizationUUID, 1, 1)
		stmt.Close()
		if err != nil {
			tx.Rollback()
			logrus.Error("insert organization error ", err.Error())
			return 0, "", 0, errors.New("insert organization failed")
		}
		newOrganizationID, _ := res.LastInsertId()
		organizationID = int(newOrganizationID)
	}
	// Add relationship between the account and organization
	stmt, err = tx.Prepare(`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, "", 0, err
	}
	res, err := stmt.Exec(organizationID, uid, accountID, 1)
	stmt.Close()
	if err != nil {
		tx.Rollback()
		logrus.Error("insert organization_accounts error ", err.Error())
		return 0, "", 0, errors.New("insert organization_accounts failed")
	}
	organizationManagentID, _ := res.LastInsertId()
	return organizationID, organizationUUID, int(organizationManagentID), nil
}

// CreateTenantManagement create the tanant managemant when creating user
func (r *accountRepositoryHandler) CreateTenantManagement(tx *sql.Tx, tenantCode string, organizationID int) (string, int, error) {
	// Check if is allowed to insert
	var tenantID int
	var organizations int
	var organizationsAllowed int
	var tenantUUID string
	var status bool
	// check the tanant exists
	stmt, err := tx.Prepare("SELECT id, organizations, organizations_allowed, tenant_uuid, status FROM tenants WHERE code = ?")
	if err != nil {
		return "", 0, err
	}
	err = stmt.QueryRow(tenantCode).Scan(&tenantID, &organizations, &organizationsAllowed, &tenantUUID, &status)
	stmt.Close()
	if err != nil {
		logrus.Error("select tenant error: ", err)
		return "", 0, fmt.Errorf("tenants of code %s not exist", tenantCode)
	}
	if !status {
		return "", 0, fmt.Errorf("tenant of code %s cannot be selected ", tenantCode)
	}
	if organizations >= organizationsAllowed {
		return "", 0, fmt.Errorf("tenant of code %s is full", tenantCode)
	}

	var tenantManagentID int
	stmt, err = tx.Prepare("SELECT id FROM tenants_management WHERE organization_id = ? AND tenant_id = ?")
	if err != nil {
		return "", 0, err
	}
	err = stmt.QueryRow(organizationID, tenantID).Scan(&tenantManagentID)
	stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("select tenants_management error ", err.Error())
		return "", 0, errors.New("select tenants_management failed")
	}
	if err == sql.ErrNoRows {
		// insert managemant
		stmt, err = tx.Prepare(`INSERT INTO tenants_management (organization_id, tenant_id, status, suspended) VALUES (?, ?, ?, ?)`)
		if err != nil {
			return "", 0, err
		}
		res, err := stmt.Exec(organizationID, tenantID, 1, 0)
		stmt.Close()
		if err != nil {
			tx.Rollback()
			logrus.Error("insert tenants_management error ", err.Error())
			return "", 0, errors.New("insert tenants_management failed")
		}
		newTenantManagentID, _ := res.LastInsertId()
		tenantManagentID = int(newTenantManagentID)
	}

	// update the organizations count in tanants
	stmt, err = tx.Prepare(`UPDATE tenants SET organizations = ? WHERE id = ?`)
	if err != nil {
		return "", 0, err
	}
	_, err = stmt.Exec(organizations+1, tenantID)
	stmt.Close()
	if err != nil {
		tx.Rollback()
		logrus.Error("update tenants error ", err.Error())
		return "", 0, errors.New("update tenants failed")
	}
	// update the tenant_id in organization
	stmt, err = tx.Prepare(`UPDATE organizations SET tenant_id = ? WHERE id = ?`)
	if err != nil {
		return "", 0, err
	}
	_, err = stmt.Exec(tenantID, organizationID)
	stmt.Close()
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
	stmt, err := tx.Prepare(`SELECT COUNT(*) FROM accounts_cards WHERE user_id = ?`)
	if err != nil {
		return err
	}
	err = stmt.QueryRow(userID).Scan(&count)
	stmt.Close()
	if err != nil {
		logrus.Error("select cards error ", err.Error())
		return errors.New("select cards failed")
	}
	isDefault := 1
	if count > 0 {
		isDefault = 0
	}
	// insert account card info
	stmt, err = tx.Prepare(`INSERT INTO accounts_cards (user_payment_gateway_id, user_id, status, is_default, source_id, account_id) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(customerID, userID, 1, isDefault, sourceID, userID)
	stmt.Close()
	if err != nil {
		tx.Rollback()
		logrus.Error("insert cards error ", err.Error())
		return errors.New("insert cards failed")
	}
	return nil
}

// CreateUser create a new account in db
func (r *accountRepositoryHandler) CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID, tenantCode string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	// insert user base info
	stmt, err := tx.Prepare(`INSERT INTO accounts (email, full_name, address, address_2, phone, city, postal_code, country, state, status, uid, org_id, verified) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(email, fullName, addressLine, addressLine2, phoneNumber, city, postalCode, country, state, 1, uid, nil, 1)
	stmt.Close()
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, errors.New("insert user failed")
	}
	userID, _ := res.LastInsertId()
	organizationID, organizationUUID, organizationManagentID, err := r.CreateOrganization(tx, organizationName, vat, organisationCountry, uid, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	err = r.CreateCard(tx, customerID, sourceID, int(userID))
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	tenantUUID, tenantManagentID, err := r.CreateTenantManagement(tx, tenantCode, organizationID)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	// create environment
	id, err := r.CreateEnvironment(tenantUUID, organizationUUID, tenantManagentID, int(userID), email, fullName, phoneNumber)
	if err != nil {
		return 1, err
	}
	stmt, err = r.db.Prepare(`UPDATE organization_accounts SET organization_user_id = ? WHERE id = ?`)
	if err != nil {

		return 1, err
	}
	_, err = stmt.Exec(id, organizationManagentID)
	stmt.Close()
	return 1, err
}

func (r *accountRepositoryHandler) SetStatusToZeroIfEnvFailed(userID, tenantManagentID int) {
	stmt, err := r.db.Prepare(`UPDATE accounts SET status=? WHERE id=?`)
	if err != nil {
		logrus.Error(err)
		return
	}
	if _, err := stmt.Exec(0, userID); err != nil {
		logrus.Error("create environment failed, update account status to 0 error: ", err.Error())
	}
	stmt.Close()
	stmt, err = r.db.Prepare(`UPDATE tenants_management SET status=? WHERE id=?`)
	if err != nil {
		logrus.Error(err)
		return
	}
	if _, err := stmt.Exec(0, tenantManagentID); err != nil {
		logrus.Error("create environment failed, update tenants_management status to 0 error: ", err.Error())
	}
	stmt.Close()
}

// CreateEnvironment create a new schema in db
func (r *accountRepositoryHandler) CreateEnvironment(tenantUUID, organizationUUID string, tenantManagentID, userID int, email, fullName, phoneNumber string) (int, error) {
	//get tanant mysql connstr from environment
	connStr := os.Getenv(tenantUUID)
	if len(connStr) == 0 {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return 0, errors.New("the tenant mysql connstr is not set")
	}
	if strings.Contains(connStr, "?") {
		connStr += "&multiStatements=true"
	} else {
		connStr += "?multiStatements=true"
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return 0, errors.New("open mysql failed")
	}
	var name string
	stmt, err := db.Prepare(`SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '` + organizationUUID + `'`)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	err = stmt.QueryRow().Scan(&name)
	if err == nil {
		logrus.Info(organizationUUID, " schema exists")
	}
	stmt.Close()
	if name == "" {
		// create database
		stmt, err := db.Prepare(fmt.Sprintf("CREATE DATABASE `%s`", organizationUUID))
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		_, err = stmt.Exec()
		stmt.Close()
		if err != nil {
			r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
			return 0, errors.New("create database failed")
		}
	}
	stmt, err = db.Prepare(fmt.Sprintf("USE `%s`", organizationUUID))
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	_, err = stmt.Exec()
	stmt.Close()
	if err != nil {
		r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
		return 0, errors.New("use database failed")
	}
	if name == "" {
		// create tables
		stmt, err = db.Prepare(model.SQL_TEMPLATE)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		_, err = stmt.Exec()
		stmt.Close()
		if err != nil {
			r.SetStatusToZeroIfEnvFailed(userID, tenantManagentID)
			return 0, errors.New("create tables failed")
		}
	}
	// create user
	stmt, err = db.Prepare(`INSERT INTO users (email, full_name, phone, status, language_preference, policy_id, theme) VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	res, err := stmt.Exec(email, fullName, phoneNumber, 1, "en", 1, "light")
	stmt.Close()
	if err != nil {
		logrus.Error("insert user error ", err.Error())
		return 0, errors.New("insert users failed")
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}
