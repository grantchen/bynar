/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package repository

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

// SelectAccount query accounts,oraginzation_accounts,organizations,tenants columns to generate idToken of
// google identify platform
func (r *accountRepositoryHandler) SelectAccount(email string) (*model.SignIn, error) {
	var querySql = `
		select a.uid,oa.organization_user_uid,os.status as oraginzation_status,
		       os.organization_uuid,t.tenant_uuid 
		from accounts a 
		join oraginzation_accounts oa on a.id = oa.organization_user_id
		join organizations os on oa.organization_id=os.id
		join tenants t on os.tenant_id = t.id
		where a.email = ? and a.status =1 and a.verified = 1`
	var signIn = model.SignIn{}
	err := r.db.QueryRow(querySql, email).Scan(&signIn.Uid, &signIn.OrganizationUserId, &signIn.OrganizationStatus, &signIn.OrganizationUuid, &signIn.TenantUuid)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &signIn, nil
}
