package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
)

// GetOrganizationAccount gets the organization account.
func (r *accountRepositoryHandler) GetOrganizationAccount(language string, accountID int, organizationUuid string) (*model.GetOrganizationAccountResponse, error) {
	stmt, err := r.db.Prepare(`
		SELECT a.email,
			   a.full_name,
			   a.country,
			   a.address,
			   a.address_2,
			   a.city,
			   a.state,
			   a.postal_code,
			   a.phone,
			   o.description,
			   o.vat_number,
			   o.country
		from accounts a
				 INNER JOIN organization_accounts oa ON oa.organization_user_uid = a.uid
				 INNER JOIN organizations o ON o.id = oa.organization_id
		WHERE a.id = ?`,
	)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var res model.GetOrganizationAccountResponse
	err = stmt.QueryRow(accountID).Scan(
		&res.Email,
		&res.FullName,
		&res.Country,
		&res.AddressLine,
		&res.AddressLine2,
		&res.City,
		&res.State,
		&res.PostalCode,
		&res.PhoneNumber,
		&res.OrganizationName,
		&res.VAT,
		&res.OrganizationCountry,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateOrganizationAccount updates the organization account.
func (r *accountRepositoryHandler) UpdateOrganizationAccount(
	db *sql.DB,
	language string,
	accountID int,
	organizationUserId int,
	organizationUuid string,
	req model.OrganizationAccountRequest,
) error {
	updateUserStmt, err := db.Prepare(`
		UPDATE users
		SET email     = ?,
			full_name = ?,
			phone     = ?
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer updateUserStmt.Close()
	_, err = updateUserStmt.Exec(req.Email, req.FullName, req.PhoneNumber, organizationUserId)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateAccountStmt, err := tx.Prepare(`
		UPDATE accounts
		SET email       = ?,
			full_name   = ?,
			country     = ?,
			address     = ?,
			address_2   = ?,
			city        = ?,
			state       = ?,
			postal_code = ?,
			phone       = ?
		WHERE id = ?;`,
	)
	if err != nil {
		return err
	}
	defer updateAccountStmt.Close()
	_, err = updateAccountStmt.Exec(
		req.Email,
		req.FullName,
		req.Country,
		req.AddressLine,
		req.AddressLine2,
		req.City,
		req.State,
		req.PostalCode,
		req.PhoneNumber,
		accountID,
	)
	if err != nil {
		return err
	}

	updateOrgStmt, err := tx.Prepare(`
		UPDATE organizations
		SET description = ?,
			vat_number  = ?,
			country     = ?
		WHERE organization_uuid = ?
		LIMIT 1;`,
	)
	if err != nil {
		return err
	}
	defer updateOrgStmt.Close()
	_, err = updateOrgStmt.Exec(
		req.OrganizationName,
		req.VAT,
		req.OrganizationCountry,
		organizationUuid,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteOrganizationAccount deletes the organization account.
func (r *accountRepositoryHandler) DeleteOrganizationAccount(db *sql.DB, language string, tenantUuid string, organizationUuid string) error {
	var organizationID int
	var tenantID int
	queryOrganizationStmt, err := r.db.Prepare(`SELECT id, tenant_id FROM organizations WHERE organization_uuid = ? LIMIT 1;`)
	if err != nil {
		return err
	}
	defer queryOrganizationStmt.Close()

	err = queryOrganizationStmt.QueryRow(organizationUuid).Scan(&organizationID, &tenantID)
	if err != nil {
		return err
	}

	err = r.deleteEnvironment(language, tenantUuid, organizationUuid)
	if err != nil {
		return err
	}

	err = r.deleteAccount(language, organizationID, tenantID)
	if err != nil {
		return err
	}

	return nil
}

// IsCanDeleteOrganizationAccount checks if the organization account can be deleted.
func (r *accountRepositoryHandler) IsCanDeleteOrganizationAccount(language string, organizationUuid string) error {
	queryInvoiceStmt, err := r.db.Prepare(`
		SELECT 1
		FROM invoices i
		WHERE account_id IN (SELECT a.id
							 FROM accounts a
									  INNER JOIN organization_accounts oa ON oa.organization_user_uid = a.uid
									  INNER JOIN organizations o ON o.id = oa.organization_id
							 WHERE o.organization_uuid = ?)
		  AND i.paid = 0
		LIMIT 1`)
	if err != nil {
		return err
	}
	defer queryInvoiceStmt.Close()

	var existFlag int
	if err = queryInvoiceStmt.QueryRow(organizationUuid).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	return i18n.TranslationI18n(language, "CannotDeleteAccountWithUnpaidInvoices", nil)
}

// GetOrganizationIdByUuid gets the organization id by uuid.
func (r *accountRepositoryHandler) GetOrganizationIdByUuid(language string, organizationUuid string) (int, error) {
	stmt, err := r.db.Prepare(`SELECT id FROM organizations WHERE organization_uuid = ? LIMIT 1;`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var organizationID int
	err = stmt.QueryRow(organizationUuid).Scan(&organizationID)
	if err != nil {
		return 0, err
	}

	return organizationID, nil
}

// GetCustomerIDsByOrganizationID gets the customer ids of checkout users by organization id.
func (r *accountRepositoryHandler) GetCustomerIDsByOrganizationID(language string, organizationID int) ([]string, error) {
	// pluck user_payment_gateway_id from accounts_cards by user_id(which associated to accounts.id by organization_accounts)
	stmt, err := r.db.Prepare(`
		SELECT DISTINCT ac.user_payment_gateway_id
		FROM accounts_cards ac
				 INNER JOIN accounts a ON a.id = ac.user_id
				 INNER JOIN organization_accounts oa ON oa.organization_user_uid = a.uid
		WHERE oa.organization_id = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerIDs []string
	for rows.Next() {
		var customerID string
		err = rows.Scan(&customerID)
		if err != nil {
			return nil, err
		}
		customerIDs = append(customerIDs, customerID)
	}

	return customerIDs, nil
}

// GetGipUserUidsByOrganizationID gets the gip user uids by organization id.
func (r *accountRepositoryHandler) GetGipUserUidsByOrganizationID(language string, organizationID int) ([]string, error) {
	// pluck organization_user_uid from organization_accounts by organization_id
	stmt, err := r.db.Prepare(`SELECT organization_user_uid FROM organization_accounts WHERE organization_id = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gipUserUids []string
	for rows.Next() {
		var gipUserUid string
		err = rows.Scan(&gipUserUid)
		if err != nil {
			return nil, err
		}
		gipUserUids = append(gipUserUids, gipUserUid)
	}

	return gipUserUids, nil
}

// deleteAccount deletes the account data of the organization in accounts_management.
func (r *accountRepositoryHandler) deleteAccount(language string, organizationID int, tenantID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// delete tenants_management by organization_id and tenant_id
	deleteTenantsManagementStmt, err := tx.Prepare(`DELETE FROM tenants_management WHERE organization_id = ? AND tenant_id = ?;`)
	if err != nil {
		return err
	}
	defer deleteTenantsManagementStmt.Close()
	_, err = deleteTenantsManagementStmt.Exec(organizationID, tenantID)
	if err != nil {
		return err
	}

	// decrease organizations of tenants by tenant_id
	decreaseOrganizationsOfTenantsStmt, err := tx.Prepare(`UPDATE tenants SET organizations = organizations - 1 WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer decreaseOrganizationsOfTenantsStmt.Close()
	_, err = decreaseOrganizationsOfTenantsStmt.Exec(tenantID)
	if err != nil {
		return err
	}

	// delete accounts_cards by user_id(account_id)
	deleteAccountsCardsStmt, err := tx.Prepare(`
		DELETE
		FROM accounts_cards
		WHERE user_id IN (SELECT a.id
						  FROM accounts a
								   INNER JOIN organization_accounts oa ON oa.organization_user_uid = a.uid
						  WHERE oa.organization_id = ?);`)
	if err != nil {
		return err
	}
	defer deleteAccountsCardsStmt.Close()
	_, err = deleteAccountsCardsStmt.Exec(organizationID)
	if err != nil {
		return err
	}

	// delete invoices by account_id
	deleteInvoicesStmt, err := tx.Prepare(`
		DELETE
		FROM invoices
		WHERE account_id IN (SELECT a.id
							 FROM accounts a
									  INNER JOIN organization_accounts oa ON oa.organization_user_uid = a.uid
							 WHERE oa.organization_id = ?);`)
	if err != nil {
		return err
	}
	defer deleteInvoicesStmt.Close()
	_, err = deleteInvoicesStmt.Exec(organizationID)
	if err != nil {
		return err
	}

	// delete accounts by uid associated to organization_accounts' organization_user_uid
	deleteAccountsStmt, err := tx.Prepare(`
		DELETE
		FROM accounts
		WHERE uid IN (SELECT oa.organization_user_uid
					  FROM organization_accounts oa
					  WHERE oa.organization_id = ?);`)
	if err != nil {
		return err
	}
	defer deleteAccountsStmt.Close()
	_, err = deleteAccountsStmt.Exec(organizationID)
	if err != nil {
		return err
	}

	// delete organization_accounts by organization_id
	deleteOrganizationAccountsStmt, err := tx.Prepare(`DELETE FROM organization_accounts WHERE organization_id = ?;`)
	if err != nil {
		return err
	}
	defer deleteOrganizationAccountsStmt.Close()
	_, err = deleteOrganizationAccountsStmt.Exec(organizationID)
	if err != nil {
		return err
	}

	// delete organizations by id
	deleteOrganizationsStmt, err := tx.Prepare(`DELETE FROM organizations WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer deleteOrganizationsStmt.Close()
	_, err = deleteOrganizationsStmt.Exec(organizationID)
	if err != nil {
		return err
	}

	// commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// deleteEnvironment deletes the environment of the organization.
func (r *accountRepositoryHandler) deleteEnvironment(language, tenantUuid, organizationUuid string) error {
	db, err := r.getEnvironmentDB(tenantUuid)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP SCHEMA `%s`", organizationUuid))
	if err != nil {
		return err
	}

	return nil
}
