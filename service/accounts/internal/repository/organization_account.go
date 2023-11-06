package repository

import (
	"database/sql"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
)

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

func (r *accountRepositoryHandler) DeleteOrganizationAccount(db *sql.DB, language string, tenantUuid string, organizationUuid string) error {
	err := r.isCanDeleteOrganizationAccount(language, organizationUuid)
	if err != nil {
		return err
	}

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

	err = r.deleteEnvironment(language, tenantUuid)
	if err != nil {
		return err
	}

	err = r.deleteAccount(language, organizationID, tenantID)
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepositoryHandler) isCanDeleteOrganizationAccount(language string, organizationUuid string) error {
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

	// commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// deleteEnvironment deletes the environment of the organization.
func (r *accountRepositoryHandler) deleteEnvironment(language string, tenantUuid string) error {
	db, err := r.getEnvironmentDB(tenantUuid)
	if err != nil {
		return err
	}
	defer db.Close()

	dropDatabaseStmt, err := db.Prepare(`DROP DATABASE IF EXISTS ?;`)
	if err != nil {
		return err
	}
	defer dropDatabaseStmt.Close()

	_, err = dropDatabaseStmt.Exec(tenantUuid)
	if err != nil {
		return err
	}

	return nil
}
