/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package repository

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

// GetOrganizationDetail query organization from database by organizationUuid
func (r *accountRepositoryHandler) GetOrganizationDetail(organizationUuid string) (*model.Organization, error) {
	var querySql = `
		select a.id,a.description,a.vat_number,a.country,a.data_sovereignty,a.organization_uuid,a.tenant_id,a.status,a.verified 
		from organizations a where a.organization_uuid = ? limit 1`
	var organization = model.Organization{}
	err := r.db.QueryRow(querySql, organizationUuid).Scan(&organization.ID,
		&organization.Description, &organization.VatNumber,
		&organization.Country, &organization.DataSovereignty,
		&organization.OrganizationUuid, &organization.TenantId, &organization.Status, &organization.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &organization, nil
}

// GetUserAccountDetail get accounts detail by uid provided
func (r *accountRepositoryHandler) GetUserAccountDetail(email string) (*model.Account, error) {
	var querySql = `
		select a.id,a.email,a.full_name,a.address,a.address_2,a.phone,a.city,a.postal_code,a.country,a.state,a.status,a.uid,a.org_id,a.verified 
		from accounts a where email = ? and status = ? and verified = ? limit 1`
	var account = model.Account{}
	err := r.db.QueryRow(querySql, email, true, true).Scan(&account.ID,
		&account.Email, &account.FullName, &account.Address, &account.Address2, &account.Phone,
		&account.City, &account.PostalCode, &account.Country, &account.State,
		&account.Status, &account.UID, &account.OrgID, &account.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &account, nil
}

// Update user language preference
func (r *accountRepositoryHandler) UpdateUserLanguagePreference(email, languagePreference string) error {
	var querySql = `update users set language_preference = ? where email = ?`
	if _, err := r.db.Exec(querySql, languagePreference, email); err != nil {
		return err
	}
	return nil
}
