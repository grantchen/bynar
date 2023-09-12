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

// SelectSignInColumns query accounts,oraginzation_accounts,organizations,tenants columns to generate idToken of
// google identify platform
func (r *accountRepositoryHandler) SelectSignInColumns(email string) (*model.SignIn, error) {
	var querySql = `
		select a.uid,
		       oa.organization_user_uid,
		       oa.oraginzation_main_account as organization_main_account,
		       os.status as organization_status,
		       os.organization_uuid,
		       t.tenant_uuid
		from accounts a 
		join oraginzation_accounts oa on a.id = oa.organization_user_id
		join organizations os on oa.organization_id=os.id
		join tenants t on os.tenant_id = t.id
		where a.email = ? and a.status = true and a.verified = true `
	var signIn = model.SignIn{}
	err := r.db.QueryRow(querySql, email).Scan(&signIn.Uid,
		&signIn.OrganizationUserId, &signIn.OrganizationMainAccount,
		&signIn.OrganizationStatus, &signIn.OrganizationUuid,
		&signIn.TenantUuid)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &signIn, nil
}

// SelectUserByUid select accounts by uid
func (r *accountRepositoryHandler) SelectUserByUid(uid string) (*model.GetUserResponse, error) {
	var querySql = `select a.email,a.full_name from accounts a where a.uid = ? limit 1`
	var user = model.GetUserResponse{}
	err := r.db.QueryRow(querySql, uid).Scan(
		&user.Email, &user.FullName)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &user, nil
}
