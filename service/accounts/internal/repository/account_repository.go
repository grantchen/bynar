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
