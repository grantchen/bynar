package repository

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
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

func (r *accountRepositoryHandler) UpdateOrganizationAccount(language string, accountID int, organizationUuid string, req model.OrganizationAccountRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	accountStmt, err := tx.Prepare(`
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

	defer accountStmt.Close()
	_, err = accountStmt.Exec(
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

	orgStmt, err := tx.Prepare(`
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

	defer orgStmt.Close()
	_, err = orgStmt.Exec(
		req.OrganizationName,
		req.VAT,
		req.OrganizationCountry,
		organizationUuid,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *accountRepositoryHandler) DeleteOrganizationAccount(language string, accountID int, organizationUuid string) error {
	//TODO implement me
	panic("implement me")
}
