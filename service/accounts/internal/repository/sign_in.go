/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package repository

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// SelectSignInColumns query accounts,organization_accounts,organizations,tenants columns to generate idToken of
// Google identify platform
func (r *accountRepositoryHandler) SelectSignInColumns(email string) (*model.SignIn, error) {
	var querySql = `
		select a.uid,
		       oa.organization_user_id,
		       oa.oraginzation_main_account as organization_main_account,
		       os.status as organization_status,
               os.verified as organization_verified,
		       os.organization_uuid,
		       t.tenant_uuid,
               tm.status as tenant_status,
               tm.suspended as tenant_suspended
		from accounts a 
		join organization_accounts oa on a.uid = oa.organization_user_uid
		join organization_accounts oa on a.id = oa.organization_user_id
		join organizations os on oa.organization_id=os.id
		join tenants t on os.tenant_id = t.id
        join tenants_management tm on tm.organization_id = os.id and tm.tenant_id = t.id
		where a.email = ? and a.status = ? and a.verified = ? `
	var signIn = model.SignIn{}
	err := r.db.QueryRow(querySql, email, true, true).Scan(&signIn.Uid,
		&signIn.OrganizationUserId, &signIn.OrganizationAccount,
		&signIn.OrganizationStatus, &signIn.OrganizationVerified, &signIn.OrganizationUuid,
		&signIn.TenantUuid, &signIn.TenantStatus, &signIn.TenantSuspended)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}

	// Query 'users' table
	db, err := sql_db.InitializeConnection(utils.GenerateOrganizationConnection(signIn.TenantUuid, signIn.OrganizationUuid))
	if err != nil {
		return nil, errors.NewUnknownError("query user language preference fail").WithInternal().WithCause(err)
	}
	defer db.Close()
	querySql = `select coalesce(language_preference,'') from users where email = ? and status = ?`
	if err = db.QueryRow(querySql, email, true).Scan(&signIn.Language, &signIn.Theme); err != nil {
		return nil, errors.NewUnknownError("query user language preference fail").WithInternal().WithCause(err)
	}

	return &signIn, nil
}
