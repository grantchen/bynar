/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package repository

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// SelectSignInColumns query accounts,organization_accounts,organizations,tenants columns to generate idToken of
// Google identify platform
func (r *accountRepositoryHandler) SelectSignInColumns(uid string) (*model.SignIn, error) {
	var querySql = `
		select oa.organization_user_uid as uid,
		       oa.organization_user_id,
		       oa.oraginzation_main_account as organization_main_account,
		       os.status as organization_status,
               os.verified as organization_verified,
		       os.organization_uuid,
		       t.tenant_uuid,
               tm.status as tenant_status,
               tm.suspended as tenant_suspended
        from organization_accounts oa
		join organizations os on oa.organization_id=os.id
		join tenants t on os.tenant_id = t.id
        join tenants_management tm on tm.organization_id = os.id and tm.tenant_id = t.id
		where oa.organization_user_uid = ? `
	var signIn = model.SignIn{}
	prepare, err := r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	err = prepare.QueryRow(uid).Scan(&signIn.Uid,
		&signIn.OrganizationUserId, &signIn.OrganizationAccount,
		&signIn.OrganizationStatus, &signIn.OrganizationVerified, &signIn.OrganizationUuid,
		&signIn.TenantUuid, &signIn.TenantStatus, &signIn.TenantSuspended)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}

	// Query 'users' table
	db, err := sqldb.InitializeConnection(utils.GenerateOrganizationConnection(signIn.TenantUuid, signIn.OrganizationUuid))
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	querySql = `select coalesce(language_preference,''),coalesce(theme,'') from users where id = ? and status = ?`
	stmt, err := db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	if err = stmt.QueryRow(signIn.OrganizationUserId, true).Scan(&signIn.Language, &signIn.Theme); err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}

	// Query 'accounts_manager.accounts' table
	querySql = `select id from accounts where uid = ?`
	stmt, err = r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	if err = stmt.QueryRow(uid).Scan(&signIn.AccountId); err != nil {
		// user not exists in accounts if user added in user list
		logrus.Error("query row: [%w]", err)
		// return nil, fmt.Errorf("query row: [%w]", err)
		err = r.db.QueryRow(`SELECT a.id FROM organization_accounts oa1 JOIN organization_accounts oa2 ON oa1.organization_id = oa2.organization_id JOIN accounts a ON oa2.organization_user_uid = a.uid WHERE oa1.organization_user_uid = ? ORDER BY oa2.id ASC LIMIT 1`, uid).Scan(&signIn.AccountId)
		if err != nil {
			logrus.Error("query row: [%w]", err)
		}
	}

	return &signIn, nil
}
