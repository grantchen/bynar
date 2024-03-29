/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package repository

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model/organization_schema"
)

// GetOrganizationDetail query organization from database by organizationUuid
func (r *accountRepositoryHandler) GetOrganizationDetail(organizationUuid string) (*model.Organization, error) {
	var querySql = `
		select a.id,coalesce(a.description,''),coalesce(a.vat_number,''),coalesce(a.country,''),
		       coalesce(a.data_sovereignty,''),coalesce(a.organization_uuid,''),a.tenant_id,a.status,a.verified
		from organizations a where a.organization_uuid = ? limit 1`
	var organization = model.Organization{}
	prepare, err := r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	err = prepare.QueryRow(organizationUuid).Scan(&organization.ID,
		&organization.Description, &organization.VatNumber,
		&organization.Country, &organization.DataSovereignty,
		&organization.OrganizationUuid, &organization.TenantId, &organization.Status, &organization.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &organization, nil
}

// GetUserAccountDetail get accounts detail by uid provided
func (r *accountRepositoryHandler) GetUserAccountDetail(uid string) (*model.Account, error) {
	var querySql = `
		select a.id,a.email,a.full_name,a.address,coalesce(a.address_2,''),a.phone,a.city,a.postal_code,a.country,a.state,a.status,a.uid,a.org_id,a.verified
		from accounts a where uid = ? and status = ? and verified = ? limit 1`
	var account = model.Account{}
	prepare, err := r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	err = prepare.QueryRow(uid, true, true).Scan(&account.ID,
		&account.Email, &account.FullName, &account.Address, &account.Address2, &account.Phone,
		&account.City, &account.PostalCode, &account.Country, &account.State,
		&account.Status, &account.UID, &account.OrgID, &account.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &account, nil
}

// UpdateUserLanguagePreference Update user language preference
func (r *accountRepositoryHandler) UpdateUserLanguagePreference(db *sql.DB, userId int, languagePreference string) error {
	var updateSql = `update users set language_preference = ? where id = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	if _, err = prepare.Exec(languagePreference, userId); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}

	return nil
}

// UpdateUserThemePreference Update user theme preference
func (r *accountRepositoryHandler) UpdateUserThemePreference(db *sql.DB, userId int, themePreference string) error {
	var updateSql = `update users set theme = ? where id = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	if _, err = prepare.Exec(themePreference, userId); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}

	return nil
}

// UpdateProfilePhotoOfUsers update column profile_photo in table users of organization_schema(uuid)
func (r *accountRepositoryHandler) UpdateProfilePhotoOfUsers(db *sql.DB, userId int, profilePhoto string) error {
	updateSql := `UPDATE users SET profile_photo = ? WHERE id = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	if _, err = prepare.Exec(profilePhoto, userId); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}
	return nil
}

// GetUserDetail get user details from organization_schema(uuid)
func (r *accountRepositoryHandler) GetUserDetail(db *sql.DB, userId int) (*organization_schema.User, error) {
	var querySql = `select a.id,
       a.email,
       coalesce(a.full_name,''),
       coalesce(a.phone,''),
       a.status,
       coalesce(a.language_preference,''),
       coalesce(a.theme,''),
       coalesce(a.profile_photo,''),
	   coalesce(a.policies,'')
		from users a
		where a.id = ? and status = ? limit 1`
	var user = organization_schema.User{}
	prepare, err := db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	err = prepare.QueryRow(userId, true).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Phone, &user.Status,
		&user.LanguagePreference, &user.Theme, &user.ProfilePhoto, &user.Policies)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &user, nil
}

// UpdateUserProfile update user profile
func (r *accountRepositoryHandler) UpdateUserProfile(db *sql.DB, userId int, uid string, req model.UpdateUserProfileRequest) error {
	updateSql := `update users set email = ?,full_name = ?,phone = ?,language_preference = ?,theme = ? where id = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer func(prepare *sql.Stmt) {
		_ = prepare.Close()
	}(prepare)
	_, err = prepare.Exec(req.Email, req.FullName, req.PhoneNumber, req.Language, req.Theme, userId)
	if err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}

	updateAccountSql := `update accounts set email = ?,full_name = ?,phone = ? where uid = ?`
	stmt, err := r.db.Prepare(updateAccountSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(req.Email, req.FullName, req.PhoneNumber, uid)
	if err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}
	return nil
}
